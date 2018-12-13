// Package sender contains everything used to send files
package sender

import (
	"bytes"
	"errors"
	"log"
	"net"
	"os"

	"../comlib"
)

// OpenFile opens the file, reads in and returns an array ready to be sent to the receiver and the file's size
func OpenFile(fileName string) []byte {
	file, err := os.Open(fileName)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer file.Close()

	fileStat, _ := file.Stat()
	fileSize := fileStat.Size()
	var data = make([]byte, fileSize)
	file.Read(data)

	return data
}

// Connect connects to the receiver, sends the password, which will be checked by the receiver
// and returns the Conn stream
func Connect(receiverAddr string, password []byte) (net.Conn, error) {
	log.Printf("Attempting a connection to receiver at: %s\n", receiverAddr)
	connToReceiver, err := net.Dial("tcp", receiverAddr)
	if err != nil {
		log.Println(err)
		return nil, nil
	}
	log.Printf("Connection success, attempting to login using password: %s\n", password)

	// Sending the password
	n, _ := connToReceiver.Write(password)
	log.Printf("Sent %d bytes\n", n)

	// Waiting for password confirmation
	var passwordConfirmation = make([]byte, len(comlib.PasswordValid))
	connToReceiver.Read(passwordConfirmation)

	// Reacting to receiver's response for password
	if bytes.Equal(comlib.PasswordUnvalid, passwordConfirmation) {
		// If password is wrong, receiver closes the connection. By safety, sender closes it too
		log.Println(comlib.UnvalidMessage)
		connToReceiver.Close()
		return nil, errors.New(comlib.UnvalidError)
	}

	log.Println(comlib.ValidMessage)
	return connToReceiver, nil
}
