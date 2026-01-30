package main

import (
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/icyfalc0n/maxcalls"
	"github.com/spf13/cobra"
)

var callerListenAddr string

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
	callerCmd.Flags().StringVarP(&callerListenAddr, "listen", "l", ":2000", "UDP address to listen on")
}

func runCaller(calls maxcalls.Calls, calltakerID string) {
	fmt.Printf("Caller ID: %s\n", calls.ExternalUserId())
	fmt.Printf("Calling: %s\n", calltakerID)

	conn, err := calls.Call(calltakerID)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Printf("Listening on %s and redirecting to call...\n", callerListenAddr)

	udpAddr, err := net.ResolveUDPAddr("udp", callerListenAddr)
	if err != nil {
		panic(err)
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		panic(err)
	}
	defer udpConn.Close()

	var remoteAddr *net.UDPAddr
	var mu sync.Mutex

	// Read from UDP and write to Call
	go func() {
		buf := make([]byte, 1500)
		for {
			n, addr, err := udpConn.ReadFromUDP(buf)
			if err != nil {
				fmt.Printf("ReadFromUDP error: %v\n", err)
				return
			}

			mu.Lock()
			remoteAddr = addr
			mu.Unlock()

			_, err = conn.Write(buf[:n])
			if err != nil {
				fmt.Printf("Connection Write error: %v\n", err)
				return
			}
		}
	}()

	// Read from Call and write to UDP
	buf := make([]byte, 1500)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("Connection Read error: %v\n", err)
			break
		}

		mu.Lock()
		addr := remoteAddr
		mu.Unlock()

		if addr != nil {
			_, err = udpConn.WriteToUDP(buf[:n], addr)
			if err != nil {
				fmt.Printf("WriteToUDP error: %v\n", err)
			}
		}
	}
}
