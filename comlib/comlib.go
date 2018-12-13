// Package comlib contains some common values used in the program
// comlib contains some byte values because it's simpler to communicate with the byte type
package comlib

import (
	"bufio"
	"fmt"
)

// PasswordValid is used when password is correct/valid
var PasswordValid = []byte("pw__valid")

// PasswordUnvalid is used when password is uncorrect/unvalid, needs to be the same length as PasswordValid
var PasswordUnvalid = []byte("pwunvalid")

// UnvalidMessage is shown when password isn't valid
var UnvalidMessage = "Password is incorrect. Terminating connection."

// ValidMessage is shown when password is valid
var ValidMessage = "Logged in successfully. Proceeding."

// UnvalidError is used in the error raised when password is wrong
var UnvalidError = "passwordError: used password is not valid, cannot login"

// RejectError is used by receiver in the error raised when received password is wrong
var RejectError = "passwordError: received password is not valid, terminating connection"

// DefaultPort is the port used by the receiver and the sender to communicate
var DefaultPort = 30000

// ConfirmationMsg is used by the receiver and the sender to confirm things
var ConfirmationMsg = []byte("OK")

// ComparePasswords compare the passwords and return true if 'received' is right
func ComparePasswords(received []byte, real []byte) bool {
	for i := 0; i < len(real) && i < len(received); i++ {
		if received[i] == 0 && real[i] == 0 {
			break
		} else {
			if received[i] != real[i] {
				return false
			}
		}
	}

	return true
}

// InputStr gets input, removes '\n' and returns a string
func InputStr(msg string, reader *bufio.Reader) string {
	fmt.Print(msg)
	data, _ := reader.ReadString('\n')
	return data[:len(data)-1]
}

// InputBytes gets input, removes '\n' and returns a slice of bytes
func InputBytes(msg string, reader *bufio.Reader) []byte {
	fmt.Print(msg)
	data, _ := reader.ReadBytes('\n')
	return data[:len(data)-1]
}

// TrunkData removes the 0 values from and returns a slice
func TrunkData(data []byte) []byte {
	var l = len(data)
	for i, b := range data {
		if b == 0 {
			return data[:l-(l-i)]
		}
	}
	return data
}
