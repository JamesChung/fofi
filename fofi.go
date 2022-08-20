package fofi

import (
	"context"
)

// Broadcast
func Broadcast[T any](ctx context.Context, in <-chan T, out ...chan T) context.CancelFunc {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		for {
			val := <-in
			select {
			case <-ctx.Done():
				return
			default:
				for _, ch := range out {
					ch <- val
				}
			}
		}
	}()
	return cancel
}

// Coalesce
func Coalesce[T any](ctx context.Context, bufferSize int, in ...chan T) (<-chan T, context.CancelFunc) {
	ch := make(chan T, bufferSize)
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		for _, c := range in {
			go func(c <-chan T) {
				for {
					select {
					case <-ctx.Done():
						return
					default:
						ch <- <-c
					}
				}
			}(c)
		}
	}()
	return ch, cancel
}
