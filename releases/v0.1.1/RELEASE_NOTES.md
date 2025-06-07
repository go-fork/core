# Release Notes - v0.1.1

## Overview
Bản phát hành này tập trung vào việc nâng cấp dependency quan trọng và cải thiện chất lượng test coverage.

## What's New
### 🔧 Improvements
- **Dependency Upgrade**: Nâng cấp go.fork.vn/log từ v0.1.4 lên v0.1.7
  - Khắc phục lỗi file handler initialization gây panic khi file path rỗng
  - Cải thiện độ ổn định và hiệu suất của hệ thống logging

### 🧪 Test Coverage Improvements
- **Complete Test Suite**: Hoàn thiện test coverage lên 95.1%
  - Sửa lỗi loader_test.go file rỗng không có test implementation
  - Thêm comprehensive test cases cho tất cả ModuleLoader functions
  - Khắc phục mock function signatures và DI container integration tests
  - Sửa các test case bị conflict do thay đổi dependency

### 🐛 Bug Fixes
- Fixed empty loader_test.go file với proper test implementations
- Fixed mock function signatures (args interface{} → mock.Arguments)
- Fixed DI container integration tests
- Fixed panic-expecting tests to use proper mock returns

## Breaking Changes
### ⚠️ Important Notes
Không có breaking changes trong bản phát hành này. Tất cả API giữ nguyên tương thích.

## Migration Guide
See [MIGRATION.md](./MIGRATION.md) for detailed migration instructions.

## Dependencies
### Updated
- go.fork.vn/log: v0.1.4 → v0.1.7
  - Khắc phục lỗi file handler initialization gây panic
  - Cải thiện độ ổn định và hiệu suất logging system

### Maintained
- go.fork.vn/config: v0.1.3 (stable)
- go.fork.vn/di: v0.1.3 (stable)
- github.com/stretchr/testify: v1.10.0 (testing framework)

## Performance
- Test execution improved với log v0.1.7 stability fixes
- Benchmark performance maintained:
  - RegisterCoreProviders: ~912.2 ns/op
  - LoadModule: ~140.0 ns/op
- Memory efficiency unchanged với dependency upgrade

## Security
- Dependency security improvements với log v0.1.7 upgrade
- Không có security vulnerabilities được phát hiện trong bản phát hành này

## Testing
- Hoàn thiện test coverage lên 95.1% statement coverage
- Thêm comprehensive test cases cho ModuleLoader functions
- Sửa lỗi loader_test.go file rỗng với proper implementations
- Tất cả tests PASS (0 failures) với log v0.1.7
- Benchmarks hoạt động ổn định

## Contributors
Thanks to all contributors who made this release possible:
- @nghiant0921
- @zinzinday

## Download
- Source code: [go.fork.vn/core@v0.1.1]
- Documentation: [pkg.go.dev/go.fork.vn/core@v0.1.1]

---
Release Date: 2025-06-07
