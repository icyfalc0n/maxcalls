# Usage Guide

This guide demonstrates how to use the MAX Calls API library to make and receive calls through the MAX messenger.

## Installation

```bash
go get github.com/icyfalc0n/maxcalls
```

## Basic Setup

### Creating a Calls Client
First, you need to create a `Calls` client with your MAX messenger authentication token:
```go
calls, err := maxcalls.NewCalls(authToken)
defer calls.Close() // Always close when done
```
The `NewCalls` function performs authentication and establishes connections to the MAX messenger APIs. 
Always call `Close()` when you're done to clean up resources.

### Getting Your User ID
You can retrieve your external user ID after authentication:
```go
userID := calls.ExternalUserId()
fmt.Printf("My user ID: %s\n", userID)
```

## Making Outgoing Calls
To initiate a call to another user, use the `Call` method with their user ID:
```go
targetUserID := "target-user-id-here"

conn, err := calls.Call(targetUserID)
defer conn.Close() // Close the connection when done

// Connection established!
```
The `Call` method:
1. Starts a conversation session
2. Establishes a signaling connection
3. Creates an ICE agent for peer-to-peer connection
4. Connects using ICE to establish the media connection

## Receiving Incoming Calls
To wait for and accept incoming calls, use the `WaitForCall` method:
```go
conn, err := calls.WaitForCall()

defer conn.Close() // Close the connection when done

// Connection established! You can now use conn.conn to send/receive media
```

## Resource Management
Always ensure proper cleanup:
1. **Close connections**: Call `conn.Close()` when done with a call
2. **Close calls client**: Call `calls.Close()` when done with the client
