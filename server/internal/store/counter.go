// Package store 存放 GORM 查询。Handler 不直接查询数据库。
package store

import (
	"context"
	"errors"

	"fullstack-app-template/server/internal/model"
	"gorm.io/gorm"
)

type CounterStore struct {
	db *gorm.DB
}

func NewCounterStore(db *gorm.DB) CounterStore {
	return CounterStore{db: db}
}

func (s CounterStore) Count(ctx context.Context, name string) (int64, error) {
	var row model.Counter
	err := s.db.WithContext(ctx).Where("name = ?", name).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return row.Count, nil
}

func (s CounterStore) Increment(ctx context.Context, name string) (int64, error) {
	result := s.db.WithContext(ctx).Model(&model.Counter{}).Where("name = ?", name).Update("count", gorm.Expr("count + ?", 1))
	if result.Error != nil {
		return 0, result.Error
	}
	if result.RowsAffected == 0 {
		return 0, gorm.ErrRecordNotFound
	}
	return s.Count(ctx, name)
}
