package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/icyfalc0n/maxcalls"
	"github.com/spf13/cobra"
)

var calltakerCmd = &cobra.Command{
	Use:   "calltaker [token-file]",
	Short: "Start the calltaker",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tokenFile := args[0]
		authTokenBytes, err := os.ReadFile(tokenFile)
		if err != nil {
			panic(err)
		}
		authToken := strings.TrimSpace(string(authTokenBytes))

		calls, err := maxcalls.NewCalls(LoggerImpl{}, authToken)
		if err != nil {
			panic(err)
		}

		runCalltaker(calls)
	},
}

func init() {
	rootCmd.AddCommand(calltakerCmd)
}

func runCalltaker(calls maxcalls.Calls) {
	fmt.Printf("Call taker ID: %s\n", calls.ExternalUserId())
	fmt.Println("Waiting for incoming calls....")

	_, err := calls.WaitForCall()
	if err != nil {
		panic(err)
	}
}
