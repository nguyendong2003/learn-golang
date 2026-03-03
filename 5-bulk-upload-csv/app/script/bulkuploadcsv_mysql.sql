CREATE DATABASE IF NOT EXISTS bulkuploadcsv
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

USE bulkuploadcsv;

CREATE TABLE categories (
    id CHAR(36) PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    deleted_at DATETIME,
    
    INDEX idx_categories_deleted_at (deleted_at)
);

CREATE TABLE products (
    id CHAR(36) PRIMARY KEY,
    sku VARCHAR(50) NOT NULL UNIQUE,
    category_id CHAR(36) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    deleted_at DATETIME,
    
    -- INDEX idx_products_category_id (category_id),
    INDEX idx_products_deleted_at (deleted_at),
    
    FOREIGN KEY(category_id) REFERENCES categories(id)
);

CREATE TABLE warehouses (
    id CHAR(36) PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    deleted_at DATETIME,
    
    INDEX idx_warehouses_deleted_at (deleted_at)
);

CREATE TABLE inventory_transactions (
    id CHAR(36) PRIMARY KEY,
    product_id CHAR(36) NOT NULL,
    category_id CHAR(36) NOT NULL,
    warehouse_id CHAR(36) NOT NULL,
    quantity INT NOT NULL,
    transaction_type VARCHAR(10) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    deleted_at DATETIME,
    
    FOREIGN KEY (product_id) REFERENCES products(id),
    FOREIGN KEY (category_id) REFERENCES categories(id),
    FOREIGN KEY (warehouse_id) REFERENCES warehouses(id),
    
    -- INDEX idx_inv_product_id (product_id),
    -- INDEX idx_inv_category_id (category_id),
    -- INDEX idx_inv_warehouse_id (warehouse_id),
    INDEX idx_inv_deleted_at (deleted_at)
);





