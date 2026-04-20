# Javascript

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
npm i brcypt
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
