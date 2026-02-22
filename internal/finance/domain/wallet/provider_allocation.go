package wallet

import (
	"sumni-finance-backend/internal/common/validator"
	"sumni-finance-backend/internal/common/valueobject"
	"sumni-finance-backend/internal/finance/domain/fundprovider"
)

type ProviderAllocation struct {
	provider  *fundprovider.FundProvider
	allocated valueobject.Money
}

func NewProviderAllocation(
	fundProvider *fundprovider.FundProvider,
	allocatedAmount int64,
) (ProviderAllocation, error) {
	v := validator.New()

	v.Check(fundProvider != nil, "fundProvider", "fundProvider is required")
	v.Check(allocatedAmount >= 0, "allocated", "allocated must be greater or equal 0")

	if err := v.Err(); err != nil {
		return ProviderAllocation{}, err
	}

	allocated, err := valueobject.NewMoney(allocatedAmount, fundProvider.Currency())
	if err != nil {
		return ProviderAllocation{}, err
	}

	return ProviderAllocation{
		provider:  fundProvider,
		allocated: allocated,
	}, nil
}

func (pa ProviderAllocation) Provider() *fundprovider.FundProvider { return pa.provider }
func (pa ProviderAllocation) Allocated() valueobject.Money         { return pa.allocated }
