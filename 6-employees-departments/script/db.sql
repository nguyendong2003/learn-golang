CREATE DATABASE IF NOT EXISTS bulk_upload_demo
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

USE bulk_upload_demo;

CREATE TABLE departments (
  id INT AUTO_INCREMENT PRIMARY KEY,
  code VARCHAR(50) NOT NULL UNIQUE,
  name VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE employees (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  employee_code VARCHAR(50) NOT NULL,
  full_name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  department_id INT NOT NULL,
  salary DECIMAL(15,2) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_employee_department
    FOREIGN KEY (department_id) REFERENCES departments(id)
);
