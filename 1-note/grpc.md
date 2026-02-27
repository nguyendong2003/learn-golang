# GRPC

[Installation]

<https://protobuf.dev/installation/>

<https://grpc.io/docs/languages/go/quickstart/>

[Fix error]

```bash
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.bashrc
source ~/.bashrc
```

## Keyword

### 1. GRPC

1️⃣ RPC là gì?
- `RPC` (Remote Procedure Call) cho phép bạn gọi một hàm ở máy khác giống như gọi hàm nội bộ (local).

- `gRPC` (Google Remote Procedure Call) là một framework RPC (Remote Procedure Call) mã nguồn mở do Google phát triển. Nó cho phép các service giao tiếp với nhau một cách nhanh, hiệu quả và đa ngôn ngữ, thường được dùng trong microservices và hệ thống phân tán.

- gRPC sử dụng:
    + `HTTP/2` làm giao thức truyền tải
    + `Protocol Buffers` (Protobuf) làm cơ chế tuần tự hóa dữ liệu

2️⃣ Kiến trúc của gRPC
- Thành phần chính:
    + Client
    + Server
    + Service definition (`.proto` file)

- File `.proto` định nghĩa:
    + Cấu trúc dữ liệu
    + Các service
    + Các method có thể gọi

- Ví dụ:
    ```proto
    syntax = "proto3";

    service UserService {
        rpc GetUser (UserRequest) returns (UserResponse);
    }

    message UserRequest {
        int32 id = 1;
    }

    message UserResponse {
        string name = 1;
        int32 age = 2;
    }
    ```

- Sau đó dùng `protoc` để generate code cho các ngôn ngữ như:
    + Java
    + Go
    + Python
    + Node.js
    + C#
    + C++

3️⃣ gRPC hoạt động như thế nào?

- gRPC sử dụng:
    + `HTTP/2`
        * Ưu điểm:
            + Multiplexing (nhiều request song song trên một kết nối)
            + Header compression
            + Server push
            + Binary framing

    + `Protocol Buffers (Protobuf)`
        * Ưu điểm:
            + Nhẹ hơn JSON
            + Nhanh hơn JSON
            + Typed schema rõ ràng
            + Backward compatible tốt

- So với REST + JSON:

    | Tiêu chí    | REST       | gRPC              |
    | ----------- | ---------- | ----------------- |
    | Data format | JSON       | Protobuf (binary) |
    | Transport   | HTTP/1.1   | HTTP/2            |
    | Performance | Trung bình | Rất cao           |
    | Streaming   | Khó        | Native support    |

4️⃣ 4 loại gRPC Communication

- gRPC hỗ trợ 4 kiểu giao tiếp:

    1️⃣ Unary RPC

    Client gửi 1 request → Server trả 1 response (Giống REST API)

    2️⃣ Server Streaming

    Client gửi 1 request → Server trả nhiều response stream

    3️⃣ Client Streaming

    Client gửi nhiều request → Server trả 1 response

    4️⃣ Bidirectional Streaming

    Client và Server gửi stream 2 chiều độc lập

    Đây là điểm cực mạnh của gRPC so với REST.

5️⃣ HTTP/1.1 vs HTTP/2 khác gì?

🚨 HTTP/1.1 vấn đề:
- 1 connection → 1 request tại 1 thời điểm
- Head-of-line blocking
- Text-based
- Header lớn

⚡ HTTP/2 cải tiến:
- Multiplexing (nhiều request cùng lúc trên 1 connection)
- Binary protocol
- Header compression
- Server push

👉 Đây là một lý do lớn khiến gRPC nhanh.

6️⃣ Khi nào dùng REST? Khi nào dùng gRPC?
- Dùng REST khi:
    + Public API cho frontend/browser
    + Cần debug dễ (curl, postman)
    + Cần compatibility rộng

- Dùng gRPC khi:
    + Microservices giao tiếp nội bộ
    + Cần performance cao
    + Streaming
    + Real-time
    + Internal system
    + Google, Netflix, Uber dùng gRPC cho internal microservices.

7️⃣ Scalability in gRPC
* `Server: Async`
    - Ý nghĩa:
        + gRPC server xử lý request bất đồng bộ
        + Mỗi RPC được xử lý song song (concurrent)
        + Không block toàn bộ server khi một request đang chạy
    - Trong Go:
        + Mỗi RPC → một goroutine riêng
        + Có thể xử lý hàng nghìn request cùng lúc

    👉 Đây là nền tảng giúp gRPC scale tốt.

* `Client: Async or Blocking`

    gRPC client có 2 kiểu gọi:

    - Blocking (Unary sync call)

        ```go
        res, err := client.GetUser(ctx, req)
        ```

        👉 Client sẽ chờ response.

    - Async (Streaming / Goroutine)

        + Client có thể:
            + Gửi nhiều request song song
            + Dùng goroutine
            + Dùng streaming không chặn
        
        👉 Điều này cho phép client tạo rất nhiều concurrent RPC.

### 2. Protoc command
- Để generate code từ file `.proto`, bạn sử dụng `protoc` (Protocol Buffers Compiler).

1️⃣ Cài đặt cần có
- Bạn cần:
    + protoc compiler
    + Plugin theo từng ngôn ngữ (grpc plugin)

- Kiểm tra protoc:
```bash
protoc --version
```

2️⃣ Cấu trúc lệnh cơ bản
```bash
protoc -I=<proto_path> --<lang>_out=<output_path> <file.proto>
```

- Nếu có gRPC:
```bash
protoc -I=. \
  --go_out=. \
  --go-grpc_out=. \
  user.proto
```

- Nếu muốn chỉ định module:
```bash
protoc -I=. \
  --go_out=paths=source_relative:. \
  --go-grpc_out=paths=source_relative:. \
  user.proto
```

- Generate nhiều file proto:
```bash
protoc -I=. \
  --go_out=. \
  --go-grpc_out=. \
  *.proto
```

- Thường trong project thực tế sẽ dùng:
```bash
protoc -I=. \
  --go_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_out=. \
  --go-grpc_opt=paths=source_relative \
  user.proto
```

6️⃣ Giải thích các flag quan trọng

| Flag                                  | Ý nghĩa                                                  |
| ------------------------------------- | -------------------------------------------------------- |
| `-I` hoặc `--proto_path`              | Chỉ định thư mục chứa file `.proto` để resolve import    |
| `--go_out`                            | Generate protobuf message code (`*.pb.go`)               |
| `--go_opt`                            | Truyền option cho `protoc-gen-go`                        |
| `--go-grpc_out`                       | Generate service gRPC code (`*_grpc.pb.go`)              |
| `--go-grpc_opt`                       | Truyền option cho `protoc-gen-go-grpc`                   |
| `paths=source_relative`               | Output file nằm cùng thư mục với file `.proto`           |
| `module=example.com/app`              | Chỉ định Go module khi generate                          |
| `Mfile.proto=go/package/path`         | Map 1 proto import sang Go package cụ thể                |
| `require_unimplemented_servers=false` | Không bắt buộc embed `UnimplementedXXXServer`            |
| `--plugin=protoc-gen-xxx=PATH`        | Chỉ định đường dẫn plugin nếu không nằm trong PATH       |
| `--descriptor_set_out=out.pb`         | Xuất descriptor file (dùng cho reflection/tooling)       |
| `--include_imports`                   | Include luôn các proto được import khi export descriptor |
| `--include_source_info`               | Include thông tin source (phục vụ tooling, doc gen)      |


7️⃣ Ví dụ hoàn chỉnh (Go project structure)
```code
proto/
  user.proto
```

- Generate:
```bash
protoc -I=proto \
  --go_out=proto \
  --go-grpc_out=proto \
  proto/user.proto
```

- Kết quả:
```code
proto/
  user.pb.go
  user_grpc.pb.go
```

9️⃣ Best practice trong project thực tế
- Thường người ta:
    + Tạo script generate.sh hoặc Makefile
    + Hoặc dùng Docker để đảm bảo version protoc đồng nhất

- Ví dụ Makefile:
```makefile
proto:
	protoc -I=. \
	--go_out=. \
	--go-grpc_out=. \
	proto/*.proto
```