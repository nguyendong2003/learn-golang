package model

import "github.com/google/uuid"

type User struct {
	BaseModel

	Email    string `gorm:"type:varchar(255);not null;uniqueIndex"`
	Username string `gorm:"type:varchar(100);not null;uniqueIndex"`
	Password string `gorm:"type:text;not null"`
	FullName string `gorm:"type:varchar(255)"`
	Avatar   string `gorm:"type:text"`
	IsActive bool   `gorm:"default:true"`
}

type GetDetailUserParams struct {
	Id *string
}

func (p GetDetailUserParams) Map() (User, error) {
	var user User
	if p.Id != nil {
		id, err := uuid.Parse(*p.Id)
		if err != nil {
			return user, err
		}

		user.ID = id
	}

	return user, nil
}
