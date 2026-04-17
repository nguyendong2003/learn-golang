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

### Câu hỏi 5: Lỡ push thừa vài file vào nhánh feature/course-coupon-manage và giờ muốn sửa lại commit gần nhất đó thì làm như sau
```
git checkout feature/course-coupon-manage   (Sau đó thay đổi)
git add .
git commit --amend    (Sau đó sửa và bấm Ctrl + O rồi Ctrl + X)
git push --force-with-lease origin feature/course-coupon-manage
```

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


```
https://excalidraw.com/#room=fe20c380a6119007ef93,XE104RxbeZZLoshZDV1Kxg
https://excalidraw.com/#json=5fbg4j_-PGL-h-p9qzXMT,KK80Exs23rZmnM774PliXw
```
