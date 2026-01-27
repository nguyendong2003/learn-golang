package model

import (
	"errors"
	"restfulapi/common"
)

var (
	ErrTitleIsBlank = errors.New("Title cannot be blank")
	ErrItemDeleted  = errors.New("Item is deleted. Cannot update it")
)

// Status dùng *ItemStatus thay vì ItemStatus vì nếu vì lí do nào đó mà Status là null thì không bị lỗi
type TodoItem struct {
	common.SQLModel
	Title       string      `json:"title" gorm:"column:title"`
	Description string      `json:"description" gorm:"column:description"`
	Status      *ItemStatus `json:"status" gorm:"column:status"`

	/*	Thay 3 field này bằng common.SQLModel bằng kỹ thuật struct embedding
		Id          int         `json:"id" gorm:"column:id"`
		CreatedAt   *time.Time  `json:"created_at" gorm:"created_at"`
		UpdatedAt   *time.Time  `json:"updated_at,omitempty" gorm:"updated_at"` // omitempty: BỎ QUA field khi marshal nếu field đó rỗng / giá trị mặc định
	*/
}

func (TodoItem) TableName() string { return "todo_items" }

type TodoItemCreation struct {
	Id          int         `json:"-" gorm:"column:id"`
	Title       string      `json:"title" gorm:"column:title"`
	Description string      `json:"description" gorm:"column:description"`
	Status      *ItemStatus `json:"status" gorm:"column:status"`
}

func (TodoItemCreation) TableName() string { return TodoItem{}.TableName() }

/*
# Lưu ý: Hoạt động của gorm:
- Truyền lên trường nào thì mới cật nhật trường đó, còn các trường không truyền thì giữ nguyên

- Nếu các trường title, description chỉ dùng string thì nếu client truyền lên title="" hoặc description="" hoặc title=nil hoặc description=nil hoặc không truyền title, description
thì gorm mặc định sẽ xem nó như không có giá trị và bỏ qua việc cập nhật database các trường đó

- Nên nếu muốn khi client truyền lên title="" hoặc description="" thì cập nhật các trường đó về chuỗi rỗng ""
thì các trường đó phải đổi từ string sang *string. Khi chuyển qua *string thì nếu trường đó truyền lên giá trị khác nil thì nó sẽ cập nhật
*/
type TodoItemUpdate struct {
	Title       *string `json:"title" gorm:"column:title"`
	Description *string `json:"description" gorm:"column:description"`
	Status      *string `json:"status" gorm:"column:status"`
}

func (TodoItemUpdate) TableName() string { return TodoItem{}.TableName() }
