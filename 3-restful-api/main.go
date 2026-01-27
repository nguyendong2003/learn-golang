package main

import (
	"log"
	"net/http"
	"os"
	"restfulapi/common"
	"restfulapi/module/item/model"
	ginitem "restfulapi/module/item/transport/ginitem"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

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
			items.POST("", ginitem.CreatItem(db))
			items.GET("", ListItem(db))
			items.GET("/:id", ginitem.GetItem(db))
			items.PATCH("/:id", ginitem.UpdateItem(db))
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
		if err := db.Table(model.TodoItem{}.TableName()).Where("id = ?", id).Updates(map[string]interface{}{"status": "Deleted"}).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))

	}
}

func ListItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var paging common.Paging

		if err := c.ShouldBind(&paging); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		paging.Process()

		//
		var result []model.TodoItem

		// Chỉ lấy các trường có status khác Deleted
		db = db.Where("status <> ?", "Deleted")

		if err := db.Table(model.TodoItem{}.TableName()).Count(&paging.Total).Error; err != nil {
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

		c.JSON(http.StatusOK, common.NewSucessResponse(result, paging, nil))
	}
}
