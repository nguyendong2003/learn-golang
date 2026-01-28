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

func ListItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var paging common.Paging

		if err := c.ShouldBind(&paging); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))

			return
		}

		paging.Process()

		//
		var filter model.Filter

		if err := c.ShouldBind(&filter); err != nil {
			c.JSON(http.StatusBadRequest, common.ErrInvalidRequest(err))

			return
		}

		//
		store := storage.NewSQLStore(db)
		biz := business.NewListItemBusiness(store)

		result, err := biz.ListItem(c.Request.Context(), &filter, &paging)

		if err != nil {
			c.JSON(http.StatusBadRequest, err)

			return
		}

		c.JSON(http.StatusOK, common.NewSucessResponse(result, paging, filter))
	}
}
