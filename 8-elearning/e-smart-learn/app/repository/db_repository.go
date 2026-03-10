package repository

import (
	"database/sql"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	MaxConnLifetime = 300
	OpenConns       = 0
	IdleConns       = 20
)

type DbRepository interface {
	InitializeDB() error
	GetDB() *gorm.DB
	GetSqlDB() *sql.DB
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

	sqlDB.SetConnMaxLifetime(MaxConnLifetime)
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

func (r *dbRepository) Close() {
	if r.sqlDB != nil {
		r.sqlDB.Close()
	}
}
