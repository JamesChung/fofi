package fofi

import (
	"context"
)

// Broadcast takes an 'in' channel of T and will broadcast that value
// to every 'out' channel of T.
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

// GenerateBroadcasters takes an 'in' channel of T and will generate 'n' number of
// channels.
//
// It returns a slice of type T channels and a CancelFunc.
func GenerateBroadcasters[T any](ctx context.Context, in <-chan T, n int) ([]chan T, context.CancelFunc) {
	chs := make([]chan T, n)
	for i := 0; i < n; i++ {
		chs[i] = make(chan T)
	}
	cancel := Broadcast(ctx, in, chs...)
	return chs, cancel
}

// Coalesce will create a channel of T with n 'bufferSize' and a slice of m channels
// of T.
//
// It returns a receive channel of T and a CancelFunc.
func Coalesce[T any](ctx context.Context, bufferSize int, in ...chan T) (<-chan T, context.CancelFunc) {
	ch := make(chan T, bufferSize)
	ctx, cancel := context.WithCancel(ctx)
	for _, c := range in {
		go func(c <-chan T) {
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
	return ch, cancel
}
