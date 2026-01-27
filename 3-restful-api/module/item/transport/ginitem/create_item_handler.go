package ginitem

import (
	"net/http"
	"restfulapi/common"
	"restfulapi/module/item/business"
	"restfulapi/module/item/model"
	"restfulapi/module/item/storage"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreatItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var data model.TodoItemCreation

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		store := storage.NewSQLStore(db)

		biz := business.NewCreateItemBusiness(store)

		if err := biz.CreateNewItem(c.Request.Context(), &data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data.Id))
	}
}
