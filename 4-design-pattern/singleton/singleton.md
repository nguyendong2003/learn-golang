# Singleton
- Singleton Pattern là một creational design pattern dùng để đảm bảo rằng một class (hoặc struct) chỉ có duy nhất một instance trong suốt vòng đời của chương trình, và cung cấp một điểm truy cập chung đến instance đó.

- Nói ngắn gọn:
👉 Tạo một lần – dùng mọi nơi.

- Chú ý:
Singleton Pattern đảm bảo một class chỉ có một instance và cung cấp global access, nhưng nên dùng cẩn thận, đặc biệt trong Go.

## Bài toán Singleton giải quyết
- Dùng khi bạn có những tài nguyên toàn hệ thống chỉ nên có một, ví dụ:
    + Logger
    + Database connection pool
    + Application config    
    + Cache

- Nếu tạo nhiều instance:
    ❌ Tốn tài nguyên
    ❌ Trạng thái không đồng bộ
    ❌ Dễ sinh bug khó truy vết

## Ý tưởng cốt lõi
- Ẩn việc khởi tạo (không cho tạo tùy ý)
- Tự quản lý instance
- Luôn trả về cùng một instance