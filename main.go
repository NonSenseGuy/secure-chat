package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"net"
	"os"

	"github.com/monnand/dhkx"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "1234"
	CONN_TYPE   = "tcp"
)

var (
	CONN_HOST = os.Getenv("CONN_HOST")
	CONN_PORT = os.Getenv("CONN_PORT")
	key       = []byte("")
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

		diffieHellman(conn)
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

	diffieHellman(conn)
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
		fmt.Println("enc>>>", base64.StdEncoding.EncodeToString(buffer[:length]))
		fmt.Println("dec>>>", string(AESDecrypt(buffer[:length])))
	}

}

func sendMessages(conn net.Conn) {

	reader := bufio.NewReader(os.Stdin)
	for {
		message, _, err := reader.ReadLine()
		if err != nil {
			panic(err)
		}

		encryptedMessage := AESEncrypt(message)
		// decryptedMessage := AESDecrypt(encryptedMessage)
		// fmt.Println("<<<", base64.StdEncoding.EncodeToString(encryptedMessage))
		// fmt.Println("<<<", string(decryptedMessage))

		_, err = conn.Write(encryptedMessage)

		if err != nil {
			panic(err)
		}
	}
}

func diffieHellman(conn net.Conn) {
	// Get a group. Use the default one would be enough.
	g, _ := dhkx.GetGroup(0)

	// Generate a private key from the group.
	// Use the default random number generator.
	priv, _ := g.GeneratePrivateKey(nil)

	// Get the public key from the private key.
	pub := priv.Bytes()

	//Send the public key
	conn.Write(pub)

	// Receive a slice of bytes from remote, which contains remote public key
	remoteKey := make([]byte, 1024)

	length, err := conn.Read(remoteKey)
	if err != nil {
		panic(err)
	}
	remoteKey = remoteKey[:length]

	// Recover remote public key
	remotePubKey := dhkx.NewPublicKey(remoteKey)

	// Compute the key
	k, _ := g.ComputeKey(remotePubKey, priv)

	// Get the key in the form of []byte
	key = k.Bytes()
	key = key[:16] // 16 bytes = 128 bits

	fmt.Println("Llave compartida:", base64.StdEncoding.EncodeToString(key))

}

func AESEncrypt(plaintext []byte) []byte {

	ciphertext := make([]byte, len(plaintext))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	stream := cipher.NewCTR(block, key[:aes.BlockSize])
	stream.XORKeyStream(ciphertext, plaintext)
	return ciphertext
}

func AESDecrypt(ciphertext []byte) []byte {
	plaintext := make([]byte, len(ciphertext))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	stream := cipher.NewCTR(block, key[:aes.BlockSize])
	stream.XORKeyStream(plaintext, ciphertext)
	return plaintext
}
