# Factory Method Pattern
- Factory Method Pattern là một creational design pattern dùng để định nghĩa một phương thức tạo object trong lớp cha (hoặc interface), nhưng để các lớp con quyết định sẽ tạo ra object cụ thể nào.

- Nói dễ hiểu:
👉 Lớp cha biết “khi nào cần tạo”, lớp con quyết định “tạo cái gì”.

## Vấn đề nó giải quyết
- Giả sử bạn có code kiểu này:

```go
if orderType == "online" {
    order = OnlineOrder{}
} else if orderType == "store" {
    order = StoreOrder{}
}
```
    ❌ Logic tạo object nằm trong client
    ❌ Mỗi lần thêm loại mới phải sửa nhiều chỗ
    ❌ Phụ thuộc chặt vào struct cụ thể

## Ý tưởng cốt lõi
- Có Product interface (object được tạo)
- Có Creator chứa factory method
- Factory Method trả về Product interface
- Concrete Creator override factory method để tạo product cụ thể

## Cấu trúc (theo GoF)
```scss
Creator
 ├─ factoryMethod() -> Product
 └─ someBusinessLogic()

ConcreteCreatorA ──> ProductA
ConcreteCreatorB ──> ProductB
```

## Điểm mấu chốt cần nhớ
- Factory Method là một method, không phải hàm tự do
- Logic tạo object nằm trong subclass
- Thường được gọi bên trong business logic, không gọi trực tiếp từ client