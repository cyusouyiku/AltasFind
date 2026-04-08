package logger

import (
	"log"
	"os"
)

var (
	base    = log.New(os.Stderr, "[gofind] ", log.LstdFlags)
	verbose bool
)

func SetVerbose(v bool) {
	verbose = v
}

func Infof(format string, args ...any) {
	base.Printf("INFO "+format, args...)
}

func Errorf(format string, args ...any) {
	base.Printf("ERROR "+format, args...)
}

func Debugf(format string, args ...any) {
	if verbose {
		base.Printf("DEBUG "+format, args...)
	}
}
