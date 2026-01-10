package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/icyfalc0n/max_calls_api/internal/api/oneme"
)

func Auth(client oneme.ApiClient) {
	reader := StdinReader{Reader: bufio.NewReader(os.Stdin)}

	fmt.Print("Enter phone number: ")
	phone := reader.Read()
	verificationToken, err := client.DoVerificationRequest(phone)
	if err != nil {
		panic(err)
	}

	fmt.Print("SMS code: ")
	code := reader.Read()
	login, err := client.DoCodeEnter(verificationToken.Token, code)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Login token: %s", login.TokenAttributes.Login.Token)
}
