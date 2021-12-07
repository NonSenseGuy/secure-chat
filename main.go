package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

const (
	// CONN_HOST = "localhost"
	CONN_PORT = "9090"
	CONN_TYPE = "tcp"
)

var (
	CONN_HOST = os.Getenv("SERVER_HOST")
)

func main() {
	fmt.Println("Escribe: \n 1 - para hostear una sesion de chat\n 2 - para unirte a una sesion de chat")
	var response int

	_, err := fmt.Scanf("%d", &response)
	if err != nil {
		panic(err)
	}

	if response == 1 {
		initAsServer()
	} else {
		initAsClient()
	}
}

func initAsServer() {
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	// Close the listener when the application closes.
	defer l.Close()

	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		fmt.Println("Connection accepted.")

		encoder := gob.NewEncoder(conn)
		decoder := gob.NewDecoder(conn)
		go receiveMessages(decoder)
		sendMessages(encoder)
	}
}

func initAsClient() {
	fmt.Println("hola soy cliente")
	con, error := net.Dial(CONN_TYPE, CONN_HOST+":"+CONN_PORT)

	// Handles eventual errors
	if error != nil {
		fmt.Println(error)
		return
	}

	fmt.Println("Connected to " + CONN_HOST + ":" + CONN_PORT)

	encoder := gob.NewEncoder(con)
	decoder := gob.NewDecoder(con)
	go receiveMessages(decoder)
	sendMessages(encoder)

}

func receiveMessages(decoder *gob.Decoder) {
	var message string

	for {
		err := decoder.Decode(&message)

		// Checks for errors
		if err != nil {
			fmt.Println(err)
			// Exit the loop
			return
		}
		fmt.Print(">>>", message)
	}
}

func sendMessages(encoder *gob.Encoder) {
	reader := bufio.NewReader(os.Stdin)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		err = encoder.Encode(message)
		if err != nil {
			panic(err)
		}
	}
}
