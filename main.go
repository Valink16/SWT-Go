package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"./comlib"
	"./receiver"
	"./sender"
)

func main() {
	inputReader := bufio.NewReader(os.Stdin)

	mode := strings.ToUpper(comlib.InputStr("(S)end or (Receive) ?: ", inputReader))

	if mode == "S" || mode == "SEND" {

		receiverAddr := comlib.InputStr("Enter receiver's IP (blank for 'localhost'): ", inputReader)
		if len(receiverAddr) == 0 {
			receiverAddr = "localhost"
		}
		receiverAddr = receiverAddr + ":" + strconv.Itoa(comlib.DefaultPort)

		attemptPassword := comlib.InputBytes("Enter password (blank for '0'): ", inputReader)
		if len(attemptPassword) == 0 {
			attemptPassword = []byte("0")
		}

		connToReceiver, err := sender.Connect(receiverAddr, attemptPassword)
		if err != nil {
			log.Println(err)
			return
		}

		fileName := comlib.InputStr("Enter name of the file to open: ", inputReader)
		data := sender.OpenFile(fileName)
		log.Printf("Successfully opened %s. Size: %d Bytes.", fileName, len(data))

		dataSize := []byte(strconv.Itoa(len(data)))
		log.Printf("Sending size of the file to receiver (%s) ...", dataSize)
		connToReceiver.Write(dataSize)

		log.Println("Waiting for confirmation to start sending data ...")
		confData := make([]byte, 2)
		_, _ = connToReceiver.Read(confData)

		if string(confData) == "OK" {
			log.Println("Starting to send data ...")
			connToReceiver.Write(data)
			log.Println("Done sending.")
			connToReceiver.Close()
		}

	} else if mode == "R" || mode == "RECEIVE" {
		password := comlib.InputBytes("Enter password (blank for '0'): ", inputReader)
		if len(password) == 0 {
			password = []byte("0")
		}

		connToSender, err := receiver.CheckLogin(comlib.DefaultPort, password)
		if err != nil {
			log.Fatalln(err)
		}

		log.Println("Waiting for sender to send data size ...")
		sizeData := make([]byte, 100)
		_, _ = connToSender.Read(sizeData)
		dataSize, err := strconv.Atoi(string(comlib.TrunkData(sizeData)))
		if err != nil {
			log.Fatalln(err)
		}

		data := make([]byte, dataSize)
		log.Printf("Allocated %d bytes.\n", dataSize)
		log.Println("Receiver is ready. sending confirmation ...")

		connToSender.Write([]byte("OK"))

		log.Println("Waiting for data ...")
		data, err = ioutil.ReadAll(connToSender)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("Received %d bytes.\n", len(data))
		log.Println("Done receiving.")
		connToSender.Close()

		fileName := comlib.InputStr("Enter name of file to save: ", inputReader)
		receiver.WriteData(fileName, data)
	}
}
