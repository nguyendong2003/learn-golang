/*
I. Kiến trúc bulk-upload-goroutine
CSV reader (1 goroutine)

	↓

jobs channel (Employee batch)

	↓

N workers (each worker = tx riêng)

Mỗi worker:

# Mở transaction riêng

# Insert batch

Không share transaction → DB mới chạy song song thật
*/
package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"os"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "dongnt1:123456@tcp(127.0.0.1:3306)/bulk_upload_demo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	r := gin.Default()

	r.GET("/generate/departments-sql", GenerateDepartmentsSQL)
	r.GET("/generate/employees-csv", GenerateEmployeesCSV)
	r.POST("/employees/bulk-upload", BulkUploadEmployees(db))
	r.POST("/employees/bulk-upload-goroutine", BulkUploadEmployeesGoroutine(db))

	r.Run(":8080")

}

// BatchSize = 5000, WorkerCount = 4 => 20.25s with 1 million records employees
// BatchSize = 10000, WorkerCount = 8 => 18.64s with 1 million records employees
// BatchSize = 20000, WorkerCount = 8 => chạy quá lâu với 1 million records employees
// BatchSize = 5000, WorkerCount = 8 => 34.62s với 1 million records employees
const (
	BatchSize   = 10000 // 5000
	WorkerCount = 8     // test: 2, 4, 8
)

type Employee struct {
	EmployeeCode string
	FullName     string
	Email        string
	DepartmentID int
	Salary       float64
}

type Department struct {
	ID   int    `gorm:"column:id"`
	Code string `gorm:"column:code"`
	Name string `gorm:"column:name"`
}

func (Department) TableName() string {
	return "departments"
}

type EmployeeCSVRow struct {
	EmployeeCode   string
	FullName       string
	Email          string
	DepartmentCode string
	Salary         float64
}

type CSVError struct {
	Row     int         `json:"row"`
	Field   string      `json:"field"`
	Value   interface{} `json:"value"`
	Message string      `json:"message"`
}

func LoadDepartmentCache(db *gorm.DB) (map[string]int, error) {
	var depts []Department

	if err := db.Select("id, code").Find(&depts).Error; err != nil {
		return nil, err
	}

	cache := make(map[string]int, len(depts))
	for _, d := range depts {
		cache[d.Code] = d.ID
	}
	return cache, nil
}

// API này mất 25.41s để upload 1 triệu records employees
func BulkUploadEmployees(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		fileHeader, err := c.FormFile("file")
		if err != nil {
			c.JSON(400, gin.H{"error": "file required"})
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer file.Close()

		reader := csv.NewReader(bufio.NewReader(file))
		reader.FieldsPerRecord = -1

		// read header
		header, err := reader.Read()
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid csv header"})
			return
		}

		expected := []string{
			"employee_code", "full_name", "email", "department_code", "salary",
		}
		if !reflect.DeepEqual(header, expected) {
			c.JSON(400, gin.H{"error": "csv header not match"})
			return
		}

		deptCache, err := LoadDepartmentCache(db)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		var errors []CSVError
		var employees []Employee
		rowIndex := 1

		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			rowIndex++

			salary, err := strconv.ParseFloat(record[4], 64)
			if err != nil || salary < 0 {
				errors = append(errors, CSVError{
					Row:     rowIndex,
					Field:   "salary",
					Value:   record[4],
					Message: "salary must be >= 0",
				})
				continue
			}

			deptID, ok := deptCache[record[3]]
			if !ok {
				errors = append(errors, CSVError{
					Row:     rowIndex,
					Field:   "department_code",
					Value:   record[3],
					Message: "department not found",
				})
				continue
			}

			employees = append(employees, Employee{
				EmployeeCode: record[0],
				FullName:     record[1],
				Email:        record[2],
				DepartmentID: deptID,
				Salary:       salary,
			})
		}

		if len(errors) > 0 {
			c.JSON(400, gin.H{
				"errors": errors,
			})
			return
		}

		// TRANSACTION + BATCH INSERT
		err = db.Transaction(func(tx *gorm.DB) error {
			batchSize := 1000
			for i := 0; i < len(employees); i += batchSize {
				end := i + batchSize
				if end > len(employees) {
					end = len(employees)
				}
				if err := tx.CreateInBatches(employees[i:end], batchSize).Error; err != nil {
					return err
				}
			}
			return nil
		})

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"message":       "bulk upload success",
			"total_records": len(employees),
		})
	}
}

// Goroutine + Channel + Batch Insert
func BulkUploadEmployeesGoroutine(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		fileHeader, err := c.FormFile("file")
		if err != nil {
			c.JSON(400, gin.H{"error": "file required"})
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer file.Close()

		reader := csv.NewReader(bufio.NewReader(file))
		reader.FieldsPerRecord = -1

		// header
		header, err := reader.Read()
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid csv header"})
			return
		}

		expected := []string{
			"employee_code", "full_name", "email", "department_code", "salary",
		}
		if !reflect.DeepEqual(header, expected) {
			c.JSON(400, gin.H{"error": "csv header not match"})
			return
		}

		deptCache, err := LoadDepartmentCache(db)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		jobs := make(chan []Employee, WorkerCount)
		errCh := make(chan error, WorkerCount)

		// start workers
		for i := 0; i < WorkerCount; i++ {
			go employeeWorker(db, jobs, errCh)
		}

		rowIndex := 1
		var errors []CSVError
		batch := make([]Employee, 0, BatchSize)

		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			rowIndex++

			salary, err := strconv.ParseFloat(record[4], 64)
			if err != nil || salary < 0 {
				errors = append(errors, CSVError{
					Row: rowIndex, Field: "salary", Value: record[4],
					Message: "salary must be >= 0",
				})
				continue
			}

			deptID, ok := deptCache[record[3]]
			if !ok {
				errors = append(errors, CSVError{
					Row: rowIndex, Field: "department_code", Value: record[3],
					Message: "department not found",
				})
				continue
			}

			batch = append(batch, Employee{
				EmployeeCode: record[0],
				FullName:     record[1],
				Email:        record[2],
				DepartmentID: deptID,
				Salary:       salary,
			})

			if len(batch) == BatchSize {
				jobs <- batch
				batch = make([]Employee, 0, BatchSize)
			}
		}

		if len(errors) > 0 {
			close(jobs)
			c.JSON(400, gin.H{"errors": errors})
			return
		}

		if len(batch) > 0 {
			jobs <- batch
		}

		close(jobs)

		// wait workers
		for i := 0; i < WorkerCount; i++ {
			if err := <-errCh; err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
		}

		c.JSON(200, gin.H{
			"message": "bulk upload goroutine finished",
			"workers": WorkerCount,
			"batch":   BatchSize,
		})
	}
}

// Worker function
func employeeWorker(db *gorm.DB, jobs <-chan []Employee, errCh chan<- error) {
	for batch := range jobs {
		err := db.Transaction(func(tx *gorm.DB) error {
			return tx.Create(&batch).Error
		})
		if err != nil {
			errCh <- err
			return
		}
	}
	errCh <- nil
}

//

func GenerateDepartmentsSQL(c *gin.Context) {
	file, err := os.Create("departments.sql")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString("USE bulk_upload_demo;\n")
	writer.WriteString("INSERT INTO departments(code, name) VALUES\n")

	for i := 1; i <= 100000; i++ {
		line := fmt.Sprintf("('DEPT_%05d','Department %05d')", i, i)
		if i < 100000 {
			line += ",\n"
		} else {
			line += ";\n"
		}
		writer.WriteString(line)
	}

	writer.Flush()

	c.FileAttachment("departments.sql", "departments.sql")
}

func GenerateEmployeesCSV(c *gin.Context) {
	file, err := os.Create("employees.csv")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// header
	writer.Write([]string{
		"employee_code",
		"full_name",
		"email",
		"department_code",
		"salary",
	})

	for i := 1; i <= 1_000_000; i++ {
		deptCode := fmt.Sprintf("DEPT_%05d", rand.Intn(100000)+1)
		row := []string{
			fmt.Sprintf("EMP_%07d", i),
			fmt.Sprintf("Employee %d", i),
			fmt.Sprintf("emp%d@company.com", i),
			deptCode,
			strconv.Itoa(rand.Intn(5000) + 1000),
		}
		writer.Write(row)
	}

	c.FileAttachment("employees.csv", "employees.csv")
}
