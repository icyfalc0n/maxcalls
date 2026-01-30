package main

import (
	"fmt"
	"net"
	"os"

	"github.com/icyfalc0n/maxcalls"
	"github.com/spf13/cobra"
)

var calltakerForwardAddr string

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
	calltakerCmd.Flags().StringVarP(&calltakerForwardAddr, "forward", "f", "127.0.0.1:3000", "UDP address to forward packets to")
}

func runCalltaker(calls maxcalls.Calls) {
	fmt.Printf("Call taker ID: %s\n", calls.ExternalUserId())
	fmt.Println("Waiting for incoming calls....")

	conn, err := calls.WaitForCall()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Printf("Call received. Redirecting packets to %s...\n", calltakerForwardAddr)

	targetAddr, err := net.ResolveUDPAddr("udp", calltakerForwardAddr)
	if err != nil {
		panic(err)
	}

	udpConn, err := net.DialUDP("udp", nil, targetAddr)
	if err != nil {
		panic(err)
	}
	defer udpConn.Close()

	// Read from Call and write to UDP
	go func() {
		buf := make([]byte, 1500)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Printf("Connection Read error: %v\n", err)
				return
			}

			_, err = udpConn.Write(buf[:n])
			if err != nil {
				fmt.Printf("Write to UDP error: %v\n", err)
				return
			}
		}
	}()

	// Read from UDP and write to Call
	buf := make([]byte, 1500)
	for {
		n, err := udpConn.Read(buf)
		if err != nil {
			fmt.Printf("Read from UDP error: %v\n", err)
			break
		}

		_, err = conn.Write(buf[:n])
		if err != nil {
			fmt.Printf("Connection Write error: %v\n", err)
			break
		}
	}
}
