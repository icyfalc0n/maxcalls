package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/icyfalc0n/max_calls_api/api/oneme"
)

type StdinReader struct {
	Reader *bufio.Reader
}

func (r *StdinReader) Read() string {
	readed, _ := r.Reader.ReadString('\n')
	return readed[:len(readed)-1]
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("expected 'auth', 'calltaker' or 'caller' subcommands")
		os.Exit(1)
	}

	client, err := oneme.NewOnemeApiClient()
	if err != nil {
		panic(err)
	}

	defer func() {
		err = client.RawClient.Close()
		if err != nil {
			panic(err)
		}
	}()

	switch os.Args[1] {
	case "auth":
		Auth(client)
	case "calltaker":
		Calltaker(client)
	case "caller":
		Caller(client)
	default:
		fmt.Println("expected 'auth', 'calltaker' or 'caller' subcommands")
		os.Exit(1)
	}
}
