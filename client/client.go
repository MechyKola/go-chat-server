package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

func read(conn *net.Conn) {
	//TODO In a continuous loop, read a message from the server and display it.
	reader := bufio.NewReader(*conn)
	for {
		msgIn, _ := reader.ReadString('\n')
		fmt.Print(msgIn)
	}
}

func write(conn *net.Conn, closer chan bool) {
	//TODO Continually get input from the user and send messages to the server.
	reader := bufio.NewReader(os.Stdin)
	for {
		msgIn, _ := reader.ReadString('\n')
		if msgIn == "exit\n" {
			fmt.Fprint(*conn, msgIn)
			fmt.Println("Disconnecting...")
			closer <- true
		} else {
			fmt.Fprint(*conn, msgIn)
		}
	}
}

func main() {
	// Get the server address and port from the commandline arguments.
	addrPtr := flag.String("ip", "127.0.0.1:8030", "IP:port string to connect to")
	flag.Parse()

	//TODO Try to connect to the server
	conn, err := net.Dial("tcp", *addrPtr)
	if err != nil {
		// handle error
		fmt.Println("Error: couldn't connect")
	} else {
		closer := make(chan bool)
		fmt.Println("Connected to the server! Type messages and press enter to send")

		//TODO Start asynchronously reading and displaying messages
		go read(&conn)

		//TODO Start getting and sending user messages.
		go write(&conn, closer)

		// exit
		<-closer
	}
}
