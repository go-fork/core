package core

import (
	"fmt"
	"reflect"

	"go.fork.vn/config"
	"go.fork.vn/di"
	"go.fork.vn/log"
)

// Application định nghĩa interface cho Fork application, extends di.Application
// với các chức năng bổ sung cần thiết cho framework.
//
// Interface này mở rộng di.Application interface với:
//   - Module loader access cho advanced bootstrapping
//   - Service provider utilities
//   - Convenience methods cho config và log access
//
// Application instance cần implement interface này để tương thích với
// toàn bộ hệ sinh thái Fork framework và các service providers.
type Application interface {
	// Embed di.Application interface để có tất cả DI functionality
	di.Application

	// Config trả về config manager instance.
	//
	// Phương thức này cung cấp truy cập trực tiếp đến configuration manager
	// mà không cần resolve từ container. Đây là convenience method cho
	// các service providers và business logic cần access config thường xuyên.
	//
	// Trả về:
	//   - config.Manager: Config manager instance
	//
	// Panics:
	//   - Nếu config manager chưa được đăng ký hoặc không đúng type
	//
	// Ví dụ:
	//   - cfg := app.Config()
	//   - dbHost, _ := cfg.GetString("database.host")
	Config() config.Manager

	// Log trả về log manager instance.
	//
	// Phương thức này cung cấp truy cập trực tiếp đến logging manager
	// mà không cần resolve từ container. Đây là convenience method cho
	// việc logging trong toàn bộ application.
	//
	// Trả về:
	//   - log.Manager: Log manager instance
	//
	// Panics:
	//   - Nếu log manager chưa được đăng ký hoặc không đúng type
	//
	// Ví dụ:
	//   - logger := app.Log()
	//   - logger.Info("Application started")
	Log() log.Manager

	// ModuleLoader trả về module loader instance để quản lý việc load và bootstrap modules.
	//
	// Module loader cung cấp các phương thức high-level để:
	//   - Đăng ký core service providers
	//   - Bootstrap application với proper dependency ordering
	//   - Load individual hoặc multiple modules
	//
	// Trả về:
	//   - di.ModuleLoaderContract: Module loader instance
	ModuleLoader() ModuleLoaderContract
}

// application là concrete implementation của Application interface.
//
// Struct này implement tất cả các phương thức cần thiết từ di.Application interface
// và các extension methods từ Application interface.
//
// Fields:
//   - container: DI container instance để quản lý dependencies
//   - providers: Slice các registered service providers
//   - booted: Flag đánh dấu providers đã được booted
//   - loader: Module loader instance
type application struct {
	container       di.Container
	providers       []di.ServiceProvider
	sortedProviders []di.ServiceProvider // Providers sorted by dependency order
	booted          bool
	loader          ModuleLoaderContract
}

// New tạo một Application instance mới với config chỉ định.
//
// Hàm này khởi tạo application với:
//   - DI container mới
//   - Empty providers slice
//   - Config được chỉ định
//   - Module loader được configured
//
// Tham số:
//   - config: map[string]interface{} - Cấu hình cho ứng dụng
//
// Trả về:
//   - Application: Interface implementation
//
// Ví dụ:
//
//	config := map[string]interface{}{
//	    "name": "myapp",
//	    "path": "./configs",
//	}
//	app := app.New(config)

func New(config map[string]interface{}) Application {
	// Validate config
	if config == nil {
		config = make(map[string]interface{})
	}

	// Tạo DI container
	container := di.New()

	// Khởi tạo app instance
	a := &application{
		container:       container,
		providers:       make([]di.ServiceProvider, 0),
		sortedProviders: make([]di.ServiceProvider, 0),
		booted:          false,
	}

	// Register app config
	a.Instance("app.config", config)

	// Tạo module loader với app instance
	a.loader = newModuleLoader(a)

	return a
}

// Container trả về DI container instance.
//
// Implement di.Application interface method.
//
// Trả về:
//   - di.Container: DI container instance
func (a *application) Container() di.Container {
	return a.container
}

// Config trả về config manager instance.
//
// Implement Application interface method.
//
// Trả về:
//   - config.Manager: Config manager instance
//
// Panics:
//   - Nếu config manager chưa được đăng ký hoặc không đúng type
func (a *application) Config() config.Manager {
	return a.container.MustMake("config").(config.Manager)
}

// Log trả về log manager instance.
//
// Implement Application interface method.
//
// Trả về:
//   - log.Manager: Log manager instance
//
// Panics:
//   - Nếu log manager chưa được đăng ký hoặc không đúng type
func (a *application) Log() log.Manager {
	return a.container.MustMake("log").(log.Manager)
}

// ModuleLoader trả về module loader instance.
//
// Implement Application interface method.
//
// Trả về:
//   - di.ModuleLoaderContract: Module loader instance
func (a *application) ModuleLoader() ModuleLoaderContract {
	return a.loader
}

// RegisterServiceProviders đăng ký tất cả service providers đã add vào application.
//
// Implement di.Application interface method.
//
// Trả về:
//   - error: Lỗi nếu có provider registration thất bại
func (a *application) RegisterServiceProviders() error {
	for _, provider := range a.providers {
		provider.Register(a)
	}
	return nil
}

// RegisterWithDependencies đăng ký providers theo thứ tự dependency.
//
// Implement di.Application interface method. Phương thức này sắp xếp
// các providers theo dependency requirements và đăng ký theo thứ tự đúng.
//
// Sử dụng topological sort để xác định thứ tự đăng ký dựa trên:
//   - Requires() method của mỗi provider
//   - Providers() method để biết provider nào cung cấp service nào
//
// Trả về:
//   - error: Lỗi nếu có circular dependency hoặc missing dependency
func (a *application) RegisterWithDependencies() error {
	// Xây dựng dependency graph
	providerMap := make(map[string]di.ServiceProvider)
	serviceToProvider := make(map[string]string)

	// Bước 1: Map providers và services
	for _, provider := range a.providers {
		// Tạo unique key cho provider (sử dụng type name)
		providerKey := getProviderKey(provider)
		providerMap[providerKey] = provider

		// Map services tới provider
		for _, service := range provider.Providers() {
			serviceToProvider[service] = providerKey
		}
	}

	// Bước 2: Topological sort
	sortedProviders, err := a.topologicalSort(providerMap, serviceToProvider)
	if err != nil {
		return err
	}

	// Bước 3: Lưu sorted providers để dùng cho boot
	a.sortedProviders = sortedProviders

	// Bước 4: Đăng ký theo thứ tự sorted
	for _, provider := range sortedProviders {
		provider.Register(a)
	}

	return nil
}

// BootServiceProviders boot tất cả service providers đã đăng ký.
//
// Implement di.Application interface method.
//
// Boot theo thứ tự dependency nếu đã có sortedProviders từ RegisterWithDependencies(),
// nếu không thì boot theo thứ tự đăng ký thông thường.
//
// Trả về:
//   - error: Lỗi nếu có provider boot thất bại
func (a *application) BootServiceProviders() error {
	if a.booted {
		return nil
	}

	// Sử dụng sorted providers nếu có, không thì dùng providers gốc
	providersToBoot := a.providers
	if len(a.sortedProviders) > 0 {
		providersToBoot = a.sortedProviders
	}

	for _, provider := range providersToBoot {
		provider.Boot(a)
	}

	a.booted = true
	return nil
}

// Register đăng ký một service provider vào application.
//
// Implement di.Application interface method.
//
// Tham số:
//   - provider: di.ServiceProvider - Provider cần đăng ký
func (a *application) Register(provider di.ServiceProvider) {
	if provider == nil {
		panic("service provider cannot be nil")
	}
	a.providers = append(a.providers, provider)
}

// Boot khởi động tất cả service providers với smart dependency handling.
//
// Implement di.Application interface method. Đây là shortcut method
// để đăng ký và boot tất cả providers trong một lần gọi.
//
// Method này sẽ:
//  1. Kiểm tra xem có dependencies phức tạp không
//  2. Nếu có dependencies, sử dụng RegisterWithDependencies()
//  3. Nếu không, sử dụng RegisterServiceProviders() đơn giản
//
// Trả về:
//   - error: Lỗi nếu registration hoặc boot thất bại
func (a *application) Boot() error {
	// Kiểm tra xem có cần dependency-aware registration không
	needsDependencyAware := a.hasDependencies()

	var err error
	if needsDependencyAware {
		err = a.RegisterWithDependencies()
	} else {
		err = a.RegisterServiceProviders()
	}

	if err != nil {
		return err
	}
	return a.BootServiceProviders()
}

// Bind đăng ký binding vào container.
//
// Implement di.Application interface method.
//
// Tham số:
//   - abstract: string - Abstract type name
//   - concrete: di.BindingFunc - Factory function
func (a *application) Bind(abstract string, concrete di.BindingFunc) {
	a.container.Bind(abstract, concrete)
}

// Singleton đăng ký singleton binding vào container.
//
// Implement di.Application interface method.
//
// Tham số:
//   - abstract: string - Abstract type name
//   - concrete: di.BindingFunc - Factory function
func (a *application) Singleton(abstract string, concrete di.BindingFunc) {
	a.container.Singleton(abstract, concrete)
}

// Instance đăng ký instance vào container.
//
// Implement di.Application interface method.
//
// Tham số:
//   - abstract: string - Abstract type name
//   - instance: interface{} - Instance object
func (a *application) Instance(abstract string, instance interface{}) {
	a.container.Instance(abstract, instance)
}

// Alias đăng ký alias cho abstract type.
//
// Implement di.Application interface method.
//
// Tham số:
//   - abstract: string - Original abstract name
//   - alias: string - Alias name
func (a *application) Alias(abstract, alias string) {
	a.container.Alias(abstract, alias)
}

// Make resolve dependency từ container.
//
// Implement di.Application interface method.
//
// Tham số:
//   - abstract: string - Abstract type name
//
// Trả về:
//   - interface{}: Resolved instance
//   - error: Lỗi nếu resolve thất bại
func (a *application) Make(abstract string) (interface{}, error) {
	return a.container.Make(abstract)
}

// MustMake resolve dependency từ container, panic nếu lỗi.
//
// Implement di.Application interface method.
//
// Tham số:
//   - abstract: string - Abstract type name
//
// Trả về:
//   - interface{}: Resolved instance
func (a *application) MustMake(abstract string) interface{} {
	return a.container.MustMake(abstract)
}

// Call gọi function với auto dependency injection.
//
// Implement di.Application interface method.
//
// Tham số:
//   - callback: interface{} - Function để gọi
//   - additionalParams: ...interface{} - Additional parameters
//
// Trả về:
//   - []interface{}: Function return values
//   - error: Lỗi nếu call thất bại
func (a *application) Call(callback interface{}, additionalParams ...interface{}) ([]interface{}, error) {
	return a.container.Call(callback, additionalParams...)
}

// hasDependencies kiểm tra xem có provider nào có dependencies không.
//
// Trả về true nếu có ít nhất một provider có requires dependencies,
// false nếu tất cả providers không có dependencies.
//
// Trả về:
//   - bool: true nếu cần dependency-aware registration
func (a *application) hasDependencies() bool {
	for _, provider := range a.providers {
		if len(provider.Requires()) > 0 {
			return true
		}
	}
	return false
}

// Helper methods for dependency ordering

// getProviderKey trả về unique key cho một service provider.
//
// Sử dụng reflection để lấy type name kết hợp với memory address để đảm bảo uniqueness.
//
// Tham số:
//   - provider: di.ServiceProvider - Provider cần tạo key
//
// Trả về:
//   - string: Unique key cho provider
func getProviderKey(provider di.ServiceProvider) string {
	return fmt.Sprintf("%s@%p", reflect.TypeOf(provider).String(), provider)
}

// topologicalSort sắp xếp providers theo dependency order.
//
// Sử dụng Kahn's algorithm để detect cycles và sort providers.
//
// Tham số:
//   - providerMap: map[string]di.ServiceProvider - Map provider key tới provider
//   - serviceToProvider: map[string]string - Map service name tới provider key
//
// Trả về:
//   - []di.ServiceProvider: Sorted providers list
//   - error: Lỗi nếu có circular dependency hoặc missing dependency
func (a *application) topologicalSort(providerMap map[string]di.ServiceProvider, serviceToProvider map[string]string) ([]di.ServiceProvider, error) {
	// Build adjacency list và in-degree count
	adjList := make(map[string][]string)
	inDegree := make(map[string]int)

	// Initialize all providers
	for providerKey := range providerMap {
		adjList[providerKey] = make([]string, 0)
		inDegree[providerKey] = 0
	}

	// Build dependency graph
	for providerKey, provider := range providerMap {
		requires := provider.Requires()
		for _, requiredService := range requires {
			// Tìm provider cung cấp required service
			if requiredProviderKey, exists := serviceToProvider[requiredService]; exists {
				// requiredProvider -> currentProvider dependency
				adjList[requiredProviderKey] = append(adjList[requiredProviderKey], providerKey)
				inDegree[providerKey]++
			} else {
				return nil, fmt.Errorf("required service '%s' not provided by any registered provider (required by %s)", requiredService, providerKey)
			}
		}
	}

	// Kahn's algorithm
	queue := make([]string, 0)
	result := make([]di.ServiceProvider, 0)

	// Find all providers với in-degree 0
	for providerKey, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, providerKey)
		}
	}

	// Process queue
	for len(queue) > 0 {
		// Dequeue
		current := queue[0]
		queue = queue[1:]

		// Add to result
		result = append(result, providerMap[current])

		// Update neighbors
		for _, neighbor := range adjList[current] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	// Check for cycles
	if len(result) != len(providerMap) {
		return nil, fmt.Errorf("circular dependency detected among service providers")
	}

	return result, nil
}
