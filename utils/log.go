package utils

import (
	"fmt"
	"io"
	"os"
)

var (
	ALL_OFF string = "\x1b[0m"
	BOLD           = "\x1b[1m"
	BLUE           = BOLD + "\x1b[34m"
	GREEN          = BOLD + "\x1b[32m"
	RED            = BOLD + "\x1b[31m"
	YELLOW         = BOLD + "\x1b[33m"

	formatSuffix = ALL_OFF + BOLD + "%s" + ALL_OFF + "\n"
)

var Stdout io.Writer = os.Stdout

func Plain(msg ...interface{}) {
	fmt.Fprintf(Stdout, formatSuffix, msg...)
}

func Msg(msg ...interface{}) {
	format := GREEN + "==> " + formatSuffix
	fmt.Fprintf(Stdout, format, msg...)
}

func Msg2(msg ...interface{}) {
	format := BLUE + "  -> " + formatSuffix
	fmt.Fprintf(Stdout, format, msg...)
}

func Warning(msg ...interface{}) {
	format := YELLOW + "==> WARNING: " + formatSuffix
	fmt.Fprintf(Stdout, format, msg...)
}

func Error(msg ...interface{}) {
	format := RED + "==> ERROR: " + formatSuffix
	fmt.Fprintf(Stdout, format, msg...)
}
