package repository

import (
	"context"
	"database/sql"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	MaxConnLifetime = 300
	OpenConns       = 30
	IdleConns       = 15
)

type DbRepository interface {
	InitializeDB() error
	GetDB() *gorm.DB
	GetSqlDB() *sql.DB
	Transaction(ctx context.Context, fn func(txDb DbRepository) error) error
	Begin(ctx context.Context) (DbRepository, error)
	Commit() error
	Rollback() error
	Close()
}

type dbRepository struct {
	dsn   string
	db    *gorm.DB
	sqlDB *sql.DB
}

func NewDbRepository(dsn string) DbRepository {
	return &dbRepository{
		dsn: dsn,
	}
}

func (r *dbRepository) InitializeDB() error {
	db, err := gorm.Open(postgres.Open(r.dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	var sqlDB *sql.DB
	if sqlDB, err = db.DB(); err != nil {
		return err
	}

	sqlDB.SetConnMaxLifetime(MaxConnLifetime * time.Second)
	sqlDB.SetMaxOpenConns(OpenConns)
	sqlDB.SetMaxIdleConns(IdleConns)
	err = Migrator(r.dsn, MigrateActionUp)
	if err != nil {
		return err
	}
	r.db = db
	r.sqlDB = sqlDB

	return nil
}

func (r *dbRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *dbRepository) GetSqlDB() *sql.DB {
	return r.sqlDB
}

func (r *dbRepository) Transaction(ctx context.Context, fn func(txDb DbRepository) error) error {
	if r.db == nil {
		return gorm.ErrInvalidDB
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &dbRepository{
			dsn:   r.dsn,
			db:    tx,
			sqlDB: r.sqlDB,
		}
		return fn(txRepo)
	})
}

func (r *dbRepository) Close() {
	if r.sqlDB != nil {
		r.sqlDB.Close()
	}
}

func (r *dbRepository) Begin(ctx context.Context) (DbRepository, error) {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &dbRepository{
		dsn:   r.dsn,
		db:    tx,
		sqlDB: r.sqlDB,
	}, nil
}

func (r *dbRepository) Commit() error {
	if r.db == nil {
		return gorm.ErrInvalidDB
	}
	return r.db.Commit().Error
}

func (r *dbRepository) Rollback() error {
	if r.db == nil {
		return gorm.ErrInvalidDB
	}
	return r.db.Rollback().Error
}
