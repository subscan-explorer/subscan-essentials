package daemons_test

import (
	"github.com/itering/subscan/internal/daemons"
	"github.com/sevlyar/go-daemon"
	"testing"
)

func TestRun(t *testing.T) {
	daemons.Run("demo", "start")
	if len(daemon.Flags()) != 1 {
		t.Errorf("Run Daemon failed")
	}
}
