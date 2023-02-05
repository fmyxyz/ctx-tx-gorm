package gormx

import (
	"context"
	"database/sql"

	tx "github.com/fmyxyz/ctx-tx"

	"gorm.io/gorm"
)

type GormDB struct {
	*gorm.DB

	instance string
}

func (g *GormDB) SavePoint(name string) error {
	return g.DB.SavePoint(name).Error
}

func (g *GormDB) RollbackTo(name string) error {
	return g.DB.RollbackTo(name).Error
}

func (g *GormDB) Commit() error {
	return g.DB.Commit().Error
}

func (g *GormDB) Rollback() error {
	return g.DB.Rollback().Error
}

func (g *GormDB) Name() string {
	return "gorm-" + g.instance
}

func (g *GormDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (tx.Tx, error) {
	db := g.DB.WithContext(ctx)
	tx0 := db.Begin(opts)
	return warp(tx0), tx0.Error
}

func warp(db *gorm.DB) *GormDB {
	return &GormDB{DB: db}
}

const defaultInstance = "default"

func Register(db *gorm.DB, opts ...GormDBOption) {
	gormDB := &GormDB{DB: db, instance: defaultInstance}
	for _, opt := range opts {
		opt(gormDB)
	}
	tx.Register(gormDB, tx.RegisterDefaultDB(gormDB.instance == defaultInstance))
}

type GormDBOption func(db *GormDB)

func Instance(instance string) GormDBOption {
	return func(db *GormDB) {
		db.instance = instance
	}
}

func FromContext(ctx context.Context, opts ...GormDBOption) *GormDB {
	gormDB := &GormDB{instance: defaultInstance}
	for _, opt := range opts {
		opt(gormDB)
	}
	name := gormDB.Name()
	txManager := tx.GetTxManager(name)
	if txManager == nil {
		panic(name + " not register in txManagers")
	}
	tx0 := txManager.TxFromContext(ctx)
	if tx0 != nil {
		return tx0.(*GormDB)
	}
	db := txManager.DBFromContext(ctx)
	if db != nil {
		return db.(*GormDB)
	}
	return nil
}
