package utils

import "testing"

func TestPrint(t *testing.T) {
	Plain("Test", "Something")
}

func TestMsg(t *testing.T) {
	Msg("Test", "Something")
}

func TestMsg2(t *testing.T) {
	Msg2("Test", "Something")
}

func TestWarning(t *testing.T) {
	Msg2("Test", "Something")
}

func TestError(t *testing.T) {
	Error("Test", "Something")
}
