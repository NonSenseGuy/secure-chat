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

		diffieHellman(conn)      // Generates a key with diffie hellman algorithm and stores it in a variable
		go receiveMessages(conn) // Creates a thread and receive messages in it
		sendMessages(conn)       // Uses main thread to send messages to chat

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

	diffieHellman(conn)      // Generates a key with diffie hellman algorithm and stores it in a variable
	go receiveMessages(conn) // Creates a thread and receive messages in it
	sendMessages(conn)       // Uses main thread to send messages to chat

}

func receiveMessages(conn net.Conn) {
	for {
		buffer := make([]byte, 8192) // Initialize buffer

		length, err := conn.Read(buffer) // Read message from tcp connection
		if err != nil {
			panic(err)
		}
		fmt.Println("enc>>>", base64.StdEncoding.EncodeToString(buffer[:length])) // Prints encrypted string received from connection
		fmt.Println("dec>>>", string(AESDecrypt(buffer[:length])))                // Decrypts messsage into readable string
	}

}

// sendMessages reads string from stdin, encrypts and send it to chat
func sendMessages(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	for {
		message, _, err := reader.ReadLine() // Read line from stdin
		if err != nil {
			panic(err)
		}

		encryptedMessage := AESEncrypt(message) // Encrypts string with diffie hellman key obtained at initialization

		_, err = conn.Write(encryptedMessage) // Writes encrypted message to tcp connection
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

	length, err := conn.Read(remoteKey) // This will wait until someone writes into the port therefore the chat will only initialize when both users run the chat
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

// AESEncrypt encrypts text with key stored at chat initialization
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

// AESDecrypt decrypt text with key stored at chat initialization - Is the same function to encrypt and decrypt since AES is simetrical
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
