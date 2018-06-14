package utils

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
)

var (
	Stdout io.Writer = os.Stdout
	Stderr io.Writer = os.Stderr
)

var (
	red    = color.New(color.FgRed).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
	blue   = color.New(color.FgBlue).SprintFunc()
	bold   = color.New(color.Bold).SprintFunc()
	Plain  = color.New().PrintlnFunc()
)

func Msg(msg interface{}) {
	fmt.Fprintln(Stdout, green("==> ") + bold(fmt.Sprintf("%+v", msg)))
}

func Msg2(msg interface{}) {
	fmt.Fprintln(Stdout, blue("  -> ") + bold(fmt.Sprintf("%+v", msg)))
}

func Warning(msg interface{}) {
	fmt.Fprintln(Stdout, yellow("==> WARNING: ") + bold(fmt.Sprintf("%+v", msg)))
}

func Error(msg interface{}) {
	fmt.Fprintln(Stdout, red("==> ERROR: ") + bold(fmt.Sprintf("%+v", msg)))
}

func UseColors(use bool) {
	color.NoColor = !use
}
