package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
)

// Message - Simple message
type Message struct {
	sender  int
	message string
}

func handleError(err error) {
	// TODO: all
	// Deal with an error event.
}

func acceptConns(ln net.Listener, conns chan net.Conn) {
	// TODO: all
	// Continuously accept a network connection from the Listener
	// and add it to the channel for handling connections.
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		} else {
			conns <- conn
		}
	}
}

func handleClient(client net.Conn, clientid int, msgs chan Message) {
	// TODO: all
	// So long as this connection is alive:
	// Read in new messages as delimited by '\n's
	// Tidy up each message and add it to the messages channel,
	// recording which client it came from.
	reader := bufio.NewReader(client)
	fmt.Fprintf(client, "You've been connected and assigned id of %d\n", clientid)
	run := true
	for run {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("An error occured, disconnecting %d\n", clientid)
			fmt.Print(msg)
			run = false
		} else if msg == "exit\n" {
			fmt.Printf("%d exiting chat...\n", clientid)
			msgs <- Message{clientid, "has left the chat\n"}
			run = false
		} else {
			msgs <- Message{clientid, msg}
		}
	}
	// maybe get some code to deal with disconnects
}

func main() {
	// Read in the network port we should listen on, from the commandline argument.
	// Default to port 8030
	portPtr := flag.String("port", ":8030", "port to listen on")
	flag.Parse()

	//TODO Create a Listener for TCP connections on the port given above.

	listener, _ := net.Listen("tcp", *portPtr)

	//Create a channel for connections
	conns := make(chan net.Conn)
	//Create a channel for messages
	msgs := make(chan Message)
	//Create a mapping of IDs to connections
	clients := make(map[int]net.Conn)
	currID := 0

	//Start accepting connections
	go acceptConns(listener, conns)
	for {
		select {
		case conn := <-conns:
			//TODO Deal with a new connection
			// - assign a client ID
			// - add the client to the clients channel
			// - start to asynchronously handle messages from this client
			clients[currID] = conn
			go handleClient(conn, currID, msgs)
			currID++
		case msg := <-msgs:
			//TODO Deal with a new message
			// Send the message to all clients that aren't the sender
			fmt.Printf("%d is sending %s", msg.sender, msg.message)
			for client := range clients {
				if client != msg.sender {
					fmt.Fprintf(clients[client], "%d: %s", msg.sender, msg.message)
				}
			}
		}
	}
}
