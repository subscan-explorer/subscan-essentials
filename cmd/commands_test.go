package main

import (
	"testing"
)

func Test_AtLeastCommands(t *testing.T) {
	// Test commands has start,install,CheckCompleteness commands
	action := []string{"start", "install", "CheckCompleteness"}
	for _, v := range action {
		var exist bool
		for _, c := range commands {
			if c.Name == v {
				exist = true
			}
		}
		if !exist {
			t.Errorf("commands should have %s command", v)
		}
	}
}
