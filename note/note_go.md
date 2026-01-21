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

1. `Pointers`: 
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
    + Hệ quả: Nếu mảng của bạn có 1 triệu phần tử, việc gán mảng sẽ rất tốn kém tài nguyên và làm chậm chương trình. Đây là lý do vì sao chúng ta thường dùng Slice hoặc truyền Con trỏ mảng (Array Pointer).
    + Duyệt mảng với `range`
        ```go
        fruits := [3]string{"Apple", "Banana", "Cherry"}

        for index, value := range fruits {
            fmt.Printf("Vị trí %d có quả: %s\n", index, value)
        }
        ```



3. `Slice`:
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

4. Cách hoạt động của `append`
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

5. Cách hoạt động của `make`, `copy`, `Full Slice Expression`
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
        - `make`:
            + Dùng khi muốn nhân bản dữ liệu ra một vùng nhớ độc lập.
            + Mục tiêu chính: An toàn dữ liệu, tránh Memory Leak.
        - `make`:
            + Dùng khi chia nhỏ slice lớn thành các slice nhỏ để xử lý riêng.
            + Mục tiêu chính: Cô lập vùng nhớ, ngăn chặn ghi đè ngoài ý muốn.

6. `Nill`, `Nil slice`
- The zero value of a slice is `nil`.
- A nil slice has a length and capacity of 0 and has no underlying array.

7. `Map`