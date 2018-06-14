package utils

import (
	"bytes"
	"os"
	"testing"
)

func TestPrint(t *testing.T) {
	Plain("Reticulating splines")
}

func TestMsg(t *testing.T) {
	Msg("Reticulating splines")
}

func TestMsg2(t *testing.T) {
	Msg2("Reticulating splines")
}

func TestWarning(t *testing.T) {
	Warning("Reticulating splines")
}

func TestError(t *testing.T) {
	Error("Reticulating splines")
}

func resetOutput() {
	Stdout = os.Stdout
	UseColors(true)
}

func TestErrorNoColor(t *testing.T) {
	var buf bytes.Buffer
	Stdout = &buf
	UseColors(false)
	defer resetOutput()
	Error("Reticulating splines")
	if buf.String() != "==> ERROR: Reticulating splines\n" {
		t.Fail()
	}
}

func TestStructNoColor(t *testing.T) {
	var buf bytes.Buffer
	type cow struct{ says string }
	Stdout = &buf
	UseColors(false)
	defer resetOutput()
	Msg2(&cow{says: "moo"})
	if buf.String() != "  -> &{says:moo}\n" {
		t.Fail()
	}
}

func TestMsgf(t *testing.T) {
	type cow struct{ says string }
	Msgf("%s %d %f %+v", "hi", 1, 0.7, &cow{says: "moo"})
}
