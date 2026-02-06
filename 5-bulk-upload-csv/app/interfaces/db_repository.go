package interfaces

import (
	"database/sql"

	"gorm.io/gorm"
)

type DbRepositoryInterface interface {
	InitializeDB() error
	GetDB() *gorm.DB
	GetSqlDB() *sql.DB
}
