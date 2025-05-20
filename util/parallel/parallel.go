package parallel

import (
	"context"
	"sync"

	"github.com/panjf2000/ants/v2"
)

type parallel struct {
	wg sync.WaitGroup
	// list slice of CallFn or CallErrFn
	list       []interface{}
	ctx        context.Context
	errHandler ErrorHandler
	stopped    chan struct{}
}

type CallFn = func(ctx context.Context)
type CallErrFn = func(ctx context.Context) error

type ErrorHandler = func(errs []error)
type Error []error

func New() *parallel {
	return &parallel{
		wg:   sync.WaitGroup{},
		list: make([]interface{}, 0),
		ctx:  context.Background(),
	}
}

func Func(fn CallFn) *parallel {
	parallel := New()
	return parallel.Func(fn)
}

func FuncE(fn CallErrFn) *parallel {
	parallel := New()
	return parallel.FuncE(fn)
}

func (p *parallel) push(fn interface{}) {
	if fn == nil {
		return
	}

	switch fn.(type) {
	case CallErrFn, CallFn:
	default:
		return
	}

	p.list = append(p.list, fn)
}

func (p *parallel) pop() interface{} {
	if len(p.list) == 0 {
		return nil
	}
	fn := p.list[0]
	p.list = p.list[1:len(p.list)]
	return fn
}

func (p *parallel) run(ctx context.Context) {
	defer close(p.stopped)

	var errs Error
	var errChan = make(chan error)
	cp, _ := ants.NewPoolWithFunc(5, func(next interface{}) {
		defer p.wg.Done()
		switch fn := next.(type) {
		case CallFn:
			fn(ctx)
		case CallErrFn:
			if err := fn(ctx); err != nil {
				errChan <- err
			}
		}
	})
	defer cp.Release()
	for {
		next := p.pop()
		if next == nil {
			break
		}

		p.wg.Add(1)
		_ = cp.Invoke(next)

		select {
		case <-ctx.Done():
			return
		default:
		}
	}

	go func() {
		p.wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		errs = append(errs, err)
	}

	if p.errHandler != nil && len(errs) > 0 {
		p.errHandler(errs)
	}
}

func (p *parallel) Func(fn CallFn) *parallel {
	p.push(fn)
	return p
}

func (p *parallel) FuncE(fn CallErrFn) *parallel {
	p.push(fn)
	return p
}

func (p *parallel) Start(ctx context.Context) *parallel {
	p.ctx = ctx
	p.stopped = make(chan struct{})
	go p.run(p.ctx)
	return p
}

func (p *parallel) ErrHandle(handler ErrorHandler) *parallel {
	p.errHandler = handler
	return p
}

func (p *parallel) Running() bool {
	select {
	case _, ok := <-p.stopped:
		return ok
	default:
		return true
	}
}

// Wait if tasks is early exit will return false
func (p *parallel) Wait() bool {
	if !p.Running() {
		return false
	}

	select {
	case <-p.ctx.Done():
		return p.ctx.Err() == nil
	case <-p.stopped:
		return true
	}
}
