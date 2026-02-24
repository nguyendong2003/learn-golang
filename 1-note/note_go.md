# Repo

[Udemy: Go - The Complete Guide - Maximilian Schwarzmüller] <https://github.com/mschwarzmueller/go-complete-guide-resources>

# Cai dat go

```bash
# Tải từ trang chủ golang -> giải nén ra lấy thư mục go
sudo mv ~/Downloads/go /usr/local/ 
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```
# Cai dat cong cu run and debug in vscode
```bash
go install -v github.com/go-delve/delve/cmd/dlv@latest
```
# Tao project hello world
```bash
go mod init learngo

go run main.go
go run .
```

# Cai dat Postman
```bash
sudo snap install postman
```


# Kiến thức Golang
- Tên package viết thường hết, tổ chức theo kiểu thư mục  (vd: fmt, math/rand) (fmt là viết tắt của format)
- `Exported Names` là biến được import từ 1 package khác, nó phải đươc viết hoa chữ cái đầu.  Nếu biến không được viết hoa chữ cái đầu thì là `Unexported Names`
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

- Mặc định `switch case` đã có câu lệnh break trong mỗi case rồi, nên không cần thêm câu lệnh `break` vào các case
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

### 1. `module, package, import`
    1. `Package (Gói)`

    - Package là cấp độ cơ bản nhất để tổ chức code. Một package bao gồm tất cả các file `.go` nằm trong cùng một thư mục.
    - Mục đích: Gom nhóm các chức năng liên quan (ví dụ: package `math` chứa các hàm tính toán).
    - Khai báo: Mọi file .go phải bắt đầu bằng dòng `package <tên_package>`
    - Quy tắc: Các file trong cùng một thư mục phải có cùng tên package.
    - Tính đóng gói (Visibility):
        + Tên hàm/biến bắt đầu bằng `chữ in hoa` (Vd: Calculate): `Public` (có thể dùng ở package khác).
        + Tên bắt đầu bằng `chữ thường` (Vd: getRate): `Private` (chỉ dùng nội bộ trong package đó).
    
    2. `Module`

    - `Module` là một tập hợp các package liên quan được phân phối cùng nhau. Đây là đơn vị quản lý phiên bản (versioning) và phụ thuộc (dependency).
    - File `go.mod`: Mỗi module phải có một file `go.mod` ở thư mục gốc. File này định nghĩa:
        + Tên của module (thường là đường dẫn URL như github.com/user/project).
        + Phiên bản Go.
        + Các thư viện bên ngoài (dependencies) mà dự án cần.
    - Khởi tạo: Dùng lệnh `go mod init <tên-module>`

    3. `Import`
    - Để sử dụng code từ một package khác, bạn dùng từ khóa `import`.

    - Cách sử dụng cơ bản:
        ```go
        package main

        import (
            "fmt"       // Package từ thư viện chuẩn của Go
            "math/rand"  // Package con
            "github.com/google/uuid" // Package từ bên ngoài (thư viện bên thứ 3)
        )

        func main() {
            fmt.Println(uuid.New())
        }
        ```
    - Các kiểu Import đặc biệt:
        + `Alias Import`: Đặt tên khác để tránh trùng lặp.

            ```go
            import m "math"
            // Sử dụng: m.Sin(0)
            ```

        + `Blank Import (_)`: Chỉ chạy hàm init() của package đó mà không gọi trực tiếp các hàm khác. Thường dùng cho các driver database.

            ```go
            import _ "github.com/lib/pq"
            ```

    4. `Build and run project`:
        ```bash
        go build
        
        ./learngo
        ```

### 2. `Pointers`: 
- Không giống C, Go KHÔNG cho phép làm toán trực tiếp trên con trỏ. (nếu trong mảng thi Go chỉ cho phép truy cập bằng index, không dùng pointer arithmetic)
    ```go
    s := []int{10, 20, 30}
    fmt.Println(s[1])
    ```

- `nil` giống với `null` trong c
2. `Arrays`
- Array là một dãy các phần tử có độ dài cố định và cùng kiểu dữ liệu
- Giá trị mặc định (Zero value): Khi bạn khai báo một mảng mà chưa gán giá trị, Go sẽ tự động lấp đầy nó bằng giá trị "không" của kiểu dữ liệu đó (ví dụ: số `0` cho `int`, "" cho `string`, `false` cho `bool`)
- Array là một "Value Type" (Kiểu giá trị)
- `Array` không có `Slice Header`
- Chú ý:
    + Trong nhiều ngôn ngữ, khi bạn truyền mảng vào hàm, nó truyền dưới dạng tham chiếu (pointer).
    + Trong Go, Array là kiểu giá trị. Khi bạn gán mảng A cho mảng B, hoặc truyền mảng vào một hàm, Go sẽ sao chép toàn bộ nội dung của mảng đó sang một vùng nhớ mới.
        ```go
        a := [3]string{"Go", "Rust", "C++"}
        b := a          // copy
        ```
    + Hệ quả: Nếu mảng của bạn có 1 triệu phần tử, việc gán mảng sẽ rất tốn kém tài nguyên và làm chậm chương trình. Đây là lý do vì sao chúng ta thường dùng `Slice` hoặc `truyền Con trỏ mảng (Array Pointer)`.
    + Duyệt mảng với `range`
        ```go
        fruits := [3]string{"Apple", "Banana", "Cherry"}

        for index, value := range fruits {
            fmt.Printf("Vị trí %d có quả: %s\n", index, value)
        }
        ```



### 3. `Slice`:
- Slice là view (cửa sổ) trỏ vào array
- Slice không phải là một mảng (array). Nó là một Descriptor (bản mô tả) nằm trên một mảng ẩn (backing array)
- Không lưu dữ liệu trực tiếp
- `Slice Header` là một struct gồm 3 trường:
    + Pointer (Con trỏ): Con trỏ trỏ đến phần tử đầu tiên của mảng ẩn mà slice có quyền truy cập.
    + Length (Độ dài - len): Số lượng phần tử hiện có trong Slice.
    + Capacity (Sức chứa - cap): Số lượng phần tử tối đa mà slice có thể chứa tính từ vị trí con trỏ mà không cần cấp phát lại bộ nhớ.
    ```go
    type slice struct {
        ptr *T   // con trỏ tới phần tử đầu tiên
        len int
        cap int
    }
    ```
- Chú ý phân biệt `array` và `slice`:
    ```go
    // có số lượng cụ thể  -> Đây là array. Nó không có Slice Header. Toàn bộ dữ liệu nằm trực tiếp trong biến a
    a := [3]int{1, 2, 3}  
        
        
    // ngoặc vuông rỗng -> Đây là Slice. Go sẽ ngầm tạo ra một mảng ẩn (backing array) kích thước 3, sau đó tạo một Slice Header cho a để quản lý mảng ẩn đó.
    a := []int{1, 2, 3}
    ```

- Ví dụ: Giả sử ta có một mảng ẩn (Backing Array) gồm 5 phần tử: [10, 20, 30, 40, 50]
    ```go
    // Bước 1: Khởi tạo Slice
    data := [5]int{10, 20, 30, 40, 50}
    s1 := data[1:4] // Lấy từ index 1 đến 3

    // Bước 2: Cắt tiếp từ slice (Reslicing)
    s2 := s1[1:3] // Lấy từ index 1 đến 2 của s1

    // Bước 3: Trường hợp 1 -> Append khi còn Capacity
    s2 = append(s2, 100)

    // Bước 4: Trường hợp 2 -> Append khi hết Capacity
    s2 = append(s2, 200)

    fmt.Println(s2)     // [30 40 100 200]
    fmt.Println(data)   // [10 20 30 40 100]
    s2[0] = 9999
    fmt.Println(s2)     // [9999 40 100 200]
    fmt.Println(data)   // [10 20 30 40 100]
    ```
    + Bước 1: s1 sẽ có:
        + Ptr: Trỏ vào địa chỉ của data[1] (giá trị 20)
        + Len = 3: Chứa {20, 30, 40}
        + Cap = 4: Vì từ vị trí data[1] đến cuối mảng data còn 4 ô (chỉ số 1, 2, 3, 4)
    + Bước 2: s2 nhìn vào mảng thông qua "cửa sổ" của s1
        + Ptr: Trỏ vào địa chỉ của s1[1] (tương đương data[2], giá trị 30).
        + Len = 2: Chứa {30, 40}
        + Cap = 3: Vì từ vị trí data[2] đến cuối mảng còn 3 ô (chỉ số 2, 3, 4)
    + Bước 3: 
        + Go thấy Cap của s2 là 3, hiện tại mới dùng Len là 2. Vẫn còn dư 1 chỗ!
        + Nó sẽ ghi đè giá trị 100 vào ô tiếp theo trong mảng ẩn (vị trí data[4])
        + Kết quả: data bây giờ là [10, 20, 30, 40, 100]. Cả s1 và s2 đều bị ảnh hưởng vì dùng chung mảng ẩn.
    + Bước 4: 
        + Lúc này s2 đã đầy (Len=3, Cap=3). Không còn chỗ trong mảng data nữa.
        + Go sẽ:
            1. Tạo một mảng ẩn mới (ví dụ kích thước 6)
            2. Copy {30, 40, 100} sang mảng mới
            3. Thêm 200 vào
        + Kết quả: s2 bây giờ trỏ sang một vùng nhớ hoàn toàn khác. Mọi thay đổi trên s2 từ nay về sau không ảnh hưởng đến data hay s1 nữa.

### 4.  Cách hoạt động của `append`
- `append` không chỉ đơn giản là "thêm một phần tử". Nó là một hàm thông minh có khả năng tự động quản lý bộ nhớ
- `appen` luôn trả về một `Slice Header`
-  Quy trình 3 bước mà append thực hiện mỗi khi được gọi:
    1. Kiểm tra Sức chứa (Capacity)
    - Khi bạn gọi append(s, x), Go sẽ so sánh độ dài hiện tại (`len`) và sức chứa (`cap`) của slice s.
        + `Trường hợp A (Còn chỗ)`: Nếu `len < cap`, Go chỉ đơn giản là đặt phần tử x vào vị trí tiếp theo trong mảng ẩn (backing array) và tăng `len` lên 1.
        + `Trường hợp B (Hết chỗ)`: Nếu `len == cap`, mảng ẩn hiện tại đã đầy. Go phải thực hiện một quá trình gọi là "Grow" (Tăng trưởng).

    2. Chiến lược Tăng trưởng (Grow) - Go 1.18+
    - Khi hết chỗ, Go không chỉ lấy thêm 1 ô nhớ (vì rất tốn kém nếu phải làm nhiều lần). Nó sẽ cấp phát một mảng ẩn mới lớn hơn theo công thức:
        + Nếu `cap cũ < 256`: cap mới sẽ gấp đôi cap cũ `(2x)`.
        + Nếu `cap cũ >= 256`: cap mới sẽ tăng theo công thức: `newcap += (newcap + 3*256) / 4`. (`Tăng khoảng 25% đến 63% tùy kích thước`, giúp việc tăng trưởng mượt mà hơn, không nhảy vọt quá lớn khi mảng đã khổng lồ).

    3. Di chuyển dữ liệu
    - Sau khi đã có mảng mới với cap mới:
        1. Go copy toàn bộ phần tử từ mảng cũ sang mảng mới.
        2. Chèn phần tử mới vào.
        3. Trả về một Slice Header mới (chứa `Ptr` mới trỏ sang mảng mới, `len` mới và `cap` mới).

- Ví dụ minh họa "Bẫy" chung mảng ẩn
    1. Ví dụ 1:
    ```go
    func main() {
        // Khởi tạo slice có len=3, cap=3
        a := []int{1, 2, 3} 
        
        // b và c cùng append vào a
        // Vì a đã hết cap (3/3), cả b và c đều được cấp phát mảng mới riêng biệt
        b := append(a, 4) 
        c := append(a, 5) 

        fmt.Println(a) // [1 2 3]
        fmt.Println(b) // [1 2 3 4]
        fmt.Println(c) // [1 2 3 5] -> b và c không liên quan nhau
    }
    ```
    - Phân tích quá trình thực thi:
        1. Dòng b := append(a, 4):
        - Go kiểm tra a thấy cap=3, len=3 (đã đầy).
        - Go cấp phát một mảng ẩn mới hoàn toàn (ví dụ mảng X, có cap=6).   
        - Go copy 1, 2, 3 sang mảng X và thêm 4.
        - b nhận được một Slice Header mới: Ptr trỏ tới mảng X, len=4, cap=6.

        2. Dòng c := append(a, 5):
        - Go lại kiểm tra a. Lưu ý: Lúc này a vẫn là [1, 2, 3] với cap=3, len=3
        - Go lại cấp phát một mảng ẩn mới khác nữa (ví dụ mảng Y, có cap=6).
        - Go copy 1, 2, 3 từ a sang mảng Y và thêm 5.
        - c nhận được một Slice Header mới: Ptr trỏ tới mảng Y, len=4, cap=6.

    - Kết quả:
        + b và c là hai Slice Header khác nhau, trỏ tới hai mảng ẩn khác nhau trong bộ nhớ. Đó là lý do tại sao thay đổi ở b không bao giờ ảnh hưởng đến c.

    2. Ví dụ 2:
    ```go
    func main() {
        // Khởi tạo slice có len=3 nhưng cap=10
        a := make([]int, 3, 10)
        a[0], a[1], a[2] = 1, 2, 3
        
        // a còn dư chỗ, append không tạo mảng mới mà dùng chung mảng của a
        b := append(a, 4) // b ghi số 4 vào vị trí index 3 của mảng ẩn
        c := append(a, 5) // c ghi đè số 5 vào đúng vị trí index 3 đó!

        fmt.Println(b) // [1 2 3 5] -> Số 4 đã bị biến thành số 5!
        fmt.Println(c) // [1 2 3 5]
    }
    ```
    - Phân tích quá trình thực thi:
        1. Bước 1: Khởi tạo `a := make([]int, 3, 10)`
        - Go cấp phát 1 mảng ẩn có 10 ô trống.
        - Trả về Slice Header a: Ptr (trỏ vào ô 0), Len = 3, Cap = 10.
        - Mảng ẩn hiện tại: [1, 2, 3, 0, 0, 0, 0, 0, 0, 0]

        2. Bước 2: Thực hiện `b := append(a, 4)`
        - Go nhìn vào Header của a, thấy Cap=10 và Len=3. Còn tận 7 chỗ trống!
        - Go không tạo mảng mới. Nó ghi thẳng số 4 vào ô tiếp theo (index 3) của mảng ẩn hiện tại.
        - Trả về Slice Header b: Ptr (vẫn trỏ vào ô 0), Len = 4, Cap = 10.
        - Mảng ẩn lúc này: [1, 2, 3, 4, 0, 0, 0, 0, 0, 0]

        3. Bước 3: Thực hiện `c := append(a, 5)`
        - Đây là chỗ then chốt: Go lại nhìn vào Header của a, không phải của b.
        - Trong Header của a, Len vẫn chỉ là 3. Ô tiếp theo đối với a vẫn là vị trí index 3.
        - Vì còn trống (Cap=10), Go ghi số 5 vào vị trí index 3 của mảng ẩn.
        - Hành động này ghi đè lên số 4 mà b vừa mới đặt vào đó.
        - Trả về Slice Header c: Ptr (vẫn trỏ vào ô 0), Len = 4, Cap = 10.
        - Mảng ẩn cuối cùng: [1, 2, 3, 5, 0, 0, 0, 0, 0, 0]

- Cảnh báo: 
    + Như ví dụ trên ta thấy rằng: Vì append có thể trả về một con trỏ tới mảng mới, bạn bắt buộc phải gán kết quả trả về cho chính slice đó: `s = append(s, element)`

- Lưu ý: `append` luôn trả về slice header mới, còn array có thể cũ hoặc mới, KHÔNG bao giờ sửa header cũ
    1. Ví dụ 1: Header mới – array CŨ (không realloc)
        ```go
        a := make([]int, 2, 4)
        a[0], a[1] = 1, 2

        b := append(a, 3)

        fmt.Println("a:", a, ", len: ", len(a), ", cap: ", cap(a))
        fmt.Println("b:", b, ", len: ", len(b), ", cap: ", cap(b))

        fmt.Printf("&a = %p\n", &a)
        fmt.Printf("&b = %p\n", &b)

        fmt.Printf("&a[0] = %p\n", &a[0])
        fmt.Printf("&b[0] = %p\n", &b[0])
        ```

        - OUTPUT:
        ```bash
        a: [1 2] , len:  2 , cap:  4
        b: [1 2 3] , len:  3 , cap:  4
        &a = 0xc0000ac030
        &b = 0xc0000ac048       ← KHÁC (header mới)
        &a[0] = 0xc0000b0000    
        &b[0] = 0xc0000b0000    ← GIỐNG (array cũ)
        ```
    2. Ví dụ 2: Header mới – array MỚI (realloc)
        ```go
        a := make([]int, 2, 2)
        a[0], a[1] = 1, 2

        b := append(a, 3)

        fmt.Println("a:", a, ", len: ", len(a), ", cap: ", cap(a))
        fmt.Println("b:", b, ", len: ", len(b), ", cap: ", cap(b))

        fmt.Printf("&a = %p\n", &a)     
        fmt.Printf("&b = %p\n", &b)

        fmt.Printf("&a[0] = %p\n", &a[0])
        fmt.Printf("&b[0] = %p\n", &b[0])
        ```

        - OUTPUT:
        ```bash
        a: [1 2] , len:  2 , cap:  2
        b: [1 2 3] , len:  3 , cap:  4
        &a = 0xc000126030       
        &b = 0xc000126048       ← KHÁC (header mới)
        &a[0] = 0xc00011c020
        &b[0] = 0xc00012c000    ← KHÁC (array mới)
        ```

    3. Ví dụ 3: Header mới – array CŨ nhưng làm “đổi” dữ liệu gốc
        ```go
        a := make([]int, 2, 3)
        a[0], a[1] = 1, 2

        b := append(a, 99)

        fmt.Printf("&a[0] = %p\n", &a[0])
        fmt.Printf("&b[0] = %p\n", &b[0])

        fmt.Println("a:", a)
        fmt.Println("b:", b)
        ```

        - OUTPUT:
        ```bash
        &a[0] = 0xc0001a8000
        &b[0] = 0xc0001a8000
        a: [1 2]
        b: [1 2 99]
        ```

        👉 Nhưng underlying array bị thay đổi  =>  a không thấy, nhưng bị ảnh hưởng ngầm
        ```text
        [1 2 99]
        ```
    4. Ví dụ 4: append mà KHÔNG gán lại → slice cũ không đổi
        ```go
        func f(s []int) {
            append(s, 100)
        }

        func main() {
            a := []int{1, 2}    // a là slice, không phải array
            f(a)
            fmt.Println(a)
        }
        ```

        - OUTPUT:
        ```bash
        [1 2]
        ```

### 5. Cách hoạt động của `make`, `copy`, `Full Slice Expression`
    1. Hàm `make`
        - Tác dụng: 
            + `make` dùng để khởi tạo một Slice và cấp phát bộ nhớ cho mảng ẩn ngay từ đầu
            + Điều này giúp tối ưu hiệu năng bằng cách tránh việc Go phải cấp phát lại bộ nhớ nhiều lần khi bạn `append`
        - Cách sử dụng
            ```go
            s := make([]T, length, capacity)
            ```
        - Cơ chế:
            + Go cấp phát một mảng ẩn có kích thước capacity. Sau đó tạo một Slice Header với `len = length` và `cap = capacity`. Các phần tử từ 0 đến length-1 sẽ được khởi tạo giá trị mặc định `(zero value)`.
        - Ví dụ:
            ```go
            s := make([]int, 3, 5) // [0 0 0] - Mảng ẩn thực tế có 5 chỗ, nhưng bạn chỉ thấy 3 số 0.
            ```
    2. Hàm `copy`
        - Tác dụng: 
            + Dùng để sao chép dữ liệu từ một slice này sang một slice khác
            + Đây là cách duy nhất để tách rời hai slice khỏi cùng một mảng ẩn, giúp tránh lỗi ghi đè dữ liệu hoặc rò rỉ bộ nhớ
        - Cách sử dụng
            ```go
            numCopied := copy(dest, src)
            ```
        - Cơ chế:
            + Go sao chép các phần tử từ `src` sang `dest`. 
            + Số lượng phần tử được copy sẽ là số nhỏ nhất giữa `len(src)` và `len(dest)`
        - Lưu ý:
            + `dest` phải được khởi tạo trước (bằng `make` hoặc `literal`) thì mới có chỗ để chứa dữ liệu

        - Ví dụ:
            ```go
            src := []int{1, 2, 3}
            dest := make([]int, len(src)) // Phải tạo dest có cùng len
            copy(dest, src)               // dest bây giờ là [1 2 3], nằm ở mảng ẩn khác hoàn toàn
            ```
    3. `Full Slice Expression`
        - Đây là kỹ thuật "nâng cao" nhất để kiểm soát slice.
        - Tác dụng:
            + Giới hạn capacity của một slice con.
            + Điều này ngăn chặn việc append vào slice con làm hỏng (ghi đè) dữ liệu của slice cha hoặc các slice con khác.
        - Cách sử dụng:
            ```go
            slice_con := slice_cha[start : end : capacity]
            ```

            ```text
            - Trong đó:
            start: Chỉ số bắt đầu.
            end: Chỉ số kết thúc
            capacity: Giá trị capacity cố định
            ```
        - Ví dụ:
            1. Trường hợp 1: Dùng cú pháp thường (Nguy hiểm)
            
            ```go
            parent := []int{1, 2, 3, 4, 5}
            child := parent[0:2]    // len=2, cap=5 (vẫn nhìn thấy 3, 4, 5 ở phía sau)
            child = append(child, 99) 
            fmt.Println(parent)     // [1 2 99 4 5] -> parent bị mất số 3!
            ```

            2. Trường hợp 2: Dùng Full Slice Expression (An toàn)
            
            ```go
            parent := []int{1, 2, 3, 4, 5}
            child := parent[0:2:2]  // len=2, cap=2 (Khóa chặt, không cho nhìn thấy 3, 4, 5)
            child = append(child, 99) 
            // Vì cap đã đầy (2/2), Go buộc phải tạo mảng ẩn mới cho child
            fmt.Println(parent)     // [1 2 3 4 5] -> parent được bảo vệ an toàn!
            fmt.Println(child)      // [1 2 99]
            ```
    4. Tổng kết
        - `make`:
            + Dùng khi bắt đầu tạo slice và biết trước kích thước.
            + Mục tiêu chính: Tối ưu hiệu năng (tránh re-allocation).
        - `copy`:
            + Dùng khi muốn nhân bản dữ liệu ra một vùng nhớ độc lập.
            + Mục tiêu chính: An toàn dữ liệu, tránh Memory Leak.
        - `Full Slice Expression`:
            + Dùng khi chia nhỏ slice lớn thành các slice nhỏ để xử lý riêng.
            + Mục tiêu chính: Cô lập vùng nhớ, ngăn chặn ghi đè ngoài ý muốn.

### 6. `Nill`, `Nil slice`
- The zero value of a slice is `nil`.
- A nil slice has a length and capacity of 0 and has no underlying array.

### 7. `Map`
### 8. `Function`, `Function closures`
    - Function Closure (hàm đóng) là một giá trị hàm (function value) mà nó tham chiếu đến các biến nằm bên ngoài phạm vi (body) của chính nó.
    - Cách thức hoạt động của `Closure`:
        + Thông thường, khi một hàm kết thúc, các biến cục bộ của nó sẽ bị xóa khỏi bộ nhớ. 
        + Tuy nhiên, với Closure, Go sẽ giữ lại các biến đó nếu chúng vẫn đang được tham chiếu bởi một hàm ẩn danh (anonymous function).
        + Mỗi closure có state riêng
            ```go
            package main

            import "fmt"

            func intSeq() func() int {
                i := 0
                return func() int {
                    i++
                    return i
                }
            }

            func main() {
                // nextInt là một closure, nó "nhớ" biến i
                nextInt := intSeq()

                fmt.Println(nextInt()) // Kết quả: 1
                fmt.Println(nextInt()) // Kết quả: 2
                fmt.Println(nextInt()) // Kết quả: 3

                // Nếu ta tạo một bộ đếm mới, nó sẽ có biến i riêng biệt
                newInts := intSeq()
                fmt.Println(newInts()) // Kết quả: 1
            }
            ```
    - Tại sao Closure lại quan trọng?
        + `Duy trì trạng thái (State)`: Như ví dụ trên, closure giúp hàm "nhớ" dữ liệu mà không cần sử dụng biến toàn cục (global variables) — điều vốn dễ gây lỗi trong các ứng dụng lớn hoặc đa luồng.
        + `Data Isolation (Cô lập dữ liệu)`: Biến i trong ví dụ trên hoàn toàn bị ẩn đi. Không ai có thể thay đổi i ngoại trừ chính hàm closure đó.
        + `Middleware & Decorators`: Trong lập trình Web (như Gin hoặc Echo), closure thường được dùng để bao bọc các handler nhằm kiểm tra quyền truy cập hoặc ghi log.

### 10. Method
- Go does not have classes. However, you can define methods on types.
- A method is a function with a special receiver argument.
- Remember: a method is just a function with a receiver argument.
- You can declare a method on non-struct types, too.
- You can only declare a method with a receiver whose type is defined in the same package as the method
- You cannot declare a method with a receiver whose type is defined in another package (which includes the built-in types such as int).
- Syntax:
    ```go
    func (receiver Type) MethodName(params) returnType {
        // code
    }
    ```
    
- Trong Go, Method (phương thức) thực chất là một hàm nhưng có một tham số đặc biệt đứng trước tên hàm, được gọi là Receiver (đối tượng nhận).
- Nó cho phép bạn định nghĩa các hành vi cho một struct hoặc bất kỳ kiểu dữ liệu nào bạn tự định nghĩa, tạo cảm giác giống như lập trình hướng đối tượng (OOP).
- Ví dụ:
    ```go
    package main

    import "fmt"

    type Rect struct {
        width, height int
    }

    // the Area method has a receiver of type Rect named r
    func (r Rect) Area() int {
        return r.width * r.height
    }

    func main() {
        r := Rect{width: 10, height: 5}
        fmt.Println("Diện tích:", r.Area()) // Gọi method giống như trong OOP
    }
    ```

- `Pointer receiver`
    + You can declare methods with pointer receivers.
    + This means the receiver type has the literal syntax *T for some type T. (Also, T cannot itself be a pointer such as *int.)
    + Methods with pointer receivers can modify the value to which the receiver points
    + There are two reasons to use a pointer receiver:
        1. The first is so that the method can modify the value that its receiver points to.
        2. The second is to avoid copying the value on each method call. This can be more efficient if the receiver is a large struct, for example.

### 11. `Interface`, `nil interface value`, `empty interface` 
- An interface type is defined as a set of method signatures
- A value of interface type can hold any value that implements those methods

    1. Interface là gì?
    - Interface là một tập hợp các chữ ký phương thức (`method signatures`). 
    - Một kiểu dữ liệu (thường là struct) được coi là "triển khai" (implement) một interface nếu nó định nghĩa tất cả các phương thức mà interface đó yêu cầu.
    - Đặc điểm đặc biệt nhất: Trong Go, bạn không cần dùng từ khóa `implements`. Việc triển khai là ngầm định (`implicit`).
    - `empty interface`:
        +  `any` là bí danh cho `interface{}`. 
        + Một `interface rỗng` là interface không có phương thức nào, điều đó có nghĩa là mọi kiểu dữ liệu đều thỏa mãn nó.
        + Một `interface rỗng` có thể chứa bất cứ `value` thuộc `concrete type` nào

    - Chú ý: 
        1. Một `biến interface` luôn gồm 2 phần:
        `(interface value) = (concrete type, value)` (concrete type là kiểu cụ thể)
        - Ví dụ:
            ```go
            type MyFloat float64
            type Abser interface {
                Abs() float64
            }
            type Vertex struct {
                X, Y float64
            }

            func (v *Vertex) Abs() float64 {
                return math.Sqrt(v.X*v.X + v.Y*v.Y)
            }
            
            func main() {
                var a Abser     // Lúc này: a = (nil, nil)
                a = MyFloat(-3) // Lúc này: a = (MyFloat, -3)
                v := Vertex{3, 4}
                a = &v          // Lúc này: a = (*Vertex, &{3,4})
            }
            ```
        2. Interface chứa nil ≠ interface là nil
        - Ví dụ:
            ```go
            var a Abser
            fmt.Println(a == nil)   // true. Vì a = (nil, nil)

            var v *Vertex = nil
            a = v
            fmt.Println(a == nil)   // false. Vì: a = (*Vertex, nil)
            ```
        3. Interface là nil đang là nil thì nếu gọi method trong interface chương trình sẽ bị panic (runtime error) bởi vì không biết concrete type bên trong interface để biết gọi method nào

    - Ví dụ 1:
        ```go
        package main

        import "fmt"

        // Định nghĩa Interface Speaker
        type Speaker interface {
            Speak() string
        }

        // Walking là interface rỗng
        type Walking interface {
        }

        // Hearing là interface rỗng
        type Hearing any

        // Struct Dog
        type Dog struct{}

        func (d Dog) Speak() string {
            return "Woof!"
        }

        // Struct Cat
        type Cat struct{}

        func (c Cat) Speak() string {
            return "Meow!"
        }

        func makeItSpeak(s Speaker) {
            fmt.Println(s.Speak())
        }

        func makeItWalking(w Walking) {
            fmt.Println(w)
        }

        func makeItHearing(h Hearing) {
            fmt.Println(h)
        }

        func main() {
            d := Dog{}
            c := Cat{}

            makeItSpeak(d) // Woof! (Chạy tốt vì Dog có method Speak() )
            makeItSpeak(c) // Meow! (Chạy tốt vì Cat có method Speak() )

            makeItWalking(d) // {}
            makeItWalking(c) // {}

            makeItHearing(d) // {}
            makeItHearing(c) // {}
        }
        ```

    - Ví dụ 2: 
        ```go
        package main

        import (
            "fmt"
            "math"
        )

        type Abser interface {
            Abs() float64
        }

        func main() {
            var a Abser                     // a là biến kiểu interface Abser.
            f := MyFloat(-math.Sqrt2)       // f là biến kiểu MyFloat
            v := Vertex{3, 4}               // v là biến kiểu Vertex

            /*
                a KHÔNG BAO GIỜ đổi kiểu (a luôn có kiểu interface Abser)
                - Khi gán a = f thì bên trong a chứa: (concrete type = MyFloat, value = f)
                - Khi gán a = &v thì bên trong a chứa: (concrete type = *Vertex, value = &v)
            */
            a = f  // a MyFloat implements Abser
            a = &v // a *Vertex implements Abser


            /*
                a = v lỗi vì:
                + v có kiểu Vertex
                + Vertex không có method Abs()
                + Chỉ *Vertex mới có method Abs()
            */
            // In the following line, v is a Vertex (not *Vertex) and does NOT implement Abser.
            a = v   // Error here

            fmt.Println(a.Abs())
        }

        type MyFloat float64

        func (f MyFloat) Abs() float64 {
            if f < 0 {
                return float64(-f)
            }
            return float64(f)
        }

        type Vertex struct {
            X, Y float64
        }

        func (v *Vertex) Abs() float64 {
            return math.Sqrt(v.X*v.X + v.Y*v.Y)
        }
        ```

### 12. `Type Assertion`
- `Type Assertion` (khẳng định kiểu) là cách để bạn lấy lại giá trị với kiểu dữ liệu cụ thể từ một biến đang được giữ dưới dạng interface.
- Vì `interface` có thể chứa bất kỳ giá trị nào, nên `Type Assertion` đóng vai trò như một lời xác nhận: "Tôi tin rằng cái interface này đang chứa kiểu dữ liệu X, hãy lấy nó ra cho tôi".

    1. Cú pháp:
        `v = i.(T)`
        - Trong đó:
            + i: Là một biến kiểu interface
            + T: Là kiểu dữ liệu mà bạn muốn khẳng định (ví dụ: int, string, hoặc một struct)
            + v: Biến mới sẽ giữ giá trị của i nhưng với kiểu T

    2. Hai cách sử dụng Type Assertion
        - Cách 1: Phản ứng dữ dội (Panic)
            + Nếu bạn khẳng định sai kiểu dữ liệu mà không kiểm tra, chương trình sẽ bị panic (dừng đột ngột).
                ```go
                var i interface{} = "Hello"

                s := i.(string) // Thành công, s là "Hello" (kiểu string)
                f := i.(float64) // PANIC: interface conversion: interface {} is string, not float64
                ```
        - Cách 2: Kiểm tra an toàn (Comma-ok)
            + Đây là cách được khuyến khích sử dụng trong thực tế. 
            + Nó trả về thêm một giá trị boolean (ok) để báo cho bạn biết việc khẳng định có thành công hay không.
                ```go
                var i interface{} = "Hello"

                s, ok := i.(string)
                if ok {
                    fmt.Println("Giá trị chuỗi là:", s)
                } else {
                    fmt.Println("Không phải là string!")
                }
                ```

- Docs:
    ```text
    - A type assertion provides access to an interface value's underlying concrete value.
        t := i.(T)

    - This statement asserts that the interface value i holds the concrete type T and assigns the underlying T value to the variable t.

    - If i does not hold a T, the statement will trigger a panic.

    - To test whether an interface value holds a specific type, a type assertion can return two values: the underlying value and a boolean value that reports whether the assertion succeeded.
        t, ok := i.(T)
    
    - If i holds a T, then t will be the underlying value and ok will be true.
    - If not, ok will be false and t will be the zero value of type T, and no panic occurs.
    - Note the similarity between this syntax and that of reading from a map.
    ```

### 13. `Type Switch`
- `Type Switch` là một cấu trúc điều khiển trong Go cho phép bạn so sánh kiểu dữ liệu của một interface với nhiều kiểu dữ liệu khác nhau trong một khối lệnh duy nhất.
- Nó giống như một lệnh switch thông thường, nhưng thay vì so sánh giá trị (ví dụ: x == 5), nó so sánh loại của dữ liệu (ví dụ: x có phải là string không?).

- Ví dụ 1:
    ```go
    package main

    import "fmt"

    func do(i interface{}) {
        switch v := i.(type) {
        case int, int32, int64:
            fmt.Printf("Twice %v is %v\n", v, v*2)
        case string:
            fmt.Printf("%q is %v bytes long\n", v, len(v))
        default:
            fmt.Printf("I don't know about type %T!\n", v)
        }
    }

    func main() {
        do(21)
        do("hello")
        do(true)
    }

    ```

### 14. `Stringer`
- `Stringer` là một trong những `interface` phổ biến và hữu ích nhất. 
- Nó được định nghĩa trong `package fmt` và cho phép bạn tự quyết định cách một đối tượng (struct) hiển thị khi được in ra dưới dạng chuỗi.
    1. Định nghĩa Interface Stringer
    - Interface Stringer cực kỳ đơn giản, nó chỉ yêu cầu một phương thức duy nhất:
        ```go
        type Stringer interface {
            String() string
        }
        ```
    - A Stringer is a type that can describe itself as a string. 
    - The fmt package (and many others) look for this interface to print values.

### 15. `Error`
- Go programs express error state with `error` values.
- The `error` type is a built-in interface similar to `fmt.Stringer`:
    ```go
    type error interface {
        Error() string
    }
    ```
- (As with `fmt.Stringer`, the fmt package looks for the `error` interface when printing values.)

- Functions often return an `error` value, and calling code should handle errors by testing whether the error equals `nil`.
    ```go
    i, err := strconv.Atoi("42")
    if err != nil {
        fmt.Printf("couldn't convert number: %v\n", err)
        return
    }
    fmt.Println("Converted integer:", i)
    ```

- A nil error denotes success; a non-nil error denotes failure.

### 16. `Generic`

2. `Type parameter`
    - `Type Parameter` là một loại tham số đặc biệt được đặt trong dấu ngoặc vuông []. Nó đóng vai trò như một "chỗ trống" sẽ được lấp đầy bởi một kiểu dữ liệu cụ thể khi hàm hoặc cấu trúc dữ liệu được sử dụng.
    - Cấu trúc cơ bản:
        ```go
        func TenHam[T AnyType](thamSo T) T {
            // code
        }
        ```

        + Trong đó:
            - `T`: Là tên của Type Parameter (bạn có thể đặt tên bất kỳ, nhưng thường dùng T, V, K).

            - `AnyType`: Là Type Constraint (Ràng buộc kiểu) – giới hạn những kiểu dữ liệu nào T có thể nhận.

    - Ví dụ generic:
        ```go
        type List[T any] struct {
            elements []T
        }

        func main() {
            // Danh sách số nguyên
            intList := List[int]{elements: []int{1, 2, 3}}
            
            // Danh sách chuỗi
            stringList := List[string]{elements: []string{"A", "B"}}
        }
        ```
    
3. `Type Constraints (Ràng buộc kiểu)`
    - Để tránh việc người dùng truyền vào những kiểu dữ liệu không phù hợp (ví dụ: truyền một struct vào hàm yêu cầu thực hiện phép toán +), 
    - Go sử dụng Constraints:
        + any: Cho phép bất kỳ kiểu dữ liệu nào (tương đương interface{}).
        + comparable: Các kiểu có thể so sánh bằng == hoặc !=.
        + Custom Interface: Bạn có thể tự định nghĩa tập hợp các kiểu cho phép.
    - Ví dụ:
        ```go
        // Chỉ chấp nhận int hoặc int64 hoặc float64
        type Number interface {
            int | int64 | float64
        }

        func Sum[T Number](a, b T) T {
            return a + b
        }
        ```
    
4. `Type Approximation (~)`
    - Nó cho phép một Type Constraint không chỉ chấp nhận một kiểu dữ liệu chính xác, mà còn chấp nhận tất cả các kiểu dữ liệu có kiểu cơ sở (underlying type) là kiểu đó.
    - Dấu ~ nói với trình biên dịch rằng: "Hãy chấp nhận bất kỳ kiểu nào có kiểu cơ sở là..."
    ```go
    type Number interface {
        ~int | ~float64 // Chấp nhận int, float64 và mọi kiểu định nghĩa dựa trên chúng
    }

    type MyID int

    func PrintInt[T ~int](v T) {
        fmt.Println(v)
    }

    func main() {
        var x MyID = 100
        PrintInt(x) // Hợp lệ!
    }
    ```
    - Quy tắc và Hạn chế:
        + `Chỉ dùng với kiểu cơ sở (Underlying types)`: Bạn không thể dùng `~` với các kiểu phức hợp như `struct` hoặc `interface` một cách tùy tiện. Nó thường được dùng với các kiểu dữ liệu cơ bản như `int, string, float64, bool`,...
            + Hợp lệ: `~int, ~string, ~[]byte`
            + Không hợp lệ: `~MyStruct` (vì `MyStruct` không phải là kiểu cơ sở của ngôn ngữ).

        + `Không dùng với struct, Interface`: Bạn không thể viết `~error` vì `error` là một `interface`, không phải là một kiểu dữ liệu cơ sở.

        + `Vị trí sử dụng`: Dấu `~` chỉ có thể xuất hiện bên trong các `Interface` dùng làm constraint.

### 18. `Struct Embedding`
- `Struct Embedding` (nhúng struct) là một kỹ thuật cho phép bạn lồng một struct này vào trong một struct khác. 
- Đây là cách Go thực hiện việc tái sử dụng mã nguồn và chia sẻ hành vi giữa các đối tượng thay vì sử dụng cơ chế kế thừa (inheritance) truyền thống như trong Java hay C++.
    1. `Cách hoạt động của Struct Embedding`
    - Khi bạn nhúng một struct vào struct khác mà `không đặt tên trường`, struct được nhúng sẽ trở thành một `Anonymous Field` (trường ẩn danh).

    ```go
    type Person struct {
        Name string
        Age  int
    }

    func (p Person) Greet() {
        fmt.Printf("Hi, I'm %s\n", p.Name)
    }

    type Employee struct {
        Person // Đây là Struct Embedding
        ID     int
    }
    ```
    - Đặc điểm quan trọng:
        + `Promoted Fields`: Các trường của `Person (như Name, Age)` được "đưa lên" và có thể truy cập trực tiếp từ `Employee`.

        + `Promoted Methods`: Các phương thức của `Person (như Greet())` cũng được đưa lên và có thể gọi trực tiếp từ biến kiểu `Employee`.
        ```go
        e := Employee{
            Person: Person{Name: "An", Age: 30},
            ID:     101,
        }

        fmt.Println(e.Name) // Truy cập trực tiếp thay vì e.Person.Name
        e.Greet()           // Gọi trực tiếp phương thức của Person
        ```
    2. `Composition over Inheritance (Hợp thành thay vì Kế thừa)`
    - Go không có từ khóa `extends`, `super`. 
    - `Struct Embedding` dựa trên nguyên lý `Composition`. Mối quan hệ ở đây không phải là `"Is-a"` (là một) mà là `"Has-a"` (có một), nhưng với khả năng truy cập thuận tiện như `"Is-a"`.

    3. Ghi đè phương thức (Method Overriding)
    - Nếu struct bên ngoài có phương thức trùng tên với struct bên trong, phương thức của struct bên ngoài sẽ được ưu tiên. Đây gọi là `Shadowing`.

    ```Go
    func (e Employee) Greet() {
        fmt.Printf("I'm employee #%d\n", e.ID)
    }
    // Khi gọi e.Greet(), nó sẽ chạy hàm của Employee thay vì Person.
    ```
    4. `Embedding Interface`
    - Bạn không chỉ nhúng được `struct` mà còn có thể nhúng cả `Interface` vào `struct`. 
    - Điều này cực kỳ hữu ích khi bạn muốn một `struct` thỏa mãn một `interface` mà không cần cài đặt lại toàn bộ các phương thức của nó.

    ```Go
    type Job struct {
        io.Reader // Nhúng interface Reader
    }
    ```
    - Bất kỳ struct nào chứa `Job` giờ đây cũng mặc định có phương thức `Read()` từ `io.Reader`
    5. `Những lưu ý quan trọng (Best Practices)`
    - `Tránh nhầm lẫn với Kế thừa`: Dù `Employee` có các trường của `Person`, nhưng bạn không thể dùng `Employee` ở nơi yêu cầu kiểu `Person`. Chúng vẫn là hai kiểu dữ liệu khác nhau hoàn toàn.
    
    - `Tránh đặt tên trùng lặp (Ambiguity)`: Nếu bạn nhúng 2 struct khác nhau vào cùng 1 struct mà cả 2 struct đó đều có trường cùng tên (ví dụ đều có trường `ID`), bạn buộc phải truy cập qua tên struct đầy đủ (ví dụ `e.Person.ID`) để tránh lỗi biên dịch.

    - `Sử dụng khi thực sự cần thiết`: Đừng lạm dụng embedding chỉ để viết code ngắn hơn. Chỉ dùng khi struct bên ngoài thực sự mang bản chất hoặc hành vi mở rộng của struct bên trong.

### 19. `Defer`, `Panic`, `Recover` (https://go.dev/blog/defer-panic-and-recover)
- Trong Go, panic và recover là cơ chế xử lý lỗi đặc biệt dùng cho các tình huống nghiêm trọng (runtime error hoặc lỗi không thể tiếp tục). 
- Tuy nhiên, trong thực tế, Go khuyến khích dùng error thay vì panic cho các lỗi thông thường.
    1. `Panic`
    - `panic` là một hàm build-in dùng để dừng luồng thực thi bình thường của chương trình.
    - `panic` dùng để:
        + Dừng chương trình ngay lập tức   
        + In ra stack trace
        + Thực thi các hàm defer trước khi thoát
        + Sau đó thoát chương trình (nếu không được recover)
    - Ví dụ:
        ```go
        package main

        import "fmt"

        func main() {
            fmt.Println("Start")
            panic("Something went wrong!")
            fmt.Println("End") // không chạy
        }

        ```
        📌 Output:
        ```bash
        Start
        panic: Something went wrong!
        ```
    2. Khi nào `Panic` xảy ra:
    * `Panic` do runtime (tự động)
        - Ví dụ: 
            Chia cho 0, Truy cập index ngoài phạm vi slice, Dereference nil pointer
        ```go
        var a []int
        fmt.Println(a[1]) // panic: index out of range
        ```    
    * `Panic` do lập trình viên chủ động
        ```go
        if user == nil {
            panic("user is nil")
        }
        ```
    3. `Defer` hoạt động thế nào khi có `panic`?
    - Khi panic xảy ra:
        + Bước 1: Hàm hiện tại dừng lại
        + Bước 2: Các defer trong hàm đó được gọi theo thứ tự ngược lại (LIFO)
        + Bước 3: Sau đó panic tiếp tục lan lên stack    
        ```go
        func main() {
            defer fmt.Println("world")
            fmt.Println("hello")
            panic("oops")
        }
        ```
        📌 Output:
        ```go
        hello
        world
        panic: oops
        ```
    4. `Recover` là gì?
    - `recover` là hàm (built-in function) dùng để giành lại quyền kiểm soát của một goroutine đang bị panic
    - `recover()` dùng để:
        + Bắt panic
        + Ngăn chương trình bị crash
        + Chỉ hoạt động bên trong `defer`
    - Ví dụ dùng `recover`:
        ```go
        package main

        import "fmt"

        func main() {
            defer func() {
                if r := recover(); r != nil {
                    fmt.Println("Recovered from:", r)
                }
            }()

            panic("something bad happened")
            fmt.Println("This will not run")
        }
        ```   
        📌 Output:
        ```go
        Recovered from: something bad happened
        ```
    5. `Panic lan truyền (Stack Unwinding)`
    - Ví dụ:
        ```go
        func a() {
            panic("error in a")
        }

        func b() {
            a()
        }

        func main() {
            b()
        }
        ```
        📌 Hoạt động:
        ```go
        Panic sẽ lan từ a() → b() → main() → crash.
        Nếu main() có recover, panic sẽ được chặn lại.
        ```
    6. Khi nào nên dùng panic?
    * ✅ Nên dùng khi:
        - Lỗi nghiêm trọng không thể tiếp tục
        - Sai logic nghiêm trọng (bug)
        - Khởi tạo thất bại (ví dụ: config sai, không load được file quan trọng)
        - Trong package nội bộ (assert-like)

        - Ví dụ:
            ```go
            func MustConnect() *DB {
                db, err := Connect()
                if err != nil {
                    panic(err)
                }
                return db
            }
            ```

    * ❌ Không nên dùng khi:
        - Lỗi nghiệp vụ (business logic)
        - Lỗi có thể xử lý được
        - Lỗi từ input người dùng
        - Thay vào đó nên dùng error:
            ```go
            func divide(a, b int) (int, error) {
                if b == 0 {
                    return 0, errors.New("division by zero")
                }
                return a / b, nil
            }
            ```

    7. `Lưu ý quan trọng`
    - recover chỉ hoạt động nếu:
        + Được gọi trong hàm defer
        + Cùng goroutine với panic
    - recover() được gọi bên ngoài defer func không lỗi, nó chỉ trả về nil và không làm gì cả, không bắt được panic
    - Nếu panic xảy ra trong goroutine khác, bạn phải recover trong goroutine đó.
    - Ví dụ sai:
        ```go
        go func() {
            panic("boom")
        }()
        ```
        👉 Nếu không recover trong chính goroutine đó → chương trình crash.


### 20. Sâu hơn về `panic`
    
1. `Panic` là gì ở mức runtime?
- `panic` là cơ chế dừng bất thường (abrupt termination) của một goroutine.
- Khi `panic(x)` xảy ra:
    + Bước 1: Goroutine hiện tại dừng thực thi bình thường
    + Bước 2: Runtime bắt đầu `stack unwinding`
    + Bước 3: Các defer được gọi theo thứ tự LIFO
    + Bước 4: Nếu không có recover() → chương trình crash
    + Bước 5: Runtime in stack trace của goroutine bị panic

- Quan trọng:
    + Panic chỉ ảnh hưởng đến goroutine đang chạy.
    + Nhưng nếu đó là main goroutine và không recover → toàn bộ chương trình kết thúc. 

2. `Stack Unwinding` chi tiết
- Ví dụ:
    ```go
    func c() {
        panic("boom")
    }

    func b() {
        defer fmt.Println("defer in b")
        c()
    }

    func a() {
        defer fmt.Println("defer in a")
        b()
    }

    func main() {
        a()
    }

    ```
- Thứ tự xảy ra:
    + B1: panic trong c
    + B2: c không có defer → quay về b
    + B3: Chạy defer trong b
    + B4: Quay về a
    + B5: Chạy defer trong a
    + Không có recover → crash
- Output:
    ```go
    defer in b
    defer in a
    panic: boom
    ```
👉 Đây gọi là `stack unwinding`

3. Panic là per-goroutine
- Ví dụ:
    ```go
    go func() {
        panic("worker crashed")
    }()
    ```
- Nếu không recover trong goroutine đó → toàn bộ chương trình crash.
- Vì: Nếu bất kỳ goroutine nào panic mà không recover → runtime kết thúc chương trình.

### 21. Sâu hơn về `recover`
1. `recover` thực sự là gì?

- `recover()` là một built-in function dùng để:
    + Chặn (intercept) một panic
    + Dừng quá trình unwinding stack
    + Trả về giá trị được truyền vào panic(...)

- Signature:
    ```go
    func recover() interface{}
    ```
- Nếu không có panic đang diễn ra → recover() trả về nil.

2. Cơ chế hoạt động bên trong (stack unwinding)
- Khi panic xảy ra:
    + B1: Hàm hiện tại dừng thực thi
    + B2: Runtime bắt đầu stack unwinding
    + B3: Các defer được gọi theo LIFO
    + B4: Nếu trong một defer có gọi recover():
        + Panic bị chặn
        + Stack unwinding dừng lại
        + Hàm chứa defer đó tiếp tục chạy sau defer
    + B5: Nếu không có recover → chương trình crash

3. Điều kiện để recover hoạt động
- Recover chỉ có tác dụng nếu:
    + ✅ 1. Được gọi trong defer
    + ✅ 2. Được gọi trực tiếp trong hàm defer (không phải hàm lồng)
    + ✅ 3. Cùng goroutine với panic

4. Ví dụ:
- Ví dụ đúng:
    ```go
    func main() {
        defer func() {
            if r := recover(); r != nil {
                fmt.Println("Recovered:", r)
            }
        }()

        panic("boom")
    }
    ```

- Ví dụ sai – không gọi trực tiếp trong defer:
    ```go
    func handle() {
        recover() // vô dụng
    }

    func main() {
        defer handle()
        panic("boom")
    }
    ```
    👉 Không hoạt động. Vì `recover()` không được gọi trực tiếp trong function literal của defer.

- Cách đúng nếu tách hàm:
    ```go
    func handle() {
        if r := recover(); r != nil {
            fmt.Println("Recovered:", r)
        }
    }

    func main() {
        defer handle() // OK vì handle được gọi bởi defer
        panic("boom")
    }
    ```
    Điểm khác biệt: recover() phải được gọi trong call stack của deferred function đang chạy do panic.

- Ví dụ sai:
    ```go
    func main() {
        defer func() {
            if r := recover(); r != nil {
                fmt.Println("Recovered:", r)
            }
        }()

        go func() {
            panic("goroutine panic")
        }()

        time.Sleep(time.Second)
    }
    ```
    👉 Chương trình vẫn crash. Vì panic xảy ra trong goroutine khác.

- Ví dụ đúng:
    ```go
    func main() {
        defer func() {
            if r := recover(); r != nil {
                fmt.Println("Recovered:", r)
            }
        }()

        go func() {
            defer func() {
                if r := recover(); r != nil {
                    fmt.Println("Recovered in goroutine:", r)
                }
            }()

            panic("goroutine panic")
        }()

        time.Sleep(time.Second)
    }
    ```
    📌 Nguyên tắc quan trọng: Goroutine nào panic → goroutine đó phải recover.

5. Recover dừng panic như thế nào?
- Khi recover thành công:
    + Stack unwinding dừng lại
    + Hàm chứa defer tiếp tục thực thi
    + Các hàm bên trên không bị ảnh hưởng
- Ví dụ:
    ```go
    func test() {
        defer func() {
            if r := recover(); r != nil {
                fmt.Println("Recovered")
            }
        }()

        panic("boom")
        fmt.Println("After panic") // không chạy
    }

    func main() {
        test()
        fmt.Println("Program continues")
    }
    ```
    Output:
    ```yaml
    Recovered
    Program continues
    ```

6. Giá trị trả về của recover
- recover() trả về đúng giá trị truyền vào panic
- Ví dụ: panic("string error"), panic(errors.New("error object")), panic(123)

- Ví dụ:
    ```go
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("Type: %T, Value: %v\n", r, r)
        }
    }()
    panic(errors.New("db failed"))
    ```
    Output:
    ```yaml
    Type: *errors.errorString, Value: db failed
    ```

### 22. Garbage Collector (GC) KHÔNG dọn dẹp Stack
1. Ai dọn dẹp Stack và dọn khi nào?

- Garbage Collector (GC) KHÔNG dọn dẹp Stack
- Stack được dọn dẹp bởi CPU (thông qua các chỉ thị được trình biên dịch tạo ra), và nó diễn ra ngay lập tức khi hàm kết thúc.

- Cơ chế: Mỗi khi một hàm được gọi, một khối bộ nhớ (gọi là Stack Frame) được "đẩy" vào Stack. Khi hàm đó chạy xong (return), toàn bộ Stack Frame đó bị "loại bỏ".

- Thao tác: Việc "dọn dẹp" này thực chất chỉ là thay đổi giá trị của một thanh ghi CPU gọi là Stack Pointer (SP). CPU chỉ cần di chuyển con trỏ này lên hoặc xuống. Nó không cần đi tìm từng biến để xóa, nó chỉ đơn giản là đánh dấu: "Vùng nhớ này bây giờ là trống, ai muốn ghi đè lên thì ghi".

2. Sự khác biệt giữa Dọn dẹp Stack và Dọn dẹp Heap

| Đặc điểm       | Stack (Ngăn xếp)                              | Heap (Đống)                                                  |
|---------------|-----------------------------------------------|--------------------------------------------------------------|
| **Người dọn** | CPU (Tự động theo cấu trúc hàm)               | Garbage Collector (Một chương trình chạy ngầm)               |
| **Thời điểm** | Ngay lập tức khi hàm return                   | Định kỳ (Khi bộ nhớ đầy hoặc theo thuật toán của Go)         |
| **Chi phí**   | Gần như bằng 0 (chỉ là 1 lệnh CPU)            | Rất đắt (phải quét toàn bộ bộ nhớ, tìm con trỏ...)           |
| **Cách dọn**  | Di chuyển Stack Pointer (ghi đè)              | Đánh dấu và thu hồi (Mark and Sweep)                         |

### 23. Garbage Collector (GC) thu dọn Heap như thế nào

- Go sử dụng thuật toán gọi là Mark and Sweep (Đánh dấu và Quét) với cơ chế Concurrent Collector (chạy song song với chương trình).

1. GC thu dọn Heap như thế nào? (Cơ chế Mark & Sweep)
- Quy trình này gồm 3 giai đoạn chính, được thiết kế để không làm "đứng hình" (Stop The World) chương trình quá lâu:

- `Giai đoạn 1`: Mark Preparation (Chuẩn bị) - Stop The World (STW)
    + Go tạm dừng tất cả Goroutines trong một khoảng thời gian cực ngắn (thường là vài micro giây) để bật Write Barrier (một cơ chế theo dõi xem có con trỏ nào thay đổi trong lúc GC đang chạy không).

- `Giai đoạn 2`: Marking (Đánh dấu) - Concurrent
    + Đây là giai đoạn tốn sức nhất. GC sẽ đi theo các con trỏ từ "Root" (biến toàn cục, các biến trên Stack) để tìm xem những vật thể nào trên Heap còn đang được sử dụng.

    + Go sử dụng mô hình Tricolor Marking (Đánh dấu 3 màu: Trắng, Xám, Đen).
        + Màu Trắng: Các đối tượng có thể là rác (chưa được quét tới).
        + Màu Xám: Các đối tượng đang được quét, nhưng các con trỏ bên trong nó chưa được kiểm tra (các đối tượng mà con trỏ trỏ tới chưa được kiểm tra.)
        + Màu Đen: Các đối tượng chắc chắn không phải là rác (chắc chắn đang được sử dụng) (đã quét xong cả nó và các con trỏ nó trỏ tới).

- `Giai đoạn 3`: Sweeping (Quét) - Concurrent
    + Sau khi đánh dấu xong, những vật thể nào vẫn là Màu Trắng chính là rác. GC sẽ giải phóng vùng nhớ này để sẵn sàng cho các lần cấp phát sau. Quá trình này diễn ra âm thầm khi chương trình của bạn vẫn đang chạy.

2. Khi nào GC tiến hành thu dọn?
- Go không dọn dẹp theo một khung giờ cố định (như "cứ mỗi 5 phút"). Thay vào đó, nó dựa trên 3 tín hiệu chính (Có 3 kịch bản chính khiến GC bắt đầu làm việc):

    + `Tín hiệu 1`: Dựa trên bộ nhớ tăng thêm (GOGC - Quy tắc chính)
        + Mặc định, Go có một biến môi trường là `GOGC=100`. Điều này có nghĩa là:
            + Nếu bộ nhớ Heap hiện tại sau lần dọn trước là 10MB.
            + Go sẽ đợi cho đến khi bộ nhớ Heap tăng lên thêm 100% nữa (tức là chạm ngưỡng 20MB) thì nó sẽ kích hoạt GC lần tiếp theo.
        + Bạn có thể chỉnh `GOGC` lên cao hơn để tiết kiệm CPU (nhưng tốn RAM) hoặc thấp hơn để tiết kiệm RAM (nhưng tốn CPU).

    + `Tín hiệu 2`: Dựa trên thời gian (2 phút)
        + Nếu trong vòng 2 phút mà bộ nhớ vẫn chưa tăng đến ngưỡng để kích hoạt GC, Go vẫn sẽ tự động chạy một lần dọn dẹp để đảm bảo tài nguyên không bị chiếm dụng quá lâu.

    + `Tín hiệu 3`: Kích hoạt thủ công
        + Bạn có thể ép Go dọn dẹp bằng cách gọi lệnh `runtime.GC()` trong code, nhưng điều này cực kỳ không khuyến khích trừ khi bạn đang làm kiểm thử hiệu năng.

3. Thay đổi biến môi trường `GOGC` ảnh hưởng như thế nào đến `CPU, RAM`

- `TH1`: Khi bạn đặt GOGC cao (Ví dụ: GOGC=200)
    + Ý nghĩa: Go sẽ đợi bộ nhớ Heap tăng thêm 200% so với mức sau khi dọn lần trước mới chạy GC tiếp.

    + Tại sao tốn RAM? 
        + Nếu bộ nhớ sau khi dọn là 10MB, nó sẽ đợi đến khi chạm 30MB mới dọn. Thay vì dọn ở mức 20MB (như mặc định 100), bạn cho phép nó "bày bừa" ra thêm 10MB nữa. Do đó, ứng dụng chiếm nhiều RAM hơn.

    + Tại sao tiết kiệm CPU?
        + Vì ngưỡng dọn dẹp cao hơn, nên khoảng thời gian giữa hai lần chạy GC sẽ dài ra. Thay vì cứ 5 phút phải đi dọn rác một lần, bây giờ 10 phút bạn mới dọn một lần. CPU không phải tốn chu kỳ để đi "đánh dấu" và "quét" rác thường xuyên, nên nó rảnh tay để chạy logic nghiệp vụ của bạn hơn.

- `TH2`: Khi bạn đặt GOGC thấp (Ví dụ: GOGC=50)
    + Ý nghĩa: Go chỉ đợi bộ nhớ tăng thêm 50% là đã lo đi dọn dẹp rồi.

    + Tại sao tiết kiệm RAM?
        + Vừa mới bày ra một chút (từ 10MB lên 15MB) là nhân viên vệ sinh đã vào quét dọn ngay. Bộ nhớ Heap luôn được giữ ở mức thấp nhất có thể.

    + Tại sao tốn CPU?
        + Nhân viên vệ sinh phải làm việc liên tục. Tần suất GC chạy sẽ rất dày đặc. Việc "đánh dấu" (Marking) hàng triệu vật thể trên Heap cực kỳ tốn CPU. CPU của bạn sẽ bị chia sẻ đáng kể cho bộ dọn rác thay vì chạy code của bạn.

4. Vậy khi nào nên chỉnh `GOGC`?

- `Kịch bản A (Server cấu hình mạnh, nhiều RAM)`: Bạn có một con server 64GB RAM nhưng ứng dụng chỉ dùng 4GB. Bạn nên tăng GOGC=200 hoặc 300. Tại sao phải để RAM trống phí phạm trong khi có thể giảm tải cho CPU để nó xử lý nhiều request hơn?

- `Kịch bản B (Server yếu, ít RAM)`: Bạn chạy Go trên một thiết bị nhúng hoặc container chỉ có 512MB RAM. Bạn nên giảm GOGC=50 để đảm bảo ứng dụng không bị hệ điều hành "giết" (OOM - Out of Memory) do chiếm quá nhiều RAM.

5. Một thông số mới: `GOMEMLIMIT` (Từ Go 1.19)

- Trước đây, chỉ có GOGC nên rất khó kiểm soát. Nếu bộ nhớ sau khi dọn là 1.1GB, GOGC=100 sẽ đợi đến 2.2GB. Nhưng nếu server (Docker Container) của bạn chỉ có 2GB RAM thì sao? Crash!

- Vì vậy, Go đã ra thêm GOMEMLIMIT (ví dụ: GOMEMLIMIT=1.8GiB cho docker container 2GB):
    + Go vẫn theo dõi GOGC để dọn dẹp định kỳ.
    + Nếu bộ nhớ tăng nhanh và sắp chạm mức GOMEMLIMIT=1.8GB, Go sẽ tự động kích hoạt GC ngay lập tức, bất chấp việc GOGC đã cho phép hay chưa => để cứu ứng dụng khỏi bị chết do hết RAM.

=> `GOMEMLIMIT Ngăn chặn lỗi OOM (Out of Memory)`.

### 24. Lập trình đồng thời (Concurrency) và song song (Parallelism)
- Cách tốt nhất để hiểu là thông qua sự khác biệt về cấu trúc và thực thi:
    + `Đồng thời (Concurrency)`: Là khả năng xử lý nhiều việc cùng một lúc. Nó giống như một đầu bếp đang nấu súp: họ bật bếp, trong lúc đợi nước sôi thì quay sang thái hành, sau đó lại quay lại kiểm tra nồi súp. Đầu bếp chỉ có một đôi tay, nhưng họ đang điều phối nhiều công việc xen kẽ nhau.

    + `Song song (Parallelism)`: Là khả năng thực hiện nhiều việc cùng một lúc. Nó giống như có hai đầu bếp: một người chuyên nấu súp và một người chuyên thái hành. Cả hai việc diễn ra tại cùng một thời điểm vật lý.

1. Lập trình đồng thời (Concurrent Programming)
- Định nghĩa
    + Lập trình đồng thời là kỹ thuật cho phép nhiều tác vụ (tasks) được xử lý chồng lấn về mặt thời gian, nhưng không nhất thiết phải chạy cùng lúc.
    + Trên hệ thống 1 CPU, hệ điều hành sẽ chuyển đổi qua lại rất nhanh giữa các tác vụ → tạo cảm giác chúng chạy đồng thời.

    👉 Mục tiêu chính:
    + Tăng khả năng phản hồi (responsive)
    + Xử lý I/O hiệu quả
    + Quản lý nhiều tác vụ cùng lúc

    📌 Ví dụ
    + Server web xử lý nhiều request cùng lúc.
    + Ứng dụng chat: vừa gửi tin nhắn, vừa nhận tin nhắn, vừa tải file.

2. Lập trình song song (Parallel Programming)
- Định nghĩa
    + Lập trình song song là kỹ thuật cho phép nhiều tác vụ chạy thực sự cùng lúc trên:
        + Nhiều CPU
        + Nhiều core
        + GPU

    👉 Mục tiêu chính:
    + Tăng tốc độ tính toán
    + Giải quyết bài toán lớn (AI, Big Data, mô phỏng)

    📌 Ví dụ
    + Chia mảng lớn thành nhiều phần và tính tổng cùng lúc.
    + Huấn luyện mô hình Deep Learning trên nhiều GPU.

3. Kết hợp đồng thời và song song
✅ Định nghĩa
- Là mô hình kết hợp cả:
    + Đồng thời → quản lý nhiều tác vụ
    + Song song → thực thi nhiều tác vụ cùng lúc

### 25. Goroutine
1. Goroutine là gì?
- Goroutine là một “luồng thực thi nhẹ” (lightweight thread) được quản lý bởi Go runtime, chứ không phải bởi hệ điều hành (OS)
- Nó cho phép bạn chạy nhiều tác vụ đồng thời (concurrency) một cách hiệu quả mà không tốn nhiều tài nguyên như thread truyền thống của hệ điều hành.

2. So sánh Goroutine với OS Thread

| Đặc điểm              | OS Thread                              | Goroutine                                    |
|-----------------------|----------------------------------------|----------------------------------------------|
| **Kích thước bộ nhớ** | Khoảng 1MB - 2MB                       | Chỉ từ 2KB                                   |
| **Khởi tạo/Hủy**      | Tốn kém tài nguyên hệ thống            | Rất nhanh và rẻ                              |
| **Context Switch**    | Chậm (do phải can thiệp vào Kernel)   | Rất nhanh (do Go Scheduler quản lý)          |
| **Số lượng**          | Giới hạn (vài nghìn là hệ thống đuối)  | Có thể chạy hàng triệu cùng lúc              |

3. Cơ chế hoạt động: `Mô hình M:N Scheduler`
- Trong các ngôn ngữ cũ, 1 luồng ứng dụng thường là 1 luồng hệ điều hành (1:1). Go thì khác, nó dùng mô hình M:N (nhiều Goroutine chạy trên ít luồng hệ điều hành).

- Go sử dụng một bộ điều phối (Scheduler) thông minh để ánh xạ `M` Goroutines vào `N` OS Threads.

- `Go Scheduler` điều phối dựa trên 3 thực thể chính:
    + `G (Goroutine)`: Đại diện cho một Goroutine.
        + Là đơn vị nhỏ nhất, chứa stack và con trỏ lệnh. 
        + Nó không tự chạy được mà cần được gán vào một P.

    + `M (Machine)`: Đại diện cho một OS Thread.
        + Là luồng thật sự của hệ điều hành. 
        + Để chạy mã Go, một M phải gắn với một P.

    + `P (Processor)`: Đại diện cho tài nguyên thực thi (ngữ cảnh), đóng vai trò trung gian điều phối G chạy trên M.
        + Đại diện cho ngữ cảnh thực thi (Context).
        + Số lượng P thường bằng số lõi CPU của máy bạn (có thể chỉnh qua biến `GOMAXPROCS`). 
        + P giữ một hàng đợi (runqueue) các Goroutine đang chờ được chạy.

- Cách thức vận hành: Chiến lược `Work Stealing`
    + Điểm thông minh nhất của Go Scheduler là nó không để bất kỳ CPU nào được nghỉ ngơi nếu vẫn còn việc.
        1. `Hàng đợi cục bộ (Local Runqueue)`: Mỗi P có một danh sách các Goroutine riêng. M sẽ lấy G từ P đang gắn với nó để xử lý.

        2. `Lấy trộm công việc (Work Stealing)`: Nếu một P đã xử lý hết sạch Goroutine trong hàng đợi của mình, nó sẽ nhìn sang các P khác. Nếu thấy P hàng xóm đang quá tải, nó sẽ "trộm" một nửa số Goroutine từ hàng đợi của hàng xóm về cho mình chạy.

        3. `Hàng đợi toàn cục (Global Runqueue)`: Nếu không trộm được từ hàng xóm, nó mới tìm đến hàng đợi chung của toàn hệ thống.

4. Giao tiếp giữa các Goroutine (Channels)
- Một triết lý nổi tiếng trong Go là: 
    ```go
    `Don't communicate by sharing memory; share memory by communicating.`
    (Đừng giao tiếp bằng cách dùng chung bộ nhớ; hãy chia sẻ bộ nhớ bằng cách giao tiếp.)
    ```

### 26. Channel
1. Channel là gì?
- Channel giống như một “đường ống” để truyền dữ liệu giữa các goroutine.
- Trong Golang, channel là một cơ chế dùng để giao tiếp và đồng bộ dữ liệu giữa các goroutine (luồng nhẹ trong Go). 
- Nó cho phép các goroutine gửi và nhận dữ liệu một cách an toàn mà không cần dùng lock thủ công như mutex.
- Channel hoạt động theo nguyên tắc `FIFO (First In, First Out)` — nghĩa là "Vào trước, Ra trước"

2. Cách khai báo channel
    ```go
    var ch chan int        // khai báo
    ch = make(chan int)    // khởi tạo

    // hoặc viết gọn
    ch := make(chan int)
    ```
    👉 Channel trên dùng để truyền dữ liệu kiểu int.

3. Gửi và nhận dữ liệu

    - Gửi dữ liệu vào channel:
        ```go
        ch <- 10
        ```
    - Nhận dữ liệu từ channel: 
        ```go
        value := <-ch
        ```

4. `Phân loại Channel`
- Trong Go, có hai loại channel chính với cơ chế hoạt động khác nhau hoàn toàn về mặt đồng bộ:
    1. `Unbuffered Channel` (Channel không đệm)
        + Đây là loại channel "giao hàng trực tiếp". Hãy tưởng tượng nó như một cuộc chuyển giao tài liệu tận tay: Người đưa phải gặp tận mặt người nhận thì việc chuyển giao mới hoàn tất.
        + Cú pháp: 
            ```go
            ch := make(chan int)
            ```
        + Cơ chế hoạt động: 
            + Người gửi (Sender) sẽ bị chặn (block) cho đến khi có Người nhận (Receiver) sẵn sàng lấy dữ liệu.
            + Ngược lại, Người nhận cũng bị chặn cho đến khi có Người gửi đẩy dữ liệu vào.

        + Mục đích: Dùng để đồng bộ hóa tuyệt đối giữa hai Goroutine.

    2. `Buffered Channel` (Channel có đệm)
        + Loại này giống như một cái "hòm thư" hoặc "kho chứa tạm". Người gửi có thể quẳng đồ vào đó rồi đi làm việc khác, miễn là hòm thư vẫn còn chỗ.
        + Cú pháp:
            ```go
            ch := make(chan int, capacity) (Ví dụ: capacity = 3)
            ```
        + Cơ chế hoạt động:
            + Người gửi: Chỉ bị chặn khi bộ đệm đã đầy (Full). Nếu vẫn còn chỗ, người gửi cứ đẩy vào và chạy tiếp mà không cần quan tâm người nhận đã lấy hay chưa.

            + Người nhận: Chỉ bị chặn khi bộ đệm đang trống (Empty). Nếu có ít nhất một phần tử trong đệm, người nhận lấy ra và chạy tiếp.

        + Mục đích: Tăng hiệu năng, cho phép người gửi tiếp tục làm việc mà không cần đợi người nhận ngay lập tức (xử lý bất đồng bộ nhẹ).

5. `Directional Channels (Channel định hướng)`
- Trong Go, mặc định khi bạn tạo một channel bằng `make(chan T)`, nó là `bi-directional (hai chiều)` — tức là bạn có thể vừa gửi vừa nhận dữ liệu trên đó.
- Tuy nhiên, `Directional Channels (Channel định hướng)` cho phép bạn giới hạn quyền hạn của một channel trong một phạm vi cụ thể (thường là trong các hàm). 

- Đây là một tính năng cực kỳ thông minh của Go để tăng tính `an toàn (Type Safety)` và `rõ ràng (Readability)` cho mã nguồn.

    1. Cú pháp và Cách phân biệt
    - Cách nhớ rất đơn giản: Mũi tên `<-` chỉ đi đâu thì dữ liệu đi đó.
        ```go
        chan T: Channel hai chiều (mặc định)

        chan<- T: Mũi tên hướng vào channel => Chỉ gửi (Send-only)

        <-chan T: Mũi tên hướng ra khỏi channel => Chỉ nhận (Receive-only).
        ```
    2. Tại sao lại cần Directional Channels?

    - Lợi ích 1: Ngăn ngừa lỗi logic (Compile-time check)
        + Nếu bạn cố tình nhận dữ liệu từ một `chan<- int` (chỉ gửi), trình biên dịch (compiler) sẽ báo lỗi ngay lập tức thay vì đợi đến lúc chương trình chạy mới bị lỗi (runtime error).

    - Lợi ích 2: Làm rõ ý đồ của lập trình viên
        + Nhìn vào chữ ký của hàm (function signature), người khác sẽ biết ngay hàm này dùng channel để làm gì:
            ```go
            func produce(out chan<- int) // "À, hàm này chỉ đổ dữ liệu vào channel thôi."
            
            func consume(in <-chan int) // "À, hàm này chỉ lấy dữ liệu ra để xử lý."

    3. Ví dụ thực tế: Mô hình Producer - Consumer
    - Go sẽ tự động chuyển đổi từ channel hai chiều sang một chiều khi bạn truyền vào hàm.
        ```go
        package main

        import "fmt"

        // Hàm này CHỈ ĐƯỢC PHÉP GỬI dữ liệu (Send-only)
        func producer(out chan<- int) {
            for i := 0; i < 5; i++ {
                out <- i // Gửi dữ liệu vào
            }
            close(out) // Đóng channel sau khi gửi xong
            // fmt.Println(<-out) // NẾU BỎ COMMENT DÒNG NÀY SẼ BỊ LỖI KHI COMPILER
        }

        // Hàm này CHỈ ĐƯỢC PHÉP NHẬN dữ liệu (Receive-only)
        func consumer(in <-chan int) {
            for v := range in {
                fmt.Println("Nhận được:", v)
            }
            // in <- 10 // NẾU BỎ COMMENT DÒNG NÀY SẼ BỊ LỖI KHI COMPILER
        }

        func main() {
            // 1. Tạo một channel hai chiều bình thường
            ch := make(chan int)

            // 2. Truyền vào các hàm, Go tự động ép kiểu sang channel một chiều
            go producer(ch)
            consumer(ch) // Đợi consumer xử lý xong
        }
        ```

    4. Những quy tắc cần nhớ
    - Chuyển đổi một chiều: Bạn có thể chuyển từ channel hai chiều sang một chiều, nhưng KHÔNG THỂ chuyển ngược lại.
    - Đóng channel: Chỉ nên đóng channel ở phía `Gửi (chan<-)`. Nếu bạn cố đóng một channel ở phía `Nhận (<-chan)`, trình biên dịch sẽ báo lỗi. Điều này cực kỳ hợp lý vì người nhận không nên có quyền đóng "vòi nước" mà họ đang dùng.

    5. Những trạng thái "nguy hiểm" cần nhớ
    - Làm việc với channel rất dễ gây ra lỗi `Deadlock` hoặc `Panic` nếu không cẩn thận:
        + Gửi vào channel đã đóng: Gây ra `panic`.
        + Đóng channel đã đóng: Gây ra `panic`.
        + Gửi/Nhận từ channel nil: Sẽ bị chặn mãi mãi (Deadlock).
        + Chỉ gửi mà không có ai nhận (hoặc ngược lại): Gây ra lỗi `all goroutines are asleep - deadlock!`.

### 27. Range và Close
1. Lệnh `close(ch)`
- Lệnh close dùng để thông báo rằng: "Sẽ không có thêm giá trị nào được gửi vào channel này nữa".
- Quy tắc "Vàng" khi đóng Channel:
    + Chỉ người gửi mới nên đóng: Người nhận không bao giờ nên đóng channel vì họ không biết khi nào dữ liệu thực sự kết thúc.

    + Không gửi vào channel đã đóng: Nếu bạn cố làm vậy (`ch <- v`), chương trình sẽ bị Panic.

    + Không đóng 2 lần: Đóng một channel đã đóng cũng gây Panic.

    + Đóng rồi vẫn nhận được: Nếu channel còn dữ liệu trong bộ đệm (buffered), người nhận vẫn lấy ra được cho đến khi hết sạch. Sau đó, họ sẽ nhận được "giá trị zero" (như `0, "", false`).

2. Lệnh `range` với Channel
- Thay vì dùng vòng lặp `for` vô tận, bạn có thể dùng `for range` để đọc dữ liệu từ channel
- Cơ chế:
    + `range chỉ dừng khi channel đã bị đóng VÀ đã đọc hết toàn bộ dữ liệu còn lại trong buffered channel`.
    + Vòng lặp range sẽ liên tục đợi (block) và lấy dữ liệu từ channel cho đến khi channel đó bị đóng và đã đọc hết toàn bộ dữ liệu còn lại trong buffered channel

- Ví dụ:
    ```go
    package main

    import "fmt"

    func main() {
        ch := make(chan int, 3)

        go func() {
            for i := 1; i <= 3; i++ {
                ch <- i
            }
            // Nếu không có dòng này, vòng range ở dưới sẽ đợi mãi -> Deadlock!
            close(ch) 
        }()

        // range sẽ chạy cho đến khi channel bị close
        for v := range ch {
            fmt.Println("Nhận:", v)
        }
        
        fmt.Println("Channel đã đóng, vòng lặp kết thúc!")
    }
    ```
3. Cách kiểm tra Channel đã đóng chưa
- Đôi khi bạn không dùng range mà dùng lệnh nhận đơn lẻ. Làm sao biết giá trị 0 nhận được là do người gửi gửi số 0, hay do channel đã đóng?
- Bạn sử dụng cú pháp này:
    ```go
    v, ok := <-ch
    if !ok {
        fmt.Println("Channel đã đóng và không còn dữ liệu!")
    }
    ```
    👉 
    
    `ok == true`: Channel vẫn đang mở hoặc vẫn còn dữ liệu trong đệm để đọc.

    `ok == false`: Channel đã đóng và không còn dữ liệu nào bên trong.

    👉 Nếu không có biến `ok` và `ch` đã bị đóng thì `v` nhận giá trị `zero value` chứ không bị lỗi chương trình

### 28. Select
1. `select` là gì?
- `select` là cấu trúc dùng để chờ và xử lý nhiều thao tác channel cùng lúc.
- Nó giống như `switch`, nhưng thay vì so sánh giá trị, nó chờ các operation trên channel (send/receive).

2. Cú pháp:
    ```go
    select {
    case v := <-ch1:
        // nhận từ ch1
    case ch2 <- 10:
        // gửi vào ch2
    default:
        // chạy nếu không case nào sẵn sàng
    }
    ```
3. Cách hoạt động:
- Bước 1: Go kiểm tra tất cả các case
- Bước 2:Nếu có nhiều case sẵn sàng → chọn ngẫu nhiên một case
- Bước 3: Nếu không có case nào sẵn sàng:
    + Có `default` → chạy `default`
    + Không có `default` → block (chờ)

4. Ví dụ:
- Ví dụ 1: Chờ nhiều channel
    ```go
    package main

    import (
        "fmt"
        "time"
    )

    func main() {
        ch1 := make(chan string)
        ch2 := make(chan string)

        go func() {
            time.Sleep(1 * time.Second)
            ch1 <- "from ch1"
        }()

        go func() {
            time.Sleep(2 * time.Second)
            ch2 <- "from ch2"
        }()

        select {
        case msg := <-ch1:
            fmt.Println(msg)
        case msg := <-ch2:
            fmt.Println(msg)
        }
    }
    ```
    👉 Sẽ in "from ch1" vì ch1 sẵn sàng trước.


- Ví dụ 2: Non-blocking với default
    ```go
    select {
    case msg := <-ch:
        fmt.Println(msg)
    default:
        fmt.Println("No data")
    }
    ```
    👉 Nếu ch chưa có dữ liệu → không block → in "No data"

- Ví dụ 3: Timeout pattern (rất hay dùng)
    ```go
    select {
    case msg := <-ch:
        fmt.Println(msg)
    case <-time.After(2 * time.Second):
        fmt.Println("timeout")
    }
    ```
    👉 Cách hoạt động:
    + select sẽ chờ một trong các case sẵn sàng.
    + Nếu ch nhận được dữ liệu trước 2 giây → in msg.
    + Nếu sau 2 giây vẫn chưa có dữ liệu từ ch → time.After(2 * time.Second) sẽ kích hoạt → in "timeout".
    + Bởi vì: `time.After()` tạo ra một channel và sẽ gửi một giá trị vào đó sau khoảng thời gian chỉ định. Vì vậy nó rất hay dùng để xử lý timeout trong `select`.

5. `Một số lưu ý quan trọng`:
    1. select không phải là vòng lặp
    - Muốn lặp thì phải viết:
        ```go
        for {
            select {
            ...
            }
        }
        ```
    2. Nếu tất cả channel đều là nil và select không có default thì select sẽ block vĩnh viễn.
        ```go
        var ch1 chan int
        var ch2 chan string

        select {
        case <-ch1:
        case <-ch2:
        }
        ```
        👉 Cả ch1 và ch2 đều là nil → chương trình sẽ deadlock (nếu ở main goroutine).

### 29. Mutex
1. Mutex là gì?
- Mutex (Mutual Exclusion) là cơ chế dùng để đảm bảo tại một thời điểm chỉ có một goroutine được truy cập vào vùng dữ liệu dùng chung (ví dụ: biến) hoặc một đoạn code nhất định.
- Trong Go, Mutex được cung cấp bởi package `sync`.

2. Khi nào cần Mutex?
- Khi có nhiều goroutine cùng đọc/ghi một biến chung, có thể xảy ra:

    ❌ Race condition

    ❌ Dữ liệu sai lệch

    ❌ Crash khó debug

- Ví dụ không dùng Mutex (bị race)
    ```go
    package main

    import (
        "fmt"
    )

    var counter int

    func main() {
        for i := 0; i < 1000; i++ {
            go func() {
                counter++
            }()
        }

        fmt.Println(counter)
    }
    ```
    👉 Kết quả gần như chắc chắn không phải 1000.

3. Cách sử dụng `sync.Mutex`:
- Gói sync cung cấp Mutex với hai phương thức cơ bản:
    + Lock(): Bắt đầu khóa. Nếu đã có goroutine khóa rồi, Goroutine hiện tại sẽ phải đứng đợi.
    + Unlock(): Mở khóa để goroutine khác vào.
- Cách hoạt động:
    + Nếu một goroutine đã Lock() thì Goroutine khác gọi Lock() sẽ bị block cho đến khi Unlock() được gọi

    + Nếu mutex đang mở (unlock) → goroutine đó chiếm lock và chạy tiếp.
    + Nếu mutex đang bị giữ bởi goroutine khác → goroutine hiện tại sẽ bị block tại dòng Lock() cho đến khi mutex được Unlock().

- Nếu một goroutine đã `mu.Lock()` mà không bao giờ `mu.Unlock()`, thì: Tất cả goroutine khác gọi `mu.Lock()` cùng mutex đó sẽ block mãi mãi tại dòng `Lock()`.
- Để tránh quên mở khóa (dẫn đến Deadlock), chúng ta luôn đặt `Unlock()` ngay sau `Lock()` bằng từ khóa `defer`
    ```go
    mu.Lock()
    defer mu.Unlock()
    ```

### 30. Các loại Mutex trong Golang
1. `sync.Mutex`
- Đây là loại khóa "độc quyền" (Exclusive Lock). 
- Một khi một Goroutine đã giữ khóa này, không có bất kỳ Goroutine khác (dù là đọc hay ghi) được phép vào cho đến khi khóa được mở.
- Trạng thái: Chỉ có 2 trạng thái là `Locked` (đã khóa) và `Unlocked` (đã mở).
- Cơ chế: Nếu Goroutine A đang Lock(), Goroutine B gọi Lock() sẽ bị chặn hoàn toàn cho đến khi A gọi Unlock().

- Đặc điểm:
    + Chỉ cho phép 1 goroutine truy cập vào critical section tại một thời điểm.
    + Goroutine khác sẽ bị block cho đến khi mutex được unlock.
    + Không phân biệt đọc/ghi.

- Khi nên dùng:
    + Khi tài nguyên bị ghi thường xuyên.
    + Khi không cần phân biệt read/write.

- Ví dụ:
    ```go
    var mu sync.Mutex

    mu.Lock()
    // critical section
    mu.Unlock()
    ```

2. `sync.RWMutex` (Reader/Writer Mutex)
- Đây là loại khóa "ưu tiên người đọc", giúp tối ưu hiệu năng cực tốt cho các hệ thống có tần suất đọc dữ liệu cao hơn ghi.
- Cho phép 1 thời điểm có thể có nhiều reader hoặc chỉ có 1 writer.

- Nó cung cấp hai bộ phương thức:
    + Cho việc Ghi (`Lock` và `Unlock`): Hoạt động y hệt `sync.Mutex`. Độc quyền hoàn toàn.
    + Cho việc Đọc (`RLock` và `RUnlock`): Nhiều Goroutine có thể cùng giữ `RLock` một lúc.
        * Nếu có ai đó đang giữ RLock, người muốn Lock (để ghi) phải đợi.
        * Nếu có ai đó đang giữ Lock (để ghi), người muốn RLock (để đọc) phải đợi.

- Đặc điểm:
    + Nhiều goroutine có thể RLock() cùng lúc.
    + Khi có Lock() (write), tất cả reader và writer khác sẽ bị block.
    + Tối ưu cho trường hợp read nhiều – write ít.

- Ví dụ:
    ```go
    var mu sync.RWMutex

    mu.RLock()
    // read
    mu.RUnlock()

    mu.Lock()
    // write
    mu.Unlock()
    ```

3. `sync.Once`
- Không phải mutex thuần túy, nhưng dùng cơ chế khóa bên trong.
- Đặc điểm:
    + Đảm bảo một function chỉ chạy 1 lần duy nhất.
    + Thread-safe.
- Khi nên dùng:
    + Lazy initialization
    + Singleton pattern

- Ví dụ:
    ```go
    var once sync.Once

    once.Do(func() {
        initConfig()
    })
    ```
4. `sync.Map`
- Không phải mutex trực tiếp nhưng là cấu trúc map đã được đồng bộ hóa.
- Dùng cơ chế lock + atomic bên trong.
- Đặc điểm:
    + Thread-safe map.
    + Tối ưu cho read-heavy workloads.
    + Không cần dùng mutex bên ngoài.
- Khi nên dùng
    + Concurrent map read-heavy
    + Không muốn tự quản lý mutex

- Ví dụ:
    ```go
    var m sync.Map
    m.Store("a", 1)
    v, ok := m.Load("a")
    ```

5. `Atomic` (không phải mutex nhưng thay thế mutex nhỏ)
- Package sync/atomic cho phép lock-free operation.
    ```go
    atomic.AddInt64(&counter, 1)
    ```
- Dùng khi:
    + Chỉ thao tác số nguyên
    + Cần hiệu năng cao
    + Tránh lock

6. Tóm tắt: Chọn loại nào?
    - Chỉ cần tăng/giảm số nguyên đơn giản? $\rightarrow$ Dùng sync/atomic.
    - Dữ liệu đọc nhiều, ghi ít (ví dụ: Cache, Config)? $\rightarrow$ Dùng sync.RWMutex.
    - Tác vụ ghi nhiều hoặc logic phức tạp? $\rightarrow$ Dùng sync.Mutex.
    - Cần đảm bảo một đoạn code chỉ chạy 1 lần duy nhất (như khởi tạo DB)? $\rightarrow$ Dùng sync.Once (nó sử dụng Mutex ngầm bên dưới).

7. Lưu ý quan trọng:
- Unlock() khi chưa Lock() → panic
- RUnlock() khi chưa RLock() → panic
- Code này sẽ bị deadlock, không phải panic.
    ```go
    mu.Lock()
    mu.Lock() // deadlock (không panic)
    ```

### 31. WaitGroup
- Trong Go (Golang), `WaitGroup` là một cấu trúc trong `package sync` dùng để đợi một nhóm goroutine hoàn thành trước khi tiếp tục thực thi chương trình.

1. `WaitGroup` hoạt động như thế nào?
- `WaitGroup`có 3 method chính:
    + Add(n): Tăng bộ đếm thêm n (số goroutine cần đợi)
    + Done(): Giảm bộ đếm đi 1 (gọi khi một goroutine xong)
    + Wait(): Chặn (block) cho đến khi bộ đếm về 0

- Ví dụ:
    ```go
    package main

    import (
        "fmt"
        "sync"
        "time"
    )

    func main() {
        var wg sync.WaitGroup

        for i := 1; i <= 3; i++ {
            wg.Add(1) // báo có thêm 1 goroutine
            go func(id int) {
                defer wg.Done() // báo hoàn thành khi xong
                fmt.Println("Worker", id, "đang chạy")
                time.Sleep(2 * time.Second)
                fmt.Println("Worker", id, "hoàn thành")
            }(i)
        }

        wg.Wait() // đợi tất cả goroutine hoàn thành
        fmt.Println("Tất cả công việc đã xong")
    }
    ```
    📌 Cách hoạt động trong ví dụ trên
    + wg.Add(1) → tăng bộ đếm lên 1 mỗi khi tạo goroutine.
    + wg.Done() → giảm bộ đếm khi goroutine hoàn thành.
    + wg.Wait() → chương trình sẽ dừng ở đây cho đến khi tất cả goroutine gọi Done().
    + Nếu không có Wait(), chương trình có thể kết thúc trước khi goroutine chạy xong.

    ⚠️ Lưu ý quan trọng
    + ❗ Luôn gọi Add() trước khi chạy goroutine.
    + ❗ Không được copy WaitGroup (nên truyền bằng pointer nếu dùng trong struct).
    + ❗ Nếu Done() nhiều hơn Add() → chương trình sẽ panic.
    + WaitGroup chỉ dùng để đợi hoàn thành, không dùng để truyền dữ liệu (dùng channel cho việc đó).

### 32. Context
- Trong Go (Golang), Context là một cơ chế dùng để:

    ✅ Quản lý timeout

    ✅ Hủy (cancel) nhiều goroutine cùng lúc

    ✅ Truyền request-scoped data (metadata) xuyên suốt call stack

- Nó nằm trong package chuẩn:
    ```go
    import "context"
    ```

1. Context là gì?
- `context.Context` là một interface giúp:
    + Kiểm soát vòng đời của request, goroutine
    + Truyền tín hiệu dừng giữa các goroutine
    + Tránh goroutine bị leak


- Đặc biệt quan trọng trong:
    + Web server
    + Microservices
    + Database call
    + API call
    + Xử lý concurrent

2. Các loại Context chính

    1️⃣ `context.Background()`
    - Context gốc
    - Thường dùng trong main() hoặc khởi tạo server
        ```go
        ctx := context.Background()
        ```

    2️⃣ `context.WithCancel()`
    - Tạo context có thể hủy thủ công.
        ```go
        ctx, cancel := context.WithCancel(parentCtx)
        cancel() // hủy
        ```
        👉 Dùng khi cần dừng nhiều goroutine cùng lúc.

    3️⃣ `context.WithTimeout()`
    - Tự động hủy sau một khoảng thời gian.
        ```go
        ctx, cancel := context.WithTimeout(parentCtx, duration)
        defer cancel()
        ```
        👉 Rất hay dùng cho: Query database, Call API, HTTP request


    4️⃣ `context.WithDeadline()`
    - Hủy tại một thời điểm cụ thể.
        ```go
        ctx, cancel := context.WithDeadline(parentCtx, deadlineTime)
        ```

    5️⃣ `context.WithValue()`
    - Truyền dữ liệu qua các layer.
        ```govalue
        ctx := context.WithValue(parentCtx, key, value)
        ```
    - ⚠️ Không nên lạm dụng. Chỉ dùng cho metadata như: request ID, user ID, auth token

3. Context hoạt động như thế nào?
- Trong Go, mọi `context.Context` đều tạo thành một cây (`tree`).
- Mỗi context mới được tạo ra sẽ:
    + Có 1 parent
    + Có thể có nhiều child
    + Khi parent bị cancel → tất cả child tự động bị cancel theo

- Gốc của cây: 
    + Cây luôn bắt đầu từ: `context.Background()` hoặc `context.TODO()`
    
        👉 Hai cái này là root context, không bao giờ bị cancel.

- 🌲 Ví dụ cấu trúc cây

    ```go
    root := context.Background()
    ctx1, cancel1 := context.WithCancel(root)
    ctx2, cancel2 := context.WithTimeout(ctx1, 5*time.Second)
    ctx3 := context.WithValue(ctx2, "userID", 123)
    ```

- Cây sẽ trông như sau:

    ```bash
    Background()
        └── ctx1 (WithCancel)
                └── ctx2 (WithTimeout 5s)
                        └── ctx3 (WithValue userID=123)
    ```

    🔥 Nguyên tắc quan trọng

    1️⃣ Cancel lan truyền từ trên xuống
    - Nếu bạn gọi: 
        ```go
        cancel1()
        ```
        👉 ctx1 bị hủy

        👉 ctx2 bị hủy

        👉 ctx3 bị hủy

    - Nhưng nếu bạn chỉ gọi:
        ```go
        cancel2()
        ```
        👉 ctx2 và ctx3 bị hủy

        👉 ctx1 vẫn sống

    ➡️ Parent không bị ảnh hưởng bởi child.

    2️⃣ Deadline cũng lan truyền

    - Nếu parent có timeout 3s, còn child đặt 10s:
        ```go
        parent, _ := context.WithTimeout(root, 3*time.Second)
        child, _ := context.WithTimeout(parent, 10*time.Second)
        ```
        ⏳ Child thực tế chỉ sống 3s.

        👉 Vì deadline của child không thể vượt quá parent.

    3️⃣ WithValue không tạo cơ chế hủy mới
    - WithValue chỉ:
        + Thêm dữ liệu
        + Không thêm cancel logic

    - Ví dụ:
        ```go
        ctxValue := context.WithValue(ctx2, "role", "admin")
        ```

        👉 Nó chỉ gắn data vào node đó trong cây.

💡 Tóm tắt 5 loại trong cây

    | Loại         | Tạo node mới? | Có thể cancel?     | Có deadline? | Có value? |
    | ------------ | ------------- | ------------------ | ------------ | --------- |
    | Background   | Root          | ❌                 | ❌           | ❌        |
    | WithCancel   | ✅            | ✅                 | ❌           | ❌        |
    | WithTimeout  | ✅            | ✅                 | ✅           | ❌        | 
    | WithDeadline | ✅            | ✅                 | ✅           | ❌        |
    | WithValue    | ✅            | ❌ (kế thừa parent)| ❌           | ✅        |

4. `Lưu ý quan trọng (Best Practices)`
- `Context là tham số đầu tiên`: Theo quy ước, `ctx context.Context` luôn đứng đầu danh sách tham số của hàm.

- `Đừng lưu Context vào Struct`: Chỉ truyền nó qua các hàm. Việc lưu vào struct khiến vòng đời của nó trở nên khó kiểm soát.

- `Luôn gọi cancel()`: Khi dùng `WithCancel` hoặc `WithTimeout`, hãy luôn gọi hàm cancel (thường dùng `defer cancel()`) để tránh rò rỉ bộ nhớ (goroutine leak).

- `WithValue chỉ dùng cho metadata`: Không dùng nó để truyền các tham số tùy chọn vào hàm (như database connection), hãy dùng nó cho những thứ như TraceID.