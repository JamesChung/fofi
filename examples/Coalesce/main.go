package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/JamesChung/fofi"
)

func main() {
	ctx := context.Background()
	ins := []chan string{make(chan string), make(chan string)}
	ch, cancel := fofi.Coalesce(ctx, 0, ins...)
	outputChanCount := len(ins)
	wg := sync.WaitGroup{}
	for i := 0; i < outputChanCount; i++ {
		go func(id int, ch chan string) {
			for i := 0; i < 10; i++ {
				ch <- fmt.Sprintf("Hello from goroutine #%d", id)
			}
		}(i, ins[i])
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		i := 0
		for v := range ch {
			if i >= 5 {
				cancel()
			}
			fmt.Println(v)
			i++
		}
	}()
	wg.Wait()
}
