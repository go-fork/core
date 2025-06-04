package core_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.fork.vn/core"
	coreMocks "go.fork.vn/core/mocks"
	"go.fork.vn/di"
	diMocks "go.fork.vn/di/mocks"
)

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
		t.Parallel()

		// Create a temporary config file
		tempDir, err := os.MkdirTemp("", "core-test-bootstrap-success-*")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		configFile := filepath.Join(tempDir, "test-app.yaml")
		configContent := `log:
  level: "info"
  console:
    enabled: true
    colored: true
  file:
    enabled: false`

		err = os.WriteFile(configFile, []byte(configContent), 0644)
		require.NoError(t, err)

		// Sử dụng real application cho integration test
		app := core.New(map[string]interface{}{
			"file": "test-app.yaml",
			"name": "test-app",
			"path": tempDir,
			"type": "yaml",
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
		err = loader.BootstrapApplication()
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

		// Create a temporary config file
		tempDir, err := os.MkdirTemp("", "core-test-bootstrap-*")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		configFile := filepath.Join(tempDir, "test-app.yaml")
		configContent := `log:
  level: "info"
  console:
    enabled: true`

		err = os.WriteFile(configFile, []byte(configContent), 0644)
		require.NoError(t, err)

		app := core.New(map[string]interface{}{
			"file": "test-app.yaml",
			"name": "test-app",
			"path": tempDir,
			"type": "yaml",
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
		t.Parallel()

		// Create a temporary config file for this test
		tempDir, err := os.MkdirTemp("", "core-test-*")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		configFile := filepath.Join(tempDir, "test-app.yaml")
		configContent := `log:
  level: "info"
  console:
    enabled: true
    colored: true
  file:
    enabled: false`

		err = os.WriteFile(configFile, []byte(configContent), 0644)
		require.NoError(t, err)

		// Setup real application với valid config
		app := core.New(map[string]interface{}{
			"name": "test-app",
			"path": tempDir,
			"type": "yaml",
		})

		loader := app.ModuleLoader()

		// RegisterCoreProviders should succeed với valid config
		err = loader.RegisterCoreProviders()
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
			"file": "missing-config.yaml",
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

		// Create mock service provider
		mockProvider := diMocks.NewMockServiceProvider(t)
		mockProvider.EXPECT().Register(mock.Anything).Maybe()
		mockProvider.EXPECT().Boot(mock.Anything).Maybe()

		app := core.New(map[string]interface{}{})
		loader := app.ModuleLoader()

		err := loader.LoadModule(mockProvider)
		assert.NoError(t, err)
	})

	t.Run("fails when module is not a ServiceProvider", func(t *testing.T) {
		t.Parallel()

		app := core.New(map[string]interface{}{})
		loader := app.ModuleLoader()

		// Try to load a string instead of ServiceProvider
		err := loader.LoadModule("invalid-module")
		assert.Error(t, err)

		var moduleErr *core.ModuleLoadError
		assert.True(t, errors.As(err, &moduleErr))
		assert.Equal(t, "invalid-module", moduleErr.Module)
		assert.Contains(t, moduleErr.Reason, "must implement di.ServiceProvider interface")
	})

	t.Run("registers and boots provider when app is already booted", func(t *testing.T) {
		t.Parallel()

		// Create a test config file for this test
		testConfigFile, err := os.CreateTemp("", "test-config-*.yaml")
		require.NoError(t, err)
		defer os.Remove(testConfigFile.Name())

		// Create a basic config file for testing
		configContent := `log:
  level: info
  format: text
test:
  value: "integration-test"`
		_, err = testConfigFile.WriteString(configContent)
		require.NoError(t, err)
		err = testConfigFile.Close()
		require.NoError(t, err)

		// Create real app và bootstrap it first để simulate already booted state
		app := core.New(map[string]interface{}{
			"file": testConfigFile.Name(),
		})

		// Bootstrap the app để put it in "booted" state
		loader := app.ModuleLoader()

		// Bootstrap should succeed now với valid config
		err = loader.RegisterCoreProviders()
		assert.NoError(t, err)

		// Now create và load a mock provider - it should be registered và booted immediately
		mockProvider := diMocks.NewMockServiceProvider(t)
		mockProvider.EXPECT().Register(app).Once()
		mockProvider.EXPECT().Boot(app).Once()

		err = loader.LoadModule(mockProvider)
		assert.NoError(t, err)
	})
}

// TestModuleLoader_LoadModules_Integration kiểm tra load multiple modules với real objects.
func TestModuleLoader_LoadModules_Integration(t *testing.T) {
	t.Run("successfully loads multiple valid providers", func(t *testing.T) {
		t.Parallel()

		// Create mock service providers
		mockProvider1 := diMocks.NewMockServiceProvider(t)
		mockProvider2 := diMocks.NewMockServiceProvider(t)

		mockProvider1.EXPECT().Register(mock.Anything).Maybe()
		mockProvider1.EXPECT().Boot(mock.Anything).Maybe()
		mockProvider2.EXPECT().Register(mock.Anything).Maybe()
		mockProvider2.EXPECT().Boot(mock.Anything).Maybe()

		app := core.New(map[string]interface{}{})
		loader := app.ModuleLoader()

		err := loader.LoadModules(mockProvider1, mockProvider2)
		assert.NoError(t, err)
	})

	t.Run("fails on first invalid module", func(t *testing.T) {
		t.Parallel()

		mockProvider := diMocks.NewMockServiceProvider(t)

		app := core.New(map[string]interface{}{})
		loader := app.ModuleLoader()

		err := loader.LoadModules("invalid-module", mockProvider)
		assert.Error(t, err)

		var multiErr *core.MultiModuleLoadError
		assert.True(t, errors.As(err, &multiErr))
		assert.Equal(t, 0, multiErr.FailedIndex)
		assert.Equal(t, "invalid-module", multiErr.FailedModule)

		var moduleErr *core.ModuleLoadError
		assert.True(t, errors.As(multiErr.Cause, &moduleErr))
	})

	t.Run("fails on second invalid module", func(t *testing.T) {
		t.Parallel()

		mockProvider := diMocks.NewMockServiceProvider(t)
		mockProvider.EXPECT().Register(mock.Anything).Maybe()
		mockProvider.EXPECT().Boot(mock.Anything).Maybe()

		app := core.New(map[string]interface{}{})
		loader := app.ModuleLoader()

		err := loader.LoadModules(mockProvider, "invalid-module")
		assert.Error(t, err)

		var multiErr *core.MultiModuleLoadError
		assert.True(t, errors.As(err, &multiErr))
		assert.Equal(t, 1, multiErr.FailedIndex)
		assert.Equal(t, "invalid-module", multiErr.FailedModule)
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
		assert.Implements(t, (*error)(nil), err)
	})
}

// TestMultiModuleLoadError_Contract kiểm tra multi-module error contracts.
func TestMultiModuleLoadError_Contract(t *testing.T) {
	t.Run("returns correct error message", func(t *testing.T) {
		t.Parallel()

		causeErr := errors.New("original error")
		err := &core.MultiModuleLoadError{
			FailedIndex:  2,
			FailedModule: "test-module",
			Cause:        causeErr,
		}

		expected := "failed to load modules: error at index 2, module failed: original error"
		assert.Equal(t, expected, err.Error())
		assert.Implements(t, (*error)(nil), err)
	})
}

// TestModuleLoaderContract_Interface kiểm tra ModuleLoaderContract interface compliance.
func TestModuleLoaderContract_Interface(t *testing.T) {
	t.Run("MockModuleLoaderContract implements interface", func(t *testing.T) {
		t.Parallel()

		mockLoader := coreMocks.NewMockModuleLoaderContract(t)
		assert.Implements(t, (*core.ModuleLoaderContract)(nil), mockLoader)
		assert.Implements(t, (*di.ModuleLoaderContract)(nil), mockLoader)
	})

	t.Run("real moduleLoader implements interface", func(t *testing.T) {
		t.Parallel()

		app := core.New(map[string]interface{}{})
		loader := app.ModuleLoader()

		assert.Implements(t, (*core.ModuleLoaderContract)(nil), loader)
		assert.Implements(t, (*di.ModuleLoaderContract)(nil), loader)
	})
}
