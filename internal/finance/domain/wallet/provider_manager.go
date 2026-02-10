package wallet

import (
	"errors"
	"fmt"
	"sumni-finance-backend/internal/common/validator"
	"sumni-finance-backend/internal/common/valueobject"
	"sumni-finance-backend/internal/finance/domain/fundprovider"

	"github.com/google/uuid"
)

type ProviderAllocation struct {
	fundProvider *fundprovider.FundProvider
	allocated    valueobject.Money
}

func (pa ProviderAllocation) FundProvider() *fundprovider.FundProvider {
	return pa.fundProvider
}

func (pa ProviderAllocation) Allocated() valueobject.Money {
	return pa.allocated
}

func NewProviderAllocation(
	fundProvider *fundprovider.FundProvider,
	allocated valueobject.Money,
) (ProviderAllocation, error) {
	v := validator.New()

	v.Check(fundProvider != nil, "fundProvider", "fundProvider is required")
	v.Check(!allocated.IsZero(), "allocated", "allocated is required")

	if err := v.Err(); err != nil {
		return ProviderAllocation{}, err
	}

	return ProviderAllocation{
		fundProvider: fundProvider,
		allocated:    allocated,
	}, nil
}

type ProviderManager struct {
	providers map[uuid.UUID]ProviderAllocation
}

func (m *ProviderManager) GetFundProviderAllocations() []ProviderAllocation {
	providerAllocations := make([]ProviderAllocation, 0, len(m.providers))

	for _, provider := range m.providers {
		providerAllocations = append(providerAllocations, provider)
	}

	return providerAllocations
}

func NewProviderManager(allocations []ProviderAllocation) (*ProviderManager, error) {
	providers := make(map[uuid.UUID]ProviderAllocation, len(allocations))

	for _, allocation := range allocations {
		if allocation.fundProvider == nil {
			return nil, errors.New("fundProvider can not be nil")
		}

		if allocation.allocated.IsZero() {
			return nil, errors.New("allocated is required")
		}

		_, exist := providers[allocation.fundProvider.ID()]
		if exist {
			return nil, fmt.Errorf("fundProvider must be unique: %s", allocation.fundProvider.ID())
		}

		providers[allocation.fundProvider.ID()] = allocation
	}

	return &ProviderManager{
		providers: providers,
	}, nil
}

func (m ProviderManager) AddAndAllocate(
	fundProvider *fundprovider.FundProvider,
	allocated valueobject.Money,
) error {
	if fundProvider == nil || allocated.IsZero() {
		return errors.New("FundProvider or allocated is required")
	}

	if m.HasFundProvider(fundProvider.ID()) {
		return ErrFundProviderAlreadyRegistered
	}

	if err := fundProvider.Allocate(allocated); err != nil {
		return err
	}

	m.providers[fundProvider.ID()] = ProviderAllocation{
		fundProvider: fundProvider,
		allocated:    allocated,
	}

	return nil
}

func (m ProviderManager) HasFundProvider(fID uuid.UUID) bool {
	_, exist := m.providers[fID]
	return exist
}

func (m ProviderManager) GetFundProvider(fID uuid.UUID) *fundprovider.FundProvider {
	if providerAllocation, exist := m.providers[fID]; exist {
		return providerAllocation.fundProvider
	}

	return nil
}

func (m ProviderManager) CalculateTotalProviderAllocated() (valueobject.Money, error) {
	var total valueobject.Money
	var err error

	for _, allocation := range m.providers {
		if !total.IsZero() {
			total, err = total.Add(allocation.allocated)
			if err != nil {
				return valueobject.Money{}, err
			}

			continue
		}

		total = allocation.allocated
	}

	return total, nil
}
