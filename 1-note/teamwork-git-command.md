### Init project
```bash
git init
git checkout -b feature/init-project
git remote add origin https://github.com/rs-kiennt2/e-smart-learn.git
git remote -v
git add .
git commit -m "feat: init project (#1)"
git push origin feature/init-project
```

### Xóa nhánh ở local
```bash
git checkout main
git branch -D feature/1-init-project
```


### Clear hết docker data
```bash
docker ps -a
docker system prune -a --volumes -f
```

### Câu hỏi 1: Sửa commit gần nhất trên nhánh feature/login-register-rbac ở remote
```bash
git checkout feature/login-register-rbac
git pull origin feature/login-register-rbac
git reset --soft HEAD~1
git commit -m "feat: login, register, setup middleware rbac (#5, #6)"
git push --force
```

### Câu hỏi 2: tôi đang ở nhánh a2 trong branch git, tôi đang code dở (chưa commit) trên branch này ở local. Giờ trên remote nhánh a1 có code mới, tôi muốn lấy code mới từ nhánh a1 remote đó về nhánh a2 local sau đó mới code tiếp thì dùng câu lệnh gì

1. Cách 1 (nếu bạn thích lịch sử git gọn hơn – dùng rebase   --> Cách này OK nhất):
```bash
git stash
git fetch origin
git rebase origin/a1
git stash pop
```

2. Cách 2 (khuyên dùng khi chưa commit)
```bash
git stash
git pull origin a1
git stash pop
```

3. Cách 3 (nếu muốn kiểm soát rõ hơn)
```bash
git stash
git fetch origin
git merge origin/a1
git stash pop
```

### Câu hỏi 3: nếu code trên nhánh a2 đã commit ở local rồi mà chưa đẩy lên remote, tôi muốn lấy code ở nhánh a1 remote về a2 trước rồi mới áp dụng tiếp cái a2 local commit thì làm sao?

- Trường hợp của bạn:
    + a2 local đã có commit mới nhưng chưa push
    + a1 remote có code mới
    + Bạn muốn lấy code mới từ a1 trước, rồi áp dụng lại commit của a2 local lên trên

    ➡️ Cách đúng nhất là dùng rebase.

- Cách làm: 
    + Đang đứng ở branch a2:
        ```bash
        git checkout a2
        git fetch origin
        git rebase origin/a1
        ```

    + Ý nghĩa

        `git fetch origin` → lấy code mới từ remote về

        `git rebase origin/a1` → Git sẽ:

            + tạm gỡ commit của a2 local ra
            + cập nhật code từ a1
            + apply lại commit của a2 lên trên
            + Lịch sử sẽ giống như bạn viết code a2 sau khi a1 đã cập nhật.

- Ví dụ:
    + Trước khi rebase
        ```
            a1:  A --- B --- C
                            \
            a2:               D --- E   (commit local)
        ```
    + Sau khi rebase:
        ```
        A --- B --- C --- D' --- E'
        ```

        (D', E' là commit được apply lại)

- Nếu có conflict
    + Git sẽ báo:
        ```
        CONFLICT ...
        ```
    + Bạn resolve rồi chạy tiếp:
        ```bash
        git add .
        git rebase --continue
        ```
    + Nếu muốn hủy rebase:
        ```bash
        git rebase --abort
        ```

### Câu hỏi 4: Đang ở nhánh feature/login-register-rbac thì kéo code từ nhánh main remote về
```
git checkout feature/login-register-rbac
git pull --rebase origin main
```
Sau đó
```
git push --force
```

###



### Migrate goose
1. Cài đặt goose
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
goose -version
```

2. Chạy migrate
```bash
goose -dir app/repository/migrations create update_courses_status_enum sql
goose -dir app/repository/migrations postgres "postgres://postgres:123456@localhost:5433/elearning?sslmode=disable" up
goose -dir app/repository/migrations postgres "postgres://postgres:123456@localhost:5433/elearning?sslmode=disable" reset
goose -dir app/repository/migrations postgres "postgres://postgres:123456@localhost:5433/elearning?sslmode=disable" down
```

### Stripe
1. Setup stripe webhook
```
stripe login
stripe listen --forward-to localhost:8080/api/v1/subscriptions/webhook/stripe
```
- Sửa `secret_key` trong `config.yaml` (optional)
- Copy webhook secret vào key `webhook_secret` trong `config.yaml`

2. Account Stripe
```
4242 4242 4242 4242
05/30
123
student1
```

###
hiện tại dự án này đang cần làm thêm tính năng chia tiền thu được giữa instructor và nền tảng theo:
1. Mua khóa học: instructor đựợc chia 70% và nền tảng thì được chia 30%
2. Subscriptions: người dùng có thể subscriptions và có thể enroll được tất cả các khóa trên nền tảng, nếu 1 người dùng enroll 1 khóa học thì chủ khóa học đó sẽ được chia tiền từ subscription của người dùng đó (cũng theo tỉ lệ 70% của chủ khóa học và 30% của nền tảng).

ngoài ra cũng cần làm chức năng thống kê tiền thu được, tiền của nền tảng, tiền của giảng viên theo tháng, quý, năm và xem thử tăng giảm bao nhiêu so với tháng, quý, năm trước.
Hãy đọc toàn bộ project của tôi hãy đề xuất cho tôi phải làm thế nào, viết api gì, sửa database như thế nào để làm toàn bộ chức năng đó

### Figma
https://www.figma.com/design/YOvNhQqxHQODgkI95ev2Op/E-Learning-Site--Community---Copy-?node-id=0-1&p=f&t=EEYv1Mlo03mTy4Zn-0

### Claude Code
1. Command
```
/init
/compact
/clear
/resume
/btw
/cost
/model
/mcp
```
https://ai-first.anhhh.workers.dev/skill_proposal/presales-skill-guideline

###
* hiện tại tôi đang làm dự án elearning giống udemy, tôi đã viết api tạo khóa học (product) trên stripe, ngoài ra tôi có phần chỉ instructor mới có thể tạo coupon trên stripe. Ngoài ra ở phần coupon có thể gán nó cho nhiều khóa học. Khi gán nó cho khóa học thì sẽ giảm giá tạm thời cho khóa học đó theo thuộc tính của coupon (ví dụ coupon đó có expires là ngày mấy đó và giảm giá 20% trong khoảng thời gian đó với số lượng sử dụng max_redemptions. Dự án tôi viết bằng golang, gorm, gin, postgres. Như vậy mỗi lần gán coupon cho khóa học thì sẽ tạo 1 giá mới cho course trên stripe hay sao

* Tôi muốn là ví dụ khi instructor áp dụng coupon cho khóa học thì ví dụ khóa học đó có giá gốc là 10 đô sẽ giảm xuống 8 đô, phần giao diện course detail sẽ thấy 10 đô (bị gạch) và giá hiện tại là 8 đô. Sau khi người dùng bấm buy now thì chuyển qua trang checkout và người dùng có thể nhập coupon do instructor public cho người dùng nhập (ví dụ coupon điền vào giảm 40%) thì giá sau giảm sẽ là 8 - 8 * 40%. Vậy tôi phải làm như thế nào

- Backend xử lý:

original_price = 10
course_discount = 20% → 8
user_coupon = 40% → 8 - 40% = 4.8

```go 
params := &stripe.CheckoutSessionParams{
  LineItems: []*stripe.CheckoutSessionLineItemParams{
    {
      PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
        Currency: stripe.String("usd"),
        UnitAmount: stripe.Int64(480), // giá sau khi discount ✅ FINAL PRICE
        Product: stripe.String(productID), // ✅ dùng product có sẵn
      },
      Quantity: stripe.Int64(1),
    },
  },
}
```

```sql
CREATE TABLE coupons (
    id BIGSERIAL PRIMARY KEY,

    code VARCHAR(50) UNIQUE,           -- mã coupon (NULL nếu auto discount)

    type VARCHAR(20) NOT NULL,         -- COURSE | GLOBAL | FLASH

    value_type VARCHAR(10) NOT NULL,   -- PERCENT | FIXED
    value NUMERIC(10,2) NOT NULL,      -- 20 (%) hoặc 5 ($)

    max_redemptions INT,               -- số lần dùng tối đa
    used_count INT DEFAULT 0,          -- đã dùng bao nhiêu lần

    per_user_limit INT DEFAULT 1,      -- mỗi user dùng tối đa bao nhiêu lần

    start_at TIMESTAMP,
    end_at TIMESTAMP,

    is_active BOOLEAN DEFAULT TRUE,

    created_by BIGINT,                 -- instructor_id hoặc admin_id

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

```sql
ALTER TABLE courses ADD COLUMN allow_global_discount BOOLEAN DEFAULT TRUE;

if course.AllowGlobalDiscount {
    applyFlashSale()
}
```
