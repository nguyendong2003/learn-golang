package handler

import (
	"bulk-upload-csv/interfaces"
	"encoding/csv"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GenerateHandler struct{}

func NewGenerateHandler() interfaces.GenerateHandlerInterface {
	return &GenerateHandler{}
}

func (h *GenerateHandler) Generate() gin.HandlerFunc {
	return func(c *gin.Context) {
		sqlPath, csvPath, err := GenerateSeedFiles()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// trả CSV (file quan trọng cho bulk upload)
		c.FileAttachment(csvPath, "inventory_transactions.csv")

		// nếu muốn trả SQL riêng → tạo endpoint khác
		_ = sqlPath
	}
}

/* =======================
   SEED STRUCTS
======================= */

type CategorySeed struct {
	ID   string
	Code string
}

type WarehouseSeed struct {
	ID   string
	Code string
}

type ProductSeed struct {
	ID           string
	SKU          string
	CategoryID   string
	CategoryCode string
}

/* =======================
   MAIN GENERATE
======================= */

func GenerateSeedFiles() (string, string, error) {
	rand.Seed(time.Now().UnixNano())

	now := time.Now().Format("2006-01-02 15:04:05")

	/* =======================
	   SQL FILE
	======================= */

	sqlPath := "./generate_data.sql"
	sqlFile, err := os.Create(sqlPath)
	if err != nil {
		return "", "", err
	}
	defer sqlFile.Close()

	/* =======================
	   1. CATEGORIES (1000)
	======================= */

	categories := make([]CategorySeed, 0, 1000)
	writeLine(sqlFile, "-- CATEGORIES")

	for i := 1; i <= 1000; i++ {
		id := uuid.New().String()
		code := fmt.Sprintf("CAT%04d", i)

		categories = append(categories, CategorySeed{
			ID:   id,
			Code: code,
		})

		if i%500 == 1 {
			writeLine(sqlFile, "INSERT INTO categories (id, code, name, created_at, updated_at) VALUES")
		}

		line := fmt.Sprintf(
			"('%s','%s','Category %d','%s','%s')",
			id, code, i, now, now,
		)

		if i%500 == 0 || i == 1000 {
			writeLine(sqlFile, line+";")
		} else {
			writeLine(sqlFile, line+",")
		}
	}

	/* =======================
	   2. WAREHOUSES (1000)
	======================= */

	warehouses := make([]WarehouseSeed, 0, 1000)
	writeLine(sqlFile, "\n-- WAREHOUSES")

	for i := 1; i <= 1000; i++ {
		id := uuid.New().String()
		code := fmt.Sprintf("WH%04d", i)

		warehouses = append(warehouses, WarehouseSeed{
			ID:   id,
			Code: code,
		})

		if i%500 == 1 {
			writeLine(sqlFile, "INSERT INTO warehouses (id, code, created_at, updated_at) VALUES")
		}

		line := fmt.Sprintf("('%s','%s','%s','%s')", id, code, now, now)

		if i%500 == 0 || i == 1000 {
			writeLine(sqlFile, line+";")
		} else {
			writeLine(sqlFile, line+",")
		}
	}

	/* =======================
	   3. PRODUCTS (10,000)
	======================= */

	products := make([]ProductSeed, 0, 10000)
	writeLine(sqlFile, "\n-- PRODUCTS")

	for i := 1; i <= 10000; i++ {
		id := uuid.New().String()
		cat := categories[rand.Intn(len(categories))]

		products = append(products, ProductSeed{
			ID:           id,
			SKU:          fmt.Sprintf("SKU%05d", i),
			CategoryID:   cat.ID,
			CategoryCode: cat.Code,
		})

		if i%1000 == 1 {
			writeLine(sqlFile, "INSERT INTO products (id, sku, category_id, created_at, updated_at) VALUES")
		}

		line := fmt.Sprintf(
			"('%s','SKU%05d','%s','%s','%s')",
			id, i, cat.ID, now, now,
		)

		if i%1000 == 0 || i == 10000 {
			writeLine(sqlFile, line+";")
		} else {
			writeLine(sqlFile, line+",")
		}
	}

	/* =======================
	   CSV FILE
	======================= */

	csvPath := "./inventory_transactions.csv"
	csvFile, err := os.Create(csvPath)
	if err != nil {
		return "", "", err
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// Header (bắt buộc đúng đề)
	writer.Write([]string{
		"product_sku",
		"category_code",
		"warehouse_code",
		"quantity",
		"transaction_type",
	})

	/* =======================
	   4. INVENTORY TRANSACTIONS (100,000)
	======================= */

	for i := 1; i <= 100000; i++ {
		p := products[rand.Intn(len(products))]
		w := warehouses[rand.Intn(len(warehouses))]

		txType := "IN"
		qty := rand.Intn(20) + 1

		if rand.Intn(2) == 0 {
			txType = "OUT"
			qty = -qty
		}

		row := []string{
			p.SKU,
			p.CategoryCode,
			w.Code,
			strconv.Itoa(qty),
			txType,
		}

		if err := writer.Write(row); err != nil {
			return "", "", err
		}
	}

	return sqlPath, csvPath, nil
}

/* =======================
   UTIL
======================= */

func writeLine(f *os.File, s string) {
	_, _ = f.WriteString(s + "\n")
}
