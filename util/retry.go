package util

import (
	"sync"
	"time"

	"golang.org/x/exp/slog"
)

func WithRetriesAndTimeout[T any](timeout time.Duration, maxRetries uint, f func() (T, error)) (T, error) {
	type result struct {
		val T
		err error
	}

	var err error
	var bad T
	ch := make(chan result, 1)

	var wg sync.WaitGroup
	for i := uint(0); i < maxRetries; i++ {
		wg.Add(1)
		go func() {
			val, err := f()
			ch <- result{val, err}
			wg.Done()
		}()
		select {
		case res := <-ch:
			if res.err != nil {
				slog.Warn("withRetriesAndTimeout: failed", "err", res.err, "try", i)
				err = res.err
				continue
			} else {
				return res.val, nil
			}
		case <-time.After(timeout):
			slog.Warn("withRetriesAndTimeout: timed out", "err", "timeout", "try", i)
			continue
		}
	}

	defer func() {
		wg.Wait()
		close(ch)
	}()

	return bad, err
}
