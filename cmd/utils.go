package main

import (
	"bufio"
	"log"
)

type StdinReader struct {
	Reader *bufio.Reader
}

func (r *StdinReader) Read() string {
	read, _ := r.Reader.ReadString('\n')
	return read[:len(read)-1]
}

type LoggerImpl struct{}

func (l LoggerImpl) Debugf(format string, args ...interface{}) {
	log.Printf(format, args...)
}
