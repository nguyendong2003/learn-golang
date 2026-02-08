package cmd

import (
	"bulk-upload-csv/handler"
	"bulk-upload-csv/repository"
	"bulk-upload-csv/service"
	"log"
)

func (server *ApiServer) dependenciesInjection() {
	server.dbRepository = repository.NewDbRepository(server.config.Database.Dsn)
	err := server.dbRepository.InitializeDB()
	if err != nil {
		log.Fatal(err)
	}

	server.mainHandler = handler.NewMainHandler()
	server.generateHandler = handler.NewGenerateHandler()

	server.categoryRepository = repository.NewCategoryRepository(server.dbRepository.GetDB())
	server.categoryService = service.NewCategoryService(server.categoryRepository)
	server.categoryHandler = handler.NewCategoryHandler(server.categoryService)

	server.warehouseRepository = repository.NewWarehouseRepository(server.dbRepository.GetDB())
	server.warehouseService = service.NewWarehouseService(server.warehouseRepository)
	server.warehouseHandler = handler.NewWarehouseHandler(server.warehouseService)

	server.productRepository = repository.NewProductRepository(server.dbRepository.GetDB())
	server.productService = service.NewProductService(server.productRepository)
	server.productHandler = handler.NewProductHandler(server.productService)

}
