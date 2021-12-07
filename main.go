package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "1234"
	CONN_TYPE   = "tcp"
)

var (
	CONN_HOST = os.Getenv("CONN_HOST")
	CONN_PORT = os.Getenv("CONN_PORT")
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
	l, err := net.Listen(CONN_TYPE, SERVER_HOST+":"+SERVER_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	// Close the listener when the application closes.
	defer l.Close()

	fmt.Println("Listening on " + SERVER_HOST + ":" + SERVER_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		fmt.Println("Connection accepted.")

		go receiveMessages(conn)
		sendMessages(conn)
	}
}

func initAsClient() {
	fmt.Println("hola soy cliente")
	conn, error := net.Dial(CONN_TYPE, SERVER_HOST+":"+SERVER_PORT)

	// Handles eventual errors
	if error != nil {
		fmt.Println(error)
		return
	}

	fmt.Println("Connected to " + SERVER_HOST + ":" + SERVER_PORT)

	go receiveMessages(conn)
	sendMessages(conn)

}

func receiveMessages(conn net.Conn) {
	for {
		buffer := make([]byte, 8192)

		length, err := conn.Read(buffer)
		if err != nil {
			panic(err)
		}
		fmt.Println(">>>", string(buffer[:length]))
	}

}

func sendMessages(conn net.Conn) {

	reader := bufio.NewReader(os.Stdin)
	for {
		message, _, err := reader.ReadLine()
		if err != nil {
			panic(err)
		}

		_, err = conn.Write(message)

		if err != nil {
			panic(err)
		}
	}
}
