package core_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.fork.vn/core"
	diMocks "go.fork.vn/di/mocks"
)

// Test trực tiếp với BootstrapApplication
func TestBootstrapApplication_RealCases(t *testing.T) {
	// Setup helper để không phải duplicate code
	prepareTestEnv := func() {
		// Tương tự như setupTestEnvironment, nhưng thực hiện trực tiếp
	}

	t.Run("success_path_coverage", func(t *testing.T) {
		t.Parallel()

		// Chuẩn bị môi trường test
		prepareTestEnv()

		// Tạo app với config hợp lệ
		app := core.New(map[string]interface{}{
			"file": "testdata/configs/console-only-simple.yaml",
		})

		// Đăng ký provider cần thiết
		provider := diMocks.NewMockServiceProvider(t)
		provider.EXPECT().Providers().Return([]string{"test.service"}).Maybe()
		provider.EXPECT().Requires().Return([]string{}).Maybe()
		provider.EXPECT().Register(app).Maybe()
		provider.EXPECT().Boot(app).Maybe()

		app.Register(provider)

		// Thực thi bootstrap và kiểm tra kết quả
		err := app.ModuleLoader().BootstrapApplication()
		assert.NoError(t, err)
	})

	t.Run("error_in_register_core_providers", func(t *testing.T) {
		t.Parallel()

		// Tạo app với config không tồn tại
		app := core.New(map[string]interface{}{
			"file": "non-existent-file.yaml", // FileNotFound error
		})

		// Không cần đăng ký provider thêm

		// Thực thi bootstrap - sẽ fail ở RegisterCoreProviders
		err := app.ModuleLoader().BootstrapApplication()
		assert.Error(t, err)
		// Kiểm tra lỗi liên quan đến config
		assert.Contains(t, err.Error(), "config")
	})

	t.Run("error_in_register_with_dependencies", func(t *testing.T) {
		t.Parallel()

		// Chuẩn bị môi trường test
		prepareTestEnv()

		// Tạo app với config hợp lệ
		app := core.New(map[string]interface{}{})

		// Đăng ký provider với dependency không tồn tại
		provider := diMocks.NewMockServiceProvider(t)
		provider.EXPECT().Providers().Return([]string{"test.service"}).Maybe()
		provider.EXPECT().Requires().Return([]string{"non-existent-service"}).Maybe()

		app.Register(provider)

		// Thực thi bootstrap - sẽ fail ở RegisterWithDependencies
		err := app.ModuleLoader().BootstrapApplication()
		assert.Error(t, err)
	})
}
