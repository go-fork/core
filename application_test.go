package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	config_mocks "go.fork.vn/config/mocks"
	"go.fork.vn/di"
	di_mocks "go.fork.vn/di/mocks"
	log_mocks "go.fork.vn/log/mocks"
)

// Mock service provider with time tracking for dependency ordering tests
type mockServiceProvider struct {
	name         string
	providers    []string
	requires     []string
	registered   bool
	booted       bool
	registerTime time.Time
	bootTime     time.Time
}

func (m *mockServiceProvider) Register(app di.Application) {
	m.registered = true
	m.registerTime = time.Now()
}

func (m *mockServiceProvider) Boot(app di.Application) {
	m.booted = true
	m.bootTime = time.Now()
}

func (m *mockServiceProvider) Requires() []string {
	return m.requires
}

func (m *mockServiceProvider) Providers() []string {
	return m.providers
}

// Test New function
func TestNew(t *testing.T) {
	t.Run("creates application with valid config", func(t *testing.T) {
		config := map[string]interface{}{
			"name": "test-app",
			"path": "./configs",
		}

		app := New(config)

		assert.NotNil(t, app)
		assert.NotNil(t, app.Container())
		assert.NotNil(t, app.ModuleLoader())

		// Verify config was registered
		appConfig, err := app.Make("app.config")
		assert.NoError(t, err)
		assert.Equal(t, config, appConfig)
	})

	t.Run("creates application with nil config", func(t *testing.T) {
		app := New(nil)

		assert.NotNil(t, app)
		assert.NotNil(t, app.Container())

		// Should have empty config
		appConfig, err := app.Make("app.config")
		assert.NoError(t, err)
		assert.NotNil(t, appConfig)
	})

	t.Run("creates application with empty config", func(t *testing.T) {
		config := make(map[string]interface{})
		app := New(config)

		assert.NotNil(t, app)
		appConfig, err := app.Make("app.config")
		assert.NoError(t, err)
		assert.Equal(t, config, appConfig)
	})
}

// Test Container method
func TestApplication_Container(t *testing.T) {
	app := New(nil)
	container := app.Container()

	assert.NotNil(t, container)
	assert.Implements(t, (*di.Container)(nil), container)
}

// Test Config method
func TestApplication_Config(t *testing.T) {
	t.Run("returns config manager when registered", func(t *testing.T) {
		app := New(nil)
		mockConfig := config_mocks.NewMockManager(t)

		// Register mock config
		app.Instance("config", mockConfig)

		config := app.Config()
		assert.Equal(t, mockConfig, config)
	})

	t.Run("panics when config not registered", func(t *testing.T) {
		app := New(nil)

		assert.Panics(t, func() {
			app.Config()
		})
	})
}

// Test Log method
func TestApplication_Log(t *testing.T) {
	t.Run("returns log manager when registered", func(t *testing.T) {
		app := New(nil)
		mockLog := log_mocks.NewMockManager(t)

		// Register mock log
		app.Instance("log", mockLog)

		log := app.Log()
		assert.Equal(t, mockLog, log)
	})

	t.Run("panics when log not registered", func(t *testing.T) {
		app := New(nil)

		assert.Panics(t, func() {
			app.Log()
		})
	})
}

// Test ModuleLoader method
func TestApplication_ModuleLoader(t *testing.T) {
	app := New(nil)
	loader := app.ModuleLoader()

	assert.NotNil(t, loader)
	assert.Implements(t, (*ModuleLoaderContract)(nil), loader)
}

// Test service provider registration
func TestApplication_ServiceProviders(t *testing.T) {
	t.Run("register service provider", func(t *testing.T) {
		app := New(nil)
		mockProvider := di_mocks.NewMockServiceProvider(t)

		// Setup expectations
		mockProvider.EXPECT().Register(app).Once()

		app.Register(mockProvider)
		err := app.RegisterServiceProviders()

		assert.NoError(t, err)
		mockProvider.AssertExpectations(t)
	})

	t.Run("register multiple service providers", func(t *testing.T) {
		app := New(nil)
		mockProvider1 := di_mocks.NewMockServiceProvider(t)
		mockProvider2 := di_mocks.NewMockServiceProvider(t)

		// Setup expectations
		mockProvider1.EXPECT().Register(app).Once()
		mockProvider2.EXPECT().Register(app).Once()

		app.Register(mockProvider1)
		app.Register(mockProvider2)
		err := app.RegisterServiceProviders()

		assert.NoError(t, err)
		mockProvider1.AssertExpectations(t)
		mockProvider2.AssertExpectations(t)
	})

	t.Run("panic on nil provider", func(t *testing.T) {
		app := New(nil)

		assert.Panics(t, func() {
			app.Register(nil)
		})
	})
}

// Test RegisterWithDependencies
func TestApplication_RegisterWithDependencies(t *testing.T) {
	t.Run("registers providers in dependency order", func(t *testing.T) {
		app := New(nil)

		// Create mock providers with dependencies
		providerA := di_mocks.NewMockServiceProvider(t)
		providerB := di_mocks.NewMockServiceProvider(t)

		// Provider A provides "serviceA", requires nothing
		providerA.EXPECT().Providers().Return([]string{"serviceA"}).Times(1)
		providerA.EXPECT().Requires().Return([]string{}).Times(1)
		providerA.EXPECT().Register(app).Once()

		// Provider B provides "serviceB", requires "serviceA"
		providerB.EXPECT().Providers().Return([]string{"serviceB"}).Times(1)
		providerB.EXPECT().Requires().Return([]string{"serviceA"}).Times(1)
		providerB.EXPECT().Register(app).Once()

		// Register providers (B first, then A - should be reordered)
		app.Register(providerB)
		app.Register(providerA)

		err := app.RegisterWithDependencies()
		assert.NoError(t, err)

		providerA.AssertExpectations(t)
		providerB.AssertExpectations(t)
	})

	t.Run("detects circular dependency", func(t *testing.T) {
		app := New(nil)

		providerA := di_mocks.NewMockServiceProvider(t)
		providerB := di_mocks.NewMockServiceProvider(t)

		// Circular dependency: A requires B, B requires A
		providerA.EXPECT().Providers().Return([]string{"serviceA"}).Times(1)
		providerA.EXPECT().Requires().Return([]string{"serviceB"}).Times(1)

		providerB.EXPECT().Providers().Return([]string{"serviceB"}).Times(1)
		providerB.EXPECT().Requires().Return([]string{"serviceA"}).Times(1)

		app.Register(providerA)
		app.Register(providerB)

		err := app.RegisterWithDependencies()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "circular dependency")
	})

	t.Run("detects missing dependency", func(t *testing.T) {
		app := New(nil)

		provider := di_mocks.NewMockServiceProvider(t)

		// Provider requires "missing-service" that no one provides
		provider.EXPECT().Providers().Return([]string{"serviceA"}).Times(1)
		provider.EXPECT().Requires().Return([]string{"missing-service"}).Times(1)

		app.Register(provider)

		err := app.RegisterWithDependencies()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "required service 'missing-service' not provided")
	})
}

// Test BootServiceProviders
func TestApplication_BootServiceProviders(t *testing.T) {
	t.Run("boots all registered providers", func(t *testing.T) {
		app := New(nil)
		mockProvider1 := di_mocks.NewMockServiceProvider(t)
		mockProvider2 := di_mocks.NewMockServiceProvider(t)

		// Setup expectations
		mockProvider1.EXPECT().Boot(app).Once()
		mockProvider2.EXPECT().Boot(app).Once()

		app.Register(mockProvider1)
		app.Register(mockProvider2)

		err := app.BootServiceProviders()
		assert.NoError(t, err)

		mockProvider1.AssertExpectations(t)
		mockProvider2.AssertExpectations(t)
	})

	t.Run("only boots once", func(t *testing.T) {
		app := New(nil)
		mockProvider := di_mocks.NewMockServiceProvider(t)

		// Should only be called once
		mockProvider.EXPECT().Boot(app).Once()

		app.Register(mockProvider)

		// Boot twice
		err1 := app.BootServiceProviders()
		err2 := app.BootServiceProviders()

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		mockProvider.AssertExpectations(t)
	})
}

// Test Boot method (combines register and boot)
func TestApplication_Boot(t *testing.T) {
	t.Run("boots providers without dependencies using simple registration", func(t *testing.T) {
		app := New(nil)
		mockProvider := di_mocks.NewMockServiceProvider(t)

		// Setup expectations - Boot() checks for dependencies first
		mockProvider.EXPECT().Requires().Return([]string{}).Once() // No dependencies
		mockProvider.EXPECT().Register(app).Once()
		mockProvider.EXPECT().Boot(app).Once()

		app.Register(mockProvider)
		err := app.Boot()

		assert.NoError(t, err)
		mockProvider.AssertExpectations(t)
	})

	t.Run("boots providers with dependencies using dependency-aware registration", func(t *testing.T) {
		app := New(nil)

		// Providers với dependencies
		provider1 := &mockServiceProvider{
			name:      "provider1",
			providers: []string{"service1"},
			requires:  []string{},
		}
		provider2 := &mockServiceProvider{
			name:      "provider2",
			providers: []string{"service2"},
			requires:  []string{"service1"}, // Depends on provider1
		}

		app.Register(provider1)
		app.Register(provider2)

		err := app.Boot()
		assert.NoError(t, err)

		// Verify providers booted trong thứ tự dependency
		assert.True(t, provider1.registered)
		assert.True(t, provider1.booted)
		assert.True(t, provider2.registered)
		assert.True(t, provider2.booted)

		// Verify boot order
		assert.True(t, provider1.bootTime.Before(provider2.bootTime))
	})
}

// Test DI container methods
func TestApplication_DIContainerMethods(t *testing.T) {
	app := New(nil)

	t.Run("Bind and Make", func(t *testing.T) {
		app.Bind("test-service", func(c di.Container) interface{} {
			return "test-value"
		})

		result, err := app.Make("test-service")
		assert.NoError(t, err)
		assert.Equal(t, "test-value", result)
	})

	t.Run("Singleton", func(t *testing.T) {
		counter := 0
		app.Singleton("counter", func(c di.Container) interface{} {
			counter++
			return counter
		})

		result1, err1 := app.Make("counter")
		result2, err2 := app.Make("counter")

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Equal(t, 1, result1)
		assert.Equal(t, 1, result2) // Same instance
	})

	t.Run("Instance", func(t *testing.T) {
		testValue := "instance-value"
		app.Instance("test-instance", testValue)

		result, err := app.Make("test-instance")
		assert.NoError(t, err)
		assert.Equal(t, testValue, result)
	})

	t.Run("Alias", func(t *testing.T) {
		app.Instance("original", "value")
		app.Alias("original", "alias")

		result, err := app.Make("alias")
		assert.NoError(t, err)
		assert.Equal(t, "value", result)
	})

	t.Run("MustMake success", func(t *testing.T) {
		app.Instance("must-make-test", "success")

		result := app.MustMake("must-make-test")
		assert.Equal(t, "success", result)
	})

	t.Run("MustMake panic", func(t *testing.T) {
		assert.Panics(t, func() {
			app.MustMake("non-existent")
		})
	})
}

// Test Call method
func TestApplication_Call(t *testing.T) {
	app := New(nil)

	t.Run("call function with DI", func(t *testing.T) {
		app.Bind("string", func(c di.Container) interface{} {
			return "injected-value"
		})

		testFunc := func(dep string) string {
			return "result-" + dep
		}

		results, err := app.Call(testFunc)
		require.NoError(t, err)
		require.Len(t, results, 1)
		assert.Equal(t, "result-injected-value", results[0])
	})

	t.Run("call function with additional params", func(t *testing.T) {
		testFunc := func(param1 string, param2 int) string {
			return param1 + "-" + string(rune(param2+'0'))
		}

		results, err := app.Call(testFunc, "test", 5)
		require.NoError(t, err)
		require.Len(t, results, 1)
		assert.Equal(t, "test-5", results[0])
	})
}

// Test helper functions
func TestGetProviderKey(t *testing.T) {
	provider := &testProvider{}
	key := getProviderKey(provider)

	assert.NotEmpty(t, key)
	assert.Contains(t, key, "testProvider")
}

// Test provider for testing
type testProvider struct{}

func (p *testProvider) Register(app di.Application) {}
func (p *testProvider) Boot(app di.Application)     {}
func (p *testProvider) Requires() []string          { return []string{} }
func (p *testProvider) Providers() []string         { return []string{"test"} }

// Trackable provider for order testing
type trackableProvider struct {
	name      string
	providers []string
	requires  []string
	bootOrder *[]string
}

func (p *trackableProvider) Register(app di.Application) {}
func (p *trackableProvider) Boot(app di.Application) {
	if p.bootOrder != nil {
		*p.bootOrder = append(*p.bootOrder, p.name)
	}
}
func (p *trackableProvider) Requires() []string  { return p.requires }
func (p *trackableProvider) Providers() []string { return p.providers }

// Integration test
func TestApplication_Integration(t *testing.T) {
	t.Run("full application lifecycle", func(t *testing.T) {
		config := map[string]interface{}{
			"name": "integration-test",
		}

		app := New(config)

		// Create a test service provider
		provider := &testProvider{}
		app.Register(provider)

		// Test the full lifecycle
		err := app.Boot()
		assert.NoError(t, err)

		// Verify app config is accessible
		appConfig, err := app.Make("app.config")
		assert.NoError(t, err)
		assert.Equal(t, config, appConfig)

		// Verify module loader is working
		loader := app.ModuleLoader()
		assert.NotNil(t, loader)
	})
}

// Benchmark tests
func BenchmarkNew(b *testing.B) {
	config := map[string]interface{}{
		"name": "benchmark-test",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		New(config)
	}
}

func BenchmarkApplication_Make(b *testing.B) {
	app := New(nil)
	app.Instance("benchmark-service", "test-value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = app.Make("benchmark-service")
	}
}

func BenchmarkApplication_RegisterAndBoot(b *testing.B) {
	for i := 0; i < b.N; i++ {
		app := New(nil)
		provider := &testProvider{}
		app.Register(provider)
		_ = app.Boot()
	}
}

// Error case tests
func TestApplication_ErrorCases(t *testing.T) {
	t.Run("topological sort with complex dependencies", func(t *testing.T) {
		app := New(nil)

		// Create a complex dependency chain: A -> B -> C
		providerA := di_mocks.NewMockServiceProvider(t)
		providerB := di_mocks.NewMockServiceProvider(t)
		providerC := di_mocks.NewMockServiceProvider(t)

		providerC.EXPECT().Providers().Return([]string{"serviceC"}).Times(1)
		providerC.EXPECT().Requires().Return([]string{}).Times(1)
		providerC.EXPECT().Register(app).Once()

		providerB.EXPECT().Providers().Return([]string{"serviceB"}).Times(1)
		providerB.EXPECT().Requires().Return([]string{"serviceC"}).Times(1)
		providerB.EXPECT().Register(app).Once()

		providerA.EXPECT().Providers().Return([]string{"serviceA"}).Times(1)
		providerA.EXPECT().Requires().Return([]string{"serviceB"}).Times(1)
		providerA.EXPECT().Register(app).Once()

		// Register in wrong order
		app.Register(providerA)
		app.Register(providerB)
		app.Register(providerC)

		err := app.RegisterWithDependencies()
		assert.NoError(t, err)

		providerA.AssertExpectations(t)
		providerB.AssertExpectations(t)
		providerC.AssertExpectations(t)
	})

	t.Run("provider with empty providers list", func(t *testing.T) {
		app := New(nil)

		provider := di_mocks.NewMockServiceProvider(t)
		provider.EXPECT().Providers().Return([]string{}).Times(1)
		provider.EXPECT().Requires().Return([]string{}).Times(1)
		provider.EXPECT().Register(app).Once()

		app.Register(provider)
		err := app.RegisterWithDependencies()

		assert.NoError(t, err)
		provider.AssertExpectations(t)
	})
}

// Test workflow consistency
func TestApplication_WorkflowConsistency(t *testing.T) {
	t.Run("BootServiceProviders uses sorted providers after RegisterWithDependencies", func(t *testing.T) {
		app := New(nil)

		// Track boot order
		var bootOrder []string

		providerA := di_mocks.NewMockServiceProvider(t)
		providerB := di_mocks.NewMockServiceProvider(t)

		// Provider A provides "serviceA", requires nothing
		providerA.EXPECT().Providers().Return([]string{"serviceA"}).Times(1)
		providerA.EXPECT().Requires().Return([]string{}).Times(1)
		providerA.EXPECT().Register(app).Once()
		providerA.EXPECT().Boot(app).Once().Run(func(args mock.Arguments) {
			bootOrder = append(bootOrder, "A")
		})

		// Provider B provides "serviceB", requires "serviceA"
		providerB.EXPECT().Providers().Return([]string{"serviceB"}).Times(1)
		providerB.EXPECT().Requires().Return([]string{"serviceA"}).Times(1)
		providerB.EXPECT().Register(app).Once()
		providerB.EXPECT().Boot(app).Once().Run(func(args mock.Arguments) {
			bootOrder = append(bootOrder, "B")
		})

		// Register providers in wrong order
		app.Register(providerB)
		app.Register(providerA)

		// Register with dependencies first
		err := app.RegisterWithDependencies()
		assert.NoError(t, err)

		// Then boot - should use dependency order
		err = app.BootServiceProviders()
		assert.NoError(t, err)

		// Verify boot order: A should be booted before B
		assert.Equal(t, []string{"A", "B"}, bootOrder)

		providerA.AssertExpectations(t)
		providerB.AssertExpectations(t)
	})

	t.Run("BootServiceProviders uses registration order without RegisterWithDependencies", func(t *testing.T) {
		app := New(nil)

		// Track boot order
		var bootOrder []string

		// Create trackable providers
		providerA := &trackableProvider{
			name:      "A",
			providers: []string{"testA"},
			requires:  []string{},
			bootOrder: &bootOrder,
		}
		providerB := &trackableProvider{
			name:      "B",
			providers: []string{"testB"},
			requires:  []string{},
			bootOrder: &bootOrder,
		}

		// Register providers in specific order
		app.Register(providerB) // B first
		app.Register(providerA) // A second

		// Boot without RegisterWithDependencies - should use registration order
		err := app.BootServiceProviders()
		assert.NoError(t, err)

		// Verify boot order follows registration: B then A
		assert.Equal(t, []string{"B", "A"}, bootOrder)
	})
}

// Test edge cases
func TestApplication_EdgeCases(t *testing.T) {
	t.Run("multiple calls to RegisterServiceProviders", func(t *testing.T) {
		app := New(nil)
		provider := di_mocks.NewMockServiceProvider(t)

		// Should be called multiple times
		provider.EXPECT().Register(app).Times(3)

		app.Register(provider)

		err1 := app.RegisterServiceProviders()
		err2 := app.RegisterServiceProviders()
		err3 := app.RegisterServiceProviders()

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NoError(t, err3)
		provider.AssertExpectations(t)
	})

	t.Run("empty providers slice", func(t *testing.T) {
		app := New(nil)

		err := app.RegisterServiceProviders()
		assert.NoError(t, err)

		err = app.BootServiceProviders()
		assert.NoError(t, err)
	})
}
