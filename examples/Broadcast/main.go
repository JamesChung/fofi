package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/JamesChung/fofi"
)

func main() {
	ctx := context.Background()
	in := make(chan int)
	out1 := make(chan int)
	out2 := make(chan int)
	out3 := make(chan int)
	outs := []chan int{out1, out2, out3}
	outputChanCount := 3
	wg := sync.WaitGroup{}
	cancel := fofi.Broadcast(ctx, in, outs...)
	go func() {
		i := 0
		for {
			in <- i
			i++
			if i >= 10 {
				cancel()
			}
			time.Sleep(time.Second)
		}
	}()

	for i := 0; i < outputChanCount; i++ {
		wg.Add(1)
		go func(ch chan int) {
			defer wg.Done()
			for v := range ch {
				fmt.Println(v)
			}
		}(outs[i])
	}
	wg.Wait()
}
