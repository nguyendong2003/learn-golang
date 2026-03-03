-- 1. Tạo database (Chạy lệnh này riêng lẻ nếu bạn chưa có db)
-- CREATE DATABASE elearning;

-- 2. Kết nối vào database elearning và chạy các lệnh dưới đây
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    email VARCHAR(255) NOT NULL UNIQUE,
    username VARCHAR(100) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    full_name VARCHAR(255),
    avatar TEXT,
    is_active BOOLEAN DEFAULT true
);

-- Tạo index cho deleted_at để tối ưu soft delete của GORM
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
