package parallel

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParallel(t *testing.T) {
	if testing.Short() {
		return
	}

	ctx := context.Background()
	p := Func(func(ctx context.Context) {
		time.Sleep(100 * time.Millisecond)
		t.Log("running 1")
	}).Func(func(ctx context.Context) {
		time.Sleep(200 * time.Millisecond)
		t.Log("running 2")
	}).Func(func(ctx context.Context) {
		t.Log("running 3")
	}).Start(ctx)
	assert.True(t, p.Running())
	t.Log("parallel start")
	assert.True(t, p.Wait())
	t.Log("should finished running")
	assert.False(t, p.Running())

	aCtx, aCancel := context.WithCancel(context.Background())
	p.Func(func(ctx context.Context) {
		time.Sleep(time.Second * 100)
	}).Start(aCtx)
	assert.True(t, p.Running())
	aCancel()
	assert.False(t, p.Wait())
	time.Sleep(time.Millisecond)
	assert.False(t, p.Running())
}

func TestErrHandler(t *testing.T) {
	if testing.Short() {
		return
	}

	var (
		err1 = errors.New("1")
		err2 = errors.New("2")
		err3 = errors.New("3")
	)

	var passed = make(chan struct{})

	FuncE(func(ctx context.Context) error {
		time.Sleep(time.Millisecond * 100)
		return err1
	}).FuncE(func(ctx context.Context) error {
		time.Sleep(time.Millisecond * 300)
		return err2
	}).FuncE(func(ctx context.Context) error {
		time.Sleep(time.Millisecond * 200)
		return err3
	}).ErrHandle(func(errs []error) {
		if !assert.Len(t, errs, 3) {
			return
		}
		assert.ErrorIs(t, errs[0], err1)
		assert.ErrorIs(t, errs[1], err3)
		assert.ErrorIs(t, errs[2], err2)
		close(passed)
	}).Start(context.Background()).Wait()

	select {
	case <-passed:
	case <-time.After(time.Second):
		assert.Fail(t, "errHandler not working")
	}
}
