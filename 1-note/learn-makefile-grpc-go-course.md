### Hiểu file script Makefile của project `grpc-go-course`

<https://github.com/Clement-Jean/grpc-go-course>

- Bạn đang có:
    + File .proto (ví dụ greet/proto/greet.proto)
    + Một Makefile tự động:
        + Generate protobuf
        + Build server
        + Build client

- Mình sẽ giải thích script chạy như thế nào và cách generate code 👇

1️⃣ Cấu trúc project của bạn (suy ra từ Makefile)

- Ví dụ:
```
greet/
 ├── proto/
 │    └── greet.proto
 ├── server/
 │    └── main.go
 └── client/
      └── main.go

bin/
```

2️⃣ Cách generate code

- Chỉ cần chạy:
```bash
make greet
```

- Hoặc:
```bash
make all
```

- Hoặc generate tất cả project:
```
make
```

- Vì:
```Makefile
.DEFAULT_GOAL := help
```

Nếu không chỉ định, nó hiện help.

Muốn build thì phải gọi target cụ thể như `make greet`.

3️⃣ Điều gì xảy ra khi chạy `make greet`?

- Makefile có đoạn:
```Makefile
$(project):
	protoc -I$@/${PROTO_DIR} \
	--go_opt=module=${PACKAGE} \
	--go_out=. \
	--go-grpc_opt=module=${PACKAGE} \
	--go-grpc_out=. \
	$@/${PROTO_DIR}/*.proto

```

- Với `greet` thì:
```
$@ = greet
```

- Nên lệnh thực tế sẽ là:
```
protoc -Igreet/proto \
  --go_opt=module=github.com/Clement-Jean/grpc-go-course \
  --go_out=. \
  --go-grpc_opt=module=github.com/Clement-Jean/grpc-go-course \
  --go-grpc_out=. \
  greet/proto/*.proto
```

4️⃣ Code được generate ở đâu?

- Vì bạn dùng:
```bash
--go_out=.
--go-grpc_out=.
```

- Nó sẽ generate file dựa theo:
```
option go_package = "github.com/Clement-Jean/grpc-go-course/greet/proto";
```
👉 Kết quả:
```
greet/proto/greet.pb.go
greet/proto/greet_grpc.pb.go
```

5️⃣ Sau khi generate xong

- Makefile tiếp tục:
```bash
go build -o bin/greet/server ./greet/server
go build -o bin/greet/client ./greet/client
```

- Nó build:
```
bin/greet/server
bin/greet/client
```

6️⃣ Điều kiện để script chạy được

- Bạn cần:

✅ Cài protoc
```bash
protoc --version
```

✅ Cài plugin Go
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Và đảm bảo:
```bash
$GOPATH/bin
```
đã nằm trong `PATH`.

- Kiểm tra:
```bash
which protoc-gen-go
```

7️⃣ Flow đầy đủ khi chạy `make greet`
```
make greet
   ↓
CHECK_DIR_CMD (check folder tồn tại)
   ↓
protoc generate pb.go
   ↓
go build server
   ↓
go build client
   ↓
Output ở bin/
```

8️⃣ Nếu bạn chỉ muốn generate code (không build)

- Bạn có thể sửa Makefile thêm target:
```Makefile
proto:
	protoc -Igreet/proto \
	--go_out=. \
	--go-grpc_out=. \
	greet/proto/*.proto
```

- Rồi chạy:
```bash
make proto
```

9️⃣ Tại sao dùng `--go_opt=module=${PACKAGE}` ?

- Biến:
```Makefile
PACKAGE = $(shell head -1 go.mod | awk '{print $$2}')
```

- Nó đọc:
```go
module github.com/Clement-Jean/grpc-go-course
```

- Rồi truyền vào protoc.

👉 Mục đích: giúp import path đúng theo module Go.

Nếu không có flag này, đôi khi sẽ bị lỗi import sai module.

🔟 Tóm tắt ngắn gọn

| Bạn muốn                | Lệnh           |
| ----------------------- | -------------- |
| Generate + build greet  | `make greet`   |
| Generate + build tất cả | `make all`     |
| Clean                   | `make clean`   |
| Rebuild                 | `make rebuild` |
| Xem info                | `make about`   |
