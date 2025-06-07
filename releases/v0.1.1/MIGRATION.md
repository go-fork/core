# Migration Guide - v0.1.1

## Overview
Hướng dẫn này giúp bạn nâng cấp từ phiên bản trước lên v0.1.1. Đây là bản cập nhật nhỏ không có breaking changes.

## Prerequisites
- Go 1.23 or later
- go.fork.vn/core v0.1.0 installed

## Quick Migration Checklist
- [x] No import statement changes required
- [x] No function signature changes
- [x] No configuration format changes
- [x] All existing code remains compatible
- [ ] Update dependency to get latest fixes

## Breaking Changes
**Không có breaking changes** trong bản phát hành này.

## Dependencies Update
Chỉ cần cập nhật dependency trong go.mod:

```bash
go get go.fork.vn/core@v0.1.1
go mod tidy
```

## What Changed
### Internal Improvements
- Upgraded internal logging dependency for better stability
- Enhanced test coverage for better code quality
- Fixed potential logging issues

### No Code Changes Required
Tất cả existing code sẽ hoạt động mà không cần thay đổi gì.
```go
// Old way (previous version)
oldFunction(param1, param2)

// New way (v0.1.1)
newFunction(param1, param2, newParam)
```

#### Removed Functions
- `removedFunction()` - Use `newAlternativeFunction()` instead

#### Changed Types
```go
// Old type definition
type OldConfig struct {
    Field1 string
    Field2 int
}

// New type definition
type NewConfig struct {
    Field1 string
    Field2 int64 // Changed from int
    Field3 bool  // New field
}
```

### Configuration Changes
If you're using configuration files:

```yaml
# Old configuration format
old_setting: value
deprecated_option: true

# New configuration format
new_setting: value
# deprecated_option removed
new_option: false
```

## Step-by-Step Migration

### Step 1: Update Dependencies
```bash
go get go.fork.vn/core@v0.1.1
go mod tidy
```

### Step 2: Update Import Statements
```go
// If import paths changed
import (
    "go.fork.vn/core" // Updated import
)
```

### Step 3: Update Code
Replace deprecated function calls:

```go
// Before
result := core.OldFunction(param)

// After
result := core.NewFunction(param, defaultValue)
```

### Step 4: Update Configuration
Update your configuration files according to the new schema.

### Step 5: Run Tests
```bash
go test ./...
```

## Common Issues and Solutions

### Issue 1: Function Not Found
**Problem**: `undefined: core.OldFunction`  
**Solution**: Replace with `core.NewFunction`

### Issue 2: Type Mismatch
**Problem**: `cannot use int as int64`  
**Solution**: Cast the value or update variable type

## Getting Help
- Check the [documentation](https://pkg.go.dev/go.fork.vn/core@v0.1.1)
- Search [existing issues](https://github.com/go-fork/core/issues)
- Create a [new issue](https://github.com/go-fork/core/issues/new) if needed

## Rollback Instructions
If you need to rollback:

```bash
go get go.fork.vn/core@previous-version
go mod tidy
```

Replace `previous-version` with your previous version tag.

---
**Need Help?** Feel free to open an issue or discussion on GitHub.
