# Strategy pattern
- Strategy Pattern (mẫu chiến lược) là một design pattern thuộc nhóm hành vi (Behavioral). Nó cho phép bạn định nghĩa nhiều thuật toán/chiến lược khác nhau, đóng gói mỗi chiến lược vào một class riêng, và hoán đổi chúng với nhau trong lúc chạy chương trình mà không cần sửa code của đối tượng sử dụng chúng.

- Nói đơn giản:
👉 Cùng một việc, nhưng có nhiều cách làm – Strategy giúp bạn thay đổi cách làm đó một cách linh hoạt.

## Ý tưởng cốt lõi
- Tách phần “làm thế nào” ra khỏi phần “ai sử dụng”
- Đối tượng chính (Context) không cần biết chi tiết thuật toán
- Chỉ cần biết: “Tôi có một Strategy, gọi nó là chạy”


## Ví dụ đời thường
Bạn đi từ nhà đến công ty:
🚶 Đi bộ
🚲 Đi xe đạp
🚗 Đi ô tô
🚌 Đi xe buýt
Mục tiêu: đến công ty
Chiến lược: cách di chuyển
→ Thay đổi chiến lược mà không thay đổi “bạn”.