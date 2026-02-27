# Proto Course

[Complete Guide to Protocol Buffers 3 [Java, Golang, Python]]
<https://www.udemy.com/course/protocol-buffers/>

<https://github.com/Clement-Jean/proto-course>

<https://github.com/Clement-Jean/proto-go-course>

<https://github.com/Clement-Jean/proto-java-course>

<https://github.com/Clement-Jean/proto-python-course>

<https://github.com/googleapis/googleapis/blob/master/google/type/dayofweek.proto>

## Keyword
### 1. `Defaults`:
- When not set:
    + A field will not be serialized
    + Populated with default value

### 2. `Scalar Types`
- `Number`:
    + Keyword:
        `int32, int64, sint32, sint64,
        uint32, unint64,
        fixed32, fixed64, sfixed32, sfixed64
        float, double`

    + Default value: 0

- `Boolean`:
    + Keyword: `bool`
    + Value: true or false
    + Default: false

- `String`
    + Keyword: `string`
    + Value: arbitrary length text
    + Default: empty string

    + An important thing: string only accepts `UTF-8 encoded` or `7-bit ASCII`

- `Bytes`
    + Keyword: `bytes`
    + Value: arbitrary length byte sequence
    + Default: empty bytes
    
    + An important thing: Up to you to interpret in code (Việc hiểu và xử lý nó trong code là tùy bạn quyết định.)

### 3. `Tags`
* Serialization (tuần tự hóa) là quá trình: Chuyển dữ liệu từ object / message → thành dạng nhị phân (binary) để gửi đi hoặc lưu trữ.

    + Ví dụ: Proto message:
    ```proto
    message User {
        int32 id = 1;
        string name = 2;
    }
    ```

    + Trong code bạn có:
    ```go
    user = User(id=1, name="An")
    ```
    👉 Khi serialize → nó biến thành chuỗi byte (binary data) như:
    ```go
    0x08 0x01 0x12 0x02 0x41 0x6E
    ```

    + Dạng này:

        Nhẹ hơn JSON

        Gửi qua network nhanh hơn

        Phù hợp cho microservices, gRPC

* Deserialization (giải tuần tự hóa) là quá trình ngược lại: Chuyển dữ liệu binary → thành object / message để sử dụng trong chương trình.
    + Ví dụ: Nhận được binary từ server → Proto sẽ parse lại thành: `User(id=1, name="An")`


* `tag` (hay còn gọi là field number) là số định danh duy nhất (ID) của mỗi field trong message.

1️⃣ Tag là gì?
- Ví dụ:
    ```proto
    message Person {
        string name = 1;
        int32 age = 2;
    }
    ```
- Ở đây:
    ```
    name có tag = 1
    age có tag = 2
    ```
    👉 Số 1 và 2 chính là tag.

2️⃣ Vì sao tag quan trọng?
- Khi Protobuf serialize dữ liệu, nó KHÔNG lưu tên field (name, age)
- Nó chỉ lưu: `(tag + wire type) + value`
- Nên khi deserialize, Protobuf dựa vào tag để biết dữ liệu thuộc field nào.

3️⃣ Tag hoạt động như thế nào?
- Ví dụ gửi:
    ```JSON
    {
        "name": "An",
        "age": 25
    }
    ```

- Protobuf sẽ encode kiểu như:
    ```go
    1: "An"
    2: 25
    ```

- Khi nhận dữ liệu:
    ```
    Thấy tag 1 → biết đó là name
    Thấy tag 2 → biết đó là age
    ```

4️⃣ Quy tắc của tag
- Phải duy nhất trong message
    + Không được:
        ```proto
        string name = 1;
        int32 age = 1; // ❌ sai
        ```

- Không nên đổi tag sau khi đã deploy
    + Ví dụ:
        + Version 1: 
            ```proto
            string name = 1;
            ```    

        + Version 2:
            ```proto
            int32 age = 1; // ❌ cực kỳ nguy hiểm
            ```
    👉 Client cũ sẽ đọc sai dữ liệu.

- Không dùng lại tag đã xóa
    + Đúng cách:
        ```proto
        message Person {
            string name = 1;
            reserved 2;
        }
        ```
- Tag nhỏ (1–15) tối ưu hơn
    + Tag 1–15 → chiếm 1 byte
    + Tag 16–2047 → chiếm 2 byte

    👉 Field dùng thường xuyên nên để tag nhỏ.


* `Tags`:
- These tags what make serialization and deserialization possible
- Field names are not important (for serialization)
- The rules:
    + Smallest tags: 1
    + Largest tag: 536870911
    + Reserved tags: 19000 to 19999

### 4. `Repeated Fields`

1️⃣ repeated là gì?
- repeated dùng để khai báo một field có thể xuất hiện nhiều lần (giống như array / list).

- Ví dụ:
```proto
message Person {
  string name = 1;
  repeated string phone_numbers = 2;
}
```

- Ở đây:
    + `phone_numbers` là một danh sách string
    + Có thể có 0, 1, hoặc nhiều phần tử

2️⃣ Ví dụ dữ liệu
- Dữ liệu dạng JSON tương đương:
    ```JSON
    {
        "name": "An",
        "phone_numbers": ["0901", "0902", "0903"]
    }
    ```

3️⃣ Cách Protobuf encode repeated
- Có 2 cách encode tùy theo kiểu dữ liệu:

🔹 Trường hợp 1: repeated của type không phải số (string, message...)

Ví dụ:
```proto
repeated string phone_numbers = 2;
```

👉 Protobuf sẽ encode mỗi phần tử như một field riêng biệt:
```go
2: "0901"
2: "0902"
2: "0903"
```

👉 Cùng tag lặp lại nhiều lần.

🔹 Trường hợp 2: repeated của số (int32, int64, bool...)

Ví dụ:
```proto
repeated int32 scores = 3;
```

Mặc định (proto3), nó dùng packed encoding:
```go
3: [10, 20, 30]
```

Thay vì:
```go
3: 10
3: 20
3: 30
```
👉 Packed giúp tiết kiệm byte.

4️⃣ Packed là gì?

- Với kiểu số:
    + Tất cả value được gói lại thành 1 block
    + Chỉ xuất hiện tag một lần

- Ví dụ:
    ```proto
    repeated int32 scores = 1;
    ```

    + Binary sẽ giống như:
        ```bash
        (tag=1, wire_type=length-delimited)
        [length]
        [value1][value2][value3]
        ```

5️⃣ So sánh repeated vs optional

| Keyword  | Số lượng giá trị |
| -------- | ---------------- |
| optional | 0 hoặc 1         |
| repeated | 0 đến N          |

6️⃣ repeated message
- Bạn cũng có thể lặp lại message:
```proto
message Phone {
  string number = 1;
  string type = 2;
}

message Person {
  repeated Phone phones = 3;
}
```

- JSON tương đương:
```JSON
{
  "phones": [
    { "number": "0901", "type": "home" },
    { "number": "0902", "type": "work" }
  ]
}
```

7️⃣ Lưu ý quan trọng

🔹 repeated luôn có default là rỗng (empty list)
- Không bao giờ là null.

🔹 Thứ tự phần tử được giữ nguyên
- Protobuf giữ thứ tự khi serialize/deserialize.

🔹 Không thể phân biệt:
- Field không gửi
- Field gửi nhưng rỗng

👉 Vì cả hai đều thành `[]`

### 5. Enum
1️⃣ Enum là gì?
- enum dùng để định nghĩa tập giá trị cố định dạng số nguyên có tên.
- Ví dụ:
    ```proto
        enum Status {
            STATUS_UNKNOWN = 0;
            ACTIVE = 1;
            INACTIVE = 2;
        }
    ```
    👉 Bên trong Protobuf, enum thực chất là int32

    👉 Nhưng khi code, bạn dùng tên (`ACTIVE`) thay vì số (`1`)

2️⃣ Dùng enum trong message
```proto
message User {
  string name = 1;
  Status status = 2;
}
```

- JSON representation (khi dùng protobuf JSON mapping):
```json
{
  "name": "An",
  "status": "ACTIVE"
}
```

- Nhưng khi serialize binary, nó lưu:
```bash
2: 1
```
👉 1 là giá trị của ACTIVE

3️⃣ Quy tắc bắt buộc trong proto3

🔹 Enum phải có giá trị đầu tiên = 0

- Ví dụ đúng:
```proto
enum Status {
  STATUS_UNKNOWN = 0;
  ACTIVE = 1;
}
```

- Sai:
```proto
enum Status {
  ACTIVE = 1; // ❌ proto3 không cho
}
```

📌 Lý do:
- Proto3 mặc định field số = 0 nếu không set
- Nên phải có giá trị đại diện cho "default"

4️⃣ Default value của enum
- Nếu không set status, thì:
```proto
status = 0
```

→ tức là `STATUS_UNKNOWN`

- Vì vậy best practice:
```proto
UNKNOWN = 0;
```

* `Enum`
- Keyword: `enum`
- Default value: first value
- An important thing:
    + First tag: 0

### 6. Comment

### Practice 1

<https://github.com/googleapis/googleapis/blob/master/google/type/dayofweek.proto>

### 7. Nested message
1️⃣ Nested message là gì?
- Là message được định nghĩa bên trong một message khác.
- Ví dụ:
```proto
message Person {
  string name = 1;

  message Address {
    string street = 1;
    string city = 2;
  }

  Address address = 2;
}
```

- Ở đây:
    + `Address` chỉ tồn tại trong phạm vi `Person`
    + Tên đầy đủ của nó là: `Person.Address`

2️⃣ Tại sao cần nested message?

🔹 1. Gom nhóm logic
- Nếu Address chỉ dùng cho Person, không cần khai báo global.

🔹 2. Tránh xung đột tên
- Bạn có thể có nhiều Address khác nhau:

```proto
message Company {
  message Address {
    string office = 1;
  }
}
```

- Lúc này có:
    + `Person.Address`
    + `Company.Address`

    👉 Không bị trùng.

3️⃣ JSON representation

- Với schema:
```proto
message Person {
  string name = 1;

  message Address {
    string city = 1;
  }

  Address address = 2;
}
```

- JSON sẽ là:
```json
{
  "name": "An",
  "address": {
    "city": "Hanoi"
  }
}
```
👉 Nested chỉ ảnh hưởng scope trong `.proto`

👉 JSON không thể hiện cấu trúc lồng của type

### 8. Import
1️⃣ import là gì?
- `import` dùng để sử dụng message / enum / service được định nghĩa ở file `.proto` khác.

- Ví dụ:

📄 common.proto:
```proto
syntax = "proto3";

package common;

message Address {
  string street = 1;
  string city = 2;
}
```

📄 user.proto:
```proto
syntax = "proto3";

package user;

import "common.proto";

message User {
  string name = 1;
  common.Address address = 2;
}
```

👉 User đang dùng Address từ file khác.


3️⃣ Phải dùng tên đầy đủ (fully-qualified name)
- Nếu có package, bạn phải gọi:
```proto
common.Address
```

- Không phải:
```proto
Address  // ❌ sai nếu khác package
```

4️⃣ Các loại import

🔹 1. Import thường
```proto
import "common.proto";
```
→ Dùng bình thường

🔹 2. public import
```proto
import public "common.proto";
```
- Ý nghĩa:

Nếu A import public B
Và C import A
→ C cũng thấy được B

- Ví dụ:
```
C → A → B
```
Nếu A dùng import public → C không cần import B trực tiếp

🔹 3. weak import (ít dùng)
```proto
import weak "old.proto";
```

- Nếu file không tồn tại vẫn compile được.
- Hiếm khi dùng trong production hiện đại.

### 9. Package
1️⃣ package là gì?

- `package` dùng để định nghĩa namespace cho các message / enum / service trong file `.proto`.

- Nếu trong file .proto không ghi package, thì: Tất cả type sẽ nằm trong global namespace

- Ví dụ:
```proto
syntax = "proto3";

package company.user;

message User {
  string name = 1;
}
```

👉 Tên đầy đủ (fully-qualified name) của User là:
```proto
company.user.User
```

2️⃣ Package dùng để làm gì?

🔹 1. Tránh trùng tên

- Bạn có thể có:
```proto
package company.user;
message Request {}
```
và
```proto
package company.payment;
message Request {}
```
👉 Hai Request không bị conflict vì khác namespace.

🔹 2. Quy định namespace khi generate code

- Tùy ngôn ngữ:

    + Go
    ```proto
    package company.user;
    ```
    → Có thể sinh ra package Go tương ứng (kết hợp với `option go_package`).

    + Java
        + Bạn thường thêm:
        ```proto
        option java_package = "com.company.user";
        ```
        + Nếu không, nó sẽ dùng `package` trong proto.

3️⃣ Package KHÔNG ảnh hưởng wire format

- Protobuf khi serialize chỉ lưu:
```bash
tag + wire_type + value
```

- Nó không lưu package name.

- Vì vậy:

👉 Đổi package không làm hỏng dữ liệu binary

👉 Nhưng có thể làm hỏng compile code

### 10. Note
- Backward compatible vs Forward compatible
- `Style Guide`: <https://protobuf.dev/programming-guides/style/>
- `Updating`: <https://protobuf.dev/programming-guides/proto3/#updating>
- Xem best practice viết protobuf: <https://github.com/protocolbuffers/protobuf/blob/main/src/google/protobuf/descriptor.proto>

### Practice 2



## Công cụ
- Dưới đây là tác dụng của từng công cụ trong workflow `Protobuf + Go`:

1️⃣ `protoc` – Trình biên dịch của Protocol Buffers

- Tác dụng chính:

    + Đọc file .proto
    + Phân tích cấu trúc message / service
    + Gọi các plugin để generate code tương ứng với từng ngôn ngữ

    👉 Hiểu đơn giản: `protoc` là compiler trung tâm, còn việc sinh code Go, Java, Python… là do `plugin` đảm nhiệm.

- Ví dụ:
```bash
protoc user.proto
```

→ Không có plugin thì nó không biết sinh ra code gì.

2️⃣ `protoc-gen-go` – Plugin sinh code Go cho message
- Cài bằng:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

- Tác dụng:
    + Convert message trong .proto
    + Sinh ra struct Go tương ứng
    + Sinh các hàm:
        + Marshal
        + Unmarshal
        + Size
        + ProtoReflect
        + Getter

- Ví dụ:
```proto
message User {
  string id = 1;
}
```

Sinh ra Go:
```go
type User struct {
    Id string
}
```

👉 Dùng khi bạn chỉ cần:
+ Serialize / Deserialize
+ Gửi message qua Kafka
+ Lưu file binary
+ Truyền data giữa service

3️⃣ `protoc-gen-go-grpc` – Plugin sinh code gRPC
- Cài bằng:
```bash
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

- Tác dụng:
    + Sinh interface server
    + Sinh client stub
    + Sinh handler đăng ký service

- Ví dụ:
```proto
service UserService {
  rpc GetUser (UserRequest) returns (User);
}
```

Nó sẽ sinh:

```go
type UserServiceServer interface {
    GetUser(context.Context, *UserRequest) (*User, error)
}

Và:

type UserServiceClient interface {
    GetUser(ctx context.Context, in *UserRequest, opts ...grpc.CallOption) (*User, error)
}
```

👉 Dùng khi:

+ Bạn build microservice
+ Làm RPC nội bộ giữa service
+ Xây hệ thống backend phân tán