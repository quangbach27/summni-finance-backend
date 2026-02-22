package wallet

import "github.com/google/uuid"

type ProviderAllocationSpec interface {
	IsSatisfiedBy(p ProviderAllocation) bool
}

type DefaultProviderAllocationSpec struct{}

func (spec DefaultProviderAllocationSpec) IsSatisfiedBy(p ProviderAllocation) bool {
	return true
}

type AllocationBelongsToAnyProviderSpec struct {
	allowed map[uuid.UUID]struct{}
}

func NewAllocationBelongsToAnyProviderSpec(providerIDs []uuid.UUID) AllocationBelongsToAnyProviderSpec {
	allowed := make(map[uuid.UUID]struct{})
	for _, id := range providerIDs {
		allowed[id] = struct{}{}
	}
	return AllocationBelongsToAnyProviderSpec{
		allowed: allowed,
	}
}

func (spec AllocationBelongsToAnyProviderSpec) IsSatisfiedBy(p ProviderAllocation) bool {
	provider := p.Provider()
	if provider == nil {
		return false
	}

	_, exists := spec.allowed[provider.ID()]
	return exists
}
