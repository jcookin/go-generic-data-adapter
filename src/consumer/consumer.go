package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	zmq "github.com/pebbe/zmq4"
)

type MessageHandler struct {
	socket *zmq.Socket
}
type DataStat struct {
	msgNumber   int64
	transitTime int64
}

var msgs [][]string
var stats []DataStat

func NewMessageHandler(endpoint string, socketType zmq.Type) (*MessageHandler, error) {
	socket, err := zmq.NewSocket(socketType)
	if err != nil {
		return nil, err
	}

	err = socket.Bind(endpoint)
	if err != nil {
		return nil, err
	}

	return &MessageHandler{
		socket: socket,
	}, nil
}

func (mh *MessageHandler) Close() {
	mh.socket.Close()
}

var msgCount int64 = 0
var max int64 = 0
var min int64 = 999999999999999999

func (mh *MessageHandler) StartHandlingMessages() {

	for {
		message, err := mh.socket.Recv(0)
		if err != nil {
			log.Println("Error receiving message:", err)
			continue
		}
		msgCount++

		// msgs = append(msgs, []string{message, string(rune(time.Now().UnixMilli()))})
		receiveT := time.Now().UnixNano()
		// fmt.Println(message)

		sendT, err := strconv.ParseInt(strings.Split(message, ":")[0], 10, 0)
		if err != nil {
			fmt.Println(err)
		}
		transitTime := (receiveT - sendT)
		fmt.Printf("Transit time: %d \n", transitTime)
		// stats = append(stats, DataStat{msgCount, transitTime})

		if transitTime > max {
			max = transitTime
		}
		if transitTime < min {
			min = transitTime
		}

		fmt.Printf("min (ns):%d  max(ms):%d\n", min, max)

		// Handle the received message here
		// fmt.Println("Received message:", message)
	}
}

func main() {
	port := os.Getenv("ADAPTER_PORT")
	host := os.Getenv("ADAPTER_HOST")

	if port == "" {
		port = "5555" // Default port if ADAPTER_PORT is not set
	}
	if host == "" {
		host = "tcp://127.0.0.1"
	}

	endpoint := fmt.Sprintf("%s:%s", host, port)
	socketType := zmq.PULL

	messageHandler, err := NewMessageHandler(endpoint, socketType)
	if err != nil {
		log.Fatal("Error creating message handler:", err)
	}
	defer messageHandler.Close()

	fmt.Println("Message handler listening on", endpoint)
	messageHandler.StartHandlingMessages()

	// Give some time for the message handling to occur
	time.Sleep(time.Minute)

	// defer func() {
	// 	fmt.Println("begin exit processing")
	// 	for c, msg := range msgs {
	// 		//strip out the send-time
	// 		sendT, err := strconv.ParseInt(strings.Split(msg[0], ":")[0], 10, 0)
	// 		if err != nil {
	// 			fmt.Println(err)
	// 		}
	// 		//strip out the receive time
	// 		receiveT, err := strconv.ParseInt(msg[1], 10, 0)
	// 		if err != nil {
	// 			fmt.Println(err)
	// 		}
	// 		//calculate the difference between the two
	// 		fmt.Printf("%d: %d ms", c, receiveT-sendT)
	// 	}
	// }()
}
