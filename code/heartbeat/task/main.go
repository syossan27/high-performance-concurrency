package main

import (
	"fmt"
	"math/rand"
)

func main() {
	done := make(chan interface{})
	defer close(done)

	heartbeat, results := doWork(done)
	for {
		select {
		case _, ok := <-heartbeat:
			if ok {
				fmt.Println("pulse")
			} else {
				fmt.Println("心臓の鼓動が停止しました・・・")
				return
			}
		case r, ok := <-results:
			if ok {
				fmt.Printf("results %v\n", r)
			} else {
				return
			}
		}
	}
}

func doWork(done <-chan interface{}) (<-chan interface{}, <-chan int) {
	heartbeatStream := make(chan interface{}, 1)
	workStream := make(chan int)
	go work(heartbeatStream, workStream, done)
	return heartbeatStream, workStream
}

func work(
	heartbeatStream chan interface{},
	workStream chan int,
	done <-chan interface{},
) {
	defer func() {
		if r := recover(); r != nil {
		}
		close(heartbeatStream)
		close(workStream)
	}()

	for i := 0; i < 10; i++ {
		sendPulse(heartbeatStream)

		if i == 3 {
			panic("foo")
		}

		select {
		case <-done:
			return
		case workStream <- rand.Intn(10):
		}
	}
}

func sendPulse(heartbeatStream chan interface{}) {
	select {
	case heartbeatStream <- struct{}{}:
	default:
	}
}
