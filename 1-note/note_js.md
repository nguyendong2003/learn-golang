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

### 19. Giới thiệu dự án

- Schema: https://dbdiagram.io/d/Ecom-67ae1af3263d6cf9a013ccd9

```ts
Table Language {
  id Int [pk, increment]
  name String
  code String [unique]

  createdById Int [ref: > User.id, null]
  updatedById Int [ref: > User.id, null]

  deletedAt DateTime [null]
  createdAt DateTime [default: `now()`]
  updatedAt DateTime
}

Table User {
  id Int [pk, increment]
  email String [unique]
  name String
  password String
  phoneNumber String
  avatar String [null]
  totpSecret String [null]
  status UserStatus [default: 'ACTIVE']
  roleId Int [ref: > Role.id]

  createdById Int [ref: > User.id, null]
  updatedById Int [ref: > User.id, null]

  deletedAt DateTime [null]
  createdAt DateTime [default: `now()`]
  updatedAt DateTime
}

Table UserTranslation {
  id Int [pk, increment]
  userId Int [ref: > User.id]
  languageId Int [ref: > Language.id]
  address String [null]
  description String [null]

  createdById Int [ref: > User.id, null]
  updatedById Int [ref: > User.id, null]

  deletedAt DateTime [null]
  createdAt DateTime [default: `now()`]
  updatedAt DateTime
}

Table VerificationCode {
  id Int [pk, increment]
  email String
  code String
  type VerificationCodeType

  expiresAt DateTime
  createdAt DateTime [default: `now()`]

  indexes {
    (email, code, type)
    expiresAt
  }
}

Table RefreshToken {
  token String [unique]
  userId Int [ref: > User.id]

  expiresAt DateTime
  createdAt DateTime [default: `now()`]

  indexes {
    expiresAt
  }
}

Table Permission {
  id Int [pk, increment]
  name String
  description String
  path String
  method HTTPMethod

  createdById Int [ref: > User.id, null]
  updatedById Int [ref: > User.id, null]

  deletedAt DateTime [null]
  createdAt DateTime [default: `now()`]
  updatedAt DateTime
}

Table Role {
  id Int [pk, increment]
  name String [unique]
  description String
  isActive Boolean [default: true]

  createdById Int [ref: > User.id, null]
  updatedById Int [ref: > User.id, null]

  deletedAt DateTime [null]
  createdAt DateTime [default: `now()`]
  updatedAt DateTime
}

Table Product {
  id Int [pk, increment]
  base_price Float
  virtual_price Float
  brandId Int [ref: > Brand.id]
  images String[]

  createdById Int [ref: > User.id, null]
  updatedById Int [ref: > User.id, null]

  deletedAt DateTime [null]
  createdAt DateTime [default: `now()`]
  updatedAt DateTime
}

Table ProductTranslation {
  id Int [pk, increment]
  productId Int [ref: > Product.id]
  languageId Int [ref: > Language.id]
  name String
  description String

  createdById Int [ref: > User.id, null]
  updatedById Int [ref: > User.id, null]

  deletedAt DateTime [null]
  createdAt DateTime [default: `now()`]
  updatedAt DateTime
}

Table Category {
  id Int [pk, increment]
  parentCategoryId Int [ref: > Category.id, null]

  createdById Int [ref: > User.id, null]
  updatedById Int [ref: > User.id, null]

  deletedAt DateTime [null]
  createdAt DateTime [default: `now()`]
  updatedAt DateTime
}

Table CategoryTranslation {
  id Int [pk, increment]
  categoryId Int [ref: > Category.id]
  languageId Int [ref: > Language.id]
  name String
  description String

  createdById Int [ref: > User.id, null]
  updatedById Int [ref: > User.id, null]

  deletedAt DateTime [null]
  createdAt DateTime [default: `now()`]
  updatedAt DateTime
}

Table Variant {
  id Int [pk, increment]
  name String
  productId Int [ref: > Product.id]

  createdById Int [ref: > User.id, null]
  updatedById Int [ref: > User.id, null]

  deletedAt DateTime [null]
  createdAt DateTime [default: `now()`]
  updatedAt DateTime
}

Table VariantOption {
  id Int [pk, increment]
  value String
  variantId Int [ref: > Variant.id]

  createdById Int [ref: > User.id, null]
  updatedById Int [ref: > User.id, null]

  deletedAt DateTime [null]
  createdAt DateTime [default: `now()`]
  updatedAt DateTime
}

Table SKU {
  id Int [pk, increment]
  value String
  price Float
  stock Int
  images String[]
  productId Int [ref: > Product.id]

  createdById Int [ref: > User.id, null]
  updatedById Int [ref: > User.id, null]

  deletedAt DateTime [null]
  createdAt DateTime [default: `now()`]
  updatedAt DateTime
}

Table Brand {
  id Int [pk, increment]
  logo String

  createdById Int [ref: > User.id, null]
  updatedById Int [ref: > User.id, null]

  deletedAt DateTime [null]
  createdAt DateTime [default: `now()`]
  updatedAt DateTime
}

Table BrandTranslation {
  id Int [pk, increment]
  brandId Int [ref: > Brand.id]
  languageId Int [ref: > Language.id]
  name String
  description String

  createdById Int [ref: > User.id, null]
  updatedById Int [ref: > User.id, null]

  deletedAt DateTime [null]
  createdAt DateTime [default: `now()`]
  updatedAt DateTime
}

Table CartItem {
  id Int [pk, increment]
  quantity Int
  skuId Int [ref: > SKU.id]
  userId Int [ref: > User.id]

  createdAt DateTime [default: `now()`]
  updatedAt DateTime
}

Table ProductSKUSnapshot {
  id Int [pk, increment]
  productName String
  price Float
  images String[]
  skuValue String
  skuId Int [ref: > SKU.id]
  orderId Int [ref: > Order.id]

  createdAt DateTime [default: `now()`]
}

Table Order {
  id Int [pk, increment]
  userId Int [ref: > User.id]
  status OrderStatus

  createdById Int [ref: > User.id, null]
  updatedById Int [ref: > User.id, null]

  deletedAt DateTime [null]
  createdAt DateTime [default: `now()`]
  updatedAt DateTime
}

Table Review {
  id Int [pk, increment]
  content String
  rating Int
  productId Int [ref: > Product.id]
  userId Int [ref: > User.id]

  deletedAt DateTime [null]
  createdAt DateTime [default: `now()`]
  updatedAt DateTime
}

Table PaymentTransaction {
  id Int [pk, increment]
  gateway String
  transactionDate DateTime [default: `now()`]
  accountNumber String
  subAccount String [null]
  amountIn Int [default: 0]
  amountOut Int [default: 0]
  accumulated Int [default: 0]
  code String [null]
  transactionContent String [null]
  referenceNumber String [null]
  body String [null]

  createdAt DateTime [default: `now()`]
}

Table Message {
   id Int [pk, increment]
   fromUserId Int [ref: > User.id]
   toUserId Int [ref: > User.id]
   content String

   readAt DateTime [null]
   createdAt DateTime [default: `now()`]
}

Enum OrderStatus {
  PENDING_CONFIRMATION
  PENDING_PICKUP
  PENDING_DELIVERY
  DELIVERED
  RETURNED
  CANCELLED
}

Enum VerificationCode {
  REGISTER
  FORGOT_PASSWORD
}

Enum UserStatus {
  ACTIVE
  BLOCKED
}

Enum HTTPMethod {
  GET
  POST
  PUT
  DELETE
  PATCH
  OPTIONS
  HEAD
}

// Many-to-Many Relationships
Table ProductsCategories {
  product_id Int [ref: > Product.id]
  category_id Int [ref: > Category.id]
}

Table PermissionsRoles {
  permission_id Int [ref: > Permission.id]
  role_id Int [ref: > Role.id]
}

Table SkusVariantOptions {
  sku_id Int [ref: > SKU.id]
  variant_option_id Int [ref: > VariantOption.id]
}
```

### 20. Prisma schema

#### 1. Self-relations

- Ví dụ 1: Category (self relation 1–n)

```schema.prisma
model Category {
  id       Int        @id @default(autoincrement())
  name     String

  parentId Int?
  parent   Category?  @relation("CategoryToCategory", fields: [parentId], references: [id])
  children Category[] @relation("CategoryToCategory")
}
```

👉 SQL tạo bảng

```sql
CREATE TABLE "Category" (
  "id" SERIAL PRIMARY KEY,
  "name" TEXT NOT NULL,
  "parentId" INTEGER,
  CONSTRAINT "Category_parentId_fkey"
    FOREIGN KEY ("parentId") REFERENCES "Category"("id")
    ON DELETE SET NULL
    ON UPDATE CASCADE
);
```

- Ví dụ 2: User (many-to-many implicit)

```schema.prisma
model User {
  id         Int     @id @default(autoincrement())
  name       String

  following  User[]  @relation("UserFollows")
  followers  User[]  @relation("UserFollows")
}
```

👉 SQL tạo bảng

```sql
-- 1. User
CREATE TABLE "User" (
  "id" SERIAL PRIMARY KEY,
  "name" TEXT NOT NULL
);

-- 2. Bảng trung gian (Prisma tự tạo)
CREATE TABLE "_UserFollows" (
  "A" INTEGER NOT NULL,
  "B" INTEGER NOT NULL,

  CONSTRAINT "_UserFollows_A_fkey"
    FOREIGN KEY ("A") REFERENCES "User"("id")
    ON DELETE CASCADE ON UPDATE CASCADE,

  CONSTRAINT "_UserFollows_B_fkey"
    FOREIGN KEY ("B") REFERENCES "User"("id")
    ON DELETE CASCADE ON UPDATE CASCADE

-- 3. Index (Prisma thường tạo thêm)
CREATE UNIQUE INDEX "_UserFollows_AB_unique" ON "_UserFollows"("A", "B");
CREATE INDEX "_UserFollows_B_index" ON "_UserFollows"("B");
);
```

- Ví dụ 3: User + Follow (explicit)

```schema.prisma
model User {
  id         Int              @id @default(autoincrement())
  name       String

  following  Follow[] @relation("following")
  followers  Follow[] @relation("followers")
}

model Follow {
  id          Int @id @default(autoincrement())

  followerId  Int
  followingId Int

  follower    User @relation("following", fields: [followerId], references: [id])
  following   User @relation("followers", fields: [followingId], references: [id])
}
```

👉 SQL tạo bảng

```sql
-- 1. User
CREATE TABLE "User" (
  "id" SERIAL PRIMARY KEY,
  "name" TEXT NOT NULL
);

-- 2. Follow
CREATE TABLE "Follow" (
  "id" SERIAL PRIMARY KEY,
  "followerId" INTEGER NOT NULL,
  "followingId" INTEGER NOT NULL,

  CONSTRAINT "Follow_followerId_fkey"
    FOREIGN KEY ("followerId") REFERENCES "User"("id")
    ON DELETE CASCADE ON UPDATE CASCADE,

  CONSTRAINT "Follow_followingId_fkey"
    FOREIGN KEY ("followingId") REFERENCES "User"("id")
    ON DELETE CASCADE ON UPDATE CASCADE
);

-- Best practice (nên thêm)
CREATE UNIQUE INDEX "Follow_unique"
ON "Follow"("followerId", "followingId");
```

### 20. Postgres

```bash
docker run -d \
  --name postgres-ecommerce-db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=123456 \
  -e POSTGRES_DB=ecom_dev \
  -p 5435:5432 \
  postgres:16-alpine


docker stop postgres-ecommerce-db
docker start postgres-ecommerce-db

docker rm -f postgres-ecommerce-db
```

### 21. Install zod (object validation schema) (https://docs.nestjs.com/pipes#object-schema-validation)

```bash
npm install --save zod
```
