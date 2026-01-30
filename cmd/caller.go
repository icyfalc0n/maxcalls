package main

import (
	"fmt"
	"os"

	"github.com/icyfalc0n/maxcalls"
	"github.com/spf13/cobra"
)

var callerCmd = &cobra.Command{
	Use:   "caller [calltaker_id]",
	Short: "Start the caller",
	Args:  cobra.ExactArgs(1),
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

		runCaller(calls, args[0])
	},
}

func init() {
	rootCmd.AddCommand(callerCmd)
}

func runCaller(calls maxcalls.Calls, calltakerID string) {
	fmt.Printf("Caller ID: %s\n", calls.ExternalUserId())
	fmt.Printf("Calling: %s\n", calltakerID)

	_, err := calls.Call(calltakerID)
	if err != nil {
		panic(err)
	}
}
