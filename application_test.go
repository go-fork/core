package core_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	configMocks "go.fork.vn/config/mocks"
	"go.fork.vn/core"
	"go.fork.vn/di"
	diMocks "go.fork.vn/di/mocks"
	logMocks "go.fork.vn/log/mocks"
)

// TestApplication_New tests Application constructor
func TestApplication_New(t *testing.T) {
	t.Run("creates_application_with_valid_config", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{
			"name": "test-app",
			"path": "./configs",
			"type": "yaml",
		}

		app := core.New(config)

		assert.NotNil(t, app)
		assert.NotNil(t, app.Container())
		assert.NotNil(t, app.ModuleLoader())
	})

	t.Run("creates_application_with_empty_config", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		assert.NotNil(t, app)
		assert.NotNil(t, app.Container())
		assert.NotNil(t, app.ModuleLoader())
	})

	t.Run("creates_application_with_nil_config", func(t *testing.T) {
		t.Parallel()

		app := core.New(nil)

		assert.NotNil(t, app)
		assert.NotNil(t, app.Container())
		assert.NotNil(t, app.ModuleLoader())
	})
}

// TestApplication_Register tests service provider registration
func TestApplication_Register(t *testing.T) {
	t.Run("successfully_registers_service_provider", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		mockProvider := diMocks.NewMockServiceProvider(t)
		mockProvider.EXPECT().Providers().Return([]string{"test.service"}).Maybe()
		mockProvider.EXPECT().Requires().Return([]string{}).Maybe()

		app.Register(mockProvider)
		assert.NotNil(t, app)
	})

	t.Run("panics_when_registering_nil_provider", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		assert.Panics(t, func() {
			app.Register(nil)
		})
	})

	t.Run("registers_multiple_providers_in_order", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		provider1 := diMocks.NewMockServiceProvider(t)
		provider1.EXPECT().Providers().Return([]string{"service1"}).Maybe()
		provider1.EXPECT().Requires().Return([]string{}).Maybe()

		provider2 := diMocks.NewMockServiceProvider(t)
		provider2.EXPECT().Providers().Return([]string{"service2"}).Maybe()
		provider2.EXPECT().Requires().Return([]string{}).Maybe()

		provider3 := diMocks.NewMockServiceProvider(t)
		provider3.EXPECT().Providers().Return([]string{"service3"}).Maybe()
		provider3.EXPECT().Requires().Return([]string{}).Maybe()

		app.Register(provider1)
		app.Register(provider2)
		app.Register(provider3)

		assert.NotNil(t, app)
	})
}

// TestApplication_RegisterServiceProviders tests RegisterServiceProviders method
func TestApplication_RegisterServiceProviders(t *testing.T) {
	t.Run("registers_all_providers_in_order", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		var registrationOrder []string

		providerA := diMocks.NewMockServiceProvider(t)
		providerA.EXPECT().Providers().Return([]string{"service.a"}).Maybe()
		providerA.EXPECT().Requires().Return([]string{}).Maybe()
		providerA.EXPECT().Register(app).Once().Run(func(args mock.Arguments) {
			registrationOrder = append(registrationOrder, "A")
		})

		providerB := diMocks.NewMockServiceProvider(t)
		providerB.EXPECT().Providers().Return([]string{"service.b"}).Maybe()
		providerB.EXPECT().Requires().Return([]string{}).Maybe()
		providerB.EXPECT().Register(app).Once().Run(func(args mock.Arguments) {
			registrationOrder = append(registrationOrder, "B")
		})

		app.Register(providerA)
		app.Register(providerB)

		err := app.RegisterServiceProviders()
		assert.NoError(t, err)
		assert.Equal(t, []string{"A", "B"}, registrationOrder)

		providerA.AssertExpectations(t)
		providerB.AssertExpectations(t)
	})

	t.Run("handles_registration_error", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		provider := diMocks.NewMockServiceProvider(t)
		provider.EXPECT().Providers().Return([]string{"service"}).Maybe()
		provider.EXPECT().Requires().Return([]string{}).Maybe()
		provider.EXPECT().Register(app).Once().Return()

		app.Register(provider)

		// This should not panic in normal circumstances
		err := app.RegisterServiceProviders()
		assert.NoError(t, err)
	})
}

// TestApplication_RegisterWithDependencies tests dependency resolution
func TestApplication_RegisterWithDependencies(t *testing.T) {
	t.Run("simple_dependency_chain_a_requires_nothing", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		providerA := diMocks.NewMockServiceProvider(t)
		providerA.EXPECT().Providers().Return([]string{"database.connection"}).Maybe()
		providerA.EXPECT().Requires().Return([]string{}).Maybe()
		providerA.EXPECT().Register(app).Once()

		app.Register(providerA)

		err := app.RegisterWithDependencies()
		assert.NoError(t, err)

		providerA.AssertExpectations(t)
	})

	t.Run("dependency_chain_b_requires_a", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		providerA := diMocks.NewMockServiceProvider(t)
		providerA.EXPECT().Providers().Return([]string{"database.connection"}).Maybe()
		providerA.EXPECT().Requires().Return([]string{}).Maybe()
		providerA.EXPECT().Register(app).Once()

		providerB := diMocks.NewMockServiceProvider(t)
		providerB.EXPECT().Providers().Return([]string{"cache.manager"}).Maybe()
		providerB.EXPECT().Requires().Return([]string{"database.connection"}).Maybe()
		providerB.EXPECT().Register(app).Once()

		app.Register(providerB) // Register B first to test sorting
		app.Register(providerA)

		err := app.RegisterWithDependencies()
		assert.NoError(t, err)

		providerA.AssertExpectations(t)
		providerB.AssertExpectations(t)
	})

	t.Run("complex_dependency_chain_c_requires_b_requires_a", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		providerA := diMocks.NewMockServiceProvider(t)
		providerA.EXPECT().Providers().Return([]string{"app.config"}).Maybe()
		providerA.EXPECT().Requires().Return([]string{}).Maybe()
		providerA.EXPECT().Register(app).Once()

		providerB := diMocks.NewMockServiceProvider(t)
		providerB.EXPECT().Providers().Return([]string{"database.connection"}).Maybe()
		providerB.EXPECT().Requires().Return([]string{"app.config"}).Maybe()
		providerB.EXPECT().Register(app).Once()

		providerC := diMocks.NewMockServiceProvider(t)
		providerC.EXPECT().Providers().Return([]string{"auth.service"}).Maybe()
		providerC.EXPECT().Requires().Return([]string{"database.connection"}).Maybe()
		providerC.EXPECT().Register(app).Once()

		// Register in reverse order to test dependency resolution
		app.Register(providerC)
		app.Register(providerB)
		app.Register(providerA)

		err := app.RegisterWithDependencies()
		assert.NoError(t, err)

		providerA.AssertExpectations(t)
		providerB.AssertExpectations(t)
		providerC.AssertExpectations(t)
	})

	t.Run("multiple_providers_requiring_same_service", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		providerA := diMocks.NewMockServiceProvider(t)
		providerA.EXPECT().Providers().Return([]string{"shared.service"}).Maybe()
		providerA.EXPECT().Requires().Return([]string{}).Maybe()
		providerA.EXPECT().Register(app).Once()

		providerB := diMocks.NewMockServiceProvider(t)
		providerB.EXPECT().Providers().Return([]string{"feature.b"}).Maybe()
		providerB.EXPECT().Requires().Return([]string{"shared.service"}).Maybe()
		providerB.EXPECT().Register(app).Once()

		providerC := diMocks.NewMockServiceProvider(t)
		providerC.EXPECT().Providers().Return([]string{"feature.c"}).Maybe()
		providerC.EXPECT().Requires().Return([]string{"shared.service"}).Maybe()
		providerC.EXPECT().Register(app).Once()

		app.Register(providerB)
		app.Register(providerC)
		app.Register(providerA)

		err := app.RegisterWithDependencies()
		assert.NoError(t, err)

		providerA.AssertExpectations(t)
		providerB.AssertExpectations(t)
		providerC.AssertExpectations(t)
	})

	t.Run("provider_with_multiple_requirements", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		providerA := diMocks.NewMockServiceProvider(t)
		providerA.EXPECT().Providers().Return([]string{"service.a"}).Maybe()
		providerA.EXPECT().Requires().Return([]string{}).Maybe()
		providerA.EXPECT().Register(app).Once()

		providerB := diMocks.NewMockServiceProvider(t)
		providerB.EXPECT().Providers().Return([]string{"service.b"}).Maybe()
		providerB.EXPECT().Requires().Return([]string{}).Maybe()
		providerB.EXPECT().Register(app).Once()

		providerC := diMocks.NewMockServiceProvider(t)
		providerC.EXPECT().Providers().Return([]string{"service.c"}).Maybe()
		providerC.EXPECT().Requires().Return([]string{"service.a", "service.b"}).Maybe()
		providerC.EXPECT().Register(app).Once()

		app.Register(providerC)
		app.Register(providerA)
		app.Register(providerB)

		err := app.RegisterWithDependencies()
		assert.NoError(t, err)

		providerA.AssertExpectations(t)
		providerB.AssertExpectations(t)
		providerC.AssertExpectations(t)
	})

	t.Run("circular_dependency_error", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		providerA := diMocks.NewMockServiceProvider(t)
		providerA.EXPECT().Providers().Return([]string{"service.a"}).Maybe()
		providerA.EXPECT().Requires().Return([]string{"service.b"}).Maybe()

		providerB := diMocks.NewMockServiceProvider(t)
		providerB.EXPECT().Providers().Return([]string{"service.b"}).Maybe()
		providerB.EXPECT().Requires().Return([]string{"service.a"}).Maybe()

		app.Register(providerA)
		app.Register(providerB)

		err := app.RegisterWithDependencies()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "circular dependency detected")
	})

	t.Run("missing_required_service_error", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		provider := diMocks.NewMockServiceProvider(t)
		provider.EXPECT().Providers().Return([]string{"service.a"}).Maybe()
		provider.EXPECT().Requires().Return([]string{"missing.service"}).Maybe()

		app.Register(provider)

		err := app.RegisterWithDependencies()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "required service 'missing.service' not provided by any registered provider")
	})
}

// TestApplication_BootServiceProviders tests service provider booting
func TestApplication_BootServiceProviders(t *testing.T) {
	t.Run("boots_providers_in_dependency_order", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		var bootOrder []string

		providerA := diMocks.NewMockServiceProvider(t)
		providerA.EXPECT().Providers().Return([]string{"service.a"}).Maybe()
		providerA.EXPECT().Requires().Return([]string{}).Maybe()
		providerA.EXPECT().Register(app).Once()
		providerA.EXPECT().Boot(app).Once().Run(func(args mock.Arguments) {
			bootOrder = append(bootOrder, "providerA")
		})

		providerB := diMocks.NewMockServiceProvider(t)
		providerB.EXPECT().Providers().Return([]string{"service.b"}).Maybe()
		providerB.EXPECT().Requires().Return([]string{"service.a"}).Maybe()
		providerB.EXPECT().Register(app).Once()
		providerB.EXPECT().Boot(app).Once().Run(func(args mock.Arguments) {
			bootOrder = append(bootOrder, "providerB")
		})

		app.Register(providerB) // Register B first
		app.Register(providerA)

		err := app.RegisterWithDependencies()
		assert.NoError(t, err)

		err = app.BootServiceProviders()
		assert.NoError(t, err)

		// Should boot A before B due to dependency
		assert.Equal(t, []string{"providerA", "providerB"}, bootOrder)

		providerA.AssertExpectations(t)
		providerB.AssertExpectations(t)
	})

	t.Run("handles_boot_error", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		provider := diMocks.NewMockServiceProvider(t)
		provider.EXPECT().Providers().Return([]string{"service"}).Maybe()
		provider.EXPECT().Requires().Return([]string{}).Maybe()
		provider.EXPECT().Register(app).Once()
		provider.EXPECT().Boot(app).Once().Return()

		app.Register(provider)

		err := app.RegisterWithDependencies()
		assert.NoError(t, err)

		err = app.BootServiceProviders()
		assert.NoError(t, err)
	})

	t.Run("handles_already_booted_scenario", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		provider := diMocks.NewMockServiceProvider(t)
		provider.EXPECT().Providers().Return([]string{"service"}).Maybe()
		provider.EXPECT().Requires().Return([]string{}).Maybe()
		provider.EXPECT().Register(app).Once()
		provider.EXPECT().Boot(app).Once()

		app.Register(provider)

		err := app.RegisterWithDependencies()
		assert.NoError(t, err)

		// First boot
		err = app.BootServiceProviders()
		assert.NoError(t, err)

		// Second boot should not call Boot again
		err = app.BootServiceProviders()
		assert.NoError(t, err)

		provider.AssertExpectations(t)
	})
}

// TestApplication_Boot tests the high-level Boot method
func TestApplication_Boot(t *testing.T) {
	t.Run("boot_with_no_dependencies_uses_simple_registration", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		provider := diMocks.NewMockServiceProvider(t)
		provider.EXPECT().Providers().Return([]string{"service"}).Maybe()
		provider.EXPECT().Requires().Return([]string{}).Maybe()
		provider.EXPECT().Register(app).Once()
		provider.EXPECT().Boot(app).Once()

		app.Register(provider)

		err := app.Boot()
		assert.NoError(t, err)

		provider.AssertExpectations(t)
	})

	t.Run("boot_with_dependencies_uses_dependency_resolution", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		providerA := diMocks.NewMockServiceProvider(t)
		providerA.EXPECT().Providers().Return([]string{"service.a"}).Maybe()
		providerA.EXPECT().Requires().Return([]string{}).Maybe()
		providerA.EXPECT().Register(app).Once()
		providerA.EXPECT().Boot(app).Once()

		providerB := diMocks.NewMockServiceProvider(t)
		providerB.EXPECT().Providers().Return([]string{"service.b"}).Maybe()
		providerB.EXPECT().Requires().Return([]string{"service.a"}).Maybe()
		providerB.EXPECT().Register(app).Once()
		providerB.EXPECT().Boot(app).Once()

		app.Register(providerB)
		app.Register(providerA)

		err := app.Boot()
		assert.NoError(t, err)

		providerA.AssertExpectations(t)
		providerB.AssertExpectations(t)
	})

	t.Run("boot_handles_registration_error", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		providerA := diMocks.NewMockServiceProvider(t)
		providerA.EXPECT().Providers().Return([]string{"service.a"}).Maybe()
		providerA.EXPECT().Requires().Return([]string{"missing.service"}).Maybe()

		app.Register(providerA)

		err := app.Boot()
		assert.Error(t, err)
	})

	t.Run("boot_handles_boot_service_providers_error", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		provider := diMocks.NewMockServiceProvider(t)
		provider.EXPECT().Providers().Return([]string{"service"}).Maybe()
		provider.EXPECT().Requires().Return([]string{}).Maybe()
		provider.EXPECT().Register(app).Once()
		provider.EXPECT().Boot(app).Once().Return()

		app.Register(provider)

		err := app.Boot()
		assert.NoError(t, err)
	})
}

// TestApplication_DI_Integration tests DI container integration
func TestApplication_DI_Integration(t *testing.T) {
	t.Run("bind_and_make_service", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		app.Bind("test.service", func(c di.Container) interface{} {
			return "test-value"
		})

		service, err := app.Make("test.service")
		assert.NoError(t, err)
		assert.Equal(t, "test-value", service)
	})

	t.Run("singleton_binding", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		counter := 0
		app.Singleton("counter.service", func(c di.Container) interface{} {
			counter++
			return counter
		})

		service1, err1 := app.Make("counter.service")
		service2, err2 := app.Make("counter.service")

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Equal(t, 1, service1)
		assert.Equal(t, 1, service2) // Same instance
		assert.Equal(t, 1, counter)  // Only instantiated once
	})

	t.Run("instance_binding", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		testValue := "instance-value"
		app.Instance("test.instance", testValue)

		service, err := app.Make("test.instance")
		assert.NoError(t, err)
		assert.Equal(t, testValue, service)
	})

	t.Run("alias_creation", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		app.Bind("original.service", func(c di.Container) interface{} {
			return "original-value"
		})

		app.Alias("original.service", "aliased.service")

		original, err1 := app.Make("original.service")
		aliased, err2 := app.Make("aliased.service")

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Equal(t, original, aliased)
	})
}

// TestApplication_MustMake tests MustMake method
func TestApplication_MustMake(t *testing.T) {
	t.Run("returns_service_when_exists", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		app.Bind("test.service", func(c di.Container) interface{} {
			return "test-value"
		})

		service := app.MustMake("test.service")
		assert.Equal(t, "test-value", service)
	})

	t.Run("panics_when_service_not_found", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		assert.Panics(t, func() {
			app.MustMake("non.existent.service")
		})
	})
}

// TestApplication_Call tests Call method
func TestApplication_Call(t *testing.T) {
	t.Run("calls_function_with_parameters", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		app.Bind("string", func(c di.Container) interface{} {
			return "injected-value"
		})

		testFunc := func(param string) string {
			return "result: " + param
		}

		result, err := app.Call(testFunc)
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "result: injected-value", result[0])
	})
}

// TestApplication_Config_Integration tests config manager integration
func TestApplication_Config_Integration(t *testing.T) {
	t.Run("returns_config_manager_when_available", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		mockConfig := configMocks.NewMockManager(t)
		app.Instance("config", mockConfig)

		configManager := app.Config()
		assert.Equal(t, mockConfig, configManager)
	})

	t.Run("panics_when_config_not_registered", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		assert.Panics(t, func() {
			app.Config()
		})
	})
}

// TestApplication_Log_Integration tests log manager integration
func TestApplication_Log_Integration(t *testing.T) {
	t.Run("returns_log_manager_when_available", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		mockLog := logMocks.NewMockManager(t)
		app.Instance("log", mockLog)

		logManager := app.Log()
		assert.Equal(t, mockLog, logManager)
	})

	t.Run("panics_when_log_not_registered", func(t *testing.T) {
		t.Parallel()

		config := map[string]interface{}{}
		app := core.New(config)

		assert.Panics(t, func() {
			app.Log()
		})
	})
}
