# Cai dat go
Tải từ trang chủ golang -> giải nén ra lấy thư mục go
sudo mv ~/Downloads/go /usr/local/ 
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Cai dat cong cu run and debug in vscode
go install -v github.com/go-delve/delve/cmd/dlv@latest

# Tao project hello world
go mod init learngo

go run main.go
go run .

# Kiến thức Golang
- Tên package viết thường hết, tổ chức theo kiểu thư mục  (vd: fmt, math/rand) (fmt là viết tắt của format)
- Exported Names là biến được import từ 1 package khác, nó phải đươc viết hoa chữ cái đầu.  Nếu biến không được viết hoa chữ cái đầu thì là Unexported Names
- Kiểu dữ liệu của biến nằm ở đằng sau biến
- Trong Go kiểu int dài 32 bit hay 64 bit tùy theo kiến trúc máy tính, 
    + Nếu máy tính 32 bit thì int dài 32 bit (int32)
    + Nếu máy tính 64 bit thì int dài 64 bit (int64)
- Hằng số không có giới hạn bit cho đến khi bị ép vào một kiểu cụ thể
    + Ví dụ:

```go
const (
    Big   = 1 << 100  // 2^100  // Không bị tràn số
    Small = Big >> 99 // 2
)

m := float64(Big)
```

- Vòng lặp trong Go:

```go
for i := 0; i < 10; i++ {}      // for truyền thống
for condition {}                // while(condition)
for {}                          // while(true) {}
for { ... if !cond { break } }  // do {} while(condition)
```

- Mặc định switch case đã có câu lệnh break trong mỗi case rồi, nên không cần thêm câu lệnh break vào các case
- Từ khóa `defer` trong Go 
    + Dùng để trì hoãn việc gọi một hàm cho đến khi hàm bao quanh kết thúc.
    + Chạy khi: return, panic, kết thúc function
    + Thứ tự thực hiện: defer chạy theo LIFO (Stack), khai báo sau được chạy trước
    + Thời điểm đánh giá tham số: Tham số được evaluate ngay khi defer được gọi, Không phải lúc hàm defer thực sự chạy
    ```go
    x := 10
    defer fmt.Println(x)
    x = 20   // vẫn in 10
    ```
    + Ứng dụng phổ biến: Đóng tài nguyên, Cleanup Code
    + defer giúp đảm bảo cleanup code luôn chạy, kể cả khi panic,
chạy theo LIFO và tham số được evaluate ngay lúc defer."

- `Pointers`: 
    + Không giống C, Go KHÔNG cho phép làm toán trực tiếp trên con trỏ. (nếu trong mảng thi Go chỉ cho phép truy cập bằng index, không dùng pointer arithmetic)
    ```go
    s := []int{10, 20, 30}
    fmt.Println(s[1])
    ```

    + `nil` giống với `null` trong c
- `Arrays`
- `Slice`: (dynamically-sized) trỏ đến 1 phần trong 1 array, thay đổi slice thì thay đổi array
