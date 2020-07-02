package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_UnknownToken(t *testing.T) {
	testSrv.UnknownToken()
	onceToken.Do(func() {
		assert.Fail(t, "call twice UnknownToken")
	})
}
