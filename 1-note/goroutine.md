1. Khái niệm cơ bản: Goroutine là gì?
Đừng nhầm lẫn Goroutine với Thread (luồng) của hệ điều hành.

Siêu nhẹ: Một Goroutine chỉ tốn khoảng 2KB bộ nhớ ban đầu, trong khi OS Thread tốn khoảng 1-2MB.

Go Scheduler: Go có một bộ điều phối riêng để chạy hàng ngàn (thậm chí hàng triệu) Goroutine trên một số ít OS Threads.

2. Channels - "Đừng chia sẻ bộ nhớ, hãy giao tiếp"
Triết lý của Go là: “Do not communicate by sharing memory; instead, share memory by communicating.”

Unbuffered Channels: Giúp các Goroutine đồng bộ hóa (chặn cho đến khi có người gửi và người nhận).

Buffered Channels: Cho phép gửi một lượng dữ liệu nhất định mà không cần người nhận sẵn sàng ngay lập tức.

Directional Channels: Giới hạn channel chỉ để gửi (chan<-) hoặc chỉ để nhận (<-chan).

3. Đồng bộ hóa (Synchronization)
Đôi khi bạn không cần truyền dữ liệu mà chỉ cần đợi các tác vụ xong xuôi.

WaitGroup (sync.WaitGroup): Cách chuẩn để đợi một nhóm Goroutine hoàn thành. Sử dụng Add(), Done(), và Wait().

Mutex (sync.Mutex): Dùng để bảo vệ dữ liệu dùng chung (shared state) khi nhiều Goroutine cùng truy cập vào một biến, tránh tình trạng Race Condition.

4. Các Pattern phổ biến
Khi đã nắm cơ bản, bạn nên tìm hiểu cách phối hợp chúng:

Select statement: Giống như switch nhưng dành cho channel. Nó cho phép một Goroutine đợi trên nhiều hoạt động channel cùng lúc.

Worker Pool: Tạo ra một số lượng Goroutine cố định để xử lý một hàng đợi công việc (giúp kiểm soát tài nguyên).

Fan-in / Fan-out: Cách phân tán công việc ra nhiều Goroutine và thu thập kết quả lại một nơi.

5. Xử lý lỗi và Hủy bỏ (Context)
Trong thực tế, bạn không muốn một Goroutine chạy mãi mãi nếu yêu cầu đã bị hủy.

Package context: Dùng để truyền tín hiệu hủy (cancel), deadline, hoặc các giá trị đi xuyên suốt các tầng xử lý. Đây là kiến thức bắt buộc nếu bạn làm Web hay Microservices với Go.

Một ví dụ đơn giản về Race Condition:
Nếu bạn chạy đoạn code sau mà không có Mutex hay Channel, kết quả có thể sẽ không như mong đợi:

Go
var count = 0
for i := 0; i < 1000; i++ {
    go func() { count++ }() // Race condition ở đây!
}