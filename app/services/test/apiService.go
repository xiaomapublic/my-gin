package test

import (
	"fmt"
	"sync"
)

func Seek(name string, match chan string, wg *sync.WaitGroup) {
	select {
	case peer := <-match:
		fmt.Printf("%s sent a message to %s.\n", peer, name)
	case match <- name:
		// 等待某个goroutine接收我的消息
	}
	wg.Done()
}
