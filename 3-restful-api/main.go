package main

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/*
| Interface                              | Ai gọi              | Khi nào                |
| -------------------------------------- | ------------------- | ---------------------- |
| `sql.Scanner` → `Scan()`               | GORM                | **SELECT từ DB**       |
| `driver.Valuer` → `Value()`            | GORM                | **INSERT / UPDATE DB** |
| `json.Marshaler` → `MarshalJSON()`     | Gin / encoding/json | **Trả response**       |
| `json.Unmarshaler` → `UnmarshalJSON()` | Gin                 | **Nhận request JSON**  |

*/

/* Luồng chạy
REQUEST JSON
   ↓
UnmarshalJSON()   ← Gin gọi
   ↓
ItemStatus (0,1,2)
   ↓
Value()           ← GORM gọi
   ↓
DB ("Doing")

----------------------------

DB ("Doing")
   ↓
Scan()            ← GORM gọi
   ↓
ItemStatus (0,1,2)
   ↓
MarshalJSON()     ← Gin gọi
   ↓
RESPONSE JSON
*/

type ItemStatus int

/*
iota là một hằng số đặc biệt trong Go, dùng để tự động tăng số khi khai báo const. Nó chỉ dùng được trong const block.
Trong const block, nếu bỏ trống giá trị, Go sẽ tự dùng lại biểu thức ở dòng trên, nhưng iota thì tự tăng.
*/
const (
	ItemStatusDoing   ItemStatus = iota // 0
	ItemStatusDone                      // iota -> 1
	ItemStatusDeleted                   // iota -> 2
)

var allItemStatuses = [3]string{"Doing", "Done", "Deleted"}

// String() – dùng khi in / trả JSON (map từ số 0,1,2 sang "Doing", "Done", "Deleted")
func (item *ItemStatus) String() string {
	return allItemStatuses[*item]
}

func parseStrToItemStatus(s string) (ItemStatus, error) {
	for index := range allItemStatuses {
		if allItemStatuses[index] == s {
			return ItemStatus(index), nil
		}
	}

	return ItemStatus(0), errors.New("Invalid status string")
}

// convert string → enum (GORM tự động gọi khi: SELECT từ DB)
func (item *ItemStatus) Scan(value interface{}) error {
	bytes, ok := value.([]byte)

	if !ok {
		return fmt.Errorf("Fail to scan data from sql: %s", value)
	}

	v, err := parseStrToItemStatus(string(bytes))

	if err != nil {
		return errors.New(err.Error())
	}

	*item = v // *item = 0,1,2

	return nil
}

func (item *ItemStatus) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	return item.String(), nil
}

// Nếu không có method này thì khi gọi api get item thì status sẽ trả về số 0,1,2 thay vì "Doing", "Done", "Deleted" (GORM gọi hàm này)
func (item *ItemStatus) MarshalJSON() ([]byte, error) {
	if item == nil {
		return nil, nil
	}

	return []byte(fmt.Sprintf("\"%s\"", item.String())), nil
}

func (item *ItemStatus) UnmarshalJSON(data []byte) error {
	str := strings.ReplaceAll(string(data), "\"", "")

	itemValue, err := parseStrToItemStatus(str)

	if err != nil {
		return err
	}

	*item = itemValue

	return nil
}

// Status dùng *ItemStatus thay vì ItemStatus vì nếu vì lí do nào đó mà Status là null thì không bị lỗi
type TodoItem struct {
	Id          int         `json:"id" gorm:"column:id"`
	Title       string      `json:"title" gorm:"column:title"`
	Description string      `json:"description" gorm:"column:description"`
	Status      *ItemStatus `json:"status" gorm:"column:status"`
	CreatedAt   *time.Time  `json:"created_at" gorm:"created_at"`
	UpdatedAt   *time.Time  `json:"updated_at,omitempty" gorm:"updated_at"` // omitempty: BỎ QUA field khi marshal nếu field đó rỗng / giá trị mặc định
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

// Phân trang (dùng tag form để parse query params từ client truyền lên)
type Paging struct {
	Page  int   `json:"page" form:"page"`
	Limit int   `json:"limit" form:"limit"`
	Total int64 `json:"total" form:"-"` // Không nhận trường total từ client truyền lên
}

// Xử lí query params client truyền lên
func (p *Paging) Process() {
	if p.Page <= 0 {
		p.Page = 1
	}

	if p.Limit <= 0 || p.Limit >= 100 {
		p.Limit = 10
	}
}

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to database
	dsn := os.Getenv("DB_CONNECTION")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	/*
		# CRUD: Create, Read, Update, Delete

		POST /v1/items (create a new item)
		GET /v1/items (list item) /v1/items?page=1
		GET /v1/items/:id (get item detail by id)
		(PUT || PATCH) /v1/items/:id (update an item by id)
		DELETE /v1/items/:id (delete item by id)
	*/
	router := gin.Default()

	v1 := router.Group("/v1")
	{
		items := v1.Group("/items")
		{
			items.POST("", CreatItem(db))
			items.GET("", ListItem(db))
			items.GET("/:id", GetItem(db))
			items.PATCH("/:id", UpdateItem(db))
			items.DELETE("/:id", DeleteItem(db))
		}
	}

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ping sucessfully",
		})
	})
	router.Run() // listens on 0.0.0.0:8080 by default // or: router.Run(":8000") // change port 8000

}

func CreatItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var data TodoItemCreation

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		if err := db.Create(&data).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": data.Id,
		})
	}
}

func GetItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		//
		var data TodoItem

		/*
			// Cách 1:
			data.Id = id
			if err := db.First(&data).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})

				return
			}

		*/

		// Cách 2:
		if err := db.Where("id = ?", id).First(&data).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})

	}
}

// Truyền trường nào lên thì mới cập nhật trường đó, còn lại giữ nguyên
func UpdateItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		//
		var data TodoItemUpdate

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		if err := db.Where("id = ?", id).Updates(&data).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": true,
		})

	}
}

func DeleteItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		/*
			// Code như này là xóa luôn row đó luôn
			if err := db.Table(TodoItem{}.TableName()).Where("id = ?", id).Delete(nil).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})

				return
			}
		*/

		// Soft delete -> chỉ thay đổi status sang Deleted
		if err := db.Table(TodoItem{}.TableName()).Where("id = ?", id).Updates(map[string]interface{}{"status": "Deleted"}).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": true,
		})

	}
}

func ListItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var paging Paging

		if err := c.ShouldBind(&paging); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		paging.Process()

		//
		var result []TodoItem

		// Chỉ lấy các trường có status khác Deleted
		db = db.Where("status <> ?", "Deleted")

		if err := db.Table(TodoItem{}.TableName()).Count(&paging.Total).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		if err := db.Order("id desc").
			Offset((paging.Page - 1) * paging.Limit).
			Limit(paging.Limit).
			Find(&result).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data":   result,
			"paging": paging,
		})

	}
}
