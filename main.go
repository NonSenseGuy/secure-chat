package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "9090"
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
	conn, error := net.Dial(CONN_TYPE, CONN_HOST+":"+CONN_PORT)

	// Handles eventual errors
	if error != nil {
		fmt.Println(error)
		return
	}

	fmt.Println("Connected to " + CONN_HOST + ":" + CONN_PORT)

	go receiveMessages(conn)
	sendMessages(conn)

}

func receiveMessages(conn net.Conn) {
	for {

		var buf bytes.Buffer
		io.Copy(&buf, conn)

		// message, err := ioutil.ReadAll(conn)
		// Checks for errors
		// if err != nil {
		// 	fmt.Println(err)
		// 	// Exit the loop
		// 	return
		// }

		// TODO PROCESS CRYPTOCRAPHIC KEY + MESSAGE
		if buf.Len() > 0 {
			fmt.Println("total size:", buf.Len())
			fmt.Println(">>>", buf.String())
		}

		// buf.Reset()
	}
}

func sendMessages(conn net.Conn) {

	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(conn)

	for {
		message, _, err := reader.ReadLine()
		if err != nil {
			panic(err)
		}

		_, err = writer.Write([]byte(message))
		if err != nil {
			panic(err)
		}
	}
}
