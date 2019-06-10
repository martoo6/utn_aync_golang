package main

import (
	"fmt"
	"time"
)

func main() {
	myChan := make(chan string)

	go miner(myChan)
	go afip(myChan)

	time.Sleep(6 * time.Second)
}

func miner(myChan chan<- string) {
	for _, item := range [3]string{"coin0", "coin1", "coin2"} {
		time.Sleep(500 * time.Millisecond)
		myChan <- item //send
	}
}

func afip(myChan <-chan string) {
	for i := 0; i < 3; i++ {
		coin := <-myChan //receive
		fmt.Println("Received:  " + coin + " from miner")
	}
}
