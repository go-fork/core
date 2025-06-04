# Workflows và Dependency Management - go.fork.vn/core

## 🎯 Giới thiệu

**Workflows** trong go.fork.vn/core là các quy trình hoạt động cốt lõi xử lý vòng đời của application, service providers, và dependency management. Hệ thống workflows đảm bảo các thành phần được khởi tạo, đăng ký và boot theo đúng thứ tự dependency.

## 🏗️ Kiến trúc Workflows

```mermaid
graph TB
    Start[Khởi tạo Application]
    RegisterProviders[Register Providers]
    DetectDeps[Detect Dependencies]
    SortDeps[Sort Dependencies]
    RegisterWithDeps[Register With Dependencies]
    RegisterNormal[Register Normal]
    BootProviders[Boot Providers]
    AppRunning[Application Running]
    
    Start --> RegisterProviders
    RegisterProviders --> DetectDeps
    
    DetectDeps -->|Has Dependencies| SortDeps
    DetectDeps -->|No Dependencies| RegisterNormal
    
    SortDeps --> RegisterWithDeps
    RegisterWithDeps --> BootProviders
    RegisterNormal --> BootProviders
    
    BootProviders --> AppRunning
    
    style Start fill:#ff9999,stroke:#333,stroke-width:2px
    style DetectDeps fill:#ffcc99,stroke:#333,stroke-width:2px
    style SortDeps fill:#99ff99,stroke:#333,stroke-width:2px
    style BootProviders fill:#99ccff,stroke:#333,stroke-width:2px
    style AppRunning fill:#cc99ff,stroke:#333,stroke-width:2px
```

## 🔄 Application Lifecycle Workflows

### 1. **Initialization Workflow**

```mermaid
sequenceDiagram
    participant Client
    participant App as Application
    participant Container as DI Container
    participant Loader as Module Loader
    
    Client->>+App: New()
    App->>+Container: di.New()
    Container-->>-App: Container instance
    App->>+Loader: newModuleLoader(app)
    Loader-->>-App: Loader instance
    App-->>-Client: Application instance
    
    Note over Client,App: Application instance created
```

### 2. **Provider Registration Workflow**

```mermaid
sequenceDiagram
    participant Client
    participant App as Application
    participant Provider as Service Provider
    
    Client->>+App: Register(provider)
    
    App->>+App: getProviderKey(provider)
    App-->>-App: unique key
    
    App->>+App: map[key] = provider
    App-->>-App: Provider stored
    
    App-->>-Client: Success
    
    Note over Client,App: Provider stored in registry
```

### 3. **Smart Registration Workflow**

```mermaid
sequenceDiagram
    participant Client
    participant App as Application
    participant Smart as Smart Dependency Manager
    participant Provider as Service Providers
    
    Client->>+App: RegisterWithDependencies()
    App->>+Smart: buildDependencyGraph()
    
    loop For each provider
        Smart->>+Provider: Requires()
        Provider-->>-Smart: []string dependencies
    end
    
    Smart->>+Smart: topologicalSort()
    Smart-->>-Smart: Sorted providers
    
    loop For each sorted provider
        Smart->>+Provider: Register(app)
        Provider-->>-Smart: Success/Error
    end
    
    Smart-->>-App: Success/Error
    App-->>-Client: Success/Error
    
    Note over Client,App: Providers registered in dependency order
```

### 4. **Boot Workflow**

```mermaid
sequenceDiagram
    participant Client
    participant App as Application
    participant Provider as Service Providers
    
    Client->>+App: Boot()
    
    alt Has Dependencies
        App->>+App: RegisterWithDependencies()
        App-->>-App: Success/Error
    else No Dependencies
        App->>+App: RegisterServiceProviders()
        App-->>-App: Success/Error
    end
    
    App->>+App: BootServiceProviders()
    
    loop For each provider
        App->>+Provider: Boot(app)
        Provider-->>-App: Success/Error
    end
    
    App-->>-Client: Success/Error
    
    Note over Client,App: Application fully booted
```

## 🧩 Dependency Management

### 1. **Dependency Graph Construction**

```mermaid
graph TD
    A[Provider A<br/>Requires: []]
    B[Provider B<br/>Requires: [A]]
    C[Provider C<br/>Requires: [A, B]]
    D[Provider D<br/>Requires: [A]]
    
    A --> B
    A --> D
    B --> C
    A -.-> C
    
    style A fill:#ff9999,stroke:#333,stroke-width:2px
    style B fill:#ffcc99,stroke:#333,stroke-width:2px
    style C fill:#99ff99,stroke:#333,stroke-width:2px
    style D fill:#99ccff,stroke:#333,stroke-width:2px
```

### 2. **Topological Sort Algorithm**

```mermaid
flowchart TD
    Start([Begin])
    CreateGraph[Create Adjacency List]
    CalcIndegrees[Calculate In-degrees]
    FindZero[Find nodes with zero in-degree]
    QueueZero[Add zero-degree nodes to queue]
    ProcessQueue{Queue empty?}
    DequeueNode[Dequeue node]
    AddResult[Add to result list]
    UpdateDeps[Update dependencies]
    CheckCircular[Check for circular deps]
    End([Finish])
    
    Start --> CreateGraph
    CreateGraph --> CalcIndegrees
    CalcIndegrees --> FindZero
    FindZero --> QueueZero
    QueueZero --> ProcessQueue
    
    ProcessQueue -->|Yes| CheckCircular
    ProcessQueue -->|No| DequeueNode
    
    DequeueNode --> AddResult
    AddResult --> UpdateDeps
    UpdateDeps --> ProcessQueue
    
    CheckCircular -->|Has remaining nodes| End
    CheckCircular -->|No remaining nodes| End
    
    style Start fill:#ff9999,stroke:#333,stroke-width:2px
    style ProcessQueue fill:#ffcc99,stroke:#333,stroke-width:2px
    style AddResult fill:#99ff99,stroke:#333,stroke-width:2px
    style CheckCircular fill:#99ccff,stroke:#333,stroke-width:2px
    style End fill:#cc99ff,stroke:#333,stroke-width:2px
```

### 3. **Dependency Resolution Code**

```go
func (a *application) RegisterWithDependencies() error {
    a.mu.Lock()
    defer a.mu.Unlock()
    
    if len(a.providers) == 0 {
        return nil
    }
    
    // Build dependency graph from providers
    graph, err := a.buildDependencyGraph()
    if err != nil {
        return err
    }
    
    // Sort providers by dependency order
    sorted, err := a.topologicalSort(graph)
    if err != nil {
        return fmt.Errorf("dependency resolution failed: %w", err)
    }
    
    // Cache sorted providers for future use
    a.sortedProviders = sorted
    
    // Register providers in dependency order
    for _, provider := range sorted {
        if err := provider.Register(a); err != nil {
            return fmt.Errorf("provider registration failed: %w", err)
        }
    }
    
    return nil
}
```

### 4. **Circular Dependency Detection**

```mermaid
graph TD
    A[Provider A<br/>Requires: [C]]
    B[Provider B<br/>Requires: [A]]
    C[Provider C<br/>Requires: [B]]
    
    A --> B
    B --> C
    C --> A
    
    style A fill:#ff9999,stroke:#333,stroke-width:2px
    style B fill:#ffcc99,stroke:#333,stroke-width:2px
    style C fill:#99ff99,stroke:#333,stroke-width:2px
    
    ERROR[Circular Dependency Detected!<br/>A → B → C → A] 
    style ERROR fill:#ff0000,stroke:#333,stroke-width:2px,color:#ffffff
```

## 🔧 Service Provider Lifecycle

### 1. **Registration Phase**

```go
// Provider Implementation Example
func (p *MyProvider) Register(app core.Application) error {
    // Step 1: Register bindings
    app.Singleton("my-service", func(c di.Container) interface{} {
        return &MyService{
            config: c.MustMake("config").(config.Manager),
            logger: c.MustMake("log").(log.Manager),
        }
    })
    
    // Step 2: Register any subproviders
    app.Register(&MySubProvider{})
    
    return nil
}
```

### 2. **Boot Phase**

```go
// Provider Implementation Example
func (p *MyProvider) Boot(app core.Application) error {
    // Step 1: Resolve service instance
    service := app.MustMake("my-service").(*MyService)
    
    // Step 2: Initialize the service
    if err := service.Initialize(); err != nil {
        return err
    }
    
    // Step 3: Register cleanup if needed
    runtime.SetFinalizer(service, func(s *MyService) {
        s.Cleanup()
    })
    
    return nil
}
```

### 3. **Dependency Declaration**

```go
// Provider Implementation Example
func (p *MyProvider) Requires() []string {
    // Declare dependencies on other providers
    return []string{
        "config",     // Require config provider
        "log",        // Require log provider
        "database",   // Require database provider
    }
}

func (p *MyProvider) Providers() []string {
    // Declare services this provider offers
    return []string{
        "my-service", // This provider offers my-service
        "my-utility", // This provider offers my-utility
    }
}
```

## 🧪 Edge Cases và Xử lý Lỗi

### 1. **Boot Errors**

```mermaid
sequenceDiagram
    participant App as Application
    participant ProviderA as Provider A
    participant ProviderB as Provider B
    participant ProviderC as Provider C
    
    App->>+ProviderA: Boot(app)
    ProviderA-->>-App: Success
    
    App->>+ProviderB: Boot(app)
    ProviderB-->>-App: Error: "Connection failed"
    
    Note over App,ProviderB: Boot process stops with error
    
    App->>App: Return error
    
    Note over App: Provider C never booted due to B's failure
```

### 2. **Missing Dependencies**

```mermaid
graph TD
    A[Provider A<br/>Requires: [X]]
    B[Provider B<br/>Requires: [A]]
    
    A --> X[X: Missing Provider]
    A --> B
    
    style A fill:#ff9999,stroke:#333,stroke-width:2px
    style B fill:#ffcc99,stroke:#333,stroke-width:2px
    style X fill:#ff0000,stroke:#333,stroke-width:2px,color:#ffffff
    
    ERROR[Error: Dependency "X" not found<br/>Required by "Provider A"] 
    style ERROR fill:#ff0000,stroke:#333,stroke-width:2px,color:#ffffff
```

### 3. **Runtime Provider Loading**

```mermaid
sequenceDiagram
    participant App as Application
    participant Loader as Module Loader
    participant NewProvider as New Provider
    participant Container as DI Container
    
    Note over App: Application already booted
    
    App->>+Loader: LoadModule(newProvider)
    Loader->>+NewProvider: typecasting check
    NewProvider-->>-Loader: Valid ServiceProvider
    
    Loader->>+App: Register(newProvider)
    App-->>-Loader: Success
    
    Loader->>+NewProvider: Register(app)
    NewProvider->>+Container: Bind services
    Container-->>-NewProvider: Success
    NewProvider-->>-Loader: Success
    
    Loader->>+NewProvider: Boot(app)
    NewProvider-->>-Loader: Success
    
    Loader-->>-App: Success
    
    Note over App: New provider registered and booted
```

## 🏆 Best Practices

### 1. **Modular Provider Structure**

```go
// Group related providers
type DatabaseProviders struct {
    Connection *DatabaseConnectionProvider
    Migration  *DatabaseMigrationProvider
    Seeder     *DatabaseSeederProvider
}

func (p *DatabaseProviders) Register(app core.Application) {
    app.Register(p.Connection)
    app.Register(p.Migration)
    app.Register(p.Seeder)
}

func (p *DatabaseProviders) Requires() []string {
    return []string{"config", "log"}
}

func (p *DatabaseProviders) Providers() []string {
    return []string{"database", "migration", "seeder"}
}
```

### 2. **Defensive Dependency Loading**

```go
func (p *MyProvider) Register(app core.Application) error {
    // Get required service safely
    configInterface, err := app.Make("config")
    if err != nil {
        return fmt.Errorf("config service not available: %w", err)
    }
    
    // Type assertion with validation
    config, ok := configInterface.(config.Manager)
    if !ok {
        return fmt.Errorf("invalid config service type: expected config.Manager, got %T", configInterface)
    }
    
    // Now use the config service safely
    dsn := config.GetString("database.dsn")
    
    // Register our service with the valid config
    app.Singleton("my-service", func(c di.Container) interface{} {
        return &MyService{config: config}
    })
    
    return nil
}
```

### 3. **Phased Registration**

```go
func (p *ComplexProvider) Register(app core.Application) error {
    // Phase 1: Register basic services
    if err := p.registerBasicServices(app); err != nil {
        return err
    }
    
    // Phase 2: Register dependent services
    if err := p.registerDependentServices(app); err != nil {
        return err
    }
    
    // Phase 3: Register advanced services
    if err := p.registerAdvancedServices(app); err != nil {
        return err
    }
    
    return nil
}
```

## 📈 Performance Optimization

### 1. **Cached Provider Sorting**

```go
func (a *application) RegisterWithDependencies() error {
    a.mu.Lock()
    defer a.mu.Unlock()
    
    // Use cached sorted providers if available
    if len(a.sortedProviders) > 0 {
        // If no providers added since last sort, use cached result
        if len(a.sortedProviders) == len(a.providers) {
            for _, provider := range a.sortedProviders {
                if err := provider.Register(a); err != nil {
                    return err
                }
            }
            return nil
        }
    }
    
    // Otherwise do full dependency sort
    // ...
}
```

### 2. **Lazy Provider Resolution**

```go
func (a *application) lazyRegisterProvider(name string) error {
    a.mu.RLock()
    defer a.mu.RUnlock()
    
    for _, provider := range a.providers {
        if slices.Contains(provider.Providers(), name) {
            if err := provider.Register(a); err != nil {
                return err
            }
            return nil
        }
    }
    
    return fmt.Errorf("no provider found for service: %s", name)
}
```

## 🔮 Roadmap

### Planned Features

- **Tự động phát hiện provider cycles** trong compile time
- **Conditional Providers** - Providers chỉ được load khi điều kiện được thỏa mãn
- **Named Provider Instances** - Support cho nhiều instances cùng provider type
- **Tagging và Grouping** - Providers và services có thể được tag và group

---

> **Next**: [Core Providers Documentation](core_providers.md) - Chi tiết về các core providers cốt lõi