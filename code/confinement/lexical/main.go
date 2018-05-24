package main

import "fmt"

func main() {
	loopData := func(handleData chan<- int) {
		defer close(handleData)
		data := []int{1, 2, 3, 4}
		for i := range data {
			handleData <- data[i]
		}
	}

	handleData := make(chan int)
	go loopData(handleData)
	for num := range handleData {
		fmt.Println(num)
	}
}
