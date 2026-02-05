package main

import (
	"log"
	"net/http"
	"os"
	ginitem "restfulapi/module/item/transport/ginitem"

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

	// router.Use(middleware.Recovery())

	// v1 := router.Group("/v1", middleware.Recovery())
	v1 := router.Group("/v1")
	{
		items := v1.Group("/items")
		{
			// items.POST("", middleware.Recovery(), ginitem.CreatItem(db))
			items.POST("", ginitem.CreatItem(db))
			items.GET("", ginitem.ListItem(db))
			items.GET("/:id", ginitem.GetItem(db))
			items.PATCH("/:id", ginitem.UpdateItem(db))
			items.DELETE("/:id", ginitem.DeleteItem(db))
		}
	}

	router.GET("/ping", func(c *gin.Context) {
		// go func() {
		// 	defer common.Recovery()
		// 	fmt.Println([]int{}[0])
		// }()

		c.JSON(http.StatusOK, gin.H{
			"message": "ping sucessfully",
		})
	})
	router.Run() // listens on 0.0.0.0:8080 by default // or: router.Run(":8000") // change port 8000

}
