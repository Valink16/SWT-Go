// Package receiver contains everything used to receive files
package receiver

import (
	"errors"
	"log"
	"net"
	"os"
	"strconv"

	"../comlib"
)

// WriteData writes received data into a file
func WriteData(fileName string, data []byte) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		log.Fatalln(err)
	}

}

// CheckLogin listens on 'port' and wait for the sender to send the password and checks it with 'password'
func CheckLogin(port int, password []byte) (net.Conn, error) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Listening on port: %d ...\n", port)

	connToSender, err := listener.Accept()
	if err != nil {
		return nil, err
	}
	log.Printf("Accepted sender's connection from %s. Waiting for password ...\n", connToSender.RemoteAddr().String())

	// Waiting for password
	var receivedPassword = make([]byte, 100)
	connToSender.Read(receivedPassword)
	log.Printf("Received password: %s\n", receivedPassword)
	log.Printf("Real password: %s\n", password)

	// Password checking
	if !comlib.ComparePasswords(receivedPassword, password) {
		log.Println("Password is incorrect, sending rejection to sender")
		connToSender.Write(comlib.PasswordUnvalid)
		// If password isn't correct, we immediately close the connection
		connToSender.Close()
		return nil, errors.New(comlib.RejectError)
	}

	log.Println("Password is correct, sending confirmation to sender")
	connToSender.Write(comlib.PasswordValid)

	return connToSender, nil
}
