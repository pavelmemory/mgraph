package main

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func T(stop chan bool, data chan string) {
	for {
		result := make(chan int)
		go func() {
			select {
			case <-stop:
				return
			case d := <-data:
				r, _ := strconv.Atoi(d)
				result <- r
			}
		}()

		select {
		case <-time.After(1 * time.Second):
			fmt.Println("timeout")
			return
		case <-stop:
			fmt.Println("stopped")
			return
		case comRes := <-result:
			stop <- false
			fmt.Println(comRes)
		}
	}
}

func TestName(t *testing.T) {
	data := make(chan string, 1)
	stop := make(chan bool, 1)
	data <- "1002"
	T(stop, data)

}
