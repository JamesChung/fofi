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
func ExampleGenerateOutputBroadcast() {
	ctx := context.Background()
	in := make(chan int)
	outputChanCount := 3
	wg := sync.WaitGroup{}
	out, cancel := fofi.GenerateOutputBroadcast(ctx, in, outputChanCount)
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
func ExampleGenerateBroadcast() {
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

// This example creates two output channels and will consolidate
// messages from both into one channel and output the value.
func ExampleCoalesce() {
	ctx := context.Background()
	ins := []chan string{make(chan string), make(chan string)}
	outputChanCount := len(ins)
	wg := sync.WaitGroup{}
	for i := 0; i < outputChanCount; i++ {
		go func(id int, ch chan string) {
			for i := 0; i < 10; i++ {
				ch <- fmt.Sprintf("Hello from goroutine #%d", id)
			}
		}(i, ins[i])
	}
	inputs := make([]<-chan string, outputChanCount)
	for i, c := range ins {
		inputs[i] = c
	}
	ch, cancel := fofi.Coalesce(ctx, 0, inputs...)
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
