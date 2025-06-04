# Core Providers - go.fork.vn/core

## 🎯 Giới thiệu

**Core Providers** là các service providers được built-in và tự động đăng ký trong go.fork.vn/core. Các providers này cung cấp chức năng cốt lõi cần thiết cho hầu hết các ứng dụng, bao gồm configuration management và logging services.

## 🏗️ Kiến trúc Core Providers

```mermaid
graph TB
    subgraph "Module Loader"
        ML[Module Loader]
    end
    
    subgraph "Core Providers"
        CP[Core Providers]
        CFG[Config Provider]
        LOG[Log Provider]
        
        CP --> CFG
        CP --> LOG
    end
    
    subgraph "Services"
        CM[Config Manager]
        LM[Log Manager]
        
        CFG --> CM
        LOG --> LM
    end
    
    ML --> CP
    
    style ML fill:#ff9999,stroke:#333,stroke-width:2px
    style CP fill:#99ccff,stroke:#333,stroke-width:2px
    style CFG fill:#ffcc99,stroke:#333,stroke-width:2px
    style LOG fill:#99ff99,stroke:#333,stroke-width:2px
```

## 🔄 Registration Workflow

```mermaid
sequenceDiagram
    participant ML as Module Loader
    participant APP as Application
    participant CFG as Config Provider
    participant LOG as Log Provider
    participant CM as Config Manager
    participant LM as Log Manager
    
    ML->>+ML: RegisterCoreProviders()
    
    ML->>+APP: Register(config.NewServiceProvider())
    APP-->>-ML: Success
    
    ML->>+APP: RegisterServiceProviders()
    APP->>+CFG: Register(app)
    CFG->>APP: Bind("config", ConfigManager)
    APP-->>-CFG: Success
    APP-->>-ML: Success
    
    ML->>+ML: applyConfig()
    ML->>+APP: Container().Make("app.config")
    APP-->>-ML: Config map
    ML->>+APP: Container().Make("config")
    APP-->>-ML: ConfigManager
    ML->>+CM: SetConfigFile, SetConfigName, etc.
    CM-->>-ML: Success
    ML->>+CM: ReadInConfig()
    CM-->>-ML: Config loaded
    ML-->>-ML: Success
    
    ML->>+APP: Register(log.NewServiceProvider())
    APP-->>-ML: Success
    
    ML->>+APP: RegisterServiceProviders()
    APP->>+LOG: Register(app)
    LOG->>APP: Bind("log", LogManager)
    APP-->>-LOG: Success
    APP-->>-ML: Success
    
    ML-->>-ML: All core providers registered
```

## 🔧 Config Provider

Config Provider cung cấp unified configuration management cho toàn bộ ứng dụng.

### 1. **ServiceProvider Implementation**

```go
// Tại go.fork.vn/config
type ServiceProvider struct{}

func NewServiceProvider() di.ServiceProvider {
    return &ServiceProvider{}
}

func (p *ServiceProvider) Register(app interface{}) error {
    application, ok := app.(core.Application)
    if !ok {
        return fmt.Errorf("invalid application instance")
    }
    
    // Register config manager as singleton
    application.Singleton("config", func(c di.Container) interface{} {
        manager := config.New()
        return manager
    })
    
    return nil
}

func (p *ServiceProvider) Boot(app interface{}) error {
    // Config provider không cần logic boot phức tạp
    return nil
}

func (p *ServiceProvider) Requires() []string {
    // Không phụ thuộc vào provider nào khác
    return []string{}
}

func (p *ServiceProvider) Providers() []string {
    return []string{"config"}
}
```

### 2. **Config Manager Interface**

```go
// Tại go.fork.vn/config
type Manager interface {
    // File configuration
    SetConfigFile(file string)
    SetConfigName(name string)
    AddConfigPath(path string)
    SetConfigType(fileType string)
    
    // Reading config
    ReadInConfig() error
    Get(key string) interface{}
    GetString(key string) string
    GetInt(key string) int
    GetBool(key string) bool
    
    // Modifying config
    Set(key string, value interface{})
    SetDefault(key string, value interface{})
    
    // Advanced operations
    Unmarshal(rawVal interface{}) error
    UnmarshalKey(key string, rawVal interface{}) error
    AllSettings() map[string]interface{}
    
    // Runtime operations
    WatchConfig()
    OnConfigChange(run func())
    Reload() error
}
```

### 3. **Config File Example**

```yaml
# app.yaml
app:
  name: "My Application"
  debug: true
  environment: "development"
  
database:
  driver: "postgres"
  dsn: "postgres://user:password@localhost:5432/mydb?sslmode=disable"
  
log:
  level: "debug"
  format: "json"
  output: "stdout"
  
http:
  port: 8080
  timeout: 30
  
cache:
  driver: "redis"
  address: "localhost:6379"
```

### 4. **Sử dụng Config Provider**

```go
// Trong provider khác
func (p *DatabaseProvider) Register(app core.Application) error {
    // Get config manager
    configManager := app.MustMake("config").(config.Manager)
    
    // Lấy database config
    dsn := configManager.GetString("database.dsn")
    driver := configManager.GetString("database.driver")
    
    // Register database singleton
    app.Singleton("database", func(c di.Container) interface{} {
        return database.NewConnection(driver, dsn)
    })
    
    return nil
}
```

## 🔧 Log Provider

Log Provider cung cấp khả năng logging mạnh mẽ, với nhiều outputs và log levels.

### 1. **ServiceProvider Implementation**

```go
// Tại go.fork.vn/log
type ServiceProvider struct{}

func NewServiceProvider() di.ServiceProvider {
    return &ServiceProvider{}
}

func (p *ServiceProvider) Register(app interface{}) error {
    application, ok := app.(core.Application)
    if !ok {
        return fmt.Errorf("invalid application instance")
    }
    
    // Register log manager as singleton
    application.Singleton("log", func(c di.Container) interface{} {
        // Get config to configure log manager
        configManager, err := c.Make("config")
        if err != nil {
            // Fallback to default config
            return log.New(log.Config{
                Level:  "info",
                Format: "text",
                Output: "stderr",
            })
        }
        
        cfg, ok := configManager.(config.Manager)
        if !ok {
            // Fallback to default config
            return log.New(log.Config{
                Level:  "info",
                Format: "text",
                Output: "stderr",
            })
        }
        
        // Create log manager with config
        return log.New(log.Config{
            Level:  cfg.GetString("log.level"),
            Format: cfg.GetString("log.format"),
            Output: cfg.GetString("log.output"),
        })
    })
    
    return nil
}

func (p *ServiceProvider) Boot(app interface{}) error {
    return nil
}

func (p *ServiceProvider) Requires() []string {
    // Phụ thuộc vào config provider
    return []string{"config"}
}

func (p *ServiceProvider) Providers() []string {
    return []string{"log"}
}
```

### 2. **Log Manager Interface**

```go
// Tại go.fork.vn/log
type Manager interface {
    // Log levels
    Debug(msg string, fields ...interface{})
    Info(msg string, fields ...interface{})
    Warning(msg string, fields ...interface{}) 
    Error(msg string, fields ...interface{})
    Fatal(msg string, fields ...interface{})
    
    // Context logs
    WithFields(fields map[string]interface{}) Logger
    WithError(err error) Logger
    
    // Output control
    SetOutput(output io.Writer)
    SetLevel(level string) error
    SetFormatter(format string) error
    
    // Filter and hook management
    AddHook(hook Hook) error
    AddFilter(filter Filter) error
}

type Logger interface {
    Debug(msg string, fields ...interface{})
    Info(msg string, fields ...interface{})
    Warning(msg string, fields ...interface{})
    Error(msg string, fields ...interface{})
    Fatal(msg string, fields ...interface{})
}
```

### 3. **Sử dụng Log Provider**

```go
// Trong provider khác
func (p *APIProvider) Boot(app core.Application) error {
    // Get log manager
    logger := app.MustMake("log").(log.Manager)
    
    // Basic logging
    logger.Info("API Provider booted")
    
    // Structured logging
    logger.WithFields(map[string]interface{}{
        "version": "1.0.0",
        "port":    8080,
    }).Info("API server starting")
    
    // Error logging
    if err := startServer(); err != nil {
        logger.WithError(err).Error("Failed to start API server")
        return err
    }
    
    return nil
}
```

## 🔄 Workflow Tích hợp

### 1. **Core Provider Registration**

```mermaid
flowchart TD
    Start([Begin])
    RegisterConfig[Register Config Provider]
    ActivateConfig[Register Config Immediately]
    ApplyConfig[Apply Config Settings]
    ReadConfig[Read Config Files]
    RegisterLog[Register Log Provider]
    ActivateLog[Register Log Immediately]
    End([Core Providers Ready])
    
    Start --> RegisterConfig
    RegisterConfig --> ActivateConfig
    ActivateConfig --> ApplyConfig
    ApplyConfig --> ReadConfig
    ReadConfig --> RegisterLog
    RegisterLog --> ActivateLog
    ActivateLog --> End
    
    style Start fill:#ff9999,stroke:#333,stroke-width:2px
    style ApplyConfig fill:#ffcc99,stroke:#333,stroke-width:2px
    style ReadConfig fill:#99ff99,stroke:#333,stroke-width:2px
    style End fill:#99ccff,stroke:#333,stroke-width:2px
```

### 2. **Hierarchical Dependencies**

```mermaid
graph TD
    A[Application]
    ML[Module Loader]
    CFG[Config Provider]
    LOG[Log Provider]
    CUS[Custom Providers]
    
    A --> ML
    ML --> CFG
    CFG --> LOG
    LOG --> CUS
    
    style A fill:#ff9999,stroke:#333,stroke-width:2px
    style ML fill:#ffcc99,stroke:#333,stroke-width:2px
    style CFG fill:#99ff99,stroke:#333,stroke-width:2px
    style LOG fill:#99ccff,stroke:#333,stroke-width:2px
    style CUS fill:#cc99ff,stroke:#333,stroke-width:2px
```

## 🏆 Best Practices

### 1. **Cấu hình Config Provider**

```go
// In main.go or bootstrap
app := core.New()
loader := app.ModuleLoader()

// Set environment-specific config settings
app.Instance("app.config", map[string]interface{}{
    "file": "config.yaml",
    "name": "config",
    "path": "configs/",
    "type": "yaml",
})

// Bootstrap application
if err := loader.BootstrapApplication(); err != nil {
    log.Fatal("Failed to bootstrap application:", err)
}
```

### 2. **Sử dụng Config & Log Providers Together**

```go
// Trong provider khác
app.Call(func(
    cfg config.Manager,
    logger log.Manager,
) {
    // Lấy giá trị từ config
    appName := cfg.GetString("app.name")
    environment := cfg.GetString("app.environment")
    
    // Log thông tin với structured logging
    logger.WithFields(map[string]interface{}{
        "app_name":    appName,
        "environment": environment,
        "version":     "1.0.0",
        "timestamp":   time.Now().Format(time.RFC3339),
    }).Info("Application started")
})
```

### 3. **Customize Log Format**

```go
// Trong config.yaml
log:
  level: "info"
  format: "json"
  output: "stdout"
  timestamp_format: "2006-01-02 15:04:05"
  include_caller: true
  caller_skip: 1

// Trong code
app.Call(func(logger log.Manager) {
    logger.Info("This will be formatted according to config")
})
```

### 4. **Dynamic Config Changes**

```go
// Trong runtime
app.Call(func(cfg config.Manager, logger log.Manager) {
    // Thay đổi config
    cfg.Set("app.debug", false)
    cfg.Set("log.level", "error")
    
    // Reload config để áp dụng thay đổi
    if err := cfg.Reload(); err != nil {
        logger.WithError(err).Error("Failed to reload config")
        return
    }
    
    logger.Info("Config reloaded with new settings")
})
```

## 📈 Performance Considerations

### 1. **Config Caching**

```go
// Enable config caching
app.Call(func(cfg config.Manager) {
    // Cache config cho các lần khởi động sau
    if err := cfg.WriteConfigToFile("config_cache.json"); err != nil {
        // Handle error
    }
})

// Sử dụng cached config
app.Instance("app.config", map[string]interface{}{
    "file": "config_cache.json",
})
```

### 2. **Log Level Filtering**

```go
// Set log level based on environment
app.Call(func(cfg config.Manager, logger log.Manager) {
    env := cfg.GetString("app.environment")
    
    switch env {
    case "production":
        logger.SetLevel("error") // Only log errors in production
    case "staging":
        logger.SetLevel("warning") // Log warnings and errors in staging
    default:
        logger.SetLevel("debug") // Log everything in development
    }
})
```

## 🛠️ Mở rộng Core Providers

### 1. **Custom Config Adapter**

```go
func CustomConfigProvider() di.ServiceProvider {
    return &customConfigProvider{}
}

type customConfigProvider struct{}

func (p *customConfigProvider) Register(app interface{}) error {
    application, ok := app.(core.Application)
    if !ok {
        return fmt.Errorf("invalid application instance")
    }
    
    // Override standard config provider
    application.Singleton("config", func(c di.Container) interface{} {
        // Implement custom config manager
        return &CustomConfigManager{}
    })
    
    return nil
}

func (p *customConfigProvider) Requires() []string {
    return []string{}
}

func (p *customConfigProvider) Providers() []string {
    return []string{"config"}
}
```

### 2. **Specialized Logger**

```go
func EnhancedLogProvider() di.ServiceProvider {
    return &enhancedLogProvider{}
}

type enhancedLogProvider struct{}

func (p *enhancedLogProvider) Register(app interface{}) error {
    application, ok := app.(core.Application)
    if !ok {
        return fmt.Errorf("invalid application instance")
    }
    
    // Register enhanced logger that augments standard logger
    application.Singleton("enhanced.log", func(c di.Container) interface{} {
        // Get standard logger
        stdLogger, err := c.Make("log")
        if err != nil {
            // Create new logger if standard one not available
            return &EnhancedLogger{logger: log.New(log.Config{})}
        }
        
        // Wrap standard logger with enhanced features
        return &EnhancedLogger{logger: stdLogger.(log.Manager)}
    })
    
    return nil
}

func (p *enhancedLogProvider) Requires() []string {
    return []string{"log"}
}

func (p *enhancedLogProvider) Providers() []string {
    return []string{"enhanced.log"}
}
```

## 🔮 Roadmap

### Planned Features

- **Distributed Configuration**: Support cho etcd, Consul configuration backends
- **Structured Logging Enhancements**: Advanced filtering, dynamic fields
- **Observability Integration**: OpenTelemetry, Prometheus metrics integration
- **Environment-aware Config**: Auto-select config sections based on environment

---

> **Next**: [Index](index.md) - Quay lại trang chủ