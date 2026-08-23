package dbutils

import (
	"context"

	"gorm.io/gorm"
)

type Manager interface {
	WithinTransaction(
		ctx context.Context,
		fn func(ctx context.Context) error,
	) error
}

type manager struct {
	db *gorm.DB
}

func NewManager(db *gorm.DB) Manager {
	return &manager{
		db: db,
	}
}

type contextKey struct{}

func (m *manager) WithinTransaction(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	return m.db.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {

			txCtx := context.WithValue(
				ctx,
				contextKey{},
				tx,
			)

			return fn(txCtx)
		},
	)
}

func DB(ctx context.Context, fallback *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(contextKey{}).(*gorm.DB); ok {
		return tx.WithContext(ctx)
	}

	return fallback.WithContext(ctx)
}
