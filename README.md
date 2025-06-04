![go.fork.vn Logo](https://go.fork.vn/assets/logo.svg)

# go.fork.vn/core

**Package nền tảng của hệ sinh thái Go.Fork với hệ thống Dependency Injection mạnh mẽ**

[![Go Reference](https://pkg.go.dev/badge/go.fork.vn/core.svg)](https://pkg.go.dev/go.fork.vn/core)
[![Go Report Card](https://goreportcard.com/badge/go.fork.vn/core)](https://goreportcard.com/report/go.fork.vn/core)
[![CI Status](https://github.com/go-fork/core/actions/workflows/test.yml/badge.svg)](https://github.com/go-fork/core/actions)
[![GitHub release](https://img.shields.io/github/v/release/go-fork/core)](https://github.com/go-fork/core/releases)
[![License](https://img.shields.io/github/license/go-fork/core)](https://github.com/go-fork/core/blob/main/LICENSE)

## 📖 Giới thiệu

**go.fork.vn/core** là package nền tảng của hệ sinh thái go.fork.vn, cung cấp các interface và implementation cốt lõi để xây dựng ứng dụng Go hiện đại với kiến trúc modular, dependency injection và quản lý lifecycle tự động.

Package này xây dựng dựa trên [go.fork.vn/di](https://github.com/go-fork/di) và mở rộng với các tính năng mạnh mẽ hơn như:

- **Smart Dependency Management**: Sắp xếp và phát hiện circular dependencies tự động
- **Module Loading System**: Quản lý bootstrap và loading của modules
- **Core Providers**: Tích hợp sẵn với Config và Log services

## ⚙️ Cài đặt

```bash
go get go.fork.vn/core
```

## 🚀 Quick Start

```go
package main

import (
    "log"
    
    "go.fork.vn/core"
    "go.fork.vn/config"
    "go.fork.vn/log"
)

func main() {
    // Tạo application instance
    app := core.New()
    
    // Lấy module loader
    loader := app.ModuleLoader()
    
    // Bootstrap application với workflow chuẩn
    if err := loader.BootstrapApplication(); err != nil {
        log.Fatal("Failed to bootstrap application:", err)
    }
    
    // Sử dụng các services
    app.Call(func(
        cfg config.Manager,
        logger log.Manager,
    ) {
        logger.Info("Application started with config", 
            "name", cfg.GetString("app.name"),
            "env", cfg.GetString("app.environment"),
        )
    })
    
    // Application sẵn sàng sử dụng
    // ...
}
```

## 🏗️ Kiến trúc

```mermaid
graph TB
    subgraph "Application Layer"
        App[Application Interface]
        Impl[Application Implementation]
        App --> Impl
    end
    
    subgraph "Core Services"
        Container[DI Container]
        Loader[Module Loader]
        Smart[Smart Dependencies]
    end
    
    subgraph "Service Providers"
        Config[Config Provider]
        Log[Log Provider]
        Custom[Custom Providers]
    end
    
    Impl --> Container
    Impl --> Loader
    Impl --> Smart
    
    Smart --> Config
    Smart --> Log
    Smart --> Custom
    
    Loader --> Config
    Loader --> Custom
    
    style App fill:#ff9999,stroke:#333,stroke-width:2px
    style Container fill:#99ccff,stroke:#333,stroke-width:2px
    style Smart fill:#ffcc99,stroke:#333,stroke-width:2px
    style Loader fill:#99ff99,stroke:#333,stroke-width:2px
```

## 📄 Tài liệu

### Core Documentation

- [**Tổng quan Hệ thống**](docs/overview.md) - Kiến trúc tổng thể và tính năng
- [**Application Interface**](docs/application.md) - Chi tiết về Application interface và implementation
- [**Module Loader**](docs/loader.md) - Hệ thống Module Loader
- [**Workflows**](docs/workflows.md) - Quy trình hoạt động và dependency management
- [**Core Providers**](docs/core_providers.md) - Chi tiết về các core providers

### Code Reference

- [GoDoc](https://pkg.go.dev/go.fork.vn/core) - API Reference
- [GitHub Releases](https://github.com/go-fork/core/releases) - Release history và changelog

## 🧩 Các thành phần chính

### 1. Application Interface

Interface trung tâm định nghĩa contract cho quản lý ứng dụng:

```go
type Application interface {
    // Container management
    Container() di.Container
    
    // Service Provider lifecycle
    Register(provider di.ServiceProvider)
    RegisterServiceProviders() error
    RegisterWithDependencies() error
    BootServiceProviders() error
    Boot() error
    
    // Dependency management shortcuts
    Bind(abstract string, concrete di.BindingFunc)
    Singleton(abstract string, concrete di.BindingFunc)
    Instance(abstract string, instance interface{})
    Alias(abstract, alias string)
    
    // Dependency resolution
    Make(abstract string) (interface{}, error)
    MustMake(abstract string) interface{}
    Call(callback interface{}, params ...interface{}) ([]interface{}, error)
    
    // Module loader access
    ModuleLoader() ModuleLoaderContract
}
```

### 2. Module Loader

Quản lý quá trình load, bootstrap và điều phối các module trong ứng dụng:

```go
type ModuleLoaderContract interface {
    // Bootstrap application với workflow chuẩn
    BootstrapApplication() error
    
    // Đăng ký core providers (config, log)
    RegisterCoreProviders() error
    
    // Load các modules và providers
    LoadModule(module interface{}) error
    LoadModules(modules ...interface{}) error
}
```

### 3. Service Provider Pattern

Interface định nghĩa một service provider:

```go
// từ go.fork.vn/di
type ServiceProvider interface {
    // Đăng ký bindings
    Register(app interface{}) error
    
    // Khởi động services
    Boot(app interface{}) error
    
    // Khai báo dependencies
    Requires() []string
    
    // Khai báo services được cung cấp
    Providers() []string
}
```

## 🧪 Ví dụ

### Custom Service Provider

```go
type MyServiceProvider struct{}

func (p *MyServiceProvider) Register(app interface{}) error {
    application, ok := app.(core.Application)
    if !ok {
        return fmt.Errorf("invalid application instance")
    }
    
    application.Singleton("my-service", func(c di.Container) interface{} {
        config := c.MustMake("config").(config.Manager)
        logger := c.MustMake("log").(log.Manager)
        
        return &MyService{
            config: config,
            logger: logger,
        }
    })
    
    return nil
}

func (p *MyServiceProvider) Boot(app interface{}) error {
    // Khởi tạo service khi cần
    return nil
}

func (p *MyServiceProvider) Requires() []string {
    return []string{"config", "log"}
}

func (p *MyServiceProvider) Providers() []string {
    return []string{"my-service"}
}
```

### Sử dụng Smart Dependencies

```go
// Đăng ký multiple providers với dependencies
app := core.New()

// Đăng ký các providers
app.Register(&ConfigProvider{})
app.Register(&LogProvider{})
app.Register(&DatabaseProvider{})  // Requires: config
app.Register(&CacheProvider{})     // Requires: config
app.Register(&AuthProvider{})      // Requires: database, config
app.Register(&APIProvider{})       // Requires: auth, cache, log

// Boot application với smart dependency resolution
if err := app.Boot(); err != nil {
    log.Fatal("Failed to boot application:", err)
}

// Providers sẽ được sắp xếp và boot theo thứ tự:
// ConfigProvider → LogProvider → DatabaseProvider → CacheProvider → AuthProvider → APIProvider
```

### Dynamic Module Loading

```go
app := core.New()
loader := app.ModuleLoader()

// Bootstrap core application
if err := loader.BootstrapApplication(); err != nil {
    log.Fatal("Failed to bootstrap application:", err)
}

// Load custom module sau khi app đã boot
if err := loader.LoadModule(&PaymentServiceProvider{}); err != nil {
    log.Fatal("Failed to load payment module:", err)
}

// Payment service đã sẵn sàng sử dụng
paymentService := app.MustMake("payment.service").(PaymentService)
```

## 📊 Performance

Package được tối ưu cho hiệu năng cao với:

```
BenchmarkRegisterProvider-8      5000000    250 ns/op    32 B/op    1 allocs/op
BenchmarkBootProvider-8          1000000   1200 ns/op   128 B/op    3 allocs/op  
BenchmarkMakeService-8          10000000    120 ns/op     0 B/op    0 allocs/op
BenchmarkCallWithDependency-8    2000000    800 ns/op    64 B/op    2 allocs/op
```

## 📝 Cấu hình

File cấu hình sample (`configs/app.sample.yaml`):

```yaml
app:
  name: "My Application"
  environment: "development"
  debug: true
  
database:
  driver: "postgres"
  dsn: "postgres://user:password@localhost:5432/mydb?sslmode=disable"
  
log:
  level: "debug"
  format: "json"
  output: "stdout"
```

## 🧠 Best Practices

- Luôn sử dụng `ModuleLoader.BootstrapApplication()` để khởi động ứng dụng
- Tổ chức providers thành các nhóm logic liên quan
- Sử dụng `Requires()` và `Providers()` để khai báo dependencies rõ ràng
- Implement error handling toàn diện trong các providers
- Sử dụng `app.Call()` để injection tự động thay vì `app.Make()`

## 🤝 Đóng góp

Contributions are welcome! Hãy follow các bước sau:

1. Fork repo
2. Tạo feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add some amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

Vui lòng đảm bảo rằng code của bạn:
- Pass tất cả các tests (`go test ./...`)
- Có test coverage phù hợp cho code mới
- Tuân thủ code style và conventions

## 📄 License

Distributed under the MIT License. See [LICENSE](LICENSE) for more information.

## 🔗 Related Packages

- [go.fork.vn/di](https://github.com/go-fork/di) - Dependency Injection core
- [go.fork.vn/config](https://github.com/go-fork/config) - Configuration management
- [go.fork.vn/log](https://github.com/go-fork/log) - Logging services
- [go.fork.vn/http](https://github.com/go-fork/http) - HTTP services framework

---

*Made with ❤️ by Go.Fork Team*
