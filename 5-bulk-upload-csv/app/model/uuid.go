package model

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ModelUUID struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`
}

func (u *ModelUUID) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID, err = uuid.NewV7()

	if err != nil {
		return errors.New("Failed to generate uuid")
	}

	return
}
