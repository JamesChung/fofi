package fofi_test

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/JamesChung/fofi"
)

// This example will create one goroutine which will count up from 0 every second
// and create output channels equal to 'outputChanCount' and each subscribed
// goroutine will print the value.
func ExampleBroadcast() {
	ctx := context.Background()
	in := make(chan int)
	outs := []chan int{make(chan int), make(chan int), make(chan int)}
	outputChanCount := len(outs)
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

// This example will create one goroutine which will count up from 0 every second
// and create output channels equal to 'outputChanCount' and each subscribed
// goroutine will print the value.
func ExampleGenerateOutputBroadcasters() {
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

// This example will create one goroutine which will count up from 0 every second
// and create output channels equal to 'outputChanCount' and each subscribed
// goroutine will print the value.
func ExampleGenerateInputOutputBroadcasters() {
	ctx := context.Background()
	outputChanCount := 3
	wg := sync.WaitGroup{}
	in, out, cancel := fofi.GenerateInputOutputBroadcasters[int](ctx, outputChanCount)
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
		go func(out <-chan int) {
			defer wg.Done()
			for v := range out {
				fmt.Println(v)
			}
		}(out[i])
	}
	wg.Wait()
}

// This example creates two output channels and will consolidate
// messages from both into one channel and output the value.
func ExampleCoalesce() {
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
