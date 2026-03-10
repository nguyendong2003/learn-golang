package model

type Blog struct {
	BaseModel
	Title     string `gorm:"type:varchar(255);not null"`
	Content   string `gorm:"type:text;not null"`
	Slug      string `gorm:"type:varchar(255);not null;uniqueIndex"`
	AuthorID  string `gorm:"type:uuid;not null;index"`
	ViewTotal int64  `gorm:"default:0"`

	Author *User `gorm:"foreignKey:AuthorID;references:ID"`
}

func (Blog) TableName() string {
	return "blogs"
}
