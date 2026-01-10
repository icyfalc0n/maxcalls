package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/icyfalc0n/maxcalls"
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

func main() {
	if len(os.Args) < 3 {
		fmt.Println("expected 'calltaker' or 'caller' subcommands and token file")
		os.Exit(1)
	}

	authTokenBytes, err := os.ReadFile(os.Args[2])
	if err != nil {
		panic(err)
	}
	authToken := strings.TrimSpace(string(authTokenBytes))

	calls, err := maxcalls.NewCalls(LoggerImpl{}, authToken)
	if err != nil {
		panic(err)
	}

	switch os.Args[1] {
	case "calltaker":
		Calltaker(calls)
	case "caller":
		Caller(calls)
	default:
		fmt.Println("expected 'calltaker' or 'caller' subcommands and token file")
		os.Exit(1)
	}
}
