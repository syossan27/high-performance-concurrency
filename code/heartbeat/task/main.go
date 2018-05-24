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
	heartbeatStream := make(chan interface{})
	workStream := make(chan int)
	go work(heartbeatStream, workStream, done)
	return heartbeatStream, workStream
}

func work(
	heartbeatStream chan interface{},
	workStream chan int,
	done <-chan interface{},
) {
	defer close(heartbeatStream)
	defer close(workStream)

	for i := 0; i < 10; i++ {
		sendPulse(heartbeatStream)

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
	}
}
