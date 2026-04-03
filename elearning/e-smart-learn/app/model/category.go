package model

type Category struct {
	BaseModel
	Name        string `gorm:"type:varchar(255);not null;uniqueIndex"`
	Description string `gorm:"type:text"`

	Courses []*Course `gorm:"foreignKey:CategoryID;references:ID"`
}

func (Category) TableName() string {
	return "categories"
}
