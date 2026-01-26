package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type TodoItem struct {
	Id          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"` // omitempty: BỎ QUA field khi marshal nếu field đó rỗng / giá trị mặc định
}

type TodoItemCreation struct {
	Id          int    `json:"-" gorm:"column:id;"`
	Title       string `json:"title" gorm:"column:title;"`
	Description string `json:"description" gorm:"column:description;"`
}

func (TodoItemCreation) TableName() string { return "todo_items" }

func ex1() {
	now := time.Now().UTC()

	item := TodoItem{
		Id:          1,
		Title:       "This is item 1",
		Description: "This is description for item 1",
		Status:      "pending",
		CreatedAt:   &now,
		UpdatedAt:   nil,
	}

	// json.Marshal(item): Biến dữ liệu Go → JSON
	jsonData, err := json.Marshal(item)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(jsonData))

	// json.Unmarshal(): Biến JSON → dữ liệu Go
	jsonStr := "{\"id\":1,\"title\":\"This is item 1\",\"description\":\"This is description for item 1\",\"status\":\"pending\",\"created_at\":\"2026-01-26T07:09:33.33976644Z\",\"updated_at\":null}"
	var item2 TodoItem
	if err := json.Unmarshal([]byte(jsonStr), &item2); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(item2)
}

// Create api
func ex2() {
	now := time.Now().UTC()

	item := TodoItem{
		Id:          1,
		Title:       "This is item 1",
		Description: "This is description for item 1",
		Status:      "pending",
		CreatedAt:   &now,
		UpdatedAt:   nil,
	}

	// Gin
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": item,
		})
	})
	router.Run() // listens on 0.0.0.0:8080 by default
	// router.Run(":8000") // change port 8000

}

// Connection database
func ex3() {
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

	fmt.Println(db)
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
			items.GET("")
			items.GET("/:id")
			items.PATCH("/:id")
			items.DELETE("/:id")
		}
	}

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ping sucessfully",
		})
	})
	router.Run() // listens on 0.0.0.0:8080 by default // or: router.Run(":8000") // change port 8000

}

func CreatItem(db *gorm.DB) func(c *gin.Context) {
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
