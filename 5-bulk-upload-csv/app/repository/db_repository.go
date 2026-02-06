package repository

import (
	"bulk-upload-csv/interfaces"
	"database/sql"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	MaxConnLifetime = 300
	OpenConns       = 0
	IdleConns       = 10
)

type DbRepository struct {
	dsn   string
	db    *gorm.DB
	sqlDB *sql.DB
}

func NewDbRepository(dsn string) interfaces.DbRepositoryInterface {
	return &DbRepository{
		dsn: dsn,
	}
}

func (r *DbRepository) InitializeDB() error {
	db, err := gorm.Open(postgres.Open(r.dsn))
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

	r.db = db
	r.sqlDB = sqlDB

	return nil
}

func (r *DbRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *DbRepository) GetSqlDB() *sql.DB {
	return r.sqlDB
}

func (r *DbRepository) Close() {
	if r.sqlDB != nil {
		r.sqlDB.Close()
	}
}
