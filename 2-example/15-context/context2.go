package contextExample

import (
	"context"
	"fmt"
	"time"
)

func cookCom(ctx context.Context, chPho chan<- string) {
	fmt.Println("Bat dau nau Com")
	select {
	case val := <-time.After(100 * time.Millisecond):
		chPho <- "Com da nau xong"
		fmt.Printf("%T, %v\n", val, val) // time.Time, 2026-02-24 14:14:38.154578472 +0700 +07 m=+0.100089429
	case <-ctx.Done():
		fmt.Println("Huy nau Com")
		return
	}
}

func cookPho(ctx context.Context, chPho chan<- string) {
	fmt.Println("Bat dau nau Pho")
	select {
	case val := <-time.After(1 * time.Second):
		chPho <- "Pho da nau xong"
		fmt.Printf("%T, %v\n", val, val) // time.Time, 2026-02-24 14:14:39.054559826 +0700 +07 m=+1.000070783
	case <-ctx.Done():
		fmt.Println("Huy nau Pho")
		return
	}
}

func cookChao(ctx context.Context, chChao chan<- string) {
	fmt.Println("Bat dau nau Chao")
	select {
	case val := <-time.After(2 * time.Second):
		chChao <- "Chao da nau xong"
		fmt.Printf("%T, %v\n", val, val)
	case <-ctx.Done():
		fmt.Println("Huy nau Chao")
		return
	}
}

func Main2() {
	chPho := make(chan string)
	chChao := make(chan string)
	chCom := make(chan string)

	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()

	go cookPho(ctx, chPho)
	go cookChao(ctx, chChao)
	go cookCom(ctx, chCom)

	for i := 1; i <= 3; i++ {
		select {
		case pho := <-chPho:
			fmt.Println("Nhan duoc: ", pho)
		case chao := <-chChao:
			fmt.Println("Nhan duoc: ", chao)
		case com := <-chCom:
			fmt.Println("Nhan duoc: ", com)
		case val := <-ctx.Done():
			fmt.Println("Timeout, khong nhan mon")
			fmt.Printf("%T, %v\n", val, val) // struct {}, {}
			return
		}
	}
}
