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
	outputChanCount := 3
	wg := sync.WaitGroup{}
	in, out, cancel := fofi.GenerateBroadcast[int](ctx, outputChanCount)
	go func() {
		i := 0
		for {
			if i >= 10 {
				cancel()
			}
			in <- i
			i++
			time.Sleep(time.Second)
		}
	}()

	for i := 0; i < outputChanCount; i++ {
		wg.Add(1)
		go func(out <-chan int) {
			defer wg.Done()
			for v := range out {
				fmt.Println(v)
			}
		}(out[i])
	}
	wg.Wait()
}
