package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/icyfalc0n/max_calls_api"
)

func Caller(calls max_calls_api.Calls) {
	reader := StdinReader{Reader: bufio.NewReader(os.Stdin)}
	fmt.Printf("Caller ID: %s\n", calls.ExternalUserId())
	fmt.Printf("Call taker ID: ")
	calltakerExternalID := reader.Read()

	_, err := calls.Call(calltakerExternalID)
	if err != nil {
		panic(err)
	}
}
