# Bookmark Management

## Note
### 1. API, Restful API là gì?
### 2. Mux là gì?
- Mux và Gorilla Mux không phải là framework dựng API server. Chúng chỉ là router (bộ định tuyến HTTP) thôi.

1️⃣ Mux trong Go là gì

- Trong Go, Mux = HTTP request multiplexer → nhiệm vụ là định tuyến request tới handler phù hợp.

- Ví dụ request:
    ```
    GET /users
    POST /login
    GET /movies/1
    ```
    Mux sẽ quyết định:

    | URL         | Handler         |
    | ----------- | --------------- |
    | `/users`    | getUsersHandler |
    | `/login`    | loginHandler    |
    | `/movies/1` | getMovieHandler |

- ✅ Kết luận
    + `Mux / Gorilla Mux = router`
    + Không phải framework
    + API server vẫn chạy bằng `net/http`
    + Framework thật sự trong Go là: `Gin, Echo, Fiber`


2️⃣ net/http ServeMux (mux mặc định của Go)

- Go có sẵn router trong standard library:
```go
import "net/http"

mux := http.NewServeMux()

mux.HandleFunc("/users", getUsers)
mux.HandleFunc("/login", login)

http.ListenAndServe(":8080", mux)
```

- Ở đây:
```
  HTTP Server
        ↓
  ServeMux (router)
        ↓
  Handler function
```

3️⃣ gorilla/mux là gì
- Gorilla Mux là thư viện router nâng cao cho Go.

- Import:
```go
import "github.com/gorilla/mux"
```

- Ví dụ:
```go
r := mux.NewRouter()

r.HandleFunc("/users", getUsers).Methods("GET")
r.HandleFunc("/users/{id}", getUser).Methods("GET")

http.ListenAndServe(":8080", r)
```

- Nó thêm các tính năng:
    + Path params /users/{id}
    + HTTP method routing
    + Middleware
    + Subrouter
    + Regex route

### 3. Framework API Server:
    https://github.com/gin-gonic/gin
    https://github.com/go-chi/chi
    https://github.com/labstack/echo
    https://github.com/gofiber/fiber

### 3. Thư viện viết unit test cho package service:
https://github.com/stretchr/testify

### 4. Thư viện viết unit test cho package handler:
https://github.com/vektra/mockery


1️⃣ `Mockery` là gì?
- `Mockery` là một CLI tool viết bằng Go dùng để tự động sinh mock từ interface.

    👉 Nó đọc interface Go của bạn

    👉 Tạo ra file mock tương ứng

    👉 Dùng trong unit test

- Nó thuộc nhóm code generation tool.

2️⃣ Mock là gì trong unit test?
- Giả sử bạn có service:
    ```go
    type UserRepository interface {
        GetUser(id int) (*User, error)
    }
    ```

- Trong test, bạn không muốn:
    + Gọi DB thật
    + Gọi API thật
    + Gọi Redis thật

    👉 Bạn tạo một fake implementation của UserRepository

    Đó gọi là mock.

3️⃣ Nếu không có mockery thì sao?
- Bạn phải tự viết:
  ```go
    type MockUserRepository struct {
    }

    func (m *MockUserRepository) GetUser(id int) (*User, error) {
        return &User{Name: "Test"}, nil
    }
    ```

    👉 Viết tay rất mệt

    👉 Interface lớn là cực hình

5️⃣ Mockery là library hay tool?

👉 Nó là tool

- Vì:
    + Bạn không import mockery trong runtime
    + Bạn chỉ chạy nó để generate code
    + Sau khi generate xong → mockery không còn liên quan

5️⃣ Cách cài đặt

- Cách 1: (Ngon nhất) Không cần cài đặt gì thêm chỉ cần dùng trong code go
    ```go
    //go:generate go run github.com/vektra/mockery/v2@latest --name Password --filename password_service.go
    ```

- Cách 2: cài đặt ở global: 
    ```bash
    go install github.com/vektra/mockery/v2@latest
    ```
    ✅ Bước 1: Kiểm tra mockery có tồn tại không
    
    Chạy:

    ```bash
    ls $HOME/go/bin
    ```

    Nếu thấy:
    ```terminal
    mockery
    ```
    → OK, đã cài rồi.

    Nếu không thấy → chạy lại:
    ```bash
    go install github.com/vektra/mockery/v2@latest
    ```

    ✅ Bước 2: Thêm vào PATH

    Mở file:
    ```bash
    nano ~/.bashrc
    ```

    Thêm dòng này vào cuối file:
    ```
    export PATH=$PATH:$HOME/go/bin
    ```

    Lưu lại:
    ```
    Ctrl + O
    Enter
    Ctrl + X
    ```
    ✅ Bước 3: Reload shell
    ```bash
    source ~/.bashrc
    ```

    ** ✅ Chú ý cách gộp bước 2 và bước 3 là:
    ```bash
    echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc
    source ~/.bashrc
    ```


    ✅ Bước 4: Kiểm tra lại
    ```bash
    which mockery
    ```
    
    Phải ra:
    ```
    /home/dongnt1/go/bin/mockery
    ```

    Rồi test:
    ``` bash
    mockery --version
    ```
    
    Nếu ra version → xong

    ✅ Bước 5: Restart VS Code

    RẤT QUAN TRỌNG.

    Tắt hẳn VS Code rồi mở lại.


### 5. Viết Integration Test: 
- Viết trong thư mục `internal/test/endpoint`

### 6. Thư viện cấu hình biến môi trường:
https://github.com/kelseyhightower/envconfig

- Run Health Check:
    ```bash
    export BOOKMARK_SERVICE_APP_PORT=9000 BOOKMARK_SERVICE_SERVICE_NAME=health-check BOOKMARK_SERVICE_INSTANCE_ID=health-123
    ```

    ```bash
    go run .
    ```

- Cách loại bỏ export:
    + Cách 1:
        ```bash
        unset BOOKMARK_SERVICE_APP_PORT
        unset BOOKMARK_SERVICE_SERVICE_NAME
        unset BOOKMARK_SERVICE_INSTANCE_ID
        ```

    + Cách 2:
        ```bash
        exit
        ```
        hoặc đóng terminal → mở terminal mới -> Biến môi trường sẽ tự mất.

- Kiểm tra lại:
    ```bash
    echo $BOOKMARK_SERVICE_APP_PORT
    ```
    Nếu không in gì → đã xoá.

### 7. Thư viện tạo API Documentation

https://github.com/swaggo/swag

1. Các bước tạo swagger document:
    + Viết comment vào cmd/api/main.go và các method trong handler (ví dụ: GenPass trong internal/handler/password.go)
    + Gõ cmd vào thư mục gốc:
        ```bash
        swag init -g cmd/api/main.go
        ```

2. Cách cài đặt thư viện swaggo

1️⃣ Cài swag bằng Go

- Chạy lệnh:
    ```bash
    go install github.com/swaggo/swag/cmd/swag@latest
    ```

- Lệnh này sẽ:
    + download source
    + build CLI tool
    + tạo binary swag

- Binary sẽ nằm ở:
    ```terminal
    ~/go/bin/swag
    ```

2️⃣ Thêm Go bin vào PATH

- Nếu chưa thêm PATH thì terminal sẽ báo lỗi:
    ```bash
    command 'swag' not found
    ```

- Thêm PATH:
    ```bash
    echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc
    ```

- Sau đó reload:
    ```bash
    source ~/.bashrc
    ```

3️⃣ Kiểm tra cài thành công

- Chạy:
    ```bash
    swag --version
    ```

- Ví dụ kết quả:
    ```
    swag version v1.16.x
    ```
4️⃣ Cài thư viện Swagger cho project Go
- Trong project Go chạy:
    ```bash
    go get github.com/swaggo/swag
    go get github.com/swaggo/gin-swagger
    go get github.com/swaggo/files
    ```

5️⃣ Viết comment swagger trong main.go
- Ví dụ:
```go
// @title Bookmark API
// @version 1.0
// @description API for bookmark service
// @host localhost:9000
// @BasePath /
func main() {
}
```

6️⃣ Generate Swagger docs
- Chạy câu lệnh sau tại thư mục gốc sẽ gen ra thư mục docs:
    ```bash
    swag init -g cmd/api/main.go
    ```

7️⃣ Test Swagger
- Sau khi chạy câu lệnh generate docs trên và run project thì vào `http://localhost:8080/swagger/index.html` để 

### 8. Tạo script chạy cmd (dùng Makefile)
1. Cài đặt make
    ```bash
    sudo apt install make
    ```
2. Tạo file Makefile và viết script vào Makefile
    ```Makefile
    .PHONY: run swagger dev-run test

    run:
        go run cmd/api/main.go

    swagger:
        swag init -g cmd/api/main.go

    dev-run: swagger run

    test:
        go test ./...
    ```
    
    Thêm `.PHONY` để nếu trường hợp project gốc có các thư mục trùng tên với các câu lệnh thì khi gọi: `make <tên câu lệnh>` vẫn gọi bình thường

3. Cách chạy:
    ```bash
    make run
    make swagger
    make dev-run
    make test
    ```

### 9. Run docker (redis)
1. Chạy terminal
```bash
docker run -d --name redis -p 6379:6379 redis:alpine
docker ps
docker exec -it redis redis-cli     # vào redis-cli để tương tác với redis container qua terminal
    set mykey hello
    get mykey
    type mykey          # Xem type của key
    set "12345" "google.com"
    keys *
    get 12345
    del mykey 12345     # Xóa mykey và 12345
    keys *

    flushdb     # Xóa toàn bộ key trong database hiện tại
    flushall    # Xóa toàn bộ Redis (tất cả database) (Trong Redis, mặc định có 16 database: 0, 1, ..., 15)

    exists mykey    # Kiểm tra key có tồn tại (1 -> tồn tại, 0 -> không tồn tại)
    scan 0          # KEYS * có thể lag Redis nếu dữ liệu lớn, nên Redis khuyên dùng
    dbsize          # Xem số lượng key trong database
    
    select 0    # chuyển sang database 0 (database 0 là mặc định)
    select 1    # chuyển sang database 1
    select 15   # chuyển sang database 15

    exit
docker rm -f redis
```

2. Cài thư viện vào project: https://github.com/redis/go-redis

```bash
go get github.com/redis/go-redis/v9
```

### 10. Viết unit test cho redis
1. Thư viện tạo mock redis: https://github.com/alicebob/miniredis
```bash
go get github.com/alicebob/miniredis/v2
```
- Thư viện miniredis là một Redis server giả lập (mock Redis) viết bằng Go, dùng chủ yếu để test code mà không cần chạy Redis thật.