package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	zmq "github.com/pebbe/zmq4"
)

type MessagePublisher struct {
	socket *zmq.Socket
}

const MSG_SIZE int = 1013

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var msgSendCounter int = 0
var startTime int64 = 0
var msg_size_calc = 0

func init() {
	rand.Seed(time.Now().UnixNano())
	startTime = time.Now().UnixMilli()
	msg_size_calc = len(GenerateMessage(RandSeq(MSG_SIZE)))
	fmt.Printf("Publishing messages of %d bytes\n", msg_size_calc)

}

func NewMessagePublisher(endpoint string, socketType zmq.Type) (*MessagePublisher, error) {
	socket, err := zmq.NewSocket(socketType)
	if err != nil {
		return nil, err
	}

	err = socket.Connect(endpoint)
	if err != nil {
		return nil, err
	}

	return &MessagePublisher{
		socket: socket,
	}, nil
}

func (mp *MessagePublisher) Close() {
	mp.socket.Close()
}

func (mp *MessagePublisher) PublishMessage(message []byte) error {
	_, err := mp.socket.SendBytes(message, 0)
	return err
}

func RandSeq(n int) string {
	b := make([]rune, n-2) // -2 leaves room for prefix to message ": "
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return ": " + string(b)
}

func GenerateMessage(msgBody string) string {
	var timestamp int64 = time.Now().UnixNano()
	message := fmt.Sprintf("%d: %s", timestamp, msgBody) // why is this 13 bytes larger than expected?
	return message
}

func PublishMessage(endpoint string, messagePublisher *MessagePublisher) {
	var now int64 = time.Now().UnixMilli()
	if now-startTime > 1000 {
		fmt.Printf("Total bytes sent: %d bytes (%2d MB)\n", msgSendCounter*msg_size_calc, msgSendCounter*msg_size_calc/1000000.0)
		startTime = now
		msgSendCounter = 0
	}
	var message string = GenerateMessage(RandSeq(MSG_SIZE))
	err := messagePublisher.PublishMessage([]byte(message))
	if err != nil {
		log.Fatal("Error publishing message:", err)
	}
	msgSendCounter++
}

func main() {
	port := os.Getenv("ADAPTER_PORT")              // defaults to 5555
	host := os.Getenv("ADAPTER_HOST")              // defaults to 127.0.0.1
	numMessagesStr := os.Getenv("NUMBER_MSG_SEND") // defaults to 50

	if port == "" {
		port = "5555" // Default port if ADAPTER_PORT is not set
	}
	if host == "" {
		host = "tcp://127.0.0.1"
	}
	var numMessages int = 0
	if numMessagesStr != "" {
		tmpA, err := strconv.ParseInt(numMessagesStr, 10, 16)
		if err != nil {
			log.Fatal(err)
		}

		numMessages = int(tmpA)
	}

	endpoint := fmt.Sprintf("%s:%s", host, port)
	socketType := zmq.PUSH

	messagePublisher, err := NewMessagePublisher(endpoint, socketType)
	if err != nil {
		log.Fatal("Error creating message publisher:", err)
	}
	defer messagePublisher.Close()

	if numMessages == 0 {
		for true {
			PublishMessage(endpoint, messagePublisher)
		}
	} else {
		for msgNum := 0; msgNum < numMessages; msgNum++ {
			PublishMessage(endpoint, messagePublisher)
		}
	}

	// Give some time for the message publishing to occur
	time.Sleep(time.Second)
}
