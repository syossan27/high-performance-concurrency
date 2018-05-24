package main

import (
	"fmt"
	"time"
)

func main() {
	done := make(chan interface{}) // 処理の終了を知らせるchannel
	time.AfterFunc(10*time.Second, func() { close(done) }) // 10s後には終了
	const timeout = 2 * time.Second // タイムアウトするまでの時間

	heartbeat, results := doWork(done, timeout/2)
	for {
		select {
		case _, ok := <-heartbeat:
			if ok == false {
				fmt.Println("心臓の鼓動が停止しました・・・")
				return
			}
			fmt.Println("pulse")
		case r, ok := <-results:
			if ok == false {
				return
			}
			fmt.Printf("results %v\n", r.Second())
		case <-time.After(timeout):
			fmt.Println("タイムアウトしました！")
			return
		}
	}
}

// workを並列で動かし、heartbeat channelとresult channelを返す
func doWork(
	done <-chan interface{},
	pulseInterval time.Duration, // heartbeatの確認パルスを送る時間間隔
) (<-chan interface{}, <-chan time.Time) {
	heartbeat := make(chan interface{})
	results := make(chan time.Time)
	go work(heartbeat, results, pulseInterval, done)
	return heartbeat, results
}

func work(
	heartbeat chan interface{},
	results chan time.Time,
	pulseInterval time.Duration,
	done <-chan interface{},
) {
	defer close(heartbeat)
	defer close(results)

	pulse := time.Tick(pulseInterval) 		// 1sごとに発火するchannel
	workGen := time.Tick(2 * pulseInterval) // 2sごとに発火するchannel

	//for {
	for i := 0; i < 2; i++ {
		select {
		case <-done:
			return
		case <-pulse:
			// 1s毎にsendPulseを実行
			sendPulse(heartbeat)
		case r := <-workGen: // 2s毎にtime.Timeを返す
			sendResult(r, done, pulse, heartbeat, results)
		}
	}
}

// heartbeat channelに何らかの値を入れる（空structがオススメ）
func sendPulse(heartbeat chan interface{}) {
	select {
	case heartbeat <- struct{}{}:
	default: // bufferが満杯の場合、このdefaultがないとblockingされてしまう
	}
}

func sendResult(
	r time.Time,
	done <-chan interface{},
	pulse <-chan time.Time,
	heartbeat chan interface{},
	results chan time.Time,
) {
	for {
		select {
		case <-done:
			return
		case <-pulse:
			sendPulse(heartbeat)
		case results <- r:
			return
		}
	}
}
