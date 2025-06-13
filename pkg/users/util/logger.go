package util

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Logger struct {
	inner *log.Logger
}

func NewLogger(prefix string) Logger {
	p := strings.TrimSpace(prefix) + " "
	inner := log.New(os.Stdout, p, log.LstdFlags|log.Lmsgprefix)

	return Logger{inner}
}

func (l *Logger) Infoln(args ...any) {
	l.inner.Println("INFO:", args)
}

func (l *Logger) Infof(format string, args ...any) {
	l.inner.Printf(fmt.Sprintf("INFO: %s", format), args)
}

func (l *Logger) Errorln(args ...any) {
	l.inner.Println("ERROR:", args)
}

func (l *Logger) Errorf(format string, args ...any) {
	l.inner.Printf(fmt.Sprintf("ERROR: %s", format), args)
}

func (l *Logger) Fatalln(args ...any) {
	l.inner.Fatalln("FATAL:", args)
}

func (l *Logger) Fatalf(format string, args ...any) {
	l.inner.Fatalf(fmt.Sprintf("FATAL: %s", format), args)
}
