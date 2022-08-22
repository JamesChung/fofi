package fofi

import (
	"context"
	"sync"
)

// Broadcast takes an 'in' channel of T and will broadcast that value
// to every 'out' channel of T. Base function that's used by higher-order
// GenerateOutputBroadcast() and GenerateBroadcast()
// functions. Ideally you should use those, this is exposed if you want
// to create the channels yourself.
//
// NOTE: Broadcast will close 'out' channels when context.CancelFunc is invoked.
//
// It returns a cancel function.
func Broadcast[T any](ctx context.Context, in <-chan T, out ...chan T) context.CancelFunc {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		defer func() {
			for _, c := range out {
				close(c)
			}
		}()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				val := <-in
				for _, ch := range out {
					ch <- val
				}
			}
		}
	}()
	return cancel
}

// GenerateOutputBroadcast takes an 'in' channel of type T and will generate 'n' number of
// channels which can be used to broadcast to via the 'in' channel. If the cancel func is invoked
// all 'out' channels will be closed.
//
// It returns a slice 'out' channels of type T and a CancelFunc.
func GenerateOutputBroadcast[T any](ctx context.Context, in <-chan T, n int) (out []<-chan T, cancel context.CancelFunc) {
	chs := make([]chan T, n)
	// Initialize channels of type T
	for i := 0; i < n; i++ {
		chs[i] = make(chan T)
	}

	cancel = Broadcast(ctx, in, chs...)

	// Convert default channels into receive only channels
	out = make([]<-chan T, n)
	for i := 0; i < n; i++ {
		out[i] = chs[i]
	}
	return
}

// GenerateBroadcast is designed to be the simplest way to create an input
// channel with a series of 'n' output channels. Useful when you don't want to create
// channels yourself. Uses Broadcast() as a base, cancelFunc will close channels for you.
//
// It returns an 'in' channel of type T and a slice of 'out' channels
// of type T length equal to the parameter 'n'. Also returns a context.CancelFunc.
func GenerateBroadcast[T any](ctx context.Context, n int) (in chan<- T, out []<-chan T, cancel context.CancelFunc) {
	tmp := make(chan T)
	out, cancel = GenerateOutputBroadcast(ctx, tmp, n)
	in = tmp
	return
}

// Coalesce will create a receive channel of T with n 'bufferSize' and a slice of m channels
// of T. Useful when you have a series of n channels which you want to join together into
// one channel to receive on. Will close the channel when CancelFunc is invoked.
//
// It returns a receive channel of T and a context.CancelFunc.
func Coalesce[T any](ctx context.Context, bufferSize int, in ...chan T) (<-chan T, context.CancelFunc) {
	ch := make(chan T, bufferSize)
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		wg := sync.WaitGroup{}
		defer close(ch)
		for _, c := range in {
			wg.Add(1)
			go func(c <-chan T) {
				defer wg.Done()
				for {
					select {
					case <-ctx.Done():
						return
					default:
						if v, ok := <-c; ok {
							ch <- v
						} else {
							return
						}
					}
				}
			}(c)
		}
		wg.Wait()
	}()
	return ch, cancel
}
