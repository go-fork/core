package core

import (
	"fmt"

	"go.fork.vn/config"
	"go.fork.vn/di"
	"go.fork.vn/log"
)

type ModuleLoaderContract interface {
	di.ModuleLoaderContract
}

// moduleLoader implement di.ModuleLoaderContract để quản lý việc load và bootstrap modules.
//
// Struct này cung cấp các phương thức high-level để:
//   - Đăng ký core service providers (config, log)
//   - Bootstrap application với proper workflow
//   - Load individual hoặc multiple modules/providers
//   - Quản lý dependency ordering (future enhancement)
type moduleLoader struct {
	app Application
}

// newModuleLoader tạo module loader instance cho application.
//
// Tham số:
//   - app: Application - Application instance
//
// Trả về:
//   - di.ModuleLoaderContract: Module loader implementation
func newModuleLoader(app Application) ModuleLoaderContract {
	return &moduleLoader{
		app: app,
	}
}

// BootstrapApplication khởi tạo application với workflow hoàn chỉnh.
//
// Implement di.ModuleLoaderContract interface method.
//
// Workflow:
//  1. Đăng ký core service providers (config, log)
//  2. Đăng ký tất cả service providers đã add
//  3. Boot tất cả service providers
//
// Trả về:
//   - error: Lỗi nếu bất kỳ bước nào thất bại
func (l *moduleLoader) BootstrapApplication() error {
	// Step 1: Register core providers
	if err := l.RegisterCoreProviders(); err != nil {
		return err
	}

	// Step 2: Register ALL providers với dependency checking
	if err := l.app.RegisterWithDependencies(); err != nil {
		return err
	}

	// Step 4: Boot all providers
	if err := l.app.BootServiceProviders(); err != nil {
		return err
	}

	return nil
}

// RegisterCoreProviders đăng ký các core service providers cần thiết.
//
// Implement di.ModuleLoaderContract interface method.
//
// Core providers bao gồm:
//   - config.ServiceProvider: Configuration management
//   - log service: Logging functionality với proper config binding
//
// Trả về:
//   - error: Lỗi nếu đăng ký core providers thất bại
func (l *moduleLoader) RegisterCoreProviders() error {
	// 1. Register config provider vào list
	configProvider := config.NewServiceProvider()
	// l.app.Register(configProvider)

	// 2. Register ngay config provider để có thể apply config
	configProvider.Register(l.app)

	// 3. Apply config sau khi config provider đã register
	if err := l.applyConfig(); err != nil {
		return err
	}

	// 4. Register log provider vào list
	l.app.Register(log.NewServiceProvider())

	return nil
}

func (l *moduleLoader) applyConfig() error {
	// Lấy config từ DI container với safe type assertion
	configInterface, err := l.app.Container().Make("app.config")
	if err != nil {
		return fmt.Errorf("app.config not found: %w", err)
	}

	cfg, ok := configInterface.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid app.config type: expected map[string]interface{}, got %T", configInterface)
	}

	// Lấy config manager với safe type assertion
	configManagerInterface, err := l.app.Container().Make("config")
	if err != nil {
		return fmt.Errorf("config manager not found: %w", err)
	}

	configManager, ok := configManagerInterface.(config.Manager)
	if !ok {
		return fmt.Errorf("invalid config manager type: expected config.Manager, got %T", configManagerInterface)
	}

	// Apply config settings safely
	if file, ok := cfg["file"].(string); ok {
		configManager.SetConfigFile(file)
		// Read config with error handling
		if err := configManager.ReadInConfig(); err != nil {
			return fmt.Errorf("config read failed: %w", err)
		}
		return nil
	} else {
		if name, ok := cfg["name"].(string); ok {
			configManager.SetConfigName(name)
		}
		if path, ok := cfg["path"].(string); ok {
			configManager.AddConfigPath(path)
		}
		if fileType, ok := cfg["type"].(string); ok {
			configManager.SetConfigType(fileType)
		}
		// Read config with error handling
		if err := configManager.ReadInConfig(); err != nil {
			return fmt.Errorf("config read failed: %w", err)
		}
	}

	return nil
}

// LoadModule tải một module/provider vào application.
//
// Implement di.ModuleLoaderContract interface method.
//
// Phương thức này:
//  1. Kiểm tra module có phải là ServiceProvider không
//  2. Đăng ký module vào application
//  3. Boot module nếu application đã được booted
//
// Tham số:
//   - module: interface{} - Module cần load (phải là di.ServiceProvider)
//
// Trả về:
//   - error: Lỗi nếu module không hợp lệ hoặc load thất bại
func (l *moduleLoader) LoadModule(module interface{}) error {
	// Kiểm tra module có phải ServiceProvider không
	provider, ok := module.(di.ServiceProvider)
	if !ok {
		return &ModuleLoadError{
			Module: module,
			Reason: "module must implement di.ServiceProvider interface",
		}
	}

	// Đăng ký provider
	l.app.Register(provider)

	// Nếu app đã booted, cần register và boot provider mới ngay
	if l.isAppBooted() {
		provider.Register(l.app)
		provider.Boot(l.app)
	}

	return nil
}

// LoadModules tải nhiều modules vào application.
//
// Implement di.ModuleLoaderContract interface method.
//
// Tham số:
//   - modules: ...interface{} - Danh sách modules cần load
//
// Trả về:
//   - error: Lỗi nếu bất kỳ module nào load thất bại
func (l *moduleLoader) LoadModules(modules ...interface{}) error {
	for i, module := range modules {
		if err := l.LoadModule(module); err != nil {
			return &MultiModuleLoadError{
				FailedIndex:  i,
				FailedModule: module,
				Cause:        err,
			}
		}
	}
	return nil
}

// isAppBooted kiểm tra xem application đã được booted chưa.
//
// Trả về true nếu app đã booted, false nếu chưa.
func (l *moduleLoader) isAppBooted() bool {
	// Chúng ta không thể truy cập trực tiếp private field booted
	// Thay vào đó, chúng ta có thể sử dụng một thuộc tính khác để xác định trạng thái booted
	// Ví dụ: Kiểm tra xem log service đã được đăng ký chưa (một dịch vụ core được đăng ký khi boot)
	_, err := l.app.Make("log")
	return err == nil
}

// ModuleLoadError represent lỗi khi load một module.
//
// Error type này cung cấp thông tin chi tiết về module nào gây lỗi
// và lý do tại sao việc load thất bại.
type ModuleLoadError struct {
	Module interface{}
	Reason string
}

// Error implement error interface.
//
// Trả về:
//   - string: Error message với thông tin chi tiết
func (e *ModuleLoadError) Error() string {
	return "failed to load module: " + e.Reason
}

// MultiModuleLoadError represent lỗi khi load multiple modules.
//
// Error type này cung cấp thông tin về module nào trong danh sách
// gây ra lỗi và nguyên nhân gốc.
type MultiModuleLoadError struct {
	FailedIndex  int
	FailedModule interface{}
	Cause        error
}

// Error implement error interface.
//
// Trả về:
//   - string: Error message với thông tin chi tiết
func (e *MultiModuleLoadError) Error() string {
	return fmt.Sprintf("failed to load modules: error at index %d, module failed: %v", e.FailedIndex, e.Cause)
}
