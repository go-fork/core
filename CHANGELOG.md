# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2024-12-19

### Added

#### Core Application Framework
- **Application Interface**: Core `Application` interface defining the contract for dependency injection container management
- **Application Implementation**: Complete implementation with support for:
  - Configuration management via `go.fork.vn/config`
  - Logging management via `go.fork.vn/log`
  - Dependency injection container via `go.fork.vn/di`
  - Service provider registration and booting

#### Smart Dependency Management
- **Dependency Ordering**: Automatic dependency resolution and ordering using topological sort
- **Smart Boot Process**: Intelligent boot workflow that:
  - Auto-detects service provider dependencies
  - Uses appropriate registration method (with or without dependencies)
  - Ensures each provider boots only once
  - Maintains backward compatibility for simple cases
- **RegisterWithDependencies()**: Advanced service provider registration with dependency ordering
- **hasDependencies()**: Helper method to detect providers with dependencies

#### Module Loading System
- **Module Loader**: Generic module loading system with support for:
  - Configuration-based module loading
  - Custom module registration
  - Error handling and validation
- **Flexible Architecture**: Support for both simple and complex dependency scenarios

#### Testing Infrastructure
- **Comprehensive Test Suite**: 93.8% test coverage including:
  - Unit tests for all public methods
  - Integration tests for service provider workflows
  - Dependency ordering validation tests
  - Error handling and edge case tests
  - Benchmark tests for performance validation
- **Mock Infrastructure**: Complete mock setup using Mockery for:
  - Service providers with time tracking
  - Trackable providers for order testing
  - Error simulation and testing

#### CI/CD and Release Management
- **GitHub Actions**: Automated workflows for:
  - Code quality checks
  - Test execution
  - Coverage reporting
  - Release automation
- **Release Scripts**: Automated release management tools:
  - Archive creation (`archive_release.sh`)
  - Release template generation (`create_release_templates.sh`)
- **Code Quality**: golangci-lint integration with comprehensive linting rules

#### Documentation and Examples
- **API Documentation**: Comprehensive Go documentation for all public APIs
- **Usage Examples**: Example configurations and usage patterns in `/docs/`
- **Architecture Documentation**: Detailed architecture and design decisions

### Technical Details

#### Dependencies
- **go.fork.vn/config v0.1.3**: Configuration management
- **go.fork.vn/di v0.1.3**: Dependency injection container
- **go.fork.vn/log v0.1.3**: Logging framework
- **github.com/stretchr/testify v1.10.0**: Testing utilities

#### Key Features
- **Zero-allocation service provider key generation** using memory addresses
- **Topological sort algorithm** for dependency resolution
- **Thread-safe operations** for concurrent environments
- **Flexible configuration** via external config files
- **Comprehensive error handling** with detailed error messages

### Breaking Changes
- This is the initial release, no breaking changes

### Migration Guide
- This is the initial release, no migration needed

### Performance
- Efficient dependency resolution with O(V + E) complexity
- Minimal memory overhead for service provider tracking
- Optimized for both simple and complex dependency scenarios

---

## Development Information

### Project Structure
```
go.fork.vn/core/
├── application.go          # Core application implementation
├── application_test.go     # Comprehensive test suite
├── loader.go              # Module loading system
├── loader_test.go         # Module loader tests
├── doc.go                 # Package documentation
├── configs/               # Configuration examples
├── docs/                  # Documentation
├── mocks/                 # Generated mocks
├── scripts/               # Release management scripts
└── .github/               # CI/CD workflows
```

### Contributing
Please read the contributing guidelines in the docs/ directory before submitting pull requests.

### License
This project is licensed under the terms specified in the LICENSE file.