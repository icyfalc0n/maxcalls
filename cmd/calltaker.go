package main

import (
	"fmt"
	"os"

	"github.com/icyfalc0n/maxcalls"
	"github.com/spf13/cobra"
)

var calltakerCmd = &cobra.Command{
	Use:   "calltaker",
	Short: "Start the calltaker",
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
