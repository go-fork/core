package core_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.fork.vn/core"
	diMocks "go.fork.vn/di/mocks"
)

// TestApplication_Boot_OrderingValidation tests Boot method call ordering
func TestApplication_Boot_OrderingValidation(t *testing.T) {
	t.Run("boot_calls_follow_dependency_order_A_then_B", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		var bootOrder []string

		// Provider A: provides service A, no dependencies
		providerA := diMocks.NewMockServiceProvider(t)
		providerA.EXPECT().Providers().Return([]string{"service.a"}).Maybe()
		providerA.EXPECT().Requires().Return([]string{}).Maybe()
		providerA.EXPECT().Register(app).Once()
		providerA.EXPECT().Boot(app).Once().Run(func(args mock.Arguments) {
			bootOrder = append(bootOrder, "providerA")
		})

		// Provider B: provides service B, requires service A
		providerB := diMocks.NewMockServiceProvider(t)
		providerB.EXPECT().Providers().Return([]string{"service.b"}).Maybe()
		providerB.EXPECT().Requires().Return([]string{"service.a"}).Maybe()
		providerB.EXPECT().Register(app).Once()
		providerB.EXPECT().Boot(app).Once().Run(func(args mock.Arguments) {
			bootOrder = append(bootOrder, "providerB")
		})

		// Register B first, then A (to test sorting)
		app.Register(providerB)
		app.Register(providerA)

		// Register with dependencies first
		err := app.RegisterWithDependencies()
		assert.NoError(t, err)

		// Then boot providers
		err = app.BootServiceProviders()
		assert.NoError(t, err)

		// Verify boot order: A should be called before B
		assert.Equal(t, []string{"providerA", "providerB"}, bootOrder)

		providerA.AssertExpectations(t)
		providerB.AssertExpectations(t)
	})

	t.Run("boot_calls_follow_complex_dependency_chain_A_B_C", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		var bootOrder []string

		// Provider A: config service, no dependencies
		providerA := diMocks.NewMockServiceProvider(t)
		providerA.EXPECT().Providers().Return([]string{"config"}).Maybe()
		providerA.EXPECT().Requires().Return([]string{}).Maybe()
		providerA.EXPECT().Register(app).Once()
		providerA.EXPECT().Boot(app).Once().Run(func(args mock.Arguments) {
			bootOrder = append(bootOrder, "config")
		})

		// Provider B: database service, requires config
		providerB := diMocks.NewMockServiceProvider(t)
		providerB.EXPECT().Providers().Return([]string{"database"}).Maybe()
		providerB.EXPECT().Requires().Return([]string{"config"}).Maybe()
		providerB.EXPECT().Register(app).Once()
		providerB.EXPECT().Boot(app).Once().Run(func(args mock.Arguments) {
			bootOrder = append(bootOrder, "database")
		})

		// Provider C: auth service, requires database
		providerC := diMocks.NewMockServiceProvider(t)
		providerC.EXPECT().Providers().Return([]string{"auth"}).Maybe()
		providerC.EXPECT().Requires().Return([]string{"database"}).Maybe()
		providerC.EXPECT().Register(app).Once()
		providerC.EXPECT().Boot(app).Once().Run(func(args mock.Arguments) {
			bootOrder = append(bootOrder, "auth")
		})

		// Register in reverse order to test dependency sorting
		app.Register(providerC)
		app.Register(providerB)
		app.Register(providerA)

		err := app.RegisterWithDependencies()
		assert.NoError(t, err)

		err = app.BootServiceProviders()
		assert.NoError(t, err)

		// Verify boot order: config -> database -> auth
		assert.Equal(t, []string{"config", "database", "auth"}, bootOrder)

		providerA.AssertExpectations(t)
		providerB.AssertExpectations(t)
		providerC.AssertExpectations(t)
	})

	t.Run("boot_calls_handle_diamond_dependency_pattern", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		var bootOrder []string

		// Diamond pattern: base -> (featureA, featureB) -> combined
		// Base service (no dependencies)
		providerBase := diMocks.NewMockServiceProvider(t)
		providerBase.EXPECT().Providers().Return([]string{"base.service"}).Maybe()
		providerBase.EXPECT().Requires().Return([]string{}).Maybe()
		providerBase.EXPECT().Register(app).Once()
		providerBase.EXPECT().Boot(app).Once().Run(func(args mock.Arguments) {
			bootOrder = append(bootOrder, "base")
		})

		// Feature A (requires base)
		providerA := diMocks.NewMockServiceProvider(t)
		providerA.EXPECT().Providers().Return([]string{"feature.a"}).Maybe()
		providerA.EXPECT().Requires().Return([]string{"base.service"}).Maybe()
		providerA.EXPECT().Register(app).Once()
		providerA.EXPECT().Boot(app).Once().Run(func(args mock.Arguments) {
			bootOrder = append(bootOrder, "featureA")
		})

		// Feature B (requires base)
		providerB := diMocks.NewMockServiceProvider(t)
		providerB.EXPECT().Providers().Return([]string{"feature.b"}).Maybe()
		providerB.EXPECT().Requires().Return([]string{"base.service"}).Maybe()
		providerB.EXPECT().Register(app).Once()
		providerB.EXPECT().Boot(app).Once().Run(func(args mock.Arguments) {
			bootOrder = append(bootOrder, "featureB")
		})

		// Combined service (requires both features)
		providerCombined := diMocks.NewMockServiceProvider(t)
		providerCombined.EXPECT().Providers().Return([]string{"combined.service"}).Maybe()
		providerCombined.EXPECT().Requires().Return([]string{"feature.a", "feature.b"}).Maybe()
		providerCombined.EXPECT().Register(app).Once()
		providerCombined.EXPECT().Boot(app).Once().Run(func(args mock.Arguments) {
			bootOrder = append(bootOrder, "combined")
		})

		// Register in random order
		app.Register(providerCombined)
		app.Register(providerA)
		app.Register(providerBase)
		app.Register(providerB)

		err := app.RegisterWithDependencies()
		assert.NoError(t, err)

		err = app.BootServiceProviders()
		assert.NoError(t, err)

		// Verify boot order: base must be first, combined must be last
		assert.Equal(t, "base", bootOrder[0])
		assert.Equal(t, "combined", bootOrder[3])
		// featureA and featureB can be in any order after base but before combined
		assert.Contains(t, bootOrder[1:3], "featureA")
		assert.Contains(t, bootOrder[1:3], "featureB")

		providerBase.AssertExpectations(t)
		providerA.AssertExpectations(t)
		providerB.AssertExpectations(t)
		providerCombined.AssertExpectations(t)
	})

	t.Run("boot_calls_follow_realistic_web_app_dependencies", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		var bootOrder []string

		// Config (no dependencies)
		configProvider := diMocks.NewMockServiceProvider(t)
		configProvider.EXPECT().Providers().Return([]string{"config"}).Maybe()
		configProvider.EXPECT().Requires().Return([]string{}).Maybe()
		configProvider.EXPECT().Register(app).Once()
		configProvider.EXPECT().Boot(app).Once().Run(func(args mock.Arguments) {
			bootOrder = append(bootOrder, "config")
		})

		// Log (depends on config)
		logProvider := diMocks.NewMockServiceProvider(t)
		logProvider.EXPECT().Providers().Return([]string{"log"}).Maybe()
		logProvider.EXPECT().Requires().Return([]string{"config"}).Maybe()
		logProvider.EXPECT().Register(app).Once()
		logProvider.EXPECT().Boot(app).Once().Run(func(args mock.Arguments) {
			bootOrder = append(bootOrder, "log")
		})

		// Database (depends on config and log)
		dbProvider := diMocks.NewMockServiceProvider(t)
		dbProvider.EXPECT().Providers().Return([]string{"database"}).Maybe()
		dbProvider.EXPECT().Requires().Return([]string{"config", "log"}).Maybe()
		dbProvider.EXPECT().Register(app).Once()
		dbProvider.EXPECT().Boot(app).Once().Run(func(args mock.Arguments) {
			bootOrder = append(bootOrder, "database")
		})

		// Cache (depends on config)
		cacheProvider := diMocks.NewMockServiceProvider(t)
		cacheProvider.EXPECT().Providers().Return([]string{"cache"}).Maybe()
		cacheProvider.EXPECT().Requires().Return([]string{"config"}).Maybe()
		cacheProvider.EXPECT().Register(app).Once()
		cacheProvider.EXPECT().Boot(app).Once().Run(func(args mock.Arguments) {
			bootOrder = append(bootOrder, "cache")
		})

		// Auth (depends on database and cache)
		authProvider := diMocks.NewMockServiceProvider(t)
		authProvider.EXPECT().Providers().Return([]string{"auth"}).Maybe()
		authProvider.EXPECT().Requires().Return([]string{"database", "cache"}).Maybe()
		authProvider.EXPECT().Register(app).Once()
		authProvider.EXPECT().Boot(app).Once().Run(func(args mock.Arguments) {
			bootOrder = append(bootOrder, "auth")
		})

		// HTTP (depends on log and auth)
		httpProvider := diMocks.NewMockServiceProvider(t)
		httpProvider.EXPECT().Providers().Return([]string{"http"}).Maybe()
		httpProvider.EXPECT().Requires().Return([]string{"log", "auth"}).Maybe()
		httpProvider.EXPECT().Register(app).Once()
		httpProvider.EXPECT().Boot(app).Once().Run(func(args mock.Arguments) {
			bootOrder = append(bootOrder, "http")
		})

		// Register in random order
		app.Register(httpProvider)
		app.Register(authProvider)
		app.Register(cacheProvider)
		app.Register(dbProvider)
		app.Register(logProvider)
		app.Register(configProvider)

		err := app.RegisterWithDependencies()
		assert.NoError(t, err)

		err = app.BootServiceProviders()
		assert.NoError(t, err)

		// Verify dependency ordering constraints
		configIndex := indexOf(bootOrder, "config")
		logIndex := indexOf(bootOrder, "log")
		dbIndex := indexOf(bootOrder, "database")
		cacheIndex := indexOf(bootOrder, "cache")
		authIndex := indexOf(bootOrder, "auth")
		httpIndex := indexOf(bootOrder, "http")

		// Config must be first
		assert.Equal(t, 0, configIndex)

		// Log depends on config
		assert.True(t, logIndex > configIndex)

		// Database depends on config and log
		assert.True(t, dbIndex > configIndex)
		assert.True(t, dbIndex > logIndex)

		// Cache depends on config
		assert.True(t, cacheIndex > configIndex)

		// Auth depends on database and cache
		assert.True(t, authIndex > dbIndex)
		assert.True(t, authIndex > cacheIndex)

		// HTTP depends on log and auth
		assert.True(t, httpIndex > logIndex)
		assert.True(t, httpIndex > authIndex)

		// Verify all providers
		configProvider.AssertExpectations(t)
		logProvider.AssertExpectations(t)
		dbProvider.AssertExpectations(t)
		cacheProvider.AssertExpectations(t)
		authProvider.AssertExpectations(t)
		httpProvider.AssertExpectations(t)
	})
}

// TestApplication_Boot_Integration tests the Boot method integration
func TestApplication_Boot_Integration(t *testing.T) {
	t.Run("boot_method_handles_registration_and_booting", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		var operationOrder []string

		// Provider A: no dependencies
		providerA := diMocks.NewMockServiceProvider(t)
		providerA.EXPECT().Providers().Return([]string{"service.a"}).Maybe()
		providerA.EXPECT().Requires().Return([]string{}).Maybe()
		providerA.EXPECT().Register(app).Once().Run(func(args mock.Arguments) {
			operationOrder = append(operationOrder, "registerA")
		})
		providerA.EXPECT().Boot(app).Once().Run(func(args mock.Arguments) {
			operationOrder = append(operationOrder, "bootA")
		})

		// Provider B: depends on A
		providerB := diMocks.NewMockServiceProvider(t)
		providerB.EXPECT().Providers().Return([]string{"service.b"}).Maybe()
		providerB.EXPECT().Requires().Return([]string{"service.a"}).Maybe()
		providerB.EXPECT().Register(app).Once().Run(func(args mock.Arguments) {
			operationOrder = append(operationOrder, "registerB")
		})
		providerB.EXPECT().Boot(app).Once().Run(func(args mock.Arguments) {
			operationOrder = append(operationOrder, "bootB")
		})

		app.Register(providerB)
		app.Register(providerA)

		// Use Boot() method which should handle everything
		err := app.Boot()
		assert.NoError(t, err)

		// Verify operation order:
		// 1. All registrations happen first (in dependency order)
		// 2. Then all boots happen (in dependency order)
		expectedOrder := []string{"registerA", "registerB", "bootA", "bootB"}
		assert.Equal(t, expectedOrder, operationOrder)

		providerA.AssertExpectations(t)
		providerB.AssertExpectations(t)
	})

	t.Run("boot_method_detects_and_handles_dependency_chains", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		var operationOrder []string

		// Create a chain: A -> B -> C
		providerA := diMocks.NewMockServiceProvider(t)
		providerA.EXPECT().Providers().Return([]string{"service.a"}).Maybe()
		providerA.EXPECT().Requires().Return([]string{}).Maybe()
		providerA.EXPECT().Register(app).Once().Run(func(args mock.Arguments) {
			operationOrder = append(operationOrder, "regA")
		})
		providerA.EXPECT().Boot(app).Once().Run(func(args mock.Arguments) {
			operationOrder = append(operationOrder, "bootA")
		})

		providerB := diMocks.NewMockServiceProvider(t)
		providerB.EXPECT().Providers().Return([]string{"service.b"}).Maybe()
		providerB.EXPECT().Requires().Return([]string{"service.a"}).Maybe()
		providerB.EXPECT().Register(app).Once().Run(func(args mock.Arguments) {
			operationOrder = append(operationOrder, "regB")
		})
		providerB.EXPECT().Boot(app).Once().Run(func(args mock.Arguments) {
			operationOrder = append(operationOrder, "bootB")
		})

		providerC := diMocks.NewMockServiceProvider(t)
		providerC.EXPECT().Providers().Return([]string{"service.c"}).Maybe()
		providerC.EXPECT().Requires().Return([]string{"service.b"}).Maybe()
		providerC.EXPECT().Register(app).Once().Run(func(args mock.Arguments) {
			operationOrder = append(operationOrder, "regC")
		})
		providerC.EXPECT().Boot(app).Once().Run(func(args mock.Arguments) {
			operationOrder = append(operationOrder, "bootC")
		})

		// Register in reverse order
		app.Register(providerC)
		app.Register(providerB)
		app.Register(providerA)

		err := app.Boot()
		assert.NoError(t, err)

		// Should detect dependencies and use RegisterWithDependencies
		// Registrations: A, B, C (dependency order)
		// Boots: A, B, C (dependency order)
		expectedOrder := []string{"regA", "regB", "regC", "bootA", "bootB", "bootC"}
		assert.Equal(t, expectedOrder, operationOrder)

		providerA.AssertExpectations(t)
		providerB.AssertExpectations(t)
		providerC.AssertExpectations(t)
	})
}

// Helper function to find index of element in slice
func indexOf(slice []string, item string) int {
	for i, v := range slice {
		if v == item {
			return i
		}
	}
	return -1
}
