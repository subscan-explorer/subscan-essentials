package observer_test

import (
	"github.com/itering/subscan/internal/observer"
	"github.com/sevlyar/go-daemon"
	"testing"
)

func TestRun(t *testing.T) {
	observer.Run("demo", "start")
	if len(daemon.Flags()) != 1 {
		t.Errorf("Run Daemon failed")
	}
}
