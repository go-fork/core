# Nhật ký Thay đổi

Tất cả các thay đổi đáng chú ý của dự án này sẽ được ghi lại trong file này.

Định dạng dựa trên [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
và dự án này tuân thủ [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Planned
- Future improvements and features

## [0.1.1] - 2025-06-07

### Changed
- **Dependency Upgrade**: Nâng cấp dependency go.fork.vn/log từ v0.1.4 lên v0.1.7
  - Khắc phục lỗi file handler initialization trong log v0.1.4 gây panic khi file path rỗng
  - Cải thiện độ ổn định và hiệu suất của hệ thống logging

### Fixed
- **Test Coverage Improvements**: Hoàn thiện test coverage lên 95.1%
  - Sửa lỗi loader_test.go file rỗng không có test implementation
  - Thêm comprehensive test cases cho tất cả ModuleLoader functions
  - Khắc phục mock function signatures và DI container integration tests
  - Sửa các test case bị conflict do thay đổi dependency

### Technical Details
- Hoàn thiện test suite với 95.1% statement coverage
- Tất cả tests đều PASS (0 failures)
- Benchmarks hoạt động ổn định
- Cải thiện error handling và edge case testing

## [0.1.0] - 2025-06-05

### Added
- **Tài liệu tiếng Việt hoàn chỉnh**: Xây dựng toàn bộ tài liệu tiếng Việt đầy đủ cho package core:
  - Tài liệu Tổng quan Hệ thống (`/docs/overview.md`): Kiến trúc, tính năng và hiệu suất
  - Tài liệu Application Interface (`/docs/application.md`): Chi tiết Application interface và implementation
  - Tài liệu Module Loader (`/docs/loader.md`): Hệ thống Module Loader
  - Tài liệu Workflows (`/docs/workflows.md`): Quy trình hoạt động và dependency management
  - Tài liệu Core Providers (`/docs/core_providers.md`): Chi tiết về các core providers
  - Trang chủ tài liệu (`/docs/index.md`): Tổng quan và hướng dẫn sử dụng
  - Sử dụng Mermaid diagram cho tất cả các workflow
  - Ví dụ mã nguồn thực tế và chi tiết kỹ thuật
  
- **README.md**: Cập nhật file README.md với định dạng Markdown chuẩn:
  - Tổng quan về package và tính năng chính
  - Quick start guide và ví dụ
  - Kiến trúc và diagram
  - Liên kết đến tài liệu chi tiết
  - Hướng dẫn cài đặt và sử dụng
  - Performance benchmark

### Changed
- **Cấu trúc tài liệu**: Tổ chức tài liệu với cấu trúc nhất quán và navigation giữa các trang
- **Định dạng tài liệu**: Chuyển đổi HTML sang Markdown chuẩn cho khả năng tương thích tốt hơn
- **Hỗ trợ đa ngôn ngữ**: Documentation hoàn chỉnh bằng tiếng Việt với thuật ngữ kỹ thuật chính xác

### Added

#### Framework Ứng dụng Cốt lõi
- **Application Interface**: Interface `Application` cốt lõi định nghĩa hợp đồng cho quản lý dependency injection container
- **Application Implementation**: Triển khai hoàn chỉnh với hỗ trợ:
  - Quản lý cấu hình thông qua `go.fork.vn/config`
  - Quản lý logging thông qua `go.fork.vn/log`
  - Dependency injection container thông qua `go.fork.vn/di`
  - Đăng ký và khởi động service provider

#### Quản lý Dependency Thông minh
- **Sắp xếp Dependency**: Tự động phân giải và sắp xếp dependency sử dụng thuật toán topological sort
- **Quy trình Boot Thông minh**: Workflow khởi động thông minh:
  - Tự động phát hiện dependency của service provider
  - Sử dụng phương thức đăng ký phù hợp (có hoặc không có dependency)
  - Đảm bảo mỗi provider chỉ boot một lần
  - Duy trì tương thích ngược cho các trường hợp đơn giản
- **RegisterWithDependencies()**: Đăng ký service provider nâng cao với sắp xếp dependency
- **hasDependencies()**: Phương thức helper để phát hiện provider có dependency

#### Hệ thống Tải Module
- **Module Loader**: Hệ thống tải module tổng quát với hỗ trợ:
  - Tải module dựa trên cấu hình
  - Đăng ký module tùy chỉnh
  - Xử lý lỗi và validation
- **Kiến trúc Linh hoạt**: Hỗ trợ cả kịch bản dependency đơn giản và phức tạp

#### Cơ sở hạ tầng Kiểm thử
- **Bộ Test Toàn diện**: 93.8% test coverage bao gồm:
  - Unit test cho tất cả public method
  - Integration test cho workflow service provider
  - Test validation sắp xếp dependency
  - Test xử lý lỗi và edge case
  - Benchmark test cho validation hiệu suất
- **Cơ sở hạ tầng Mock**: Thiết lập mock hoàn chỉnh sử dụng Mockery cho:
  - Service provider với theo dõi thời gian
  - Trackable provider cho test thứ tự
  - Mô phỏng lỗi và testing

#### CI/CD và Quản lý Phát hành
- **GitHub Actions**: Workflow tự động cho:
  - Kiểm tra chất lượng code
  - Thực thi test
  - Báo cáo coverage
  - Tự động hóa phát hành
- **Script Phát hành**: Công cụ quản lý phát hành tự động:
  - Tạo archive (`archive_release.sh`)
  - Tạo template phát hành (`create_release_templates.sh`)
- **Chất lượng Code**: Tích hợp golangci-lint với quy tắc linting toàn diện

#### Tài liệu và Ví dụ
- **Tài liệu API**: Tài liệu Go toàn diện cho tất cả public API
- **Ví dụ Sử dụng**: Cấu hình và pattern sử dụng mẫu trong `/docs/`
- **Tài liệu Kiến trúc**: Kiến trúc chi tiết và quyết định thiết kế

### Technical Details

#### Dependency
- **go.fork.vn/config v0.1.3**: Quản lý cấu hình
- **go.fork.vn/di v0.1.3**: Dependency injection container
- **go.fork.vn/log v0.1.3**: Framework logging
- **github.com/stretchr/testify v1.10.0**: Tiện ích testing

#### Key Features
- **Tạo key service provider zero-allocation** sử dụng địa chỉ bộ nhớ
- **Thuật toán topological sort** cho phân giải dependency
- **Thao tác thread-safe** cho môi trường concurrent
- **Cấu hình linh hoạt** thông qua file config bên ngoài
- **Xử lý lỗi toàn diện** với thông báo lỗi chi tiết

### Breaking Changes
- Đây là phát hành đầu tiên, không có thay đổi phá vỡ

### Migration Guide
- Đây là phát hành đầu tiên, không cần migration

### Performance
- Phân giải dependency hiệu quả với độ phức tạp O(V + E)
- Overhead bộ nhớ tối thiểu cho theo dõi service provider
- Tối ưu hóa cho cả kịch bản dependency đơn giản và phức tạp

---

## Development Information

### Project Structure
```
go.fork.vn/core/
├── application.go          # Triển khai application cốt lõi
├── application_test.go     # Bộ test toàn diện
├── loader.go              # Hệ thống tải module
├── loader_test.go         # Test module loader
├── doc.go                 # Tài liệu package
├── configs/               # Ví dụ cấu hình
├── docs/                  # Tài liệu
├── mocks/                 # Mock được tạo
├── scripts/               # Script quản lý phát hành
└── .github/               # Workflow CI/CD
```

### Contributing
Vui lòng đọc hướng dẫn đóng góp trong thư mục docs/ trước khi gửi pull request.

### License
Dự án này được cấp phép theo các điều khoản được chỉ định trong file LICENSE.