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
	outputChanCount := 3
	wg := sync.WaitGroup{}
	out, cancel := fofi.GenerateOutputBroadcasters(ctx, in, outputChanCount)
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
		go func(ch <-chan int) {
			defer wg.Done()
			for v := range ch {
				fmt.Println(v)
			}
		}(out[i])
	}
	wg.Wait()
}
