// Package database 打开 GORM 连接，并在启动时集中迁移。
// 数据库固定为 PostgreSQL，驱动是 gorm.io/driver/postgres。
package database

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/url"
	"strings"
	"time"

	"fullstack-app-template/server/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Open(databaseURL string) (*gorm.DB, error) {
	databaseURL = strings.TrimSpace(databaseURL)
	if databaseURL == "" {
		return nil, errors.New("未设置数据库连接串")
	}
	lower := strings.ToLower(databaseURL)
	if !strings.HasPrefix(lower, "postgres://") && !strings.HasPrefix(lower, "postgresql://") {
		return nil, errors.New("数据库连接串必须是 PostgreSQL")
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  databaseURL,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}
	return db, nil
}

func Migrate(db *gorm.DB) error {
	if db == nil {
		return errors.New("数据库未连接")
	}
	if err := db.AutoMigrate(&model.Counter{}); err != nil {
		return err
	}
	var row model.Counter
	return db.Where(model.Counter{Name: model.WoodfishName}).Attrs(model.Counter{Count: 0}).FirstOrCreate(&row).Error
}

func Close(db *gorm.DB) error {
	if db == nil {
		return nil
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// OpenIsolated 从已有的 PostgreSQL 创建一个独立数据库，并完成迁移。
// cleanup 关闭连接并删除该数据库。
func OpenIsolated(adminURL string) (*gorm.DB, func() error, error) {
	adminURL = strings.TrimSpace(adminURL)
	if adminURL == "" {
		return nil, nil, errors.New("未设置数据库连接串")
	}
	admin, err := Open(adminURL)
	if err != nil {
		return nil, nil, err
	}
	name, err := isolatedDatabaseName()
	if err != nil {
		_ = Close(admin)
		return nil, nil, err
	}
	if err := admin.Exec("CREATE DATABASE " + name).Error; err != nil {
		_ = Close(admin)
		return nil, nil, err
	}
	testURL, err := withDatabaseName(adminURL, name)
	if err != nil {
		_ = admin.Exec("DROP DATABASE IF EXISTS " + name + " WITH (FORCE)")
		_ = Close(admin)
		return nil, nil, err
	}
	db, err := Open(testURL)
	if err != nil {
		_ = admin.Exec("DROP DATABASE IF EXISTS " + name + " WITH (FORCE)")
		_ = Close(admin)
		return nil, nil, err
	}
	if err := Migrate(db); err != nil {
		_ = Close(db)
		_ = admin.Exec("DROP DATABASE IF EXISTS " + name + " WITH (FORCE)")
		_ = Close(admin)
		return nil, nil, err
	}
	cleanup := func() error {
		return errors.Join(
			Close(db),
			admin.Exec("DROP DATABASE IF EXISTS "+name+" WITH (FORCE)").Error,
			Close(admin),
		)
	}
	return db, cleanup, nil
}

func isolatedDatabaseName() (string, error) {
	var buf [8]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return "", err
	}
	return "t" + hex.EncodeToString(buf[:]), nil
}

func withDatabaseName(raw, name string) (string, error) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	parsed.Path = "/" + name
	return parsed.String(), nil
}
