package database

import (
	"testing"

	"fullstack-app-template/server/internal/config"
	"fullstack-app-template/server/internal/model"
	"gorm.io/gorm"
)

func openTest(t *testing.T) *gorm.DB {
	t.Helper()
	adminURL, err := config.DatabaseURL()
	if err != nil {
		t.Fatalf("读取数据库配置失败: %v", err)
	}
	db, cleanup, err := OpenIsolated(adminURL)
	if err != nil {
		t.Fatalf("打开测试库失败: %v", err)
	}
	t.Cleanup(func() {
		if err := cleanup(); err != nil {
			t.Errorf("清理测试库失败: %v", err)
		}
	})
	return db
}

func TestOpenRejectsInvalidURL(t *testing.T) {
	if _, err := Open(""); err == nil {
		t.Fatal("empty URL returned a database")
	}
	if _, err := Open("data/app.db"); err == nil {
		t.Fatal("non-postgres URL returned a database")
	}
}

func TestMigrateAndReadBack(t *testing.T) {
	db := openTest(t)
	if err := Migrate(db); err != nil {
		t.Fatalf("migrate again: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	if sqlDB.Stats().MaxOpenConnections != 10 {
		t.Fatalf("max open connections = %d, want 10", sqlDB.Stats().MaxOpenConnections)
	}

	var row model.Counter
	if err := db.Where("name = ?", model.WoodfishName).First(&row).Error; err != nil {
		t.Fatal(err)
	}
	if row.Count != 0 {
		t.Fatalf("count = %d, want 0", row.Count)
	}
	if err := db.Model(&row).Update("count", int64(3)).Error; err != nil {
		t.Fatal(err)
	}
	var got model.Counter
	if err := db.First(&got, row.ID).Error; err != nil {
		t.Fatal(err)
	}
	if got.Count != 3 {
		t.Fatalf("count = %d, want 3", got.Count)
	}
}
