package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	myChan := make(chan string)
	go miner(myChan)
	go afip(ctx, myChan)
	now := time.Now().UnixNano()
	if now%2 == 0 {
		fmt.Println(now)
		cancel()
	}

	time.Sleep(4 * time.Second)
}

func miner(myChan chan<- string) {
	for _, item := range [3]string{"coin0", "coin1", "coin2"} {
		time.Sleep(500 * time.Millisecond)
		myChan <- item //send
	}
}

func afip(ctx context.Context, myChan <-chan string) {
	for i := 0; i < 3; {
		select {
		case <-ctx.Done():
			fmt.Println("Cancelled!")
			return
		case coin := <-myChan: //receive
			fmt.Println("Received:  " + coin + " from miner")
			i++
		default:
			fmt.Println("Searching more miners !!")
			time.Sleep(100 * time.Millisecond)
		}
	}
}
