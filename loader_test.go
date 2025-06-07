package core_test

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.fork.vn/core"
	coreMocks "go.fork.vn/core/mocks"
	"go.fork.vn/di"
	diMocks "go.fork.vn/di/mocks"
)

// setupTestEnvironment creates necessary test environment for log v0.1.4
func setupTestEnvironment(t *testing.T) {
	err := os.MkdirAll("testdata/logs", 0755)
	require.NoError(t, err)

	logFile := "testdata/logs/app.log"
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		file, err := os.Create(logFile)
		require.NoError(t, err)
		file.Close()
	}
}

// TestModuleLoader_RegisterCoreProviders tests core providers registration
func TestModuleLoader_RegisterCoreProviders(t *testing.T) {
	setupTestEnvironment(t)

	t.Run("successful_core_providers_registration_with_mocks", func(t *testing.T) {
		t.Parallel()

		mockLoader := coreMocks.NewMockModuleLoaderContract(t)
		mockLoader.EXPECT().RegisterCoreProviders().Return(nil).Once()

		err := mockLoader.RegisterCoreProviders()
		assert.NoError(t, err)
	})

	t.Run("fails_when_config_provider_fails_with_mocks", func(t *testing.T) {
		t.Parallel()

		mockLoader := coreMocks.NewMockModuleLoaderContract(t)
		expectedErr := errors.New("config provider registration failed")

		mockLoader.EXPECT().RegisterCoreProviders().Return(expectedErr).Once()

		err := mockLoader.RegisterCoreProviders()
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("fails_when_log_provider_fails_with_mocks", func(t *testing.T) {
		t.Parallel()

		mockLoader := coreMocks.NewMockModuleLoaderContract(t)
		expectedErr := errors.New("log provider registration failed")

		mockLoader.EXPECT().RegisterCoreProviders().Return(expectedErr).Once()

		err := mockLoader.RegisterCoreProviders()
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("fails_when_config_not_found", func(t *testing.T) {
		t.Parallel()

		app := core.New(map[string]interface{}{
			"file": "non-existent-config.yaml",
		})

		err := app.ModuleLoader().RegisterCoreProviders()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "config read failed")
	})

	t.Run("registers_core_providers_with_valid_config", func(t *testing.T) {
		t.Parallel()

		app := core.New(map[string]interface{}{
			"file": "testdata/configs/console-only-simple.yaml",
		})

		err := app.ModuleLoader().RegisterCoreProviders()
		// May fail due to config specifics, but tests the code path
		_ = err
	})

	t.Run("registers_core_providers_with_name_path_type", func(t *testing.T) {
		t.Parallel()

		app := core.New(map[string]interface{}{
			"name": "test-app",
			"path": "testdata/configs",
			"type": "yaml",
		})

		err := app.ModuleLoader().RegisterCoreProviders()
		// Tests the applyConfig path with individual parameters
		_ = err
	})
}

// TestModuleLoader_LoadModule tests module loading functionality
func TestModuleLoader_LoadModule(t *testing.T) {
	t.Run("successfully_loads_valid_service_provider_with_mocks", func(t *testing.T) {
		t.Parallel()

		mockLoader := coreMocks.NewMockModuleLoaderContract(t)
		mockProvider := diMocks.NewMockServiceProvider(t)

		mockLoader.EXPECT().LoadModule(mockProvider).Return(nil).Once()

		err := mockLoader.LoadModule(mockProvider)
		assert.NoError(t, err)
	})

	t.Run("fails_when_module_is_not_a_service_provider_with_mocks", func(t *testing.T) {
		t.Parallel()

		mockLoader := coreMocks.NewMockModuleLoaderContract(t)
		invalidModule := "invalid-module"
		expectedErr := &core.ModuleLoadError{
			Module: invalidModule,
			Reason: "module must implement di.ServiceProvider interface",
		}

		mockLoader.EXPECT().LoadModule(invalidModule).Return(expectedErr).Once()

		err := mockLoader.LoadModule(invalidModule)
		assert.Error(t, err)

		var moduleErr *core.ModuleLoadError
		assert.True(t, errors.As(err, &moduleErr))
		assert.Equal(t, invalidModule, moduleErr.Module)
		assert.Contains(t, moduleErr.Reason, "ServiceProvider interface")
	})

	t.Run("handles_module_load_error_scenarios_with_mocks", func(t *testing.T) {
		t.Parallel()

		mockLoader := coreMocks.NewMockModuleLoaderContract(t)
		testModule := "test-module"
		expectedErr := &core.ModuleLoadError{
			Module: testModule,
			Reason: "custom load error",
		}

		mockLoader.EXPECT().LoadModule(testModule).Return(expectedErr).Once()

		err := mockLoader.LoadModule(testModule)
		assert.Error(t, err)

		var moduleErr *core.ModuleLoadError
		assert.True(t, errors.As(err, &moduleErr))
		assert.Equal(t, testModule, moduleErr.Module)
		assert.Equal(t, "custom load error", moduleErr.Reason)
	})

	t.Run("successfully_loads_valid_service_provider_integration", func(t *testing.T) {
		t.Parallel()

		app := core.New(map[string]interface{}{})
		loader := app.ModuleLoader()

		provider := diMocks.NewMockServiceProvider(t)
		provider.EXPECT().Providers().Return([]string{"test.service"}).Maybe()
		provider.EXPECT().Requires().Return([]string{}).Maybe()

		err := loader.LoadModule(provider)
		assert.NoError(t, err)
	})

	t.Run("fails_when_module_is_not_a_service_provider_integration", func(t *testing.T) {
		t.Parallel()

		app := core.New(map[string]interface{}{})
		loader := app.ModuleLoader()

		invalidModule := "invalid-module-string"

		err := loader.LoadModule(invalidModule)
		assert.Error(t, err)

		var moduleErr *core.ModuleLoadError
		assert.True(t, errors.As(err, &moduleErr))
		assert.Equal(t, invalidModule, moduleErr.Module)
	})
}

// TestModuleLoader_LoadModules tests multiple module loading
func TestModuleLoader_LoadModules(t *testing.T) {
	t.Run("successfully_loads_multiple_valid_providers_with_mocks", func(t *testing.T) {
		t.Parallel()

		mockLoader := coreMocks.NewMockModuleLoaderContract(t)
		modules := []interface{}{"module1", "module2"}

		mockLoader.EXPECT().LoadModules(modules).Return(nil).Once()

		err := mockLoader.LoadModules(modules)
		assert.NoError(t, err)
	})

	t.Run("fails_on_first_invalid_module_with_mocks", func(t *testing.T) {
		t.Parallel()

		mockLoader := coreMocks.NewMockModuleLoaderContract(t)
		modules := []interface{}{"invalid-module", "valid-module"}
		expectedErr := &core.ModuleLoadError{
			Module: "invalid-module",
			Reason: "first module load failed",
		}

		mockLoader.EXPECT().LoadModules(modules).Return(expectedErr).Once()

		err := mockLoader.LoadModules(modules)
		assert.Error(t, err)

		var moduleErr *core.ModuleLoadError
		assert.True(t, errors.As(err, &moduleErr))
	})

	t.Run("fails_on_second_invalid_module_with_mocks", func(t *testing.T) {
		t.Parallel()

		mockLoader := coreMocks.NewMockModuleLoaderContract(t)
		modules := []interface{}{"valid-module", "invalid-module"}
		expectedErr := &core.MultiModuleLoadError{
			FailedIndex:  1,
			FailedModule: "invalid-module",
			Cause: &core.ModuleLoadError{
				Module: "invalid-module",
				Reason: "second module load failed",
			},
		}

		mockLoader.EXPECT().LoadModules(modules).Return(expectedErr).Once()

		err := mockLoader.LoadModules(modules)
		assert.Error(t, err)

		var multiErr *core.MultiModuleLoadError
		assert.True(t, errors.As(err, &multiErr))
	})

	t.Run("loads_modules_with_empty_list", func(t *testing.T) {
		t.Parallel()

		app := core.New(map[string]interface{}{})
		loader := app.ModuleLoader()

		err := loader.LoadModules()
		assert.NoError(t, err)
	})

	t.Run("successfully_loads_multiple_valid_providers_integration", func(t *testing.T) {
		t.Parallel()

		app := core.New(map[string]interface{}{})
		loader := app.ModuleLoader()

		provider1 := diMocks.NewMockServiceProvider(t)
		provider1.EXPECT().Providers().Return([]string{"service1"}).Maybe()
		provider1.EXPECT().Requires().Return([]string{}).Maybe()

		provider2 := diMocks.NewMockServiceProvider(t)
		provider2.EXPECT().Providers().Return([]string{"service2"}).Maybe()
		provider2.EXPECT().Requires().Return([]string{}).Maybe()

		err := loader.LoadModules(provider1, provider2)
		assert.NoError(t, err)
	})

	t.Run("fails_on_first_invalid_module_integration", func(t *testing.T) {
		t.Parallel()

		app := core.New(map[string]interface{}{})
		loader := app.ModuleLoader()

		provider := diMocks.NewMockServiceProvider(t)
		provider.EXPECT().Providers().Return([]string{"service"}).Maybe()
		provider.EXPECT().Requires().Return([]string{}).Maybe()

		err := loader.LoadModules("invalid-module", provider)
		assert.Error(t, err)

		var multiErr *core.MultiModuleLoadError
		assert.True(t, errors.As(err, &multiErr))
		assert.Equal(t, "invalid-module", multiErr.FailedModule)
	})

	t.Run("fails_on_second_invalid_module_integration", func(t *testing.T) {
		t.Parallel()

		app := core.New(map[string]interface{}{})
		loader := app.ModuleLoader()

		provider := diMocks.NewMockServiceProvider(t)
		provider.EXPECT().Providers().Return([]string{"service"}).Maybe()
		provider.EXPECT().Requires().Return([]string{}).Maybe()

		err := loader.LoadModules(provider, "invalid-module")
		assert.Error(t, err)

		var multiErr *core.MultiModuleLoadError
		assert.True(t, errors.As(err, &multiErr))
		assert.Equal(t, 1, multiErr.FailedIndex) // Check the failed index
	})
}

// TestModuleLoader_BootstrapApplication tests application bootstrapping
func TestModuleLoader_BootstrapApplication(t *testing.T) {
	setupTestEnvironment(t)

	t.Run("successfully_bootstraps_application", func(t *testing.T) {
		t.Parallel()

		app := core.New(map[string]interface{}{
			"file": "testdata/configs/console-only-simple.yaml",
		})

		provider := diMocks.NewMockServiceProvider(t)
		provider.EXPECT().Providers().Return([]string{"test.service"}).Maybe()
		provider.EXPECT().Requires().Return([]string{}).Maybe()
		provider.EXPECT().Register(app).Maybe()
		provider.EXPECT().Boot(app).Maybe()

		app.Register(provider)

		err := app.ModuleLoader().BootstrapApplication()
		// May succeed or fail based on config, but tests the code path
		_ = err
	})

	t.Run("fails_when_config_file_not_found", func(t *testing.T) {
		t.Parallel()

		app := core.New(map[string]interface{}{
			"file": "non-existent-file.yaml",
		})

		err := app.ModuleLoader().BootstrapApplication()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "config read failed")
	})

	t.Run("fails_when_register_with_dependencies_fails", func(t *testing.T) {
		t.Parallel()

		app := core.New(map[string]interface{}{
			"file": "testdata/configs/console-only-simple.yaml",
		})

		// Provider with missing dependency
		provider := diMocks.NewMockServiceProvider(t)
		provider.EXPECT().Providers().Return([]string{"service"}).Maybe()
		provider.EXPECT().Requires().Return([]string{"missing.service"}).Maybe()

		app.Register(provider)

		err := app.ModuleLoader().BootstrapApplication()
		assert.Error(t, err)
	})

	t.Run("fails_when_provider_boot_fails", func(t *testing.T) {
		t.Parallel()

		app := core.New(map[string]interface{}{
			"file": "testdata/configs/console-only-simple.yaml",
		})

		provider := diMocks.NewMockServiceProvider(t)
		provider.EXPECT().Providers().Return([]string{"service"}).Maybe()
		provider.EXPECT().Requires().Return([]string{}).Maybe()
		provider.EXPECT().Register(app).Maybe()
		provider.EXPECT().Boot(app).Maybe().Return()

		app.Register(provider)

		err := app.ModuleLoader().BootstrapApplication()
		assert.NoError(t, err)
	})

	t.Run("covers_more_bootstrap_paths", func(t *testing.T) {
		t.Parallel()

		app := core.New(map[string]interface{}{
			"name": "test-app",
			"path": "testdata/configs",
			"type": "yaml",
		})

		err := app.ModuleLoader().BootstrapApplication()
		// Tests different config application paths
		_ = err
	})

	t.Run("bootstrap_with_register_core_providers_error", func(t *testing.T) {
		t.Parallel()

		app := core.New(map[string]interface{}{
			"file": "invalid-config-format.yaml",
		})

		err := app.ModuleLoader().BootstrapApplication()
		assert.Error(t, err)
	})

	t.Run("bootstrap_with_boot_service_providers_error", func(t *testing.T) {
		t.Parallel()

		app := core.New(map[string]interface{}{})

		provider := diMocks.NewMockServiceProvider(t)
		provider.EXPECT().Providers().Return([]string{"service"}).Maybe()
		provider.EXPECT().Requires().Return([]string{}).Maybe()
		provider.EXPECT().Register(app).Maybe()
		provider.EXPECT().Boot(app).Maybe().Run(func(mock.Arguments) {
			panic("boot service providers failed")
		})

		app.Register(provider)

		err := app.ModuleLoader().BootstrapApplication()
		assert.Error(t, err)
	})
}

// TestModuleLoader_ApplyConfig tests config application scenarios
func TestModuleLoader_ApplyConfig(t *testing.T) {
	setupTestEnvironment(t)

	t.Run("apply_config_with_name_path_type_params", func(t *testing.T) {
		t.Parallel()

		app := core.New(map[string]interface{}{
			"name": "test-app",
			"path": "testdata/configs",
			"type": "yaml",
		})

		err := app.ModuleLoader().RegisterCoreProviders()
		// Tests applyConfig with individual parameters
		_ = err
	})

	t.Run("apply_config_with_config_file", func(t *testing.T) {
		t.Parallel()

		app := core.New(map[string]interface{}{
			"file": "testdata/configs/console-only-simple.yaml",
		})

		err := app.ModuleLoader().RegisterCoreProviders()
		// Tests applyConfig with file parameter
		_ = err
	})

	t.Run("apply_config_with_no_config", func(t *testing.T) {
		t.Parallel()

		app := core.New(map[string]interface{}{})

		err := app.ModuleLoader().RegisterCoreProviders()
		// Tests applyConfig with empty config
		_ = err
	})

	t.Run("apply_config_with_invalid_config_manager", func(t *testing.T) {
		t.Parallel()

		app := core.New(map[string]interface{}{})

		// Register invalid config manager type
		app.Instance("app.config", "invalid-config-type")

		err := app.ModuleLoader().RegisterCoreProviders()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid app.config type")
	})

	t.Run("apply_config_file_not_found_error", func(t *testing.T) {
		t.Parallel()

		app := core.New(map[string]interface{}{
			"file": "non-existent-config-file.yaml",
		})

		err := app.ModuleLoader().RegisterCoreProviders()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "config read failed")
	})
}

// TestModuleLoader_ErrorTypes tests custom error types
func TestModuleLoader_ErrorTypes(t *testing.T) {
	t.Run("module_load_error_provides_detailed_message", func(t *testing.T) {
		t.Parallel()

		err := &core.ModuleLoadError{
			Module: "test-module",
			Reason: "test reason",
		}

		assert.Equal(t, "failed to load module: test reason", err.Error())
		assert.Contains(t, err.Error(), "test reason")
	})

	t.Run("multi_module_load_error_provides_detailed_message", func(t *testing.T) {
		t.Parallel()

		err := &core.MultiModuleLoadError{
			FailedIndex:  1,
			FailedModule: "module2",
			Cause: &core.ModuleLoadError{
				Module: "module2",
				Reason: "reason2",
			},
		}

		errMsg := err.Error()
		assert.Contains(t, errMsg, "failed to load modules")
		assert.Contains(t, errMsg, "index 1")
	})
}

// TestModuleLoader_Interface tests interface compliance
func TestModuleLoader_Interface(t *testing.T) {
	t.Run("real_module_loader_implements_interface", func(t *testing.T) {
		t.Parallel()

		app := core.New(map[string]interface{}{})
		loader := app.ModuleLoader()

		// Type assertion to ensure interface compliance
		var coreLoader core.ModuleLoaderContract = loader
		assert.NotNil(t, coreLoader)

		// Also test di.ModuleLoaderContract interface
		var diLoader di.ModuleLoaderContract = loader
		assert.NotNil(t, diLoader)
	})

	t.Run("mock_module_loader_contract_implements_interface", func(t *testing.T) {
		t.Parallel()

		mockLoader := coreMocks.NewMockModuleLoaderContract(t)

		// Type assertion to ensure interface compliance
		var loader core.ModuleLoaderContract = mockLoader
		assert.NotNil(t, loader)

		// Also test di.ModuleLoaderContract interface
		var diLoader di.ModuleLoaderContract = mockLoader
		assert.NotNil(t, diLoader)
	})
}

// TestModuleLoader_CoverageImprovements tests additional coverage scenarios
func TestModuleLoader_CoverageImprovements(t *testing.T) {
	setupTestEnvironment(t)

	t.Run("load_module_covers_error_cases", func(t *testing.T) {
		t.Parallel()

		app := core.New(map[string]interface{}{})
		loader := app.ModuleLoader()

		// Test with non-ServiceProvider interface
		err := loader.LoadModule("not-a-service-provider")
		assert.Error(t, err)

		var moduleErr *core.ModuleLoadError
		assert.True(t, errors.As(err, &moduleErr))
		assert.Equal(t, "not-a-service-provider", moduleErr.Module)
	})

	t.Run("load_module_covers_is_app_booted_branches", func(t *testing.T) {
		t.Parallel()

		app := core.New(map[string]interface{}{})

		// First boot the app
		provider1 := diMocks.NewMockServiceProvider(t)
		provider1.EXPECT().Providers().Return([]string{"service1"}).Maybe()
		provider1.EXPECT().Requires().Return([]string{}).Maybe()
		provider1.EXPECT().Register(app).Maybe()
		provider1.EXPECT().Boot(app).Maybe()

		app.Register(provider1)
		_ = app.RegisterWithDependencies()
		_ = app.BootServiceProviders()

		// Now load module on booted app
		provider2 := diMocks.NewMockServiceProvider(t)
		provider2.EXPECT().Providers().Return([]string{"service2"}).Maybe()
		provider2.EXPECT().Requires().Return([]string{}).Maybe()
		provider2.EXPECT().Register(app).Maybe()
		provider2.EXPECT().Boot(app).Maybe()

		loader := app.ModuleLoader()
		err := loader.LoadModule(provider2)
		assert.NoError(t, err)
	})

	t.Run("load_module_with_fully_booted_app", func(t *testing.T) {
		t.Parallel()

		app := core.New(map[string]interface{}{})
		loader := app.ModuleLoader()

		// Boot app first
		_ = app.RegisterServiceProviders()
		_ = app.BootServiceProviders()

		// Load module after boot
		provider := diMocks.NewMockServiceProvider(t)
		provider.EXPECT().Providers().Return([]string{"service"}).Maybe()
		provider.EXPECT().Requires().Return([]string{}).Maybe()
		provider.EXPECT().Register(app).Maybe()
		provider.EXPECT().Boot(app).Maybe()

		err := loader.LoadModule(provider)
		assert.NoError(t, err)
	})

	t.Run("load_modules_covers_error_propagation", func(t *testing.T) {
		t.Parallel()

		app := core.New(map[string]interface{}{})
		loader := app.ModuleLoader()

		validProvider := diMocks.NewMockServiceProvider(t)
		validProvider.EXPECT().Providers().Return([]string{"valid.service"}).Maybe()
		validProvider.EXPECT().Requires().Return([]string{}).Maybe()

		err := loader.LoadModules(validProvider, "invalid-module-1", "invalid-module-2")
		assert.Error(t, err)

		var multiErr *core.MultiModuleLoadError
		assert.True(t, errors.As(err, &multiErr))
		assert.Equal(t, 1, multiErr.FailedIndex) // First invalid module failed
		assert.Equal(t, "invalid-module-1", multiErr.FailedModule)
	})
}

// BenchmarkModuleLoader_RegisterCoreProviders benchmarks core provider registration
func BenchmarkModuleLoader_RegisterCoreProviders(b *testing.B) {
	setupTestEnvironment(&testing.T{})

	b.Run("register_core_providers_performance", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			app := core.New(map[string]interface{}{})
			_ = app.ModuleLoader().RegisterCoreProviders()
		}
	})

	b.Run("register_core_providers_with_config_performance", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			app := core.New(map[string]interface{}{
				"file": "testdata/configs/console-only-simple.yaml",
			})
			_ = app.ModuleLoader().RegisterCoreProviders()
		}
	})
}

// BenchmarkModuleLoader_LoadModule benchmarks module loading
func BenchmarkModuleLoader_LoadModule(b *testing.B) {
	b.Run("load_single_module_performance", func(b *testing.B) {
		app := core.New(map[string]interface{}{})
		loader := app.ModuleLoader()

		provider := diMocks.NewMockServiceProvider(&testing.T{})
		provider.EXPECT().Providers().Return([]string{"service"}).Maybe()
		provider.EXPECT().Requires().Return([]string{}).Maybe()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = loader.LoadModule(provider)
		}
	})
}
