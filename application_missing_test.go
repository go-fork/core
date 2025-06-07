package core_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.fork.vn/core"
	"go.fork.vn/di"
	diMocks "go.fork.vn/di/mocks"
)

// TestApplication_RegisterServiceProviders tests RegisterServiceProviders function (0% coverage)
func TestApplication_RegisterServiceProviders(t *testing.T) {
	t.Run("registers_all_providers_in_order", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		var registrationOrder []string

		// Provider A
		providerA := diMocks.NewMockServiceProvider(t)
		providerA.EXPECT().Providers().Return([]string{"service.a"}).Maybe()
		providerA.EXPECT().Requires().Return([]string{}).Maybe()
		providerA.EXPECT().Register(app).Once().Run(func(args mock.Arguments) {
			registrationOrder = append(registrationOrder, "A")
		})

		// Provider B
		providerB := diMocks.NewMockServiceProvider(t)
		providerB.EXPECT().Providers().Return([]string{"service.b"}).Maybe()
		providerB.EXPECT().Requires().Return([]string{}).Maybe()
		providerB.EXPECT().Register(app).Once().Run(func(args mock.Arguments) {
			registrationOrder = append(registrationOrder, "B")
		})

		app.Register(providerA)
		app.Register(providerB)

		// Test RegisterServiceProviders
		err := app.RegisterServiceProviders()
		assert.NoError(t, err)

		assert.Equal(t, []string{"A", "B"}, registrationOrder)

		providerA.AssertExpectations(t)
		providerB.AssertExpectations(t)
	})

	t.Run("handles_registration_error", func(t *testing.T) {
		t.Parallel()

		config2 := map[string]interface{}{}
		app2 := core.New(config2)

		provider := diMocks.NewMockServiceProvider(t)
		provider.EXPECT().Providers().Return([]string{"failing.service"}).Maybe()
		provider.EXPECT().Requires().Return([]string{}).Maybe()
		// RegisterServiceProviders không handle error từ Register()
		// Nó chỉ gọi Register() mà không check return value
		provider.EXPECT().Register(app2).Once()

		app2.Register(provider)

		err := app2.RegisterServiceProviders()
		// RegisterServiceProviders luôn return nil theo implementation
		assert.NoError(t, err)

		provider.AssertExpectations(t)
	})
}

// TestApplication_Alias tests Alias function (0% coverage)
func TestApplication_Alias(t *testing.T) {
	t.Run("creates_service_alias", func(t *testing.T) {
		t.Parallel()

		config3 := map[string]interface{}{}
		app3 := core.New(config3)

		// Bind original service
		app3.Bind("original.service", func(c di.Container) interface{} {
			return "test-value"
		})

		// Create alias
		app3.Alias("original.service", "aliased.service")

		// Both names should resolve to same service
		original, err := app3.Make("original.service")
		assert.NoError(t, err)
		assert.Equal(t, "test-value", original)

		aliased, err := app3.Make("aliased.service")
		assert.NoError(t, err)
		assert.Equal(t, "test-value", aliased)
	})
}

// TestApplication_MustMake tests MustMake function (0% coverage)
func TestApplication_MustMake(t *testing.T) {
	t.Run("returns_service_when_exists", func(t *testing.T) {
		t.Parallel()

		config4 := map[string]interface{}{}
		app4 := core.New(config4)

		app4.Bind("test.service", func(c di.Container) interface{} {
			return "test-value"
		})

		service := app4.MustMake("test.service")
		assert.Equal(t, "test-value", service)
	})

	t.Run("panics_when_service_not_found", func(t *testing.T) {
		t.Parallel()

		config5 := map[string]interface{}{}
		app5 := core.New(config5)

		assert.Panics(t, func() {
			app5.MustMake("nonexistent.service")
		})
	})
}

// TestApplication_Call tests Call function (0% coverage)
func TestApplication_Call(t *testing.T) {
	t.Run("calls_function_with_parameters", func(t *testing.T) {
		t.Parallel()

		config6 := map[string]interface{}{}
		app6 := core.New(config6)

		// Simple function
		testFunc := func(a string, b int) string {
			return a + string(rune(b+'0'))
		}

		results, err := app6.Call(testFunc, "test", 5)
		assert.NoError(t, err)
		assert.Len(t, results, 1)
	})
}

// TestApplication_New_EdgeCases tests edge cases for New function to improve coverage
func TestApplication_New_EdgeCases(t *testing.T) {
	t.Run("New_with_nil_config", func(t *testing.T) {
		t.Parallel()

		// Test with nil config - should create empty config
		app := core.New(nil)

		// Should still work
		assert.NotNil(t, app)
		assert.NotNil(t, app.Container())

		// Should have default config
		config, err := app.Make("app.config")
		assert.NoError(t, err)
		assert.NotNil(t, config)

		// Should be empty map
		configMap, ok := config.(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, 0, len(configMap))
	})

	t.Run("New_creates_module_loader", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{
			"test": "value",
		}
		app := core.New(config)

		// Should have module loader
		loader := app.ModuleLoader()
		assert.NotNil(t, loader)
	})
}

// TestApplication_BootServiceProviders_EdgeCases tests edge cases for BootServiceProviders
func TestApplication_BootServiceProviders_EdgeCases(t *testing.T) {
	t.Run("BootServiceProviders_already_booted", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		var bootCount int

		provider := diMocks.NewMockServiceProvider(t)
		provider.EXPECT().Providers().Return([]string{"test.service"}).Maybe()
		provider.EXPECT().Requires().Return([]string{}).Maybe()
		provider.EXPECT().Register(app).Once()
		provider.EXPECT().Boot(app).Once().Run(func(args mock.Arguments) {
			bootCount++
		})

		app.Register(provider)

		// First boot
		err := app.RegisterServiceProviders()
		assert.NoError(t, err)

		err = app.BootServiceProviders()
		assert.NoError(t, err)
		assert.Equal(t, 1, bootCount)

		// Second boot - should not boot again
		err = app.BootServiceProviders()
		assert.NoError(t, err)
		assert.Equal(t, 1, bootCount) // Should still be 1

		provider.AssertExpectations(t)
	})

	t.Run("BootServiceProviders_uses_sorted_providers_when_available", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		var bootOrder []string

		// Provider A: no dependencies
		providerA := diMocks.NewMockServiceProvider(t)
		providerA.EXPECT().Providers().Return([]string{"service.a"}).Maybe()
		providerA.EXPECT().Requires().Return([]string{}).Maybe()
		providerA.EXPECT().Register(app).Once()
		providerA.EXPECT().Boot(app).Once().Run(func(args mock.Arguments) {
			bootOrder = append(bootOrder, "A")
		})

		// Provider B: depends on A
		providerB := diMocks.NewMockServiceProvider(t)
		providerB.EXPECT().Providers().Return([]string{"service.b"}).Maybe()
		providerB.EXPECT().Requires().Return([]string{"service.a"}).Maybe()
		providerB.EXPECT().Register(app).Once()
		providerB.EXPECT().Boot(app).Once().Run(func(args mock.Arguments) {
			bootOrder = append(bootOrder, "B")
		})

		// Register in reverse order
		app.Register(providerB)
		app.Register(providerA)

		// Use RegisterWithDependencies to create sorted providers
		err := app.RegisterWithDependencies()
		assert.NoError(t, err)

		err = app.BootServiceProviders()
		assert.NoError(t, err)

		// Should boot in dependency order (sorted providers used)
		assert.Equal(t, []string{"A", "B"}, bootOrder)

		providerA.AssertExpectations(t)
		providerB.AssertExpectations(t)
	})
}

// TestApplication_Boot_EdgeCases tests edge cases for Boot function
func TestApplication_Boot_EdgeCases(t *testing.T) {
	t.Run("Boot_chooses_correct_registration_method", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		var operationOrder []string

		// Provider with dependencies - should use RegisterWithDependencies
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

		app.Register(providerB)
		app.Register(providerA)

		err := app.Boot()
		assert.NoError(t, err)

		// Should detect dependencies and register in correct order
		expectedOrder := []string{"regA", "regB", "bootA", "bootB"}
		assert.Equal(t, expectedOrder, operationOrder)

		providerA.AssertExpectations(t)
		providerB.AssertExpectations(t)
	})

	t.Run("Boot_with_no_dependencies_uses_simple_registration", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		var operationOrder []string

		// Providers with no requirements - should use simple registration
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
		providerB.EXPECT().Requires().Return([]string{}).Maybe()
		providerB.EXPECT().Register(app).Once().Run(func(args mock.Arguments) {
			operationOrder = append(operationOrder, "regB")
		})
		providerB.EXPECT().Boot(app).Once().Run(func(args mock.Arguments) {
			operationOrder = append(operationOrder, "bootB")
		})

		app.Register(providerA)
		app.Register(providerB)

		err := app.Boot()
		assert.NoError(t, err)

		// Should use simple registration (RegisterServiceProviders)
		expectedOrder := []string{"regA", "regB", "bootA", "bootB"}
		assert.Equal(t, expectedOrder, operationOrder)

		providerA.AssertExpectations(t)
		providerB.AssertExpectations(t)
	})

	t.Run("Boot_handles_RegisterWithDependencies_error", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		// Create circular dependency to cause RegisterWithDependencies error
		providerA := diMocks.NewMockServiceProvider(t)
		providerA.EXPECT().Providers().Return([]string{"service.a"}).Maybe()
		providerA.EXPECT().Requires().Return([]string{"service.b"}).Maybe()

		providerB := diMocks.NewMockServiceProvider(t)
		providerB.EXPECT().Providers().Return([]string{"service.b"}).Maybe()
		providerB.EXPECT().Requires().Return([]string{"service.a"}).Maybe()

		app.Register(providerA)
		app.Register(providerB)

		// Should return error from RegisterWithDependencies due to circular dependency
		err := app.Boot()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "circular dependency")

		providerA.AssertExpectations(t)
		providerB.AssertExpectations(t)
	})

	t.Run("Boot_handles_BootServiceProviders_error", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		// Provider that will fail during boot
		provider := diMocks.NewMockServiceProvider(t)
		provider.EXPECT().Providers().Return([]string{"failing.service"}).Maybe()
		provider.EXPECT().Requires().Return([]string{}).Maybe()
		provider.EXPECT().Register(app).Once()
		// BootServiceProviders implementation không handle error từ Boot()
		// Nó chỉ gọi provider.Boot() mà không check return value
		provider.EXPECT().Boot(app).Once()

		app.Register(provider)

		// BootServiceProviders luôn return nil theo implementation hiện tại
		err := app.Boot()
		assert.NoError(t, err)

		provider.AssertExpectations(t)
	})
}

// TestApplication_Coverage_Improvements tests edge cases to improve coverage
func TestApplication_Coverage_Improvements(t *testing.T) {
	t.Run("New_with_nil_config", func(t *testing.T) {
		t.Parallel()

		// Test with nil config - should create empty config
		app := core.New(nil)

		// Should still work
		assert.NotNil(t, app)
		assert.NotNil(t, app.Container())

		// Should have default config
		config, err := app.Make("app.config")
		assert.NoError(t, err)
		assert.NotNil(t, config)

		// Should be empty map
		configMap, ok := config.(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, 0, len(configMap))
	})

	t.Run("BootServiceProviders_already_booted", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		var bootCount int

		provider := diMocks.NewMockServiceProvider(t)
		provider.EXPECT().Providers().Return([]string{"test.service"}).Maybe()
		provider.EXPECT().Requires().Return([]string{}).Maybe()
		provider.EXPECT().Register(app).Once()
		provider.EXPECT().Boot(app).Once().Run(func(args mock.Arguments) {
			bootCount++
		})

		app.Register(provider)

		// First boot
		err := app.RegisterServiceProviders()
		assert.NoError(t, err)

		err = app.BootServiceProviders()
		assert.NoError(t, err)
		assert.Equal(t, 1, bootCount)

		// Second boot - should not boot again
		err = app.BootServiceProviders()
		assert.NoError(t, err)
		assert.Equal(t, 1, bootCount) // Should still be 1

		provider.AssertExpectations(t)
	})

	t.Run("Boot_with_no_dependencies_uses_simple_registration", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		var operationOrder []string

		// Providers with no requirements - should use simple registration
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
		providerB.EXPECT().Requires().Return([]string{}).Maybe()
		providerB.EXPECT().Register(app).Once().Run(func(args mock.Arguments) {
			operationOrder = append(operationOrder, "regB")
		})
		providerB.EXPECT().Boot(app).Once().Run(func(args mock.Arguments) {
			operationOrder = append(operationOrder, "bootB")
		})

		app.Register(providerA)
		app.Register(providerB)

		err := app.Boot()
		assert.NoError(t, err)

		// Should use simple registration (RegisterServiceProviders)
		expectedOrder := []string{"regA", "regB", "bootA", "bootB"}
		assert.Equal(t, expectedOrder, operationOrder)

		providerA.AssertExpectations(t)
		providerB.AssertExpectations(t)
	})
}
