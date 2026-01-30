package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/icyfalc0n/maxcalls"
	"github.com/spf13/cobra"
)

var callerCmd = &cobra.Command{
	Use:   "caller [token-file]",
	Short: "Start the caller",
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

		runCaller(calls)
	},
}

func init() {
	rootCmd.AddCommand(callerCmd)
}

func runCaller(calls maxcalls.Calls) {
	reader := StdinReader{Reader: bufio.NewReader(os.Stdin)}
	fmt.Printf("Caller ID: %s\n", calls.ExternalUserId())
	fmt.Printf("Call taker ID: ")
	calltakerExternalID := reader.Read()

	_, err := calls.Call(calltakerExternalID)
	if err != nil {
		panic(err)
	}
}
