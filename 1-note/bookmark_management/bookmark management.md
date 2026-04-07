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

docker logs redis    # Hiển thị toàn bộ log đã được ghi lại của container (Sau khi in ra log hiện có, lệnh sẽ kết thúc ngay)
docker logs -f redis # Xem log của container redis theo thời gian thực (Lệnh sẽ tiếp tục hiển thị các log mới khi container ghi thêm log)

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

### 11. Thư viện viết log (Zerolog)
- Thư viện viết log (Zerolog): https://github.com/rs/zerolog
```bash
go get github.com/rs/zerolog/log
```
- Lưu ý: tất cẩ các lỗi ở phía server (`http.StatusInternalServerError`) thì đều cần phải ghi log
    ```
    log.Error().Err(err).Msg("Failed to get URL")
    ```

### 12. Thư viện cấu hình CORS: https://github.com/gin-contrib/cors
```bash
go get github.com/gin-contrib/cors
```


### 13. Viết Dockerfile
1. Tạo Dockerfile
```Dockerfile
FROM golang:1.25.6-alpine

RUN mkdir -p /opt/app

WORKDIR /opt/app

COPY . .

RUN apk add build-base

RUN go mod download && go build -o bookmark_service cmd/api/main.go

CMD ["/opt/app/bookmark_service"]
```

2. Build
```bash
docker build -t bookmark_service:test .
```
- Nếu vẫn gặp lỗi DNS khi build
    + Build với network host:
        ```bash
        docker build --network=host -t bookmark_service:test .
        ```

3. Kiểm tra
```bash
docker images
```

4. Run container
```
docker run -d --name bookmark -p 8080:8080 bookmark_service:test
```

- Thêm biến môi trường
```bash
docker run -d --name bookmark --env APP_PORT=8082 -p 8082:8082 bookmark_service:test
```

- Như vậy vẫn chưa kết nối được app với redis => phải thêm `--network host` (hoặc `--network=host`) vào mới chạy được (Đoạn `Lecture 3 - 52:00 ` có nói phần này)

```bash
docker run -d --name bookmark --network host bookmark_service:test
```

```bash
docker run -d --name bookmark --network host --env APP_PORT=8082 -p 8082:8082 bookmark_service:test
```

5. Xóa container
```bash
docker rm -f bookmark
```

6. Tại sao không nên dùng network host mà nên dùng network bridge?
- Video `Lecture 3 - 53:15` có nói phần này

7. Cách fix lỗi container không build được Dockerfile và dừng chỗ go mod download
- Lỗi đó là DNS + IPv6 timeout khi Docker container truy cập internet.
- Cách fix:
    + B1:
        ```bash
        sudo nano /etc/docker/daemon.json
        ```
    + B2: Thêm dòng này vào và lưu
        ```json
        {
            "dns": ["8.8.8.8", "1.1.1.1"]
        }
        ```
    + B3: 
        ```bash
        sudo systemctl restart docker
        ```

### 14. Sử dụng docker-compose (Lecture 13)
1. Cài đặt docker-compose
```bash
apt install docker-compose
```
2. Viết file docker-compose.yaml:

3. Run docker-compose
```bash
docker compose up
```
- Chạy ngầm không chiểm terminal
```bash
docker compose up -d
```
- Cách xịn nhất (xóa hết và build lại image nếu code có thay đổi)
```bash
docker compose down
docker compose up --build -d
```

4. Check
```bash
docker compose ps
```
```bash
docker ps
```
5. Delete
```bash
docker compose down
```

### 15. Nginx, Load Balancer (Lecture 4)
1. Tại sao cấu hình Nginx hay listen cổng 80, 443
- Tại vì khi gõ URL trên trình duyệt thì không cần gõ port 80, 443
- Cổng 80 là mặc định cho http và cổng 443 là mặc định cho https

- Lưu ý: 
    + Trong thực tế phải tạo project riêng tên là deployment. Còn viết chung như kia cho dễ coi thôi 
    + Khi đặt REDIS_PASSWORD thì khi thao tác terminal redis phải đăng nhập trước:
        docker exec -it redis redis-cli
        auth 123456789

### 16. Lecture 4
- Muốn xem thì vào browser gõ:
    + Frontend: http://localhost/
    + Swgger: http://localhost/api/bookmark_service/swagger/index.html

- Login Dockerhub và đẩy image lên dockerhub
    + Tạo tài khoản và tạo personal access token trên trang dockerhub
    + Run
        ```
        docker login -u dongcoi14122003
        ```
        Sau đó dán token đó vào -> enter
    + Vào project sau đó build image:
        ```
        docker build -t dongcoi14122003/bookmark_service:dev .
        docker push donggcoi14122003/bookmark_service:dev
        ```
- Login vào VPS và setup trên VPS (cách cài đặt: https://docs.docker.com/engine/install/ubuntu/)
    ```bash
    ssh root@103.118.29.123

    # Run the following command to uninstall all conflicting packages
    sudo apt remove $(dpkg --get-selections docker.io docker-compose docker-compose-v2 docker-doc podman-docker containerd runc | cut -f1)

    # Install using the apt repository
    # 1. Set up Docker's apt repository.
    # Add Docker's official GPG key:
    sudo apt update
    sudo apt install ca-certificates curl
    sudo install -m 0755 -d /etc/apt/keyrings
    sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
    sudo chmod a+r /etc/apt/keyrings/docker.asc

    # Add the repository to Apt sources:
    sudo tee /etc/apt/sources.list.d/docker.sources <<EOF
    Types: deb
    URIs: https://download.docker.com/linux/ubuntu
    Suites: $(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}")
    Components: stable
    Architectures: $(dpkg --print-architecture)
    Signed-By: /etc/apt/keyrings/docker.asc
    EOF

    sudo apt update

    # 2. Install the Docker packages.
    sudo apt install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin


    apt install docker-compose
    apt install git
    git clone <địa chỉ repo project deployment>
    make up
    ```
- VM:
```
IP Address: 103.118.29.123
User: root
Password: Qwerty1234zxcv#
Port ssh: 22
```

### 17. Link API
https://ebvn.top/api/bookmark_service/swagger/index.html

### 18. Clear all data in docker
```
docker rm -f $(docker ps -aq)
docker rmi -f $(docker images -q)
docker volume rm $(docker volume ls -q)
docker network prune
```

### 19. Lecture 5
- Phần deploy nếu môi trường server có cho cài đặt `go` thì chỉ cần dùng file `ci-native.yaml` và `cd.yaml`
- Nếu môi trường server không cho cài đặt `go` thì cần dùng file `ci-pure.yaml` và `cd.yaml`
- Lưu ý khi `Merge pull request` thì chuyền sang `Squash and merge` để ví dụ nhánh đó có nhiều commit thì chỉ lấy 1 commit mới nhất của nhánh khác vào nhánh `main` thôi

- Setup self hosted runners
1. Bước 1: SSH vào VPS với quyền root
```bash
ssh root@103.118.29.123
```
2. Bước 2: Thiết lập User mới (Nếu chưa làm)
- Để bảo mật, GitHub không cho chạy Runner bằng root. Chúng ta tạo user nguyendong và cấp quyền sudo cho nó:

```bash
# Tạo user và đặt mật khẩu
adduser nguyendong

# Thêm vào nhóm sudo
usermod -aG sudo nguyendong

# Thêm vào group docker
sudo usermod -aG docker nguyendong

# Apply group (👉 hoặc logout/login lại)
newgrp docker

# Chuyển sang dùng user này luôn
su - nguyendong
```

3. Bước 3: Cài đặt Runner
- Bây giờ bạn đang ở thư mục `/home/nguyendong`, hãy bắt đầu tải Runner:
```bash
# Create a folder
mkdir actions-runner && cd actions-runner

# Download the latest runner package
curl -o actions-runner-linux-x64-2.333.1.tar.gz -L https://github.com/actions/runner/releases/download/v2.333.1/actions-runner-linux-x64-2.333.1.tar.gz

# Optional: Validate the hash
echo "18f8f68ed1892854ff2ab1bab4fcaa2f5abeedc98093b6cb13638991725cab74  actions-runner-linux-x64-2.333.1.tar.gz" | shasum -a 256 -c

# Extract the installer
tar xzf ./actions-runner-linux-x64-2.333.1.tar.gz
```

4. Bước 4: Cấu hình kết nối với GitHub (Config)
- Đây là bước quan trọng nhất. Bạn cần quay lại trình duyệt (phần Settings > Actions > Runners trên GitHub) để lấy Token mới (vì Token cũ thường hết hạn sau vài phút).

```bash
./config.sh --url https://github.com/nguyendong2003/bookmark-management --token <TOKEN_MOI_NHAT>
```

5. Bước 5: Chạy Runner
```
./run.sh
```

### 20. Cách thao tác trong thực tế và Luồng hoạt động CI, CD
1. Bước 1: Đứng ở main → tạo branch feature
```bash
# đang ở main
git checkout main

# cập nhật code mới nhất
git pull origin main

# tạo branch mới
git checkout -b feature/login
```

2. Bước 2: Code + commit + push
```bash
# đang ở feature/login
git add .
git commit -m "feat: add login feature"
git push origin feature/login
```

3. Bước 3: Tạo Pull Request (trên GitHub)
- Từ: `feature/login`
- Sang: `main`

👉 Lúc này: CI chạy (do PR vào main)

4. Bước 4: Merge PR → quay lại main
- Sau khi merge trên GitHub:
```bash
# chuyển về main
git checkout main

# kéo code mới nhất
git pull origin main
```

5. Bước 5: Tạo tag (QUAN TRỌNG)
```bash
# đang ở main
git tag v1.1.0
git push origin v1.1.0
```
- Sau bước này:

    ```
    👉 Tag được tạo từ main
    👉 CI chạy với tag
    👉 Docker image: v1.1.0
    ```
### 21. DevSecOps
- Dùng `SonarQube` để cài đặt tool thêm Security vào luồng CI-CD

