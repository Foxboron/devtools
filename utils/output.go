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

func UseColors(use bool) {
	color.NoColor = !use
}

// Functions where the format string is supplied

func Msgf(format string, values ...interface{}) {
	fmt.Fprintln(Stdout, green("==> ")+bold(fmt.Sprintf(format, values...)))
}

func Msg2f(format string, values ...interface{}) {
	fmt.Fprintln(Stdout, blue("  -> ")+bold(fmt.Sprintf(format, values...)))
}

func Warningf(format string, values ...interface{}) {
	fmt.Fprintln(Stdout, yellow("==> WARNING: ")+bold(fmt.Sprintf(format, values...)))
}

func Errorf(format string, values ...interface{}) {
	fmt.Fprintln(Stdout, red("==> ERROR: ")+bold(fmt.Sprintf(format, values...)))
}

// Functions where only one argument is given

func Msg(msg interface{}) {
	Msgf("%+v", msg)
}

func Msg2(msg interface{}) {
	Msg2f("%+v", msg)
}

func Warning(msg interface{}) {
	Warningf("%+v", msg)
}

func Error(msg interface{}) {
	Errorf("%+v", msg)
}
