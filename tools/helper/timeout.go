package helper

import (
	"context"
	"time"
)

func TimeoutGuardWithFunc(ctx context.Context, timeout time.Duration, onBegin, onEnd, onTimeout func(...interface{}), args ...interface{}) func() {
	if onBegin != nil {
		onBegin(args...)
	}

	timer := time.NewTimer(timeout)

	ctxG, canceler := context.WithCancel(ctx)
	go func() {
		select {
		case <-timer.C:
			if onTimeout != nil {
				onTimeout(args...)
			}
		case <-ctxG.Done():
			timer.Stop()
			if onEnd != nil {
				onEnd(args...)
			}
		}
	}()

	return canceler
}
