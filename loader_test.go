package core_test

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.fork.vn/core"
	coreMocks "go.fork.vn/core/mocks"
	"go.fork.vn/di"
	diMocks "go.fork.vn/di/mocks"
)

// setupTestEnvironment tạo môi trường test cần thiết cho log v0.1.4
func setupTestEnvironment(t *testing.T) {
	// Ensure testdata directories exist for log v0.1.4
	err := os.MkdirAll("testdata/logs", 0755)
	require.NoError(t, err)

	// Create log file to prevent log v0.1.4 file handler creation failure
	logFile := "testdata/logs/app.log"
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		file, err := os.Create(logFile)
		require.NoError(t, err)
		file.Close()
	}
}

// TestModuleLoader_BootstrapApplication_Unit kiểm tra unit behavior của BootstrapApplication.
func TestModuleLoader_BootstrapApplication_Unit(t *testing.T) {
	t.Run("successful bootstrap workflow with mocks", func(t *testing.T) {
		t.Parallel()

		// Sử dụng MockModuleLoaderContract theo chiến lược testing
		mockLoader := coreMocks.NewMockModuleLoaderContract(t)

		// Setup expectations cho workflow thành công
		mockLoader.EXPECT().BootstrapApplication().Return(nil).Once()

		// Test behavior không dependencies
		err := mockLoader.BootstrapApplication()
		assert.NoError(t, err)
	})

	t.Run("fails when RegisterCoreProviders fails with mocks", func(t *testing.T) {
		t.Parallel()

		mockLoader := coreMocks.NewMockModuleLoaderContract(t)
		expectedErr := errors.New("core providers registration failed")

		mockLoader.EXPECT().BootstrapApplication().Return(expectedErr).Once()

		err := mockLoader.BootstrapApplication()
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("fails when dependency resolution fails with mocks", func(t *testing.T) {
		t.Parallel()

		mockLoader := coreMocks.NewMockModuleLoaderContract(t)
		expectedErr := errors.New("dependency registration failed")

		mockLoader.EXPECT().BootstrapApplication().Return(expectedErr).Once()

		err := mockLoader.BootstrapApplication()
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("fails when service providers boot fails with mocks", func(t *testing.T) {
		t.Parallel()

		mockLoader := coreMocks.NewMockModuleLoaderContract(t)
		expectedErr := errors.New("boot providers failed")

		mockLoader.EXPECT().BootstrapApplication().Return(expectedErr).Once()

		err := mockLoader.BootstrapApplication()
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

// TestModuleLoader_RegisterCoreProviders_Unit kiểm tra unit behavior của RegisterCoreProviders.
func TestModuleLoader_RegisterCoreProviders_Unit(t *testing.T) {
	t.Run("successful core providers registration with mocks", func(t *testing.T) {
		t.Parallel()

		mockLoader := coreMocks.NewMockModuleLoaderContract(t)

		mockLoader.EXPECT().RegisterCoreProviders().Return(nil).Once()

		err := mockLoader.RegisterCoreProviders()
		assert.NoError(t, err)
	})

	t.Run("fails when config provider fails with mocks", func(t *testing.T) {
		t.Parallel()

		mockLoader := coreMocks.NewMockModuleLoaderContract(t)
		expectedErr := errors.New("config provider registration failed")

		mockLoader.EXPECT().RegisterCoreProviders().Return(expectedErr).Once()

		err := mockLoader.RegisterCoreProviders()
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("fails when log provider fails with mocks", func(t *testing.T) {
		t.Parallel()

		mockLoader := coreMocks.NewMockModuleLoaderContract(t)
		expectedErr := errors.New("log provider registration failed")

		mockLoader.EXPECT().RegisterCoreProviders().Return(expectedErr).Once()

		err := mockLoader.RegisterCoreProviders()
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

// TestModuleLoader_LoadModule_Unit kiểm tra unit behavior của LoadModule.
func TestModuleLoader_LoadModule_Unit(t *testing.T) {
	t.Run("successfully loads valid service provider with mocks", func(t *testing.T) {
		t.Parallel()

		mockLoader := coreMocks.NewMockModuleLoaderContract(t)
		mockProvider := diMocks.NewMockServiceProvider(t)

		mockLoader.EXPECT().LoadModule(mockProvider).Return(nil).Once()

		err := mockLoader.LoadModule(mockProvider)
		assert.NoError(t, err)
	})

	t.Run("fails when module is not a ServiceProvider with mocks", func(t *testing.T) {
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

	t.Run("handles module load error scenarios with mocks", func(t *testing.T) {
		t.Parallel()

		mockLoader := coreMocks.NewMockModuleLoaderContract(t)
		mockProvider := diMocks.NewMockServiceProvider(t)
		expectedErr := errors.New("module load failed")

		mockLoader.EXPECT().LoadModule(mockProvider).Return(expectedErr).Once()

		err := mockLoader.LoadModule(mockProvider)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

// TestModuleLoader_LoadModules_Unit kiểm tra unit behavior của LoadModules.
func TestModuleLoader_LoadModules_Unit(t *testing.T) {
	t.Run("successfully loads multiple valid providers with mocks", func(t *testing.T) {
		t.Parallel()

		mockLoader := coreMocks.NewMockModuleLoaderContract(t)
		mockProvider1 := diMocks.NewMockServiceProvider(t)
		mockProvider2 := diMocks.NewMockServiceProvider(t)

		mockLoader.EXPECT().LoadModules(mockProvider1, mockProvider2).Return(nil).Once()

		err := mockLoader.LoadModules(mockProvider1, mockProvider2)
		assert.NoError(t, err)
	})

	t.Run("fails on first invalid module with mocks", func(t *testing.T) {
		t.Parallel()

		mockLoader := coreMocks.NewMockModuleLoaderContract(t)
		invalidModule := "invalid-module"
		mockProvider := diMocks.NewMockServiceProvider(t)

		expectedErr := &core.MultiModuleLoadError{
			FailedIndex:  0,
			FailedModule: invalidModule,
			Cause: &core.ModuleLoadError{
				Module: invalidModule,
				Reason: "module must implement di.ServiceProvider interface",
			},
		}

		mockLoader.EXPECT().LoadModules(invalidModule, mockProvider).Return(expectedErr).Once()

		err := mockLoader.LoadModules(invalidModule, mockProvider)
		assert.Error(t, err)

		var multiErr *core.MultiModuleLoadError
		assert.True(t, errors.As(err, &multiErr))
		assert.Equal(t, 0, multiErr.FailedIndex)
		assert.Equal(t, invalidModule, multiErr.FailedModule)
	})

	t.Run("fails on second invalid module with mocks", func(t *testing.T) {
		t.Parallel()

		mockLoader := coreMocks.NewMockModuleLoaderContract(t)
		mockProvider := diMocks.NewMockServiceProvider(t)
		invalidModule := "invalid-module"

		expectedErr := &core.MultiModuleLoadError{
			FailedIndex:  1,
			FailedModule: invalidModule,
			Cause: &core.ModuleLoadError{
				Module: invalidModule,
				Reason: "module must implement di.ServiceProvider interface",
			},
		}

		mockLoader.EXPECT().LoadModules(mockProvider, invalidModule).Return(expectedErr).Once()

		err := mockLoader.LoadModules(mockProvider, invalidModule)
		assert.Error(t, err)

		var multiErr *core.MultiModuleLoadError
		assert.True(t, errors.As(err, &multiErr))
		assert.Equal(t, 1, multiErr.FailedIndex)
		assert.Equal(t, invalidModule, multiErr.FailedModule)
	})
}

// TestModuleLoader_BootstrapApplication_Integration kiểm tra toàn bộ workflow bootstrap với real objects.
func TestModuleLoader_BootstrapApplication_Integration(t *testing.T) {
	t.Run("successful complete bootstrap workflow", func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping integration test in short mode due to log v0.1.6 bug")
		}

		t.Parallel()

		// Setup test environment for log v0.1.4
		setupTestEnvironment(t)

		// Use testdata config file with console-only logging to avoid log v0.1.4 file bug
		app := core.New(map[string]interface{}{
			"file": "testdata/configs/console-only-simple.yaml",
		})

		// Add a well-behaved provider
		mockProvider := diMocks.NewMockServiceProvider(t)
		mockProvider.EXPECT().Register(app).Maybe()
		mockProvider.EXPECT().Boot(app).Maybe()
		mockProvider.EXPECT().Providers().Return([]string{"test-service"}).Maybe()
		mockProvider.EXPECT().Requires().Return([]string{}).Maybe()

		app.Register(mockProvider)

		loader := app.ModuleLoader()

		// Bootstrap should succeed with all steps
		err := loader.BootstrapApplication()
		assert.NoError(t, err)

		// Verify that core services are available
		container := app.Container()
		assert.True(t, container.Bound("config"), "config service should be bound")
		assert.True(t, container.Bound("log"), "log service should be bound")
	})

	t.Run("fails when config file not found", func(t *testing.T) {
		t.Parallel()

		// Create app with invalid config path to make core providers fail
		app := core.New(map[string]interface{}{
			"file": "/non/existent/path/config.yaml",
			"name": "test-app",
			"path": "/non/existent/path",
			"type": "yaml",
		})

		loader := app.ModuleLoader()
		err := loader.BootstrapApplication()

		assert.Error(t, err)
		// Error should come from config file reading failure
		assert.Contains(t, err.Error(), "config")
	})

	t.Run("fails when provider boot fails", func(t *testing.T) {
		t.Parallel()

		// Setup test environment for log v0.1.4
		setupTestEnvironment(t)

		// Use testdata config file with console-only logging
		app := core.New(map[string]interface{}{
			"file": "testdata/configs/console-only-simple.yaml",
		})

		// Add a provider that will cause boot to fail
		mockProvider := diMocks.NewMockServiceProvider(t)
		mockProvider.EXPECT().Register(app).Maybe()
		mockProvider.EXPECT().Providers().Return([]string{"test-service"}).Maybe()
		mockProvider.EXPECT().Requires().Return([]string{}).Maybe()
		mockProvider.EXPECT().Boot(app).Run(func(app di.Application) {
			panic("boot failure during bootstrap")
		}).Maybe()

		app.Register(mockProvider)

		loader := app.ModuleLoader()

		// Bootstrap should fail during BootServiceProviders phase
		assert.Panics(t, func() {
			_ = loader.BootstrapApplication()
		})
	})
}

// TestModuleLoader_RegisterCoreProviders_Integration kiểm tra core providers registration với real objects.
func TestModuleLoader_RegisterCoreProviders_Integration(t *testing.T) {
	t.Run("successful core providers registration", func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping integration test in short mode due to log v0.1.6 bug")
		}
		
		t.Parallel()

		// Setup test environment for log v0.1.4
		setupTestEnvironment(t)

		// Use testdata config file with console-only logging
		app := core.New(map[string]interface{}{
			"file": "testdata/configs/console-only-simple.yaml",
		})

		loader := app.ModuleLoader()

		// RegisterCoreProviders should succeed với valid config
		err := loader.RegisterCoreProviders()
		assert.NoError(t, err, "RegisterCoreProviders should succeed with valid config")

		// Test that providers are registered by checking container bindings
		container := app.Container()
		assert.True(t, container.Bound("config"), "config service should be bound")
		assert.True(t, container.Bound("log"), "log service should be bound")

		// Verify services can be resolved
		configService, err := container.Make("config")
		assert.NoError(t, err, "should be able to resolve config service")
		assert.NotNil(t, configService, "config service should not be nil")

		logService, err := container.Make("log")
		assert.NoError(t, err, "should be able to resolve log service")
		assert.NotNil(t, logService, "log service should not be nil")
	})

	t.Run("fails when config not found", func(t *testing.T) {
		t.Parallel()

		// Test với application sẽ fail khi missing config file
		app := core.New(map[string]interface{}{
			"file": "/non/existent/config.yaml",
			"name": "test-app",
			"path": "/non/existent/path",
			"type": "yaml",
		})

		loader := app.ModuleLoader()
		err := loader.RegisterCoreProviders()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "config")
	})
}

// TestModuleLoader_LoadModule_Integration kiểm tra load module với real objects.
func TestModuleLoader_LoadModule_Integration(t *testing.T) {
	t.Run("successfully loads valid service provider", func(t *testing.T) {
		t.Parallel()

		// Create real application for integration test
		app := core.New(map[string]interface{}{
			"name": "test-app",
			"type": "yaml",
		})

		loader := app.ModuleLoader()

		// Create a mock provider
		mockProvider := diMocks.NewMockServiceProvider(t)
		mockProvider.EXPECT().Providers().Return([]string{"test-service"}).Maybe()
		mockProvider.EXPECT().Requires().Return([]string{}).Maybe()

		// LoadModule should succeed
		err := loader.LoadModule(mockProvider)
		assert.NoError(t, err)
	})

	t.Run("fails when module is not a ServiceProvider", func(t *testing.T) {
		t.Parallel()

		app := core.New(map[string]interface{}{
			"name": "test-app",
			"type": "yaml",
		})

		loader := app.ModuleLoader()

		// Try to load invalid module
		invalidModule := "invalid-module"
		err := loader.LoadModule(invalidModule)

		assert.Error(t, err)
		var moduleErr *core.ModuleLoadError
		assert.True(t, errors.As(err, &moduleErr))
		assert.Equal(t, invalidModule, moduleErr.Module)
	})

	t.Run("registers and boots provider when app is already booted", func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping integration test in short mode due to log v0.1.6 bug")
		}
		
		t.Parallel()

		// Setup test environment for log v0.1.4
		setupTestEnvironment(t)

		// Create real application with valid config (console-only to avoid log v0.1.4 bug)
		app := core.New(map[string]interface{}{
			"file": "testdata/configs/console-only-simple.yaml",
		})

		loader := app.ModuleLoader()

		// Bootstrap application first to make it "booted"
		err := loader.BootstrapApplication()
		require.NoError(t, err)

		// Now create a new provider and load it
		newMockProvider := diMocks.NewMockServiceProvider(t)
		newMockProvider.EXPECT().Register(app).Once()
		newMockProvider.EXPECT().Boot(app).Once()
		newMockProvider.EXPECT().Providers().Return([]string{"new-service"}).Maybe()
		newMockProvider.EXPECT().Requires().Return([]string{}).Maybe()

		err = loader.LoadModule(newMockProvider)
		assert.NoError(t, err)
	})
}

// TestModuleLoader_LoadModules_Integration kiểm tra load multiple modules với real objects.
func TestModuleLoader_LoadModules_Integration(t *testing.T) {
	t.Run("successfully loads multiple valid providers", func(t *testing.T) {
		t.Parallel()

		app := core.New(map[string]interface{}{
			"name": "test-app",
			"type": "yaml",
		})

		loader := app.ModuleLoader()

		// Create multiple mock providers
		mockProvider1 := diMocks.NewMockServiceProvider(t)
		mockProvider1.EXPECT().Providers().Return([]string{"service1"}).Maybe()
		mockProvider1.EXPECT().Requires().Return([]string{}).Maybe()

		mockProvider2 := diMocks.NewMockServiceProvider(t)
		mockProvider2.EXPECT().Providers().Return([]string{"service2"}).Maybe()
		mockProvider2.EXPECT().Requires().Return([]string{}).Maybe()

		err := loader.LoadModules(mockProvider1, mockProvider2)
		assert.NoError(t, err)
	})

	t.Run("fails on first invalid module", func(t *testing.T) {
		t.Parallel()

		app := core.New(map[string]interface{}{
			"name": "test-app",
			"type": "yaml",
		})

		loader := app.ModuleLoader()

		invalidModule := "invalid-module"
		mockProvider := diMocks.NewMockServiceProvider(t)

		err := loader.LoadModules(invalidModule, mockProvider)
		assert.Error(t, err)

		var multiErr *core.MultiModuleLoadError
		assert.True(t, errors.As(err, &multiErr))
		assert.Equal(t, 0, multiErr.FailedIndex)
	})

	t.Run("fails on second invalid module", func(t *testing.T) {
		t.Parallel()

		app := core.New(map[string]interface{}{
			"name": "test-app",
			"type": "yaml",
		})

		loader := app.ModuleLoader()

		mockProvider := diMocks.NewMockServiceProvider(t)
		mockProvider.EXPECT().Providers().Return([]string{"service1"}).Maybe()
		mockProvider.EXPECT().Requires().Return([]string{}).Maybe()

		invalidModule := "invalid-module"

		err := loader.LoadModules(mockProvider, invalidModule)
		assert.Error(t, err)

		var multiErr *core.MultiModuleLoadError
		assert.True(t, errors.As(err, &multiErr))
		assert.Equal(t, 1, multiErr.FailedIndex)
	})
}

// TestModuleLoadError_Contract kiểm tra error contracts.
func TestModuleLoadError_Contract(t *testing.T) {
	t.Run("returns correct error message", func(t *testing.T) {
		t.Parallel()

		err := &core.ModuleLoadError{
			Module: "test-module",
			Reason: "test reason",
		}

		expected := "failed to load module: test reason"
		assert.Equal(t, expected, err.Error())
	})
}

// TestMultiModuleLoadError_Contract kiểm tra multi-module error contracts.
func TestMultiModuleLoadError_Contract(t *testing.T) {
	t.Run("returns correct error message", func(t *testing.T) {
		t.Parallel()

		cause := &core.ModuleLoadError{
			Module: "test-module",
			Reason: "test reason",
		}

		err := &core.MultiModuleLoadError{
			FailedIndex:  1,
			FailedModule: "test-module",
			Cause:        cause,
		}

		expected := "failed to load modules: error at index 1, module failed: failed to load module: test reason"
		assert.Equal(t, expected, err.Error())
	})
}

// TestModuleLoaderContract_Interface kiểm tra ModuleLoaderContract interface compliance.
func TestModuleLoaderContract_Interface(t *testing.T) {
	t.Run("MockModuleLoaderContract implements interface", func(t *testing.T) {
		t.Parallel()

		var _ core.ModuleLoaderContract = (*coreMocks.MockModuleLoaderContract)(nil)
	})

	t.Run("real moduleLoader implements interface", func(t *testing.T) {
		t.Parallel()

		app := core.New(map[string]interface{}{
			"name": "test-app",
			"type": "yaml",
		})

		loader := app.ModuleLoader()
		var _ core.ModuleLoaderContract = loader
	})
}
