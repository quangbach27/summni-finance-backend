package config_test

import (
	"sumni-finance-backend/internal/config"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig_Singleton(t *testing.T) {
	const goroutines = 20
	var wg sync.WaitGroup
	wg.Add(goroutines)

	instances := make([]*config.Config, goroutines)

	for i := 0; i < goroutines; i++ {
		go func(index int) {
			defer wg.Done()
			instances[index] = config.GetConfig()
		}(i)
	}

	wg.Wait()

	// All instance must point to the same object
	first := instances[0]
	for i, cfg := range instances {
		assert.Equal(t, first, cfg, "Instance %d is not equal to the first instance", i)
	}
}
