package log

import "testing"

func TestInfo(t *testing.T) {
	log := NewLog("debug", true)
	log.Warn("nnnnnnnnnn")
}
