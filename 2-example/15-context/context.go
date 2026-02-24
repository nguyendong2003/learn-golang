package contextExample

import (
	"context"
	"fmt"
	"time"
)

type contextKey string

func Main() {

	// 1️⃣ Background - root context
	rootCtx := context.Background()

	// 2️⃣ WithValue - gắn requestID
	ctxWithValue := context.WithValue(rootCtx, contextKey("requestID"), "req-456")

	// 3️⃣ WithCancel - cho phép hủy thủ công
	ctxWithCancel, cancel := context.WithCancel(ctxWithValue)

	// 4️⃣ WithDeadline - hết hạn sau 5 giây (mốc thời gian cụ thể)
	deadline := time.Now().Add(5 * time.Second)
	ctxWithDeadline, deadlineCancel := context.WithDeadline(ctxWithCancel, deadline)
	defer deadlineCancel()

	// 5️⃣ WithTimeout - timeout 3 giây
	ctxWithTimeout, timeoutCancel := context.WithTimeout(ctxWithDeadline, 3*time.Second)
	defer timeoutCancel()

	// 6️⃣ TODO - giả sử một chỗ chưa quyết định context
	todoCtx := context.TODO()
	fmt.Println("TODO ctx:", todoCtx)

	// Chạy công việc
	go doWork(ctxWithTimeout)

	// Sau 2 giây thì cancel thủ công
	time.Sleep(2 * time.Second)
	fmt.Println("Manual cancel() called")
	cancel()

	// Đợi context kết thúc
	<-ctxWithTimeout.Done()
	fmt.Println("Main stopped:", ctxWithTimeout.Err())
}

func doWork(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Work stopped:", ctx.Err())
			return
		default:
			requestID := ctx.Value(contextKey("requestID"))
			fmt.Println("Working... requestID =", requestID)
			time.Sleep(1 * time.Second)
		}
	}
}
