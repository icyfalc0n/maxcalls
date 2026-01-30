package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/icyfalc0n/maxcalls"
	"github.com/spf13/cobra"
)

var callerCmd = &cobra.Command{
	Use:   "caller",
	Short: "Start the caller",
	Run: func(cmd *cobra.Command, args []string) {
		authToken := os.Getenv("MAXCALLS_TOKEN")
		if authToken == "" {
			fmt.Println("Error: MAXCALLS_TOKEN environment variable is required")
			os.Exit(1)
		}

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
