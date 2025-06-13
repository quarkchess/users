package util

import (
	"log"
	"os"
	"strings"
)

func NewLogger(prefix string) *log.Logger {
	p := strings.TrimSpace(prefix) + "\t"

	return log.New(os.Stdout, p, log.LstdFlags)
}
