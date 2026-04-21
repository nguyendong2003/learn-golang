# Javascript

## Kiến thức NestJS

### 1. @Module

- Trong NestJS, `@Module` chỉ làm 1 việc: quản lý dependency (ai dùng được ai)
- `@Module` bao gồm:

  | key         | câu hỏi nó trả lời                 |
  | ----------- | ---------------------------------- |
  | providers   | Module này **tạo ra cái gì?**      |
  | controllers | Module này **có API nào?**         |
  | imports     | Module này **muốn dùng ai?**       |
  | exports     | Module này **cho ai dùng cái gì?** |

- Cấu trúc cơ bản:

```ts
@Module({
  imports: [],
  controllers: [],
  providers: [],
  exports: [],
})
export class UserModule {}
```

🧩 `providers`

```ts
providers: [UserService];
```

👉 Ý nghĩa:

- Khai báo những class được tạo ra trong module
- Thường là:
  - Service
  - Repository
  - Helper

👉 Quan trọng:

- NestJS sẽ quản lý và inject các provider này
- Mặc định: chỉ dùng trong module này

🎮 `controllers`

```ts
controllers: [UserController];
```

👉 Ý nghĩa:

- Khai báo các controller xử lý request

👉 Controller dùng để:

- Nhận request từ client
- Gọi service xử lý

👉 Ví dụ:

```ts
@Get('/users')
```

🔗 `imports`

```ts
imports: [UserModule];
```

👉 Ý nghĩa:

- Module này muốn dùng tài nguyên từ module khác

⚠️ Lưu ý:

- Chỉ import thôi chưa đủ
- Phải kết hợp với exports bên module kia

📤 `exports`

```ts
exports: [UserService];
```

👉 Ý nghĩa:

- Cho phép module khác sử dụng provider của mình

👉 Nếu không export:

- Module khác không thể inject được

- Ví dụ thực tế:

```ts
// src/user/user.service.ts
import { Injectable } from "@nestjs/common";

@Injectable()
export class UserService {
  getUser() {
    return { id: 1, name: "Dong" };
  }
}

// src/user/user.module.ts
import { Module } from "@nestjs/common";
import { UserService } from "./user.service";

@Module({
  providers: [UserService],
  exports: [UserService], // 👈 cho module khác dùng
})
export class UserModule {}

// src/auth/auth.service.ts
import { Injectable } from "@nestjs/common";
import { UserService } from "../user/user.service";

@Injectable()
export class AuthService {
  constructor(private userService: UserService) {}

  login() {
    const user = this.userService.getUser();
    return {
      message: "Login success",
      user,
    };
  }
}

// src/auth/auth.module.ts
import { Module } from "@nestjs/common";
import { AuthService } from "./auth.service";
import { UserModule } from "../user/user.module";

@Module({
  imports: [UserModule], // 👈 import để dùng
  providers: [AuthService],
})
export class AuthModule {}

// src/app.module.ts
import { Module } from "@nestjs/common";
import { AuthModule } from "./auth/auth.module";

@Module({
  imports: [AuthModule],
})
export class AppModule {}

// src/main.ts
import { NestFactory } from "@nestjs/core";
import { AppModule } from "./app.module";
import { AuthService } from "./auth/auth.service";

async function bootstrap() {
  const app = await NestFactory.createApplicationContext(AppModule);

  // 👇 test thử DI
  const authService = app.get(AuthService);

  const result = authService.login();
  console.log(result);
}

bootstrap();
```

### 2. Project structure

```bash
src/
├── modules/
│   ├── user/
│   │   ├── user.module.ts
│   │   ├── user.controller.ts
│   │   ├── user.service.ts
│   │   ├── dto/
│   │   │    ├── create-user.dto.ts
│   │   │    └── update-user.dto.ts
│   │   ├── entities/
│   │   │    └── user.entity.ts
│   │
│   ├── auth/
│   │   ├── auth.module.ts
│   │   ├── auth.controller.ts
│   │   ├── auth.service.ts
│   │
│   └── ...
│
├── common/
│   ├── guards/
│   ├── interceptors/
│   ├── pipes/
│   ├── decorators/
│   ├── filters/
│
├── prisma/
│   ├── prisma.service.ts
│   └── prisma.module.ts
│
├── config/
│   ├── env.config.ts
│   └── app.config.ts
│
├── main.ts
└── app.module.ts
```

## Install

### 1. Install NVM + Nodejs (https://nodejs.org/en/download)

```bash
# Download and install nvm:
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.4/install.sh | bash
# in lieu of restarting the shell
\. "$HOME/.nvm/nvm.sh"
# Download and install Node.js:
nvm install 24
# Verify the Node.js version:
node -v # Should print "v24.15.0".
# Verify npm version:
npm -v # Should print "11.12.1".
```

### 2. Install nestjs + Setup project (https://docs.nestjs.com/first-steps)

```bash
npm i -g @nestjs/cli
nest new project-name

cd project-name
npm i --save-dev @types/jest
```

### 3. Init module, service, controller

```bash
nest
nest g mo posts
nest g s posts
nest g co posts
```

### 4. Cài đặt extention Prima trên vscode

### 5. Install prisma on project (https://docs.nestjs.com/recipes/prisma)

```bash
npm install prisma --save-dev
npx prisma
npx prisma init
```

### 6. Apply current schema change to database

```bash
npx prisma db push
```

### 7. Xem database

```bash
npx prisma studio
```

- Sau đó vào trang `http://localhost:51212` để xem database

### 8. Install and generate Prisma Client (https://docs.nestjs.com/recipes/prisma)

```bash
npm install @prisma/client
npx prisma generate
```

### 9. Install driver for SQLite

```bash
npm install @prisma/adapter-better-sqlite3
```

### 10

```
cd crud
nest g mo shared
```

- Sau câu lệnh trên sẽ tạo ra 1 thư mục `shared` trong thư mục `src` và gen ra file `shared.module.ts` trong thư mục `shared` và tự import `SharedModule` vào file `app.module.ts`

```bash
cd crud/src/shared/services
nest g s prisma --flat --no-spec
```

- Sau câu lệnh này sẽ tạo ra file `shared.module.ts` trong thư mục `crud/src/shared/services`

### 11. Câu lệnh migrate database

- Xem trong file `crud/prisma-cli.md`

### 12. Validation trong nestjs (https://docs.nestjs.com/techniques/validation)

```bash
npm i --save class-validator class-transformer
```

- Validation trong nestjs dùng thư viện (`https://github.com/typestack/class-validator`)

- Thêm `app.useGlobalPipes(new ValidationPipe())` vào `main.ts` để kích hoạt auto validation

### 13. Hash password

- Sau khi chạy câu lệnh dưới thì tạo ra file `hashing.service.ts`

```bash
nest g s shared/services/hashing --flat --no-spec
```

- Cài đặt thư viện `bcrypt`

```bash
npm i bcrypt
npm i @types/bcrypt -D
```

- Xem `prisma error code`: `https://www.prisma.io/docs/orm/reference/error-reference`

### 14. Serialization (https://docs.nestjs.com/techniques/serialization)

- Thêm code này vào `app.module.ts` để serialize global

```ts
  providers: [
    AppService,
    {
      provide: APP_INTERCEPTOR,
      useClass: ClassSerializerInterceptor,
    },
  ],
```

### 15. Interceptors (https://docs.nestjs.com/interceptors)

- Thêm code này vào `main.ts` để chèn thêm `interceptor` ở global

```ts
app.useGlobalInterceptors(new LoggingInterceptor());
```

### 16. JWT (https://docs.nestjs.com/security/authentication#jwt-token)

```bash
npm install --save @nestjs/jwt
```

```bash
nest g s shared/services/token --flat --no-spec
```

### 17. Guards (https://docs.nestjs.com/guards)

- Trong NestJS, Guard là một cơ chế dùng để kiểm soát việc truy cập (authorization) vào route (API endpoint).
- Guard quyết định request có được phép đi vào controller hay không.

1. 🧠 Guard dùng để làm gì?

- Guard chạy trước controller, sau middleware. Nó thường dùng để:
  - Kiểm tra authentication (đã đăng nhập chưa)
  - Kiểm tra authorization (có quyền không, ví dụ: admin vs user)
  - Chặn request nếu không hợp lệ

### 18. Authentication (https://docs.nestjs.com/security/authentication)
