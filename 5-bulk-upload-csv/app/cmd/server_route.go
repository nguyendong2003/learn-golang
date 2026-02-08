package cmd

func (server *ApiServer) route() {
	server.router.GET("/api/v1/healthy", server.mainHandler.HealthCheck())
	server.router.GET("/api/generate", server.generateHandler.Generate())
	{
		categoryGroup := server.router.Group("/api/v1/categories")
		categoryGroup.GET("", server.categoryHandler.GetList())
		categoryGroup.GET("/:id", server.categoryHandler.GetDetail())
		categoryGroup.POST("", server.categoryHandler.Create())
	}

	{
		warehouseGroup := server.router.Group("/api/v1/warehouses")
		warehouseGroup.GET("", server.warehouseHandler.GetList())
		warehouseGroup.GET("/:id", server.warehouseHandler.GetDetail())
		warehouseGroup.POST("", server.warehouseHandler.Create())
	}

	{
		productGroup := server.router.Group("/api/v1/products")
		productGroup.GET("", server.productHandler.GetList())
		productGroup.GET("/:id", server.productHandler.GetDetail())
		productGroup.POST("", server.productHandler.Create())
	}
}
