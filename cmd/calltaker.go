package main

import (
	"fmt"

	"github.com/icyfalc0n/maxcalls"
)

func Calltaker(calls maxcalls.Calls) {
	fmt.Printf("Call taker ID: %s\n", calls.ExternalUserId())
	fmt.Println("Waiting for incoming calls....")

	_, err := calls.WaitForCall()
	if err != nil {
		panic(err)
	}
}
