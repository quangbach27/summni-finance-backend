package fundprovider

import (
	"errors"
	"fmt"
	"sumni-finance-backend/internal/common/validator"
	"sumni-finance-backend/internal/common/valueobject"

	"github.com/google/uuid"
)

var (
	ErrInvalidAllocationEntry = errors.New("target wallet and a positive amount are required")

	ErrWalletAlreadyAllocated  = errors.New("this wallet is already allocated to this provider")
	ErrWalletNotAllocated      = errors.New("this wallet is not assigned to this provider")
	ErrAllocationLimitExceeded = errors.New("amount exceeds the wallet's current allocation")
)

// AllocationEntry
type AllocationEntry struct {
	WalletID uuid.UUID
	Amount   valueobject.Money
}

func (e AllocationEntry) IsZero() bool { return e == AllocationEntry{} }

func (e AllocationEntry) IsValid() error {
	v := validator.New()

	v.Check(e.WalletID != uuid.Nil, "WalletID", "WalletID is required")
	v.Check(!e.Amount.IsZero(), "Amount", "Amount is required")
	v.Check(e.Amount.Amount() >= 0, "amount", "Amount must be positive")
	v.Check(!e.Amount.Currency().IsZero(), "currency", "currency is required")

	if err := v.Err(); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidAllocationEntry, err)
	}

	return nil
}

// Allocation
type Allocation struct {
	totalAllocated valueobject.Money
	entries        map[uuid.UUID]AllocationEntry
}

func NewAllocation(
	currency valueobject.Currency,
	entries ...AllocationEntry,
) (Allocation, error) {
	if currency.IsZero() {
		return Allocation{}, errors.New("currency is required")
	}

	total, err := valueobject.NewMoney(0, currency)
	if err != nil {
		return Allocation{}, err
	}

	allocation := make(map[uuid.UUID]AllocationEntry, len(entries))
	for _, entry := range entries {
		if err := entry.IsValid(); err != nil {
			return Allocation{}, err
		}

		allocation[entry.WalletID] = entry
		total, err = total.Add(entry.Amount)
		if err != nil {
			return Allocation{}, err
		}
	}

	return Allocation{
		entries:        allocation,
		totalAllocated: total,
	}, nil
}

func (w Allocation) Entries() []AllocationEntry {
	walletLocationEntries := make([]AllocationEntry, 0, len(w.entries))

	for _, entry := range w.entries {
		walletLocationEntries = append(walletLocationEntries, entry)
	}

	return walletLocationEntries
}

func (w Allocation) TotalAllocated() valueobject.Money { return w.totalAllocated }

func (w Allocation) EntryOf(walletID uuid.UUID) (AllocationEntry, bool) {
	entry, exist := w.entries[walletID]
	return entry, exist
}

func (w Allocation) Allocate(
	walletID uuid.UUID,
	amount valueobject.Money,
) (Allocation, error) {
	if _, exist := w.EntryOf(walletID); exist {
		return Allocation{}, ErrWalletAlreadyAllocated
	}

	newAllocationEntry := AllocationEntry{WalletID: walletID, Amount: amount}
	if err := newAllocationEntry.IsValid(); err != nil {
		return Allocation{}, err
	}

	newTotalAllocated, err := w.totalAllocated.Add(amount)
	if err != nil {
		return Allocation{}, err
	}

	cloneEntries := w.cloneEntries()
	cloneEntries[walletID] = newAllocationEntry

	return Allocation{
		totalAllocated: newTotalAllocated,
		entries:        cloneEntries,
	}, nil
}

func (w Allocation) DecreaseAllocation(
	walletID uuid.UUID,
	amount valueobject.Money,
) (Allocation, error) {
	currentEntry, exist := w.EntryOf(walletID)
	if !exist {
		return Allocation{}, ErrWalletNotAllocated
	}

	// Validate amount before performing operations
	if amount.IsZero() || amount.Currency().IsZero() {
		return Allocation{}, ErrInvalidAllocationEntry
	}

	// Calculate new totalAllocated
	newTotalAllocated, err := w.totalAllocated.Subtract(amount)
	if err != nil {
		return Allocation{}, err
	}

	newEntryAmount, err := currentEntry.Amount.Subtract(amount)
	if err != nil {
		return Allocation{}, err
	}

	if newEntryAmount.IsNegative() {
		return Allocation{}, ErrAllocationLimitExceeded
	}

	newAllocationEntry := AllocationEntry{
		WalletID: walletID,
		Amount:   newEntryAmount,
	}
	if err = newAllocationEntry.IsValid(); err != nil {
		return Allocation{}, err
	}

	cloneEntries := w.cloneEntries()
	cloneEntries[walletID] = newAllocationEntry

	return Allocation{
		totalAllocated: newTotalAllocated,
		entries:        cloneEntries,
	}, nil
}

func (w Allocation) IncreaseAllocation(
	walletID uuid.UUID,
	amount valueobject.Money,
) (Allocation, error) {
	currentEntry, exist := w.EntryOf(walletID)
	if !exist {
		return Allocation{}, ErrWalletNotAllocated
	}

	newTotalAllocated, err := w.totalAllocated.Add(amount)
	if err != nil {
		return Allocation{}, err
	}

	newAmount, err := currentEntry.Amount.Add(amount)
	if err != nil {
		return Allocation{}, err
	}

	newAllocationEntry := AllocationEntry{
		WalletID: walletID,
		Amount:   newAmount,
	}
	if err := newAllocationEntry.IsValid(); err != nil {
		return Allocation{}, err
	}

	cloneEntries := w.cloneEntries()
	cloneEntries[walletID] = newAllocationEntry

	return Allocation{
		totalAllocated: newTotalAllocated,
		entries:        cloneEntries,
	}, nil
}

func (w Allocation) cloneEntries() map[uuid.UUID]AllocationEntry {
	newEntries := make(map[uuid.UUID]AllocationEntry, len(w.entries))
	for k, v := range w.entries {
		newEntries[k] = v
	}

	return newEntries
}
