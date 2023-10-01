package zmqmessagehandler

import (
	"fmt"
	"log"

	zmq "github.com/pebbe/zmq4"
)

// TODO: channels as FIFO queues

type Message struct {
	message string
}

type ZmqConnection struct {
	socket *zmq.Socket
}

var conn *ZmqConnection
var err error

// TODO: abstract socketType to general purpose processor
func NewConnection(address string, connType string) *ZmqConnection {

	if connType == "publisher" {
		conn, err = NewMessagePublisher(address)
	}
	if connType == "subscriber" {
		conn, err = NewMessageSubscriber(address)
	}

	if err != nil {
		log.Println(err)
	}

	return conn
}

func NewMessageSubscriber(address string) (*ZmqConnection, error) {
	socket, err := zmq.NewSocket(zmq.SUB)
	if err != nil {
		return nil, err
	}

	err = socket.Bind(address)
	if err != nil {
		return nil, err
	}

	return &ZmqConnection{
		socket: socket,
	}, nil
}

func NewMessagePublisher(address string) (*ZmqConnection, error) {
	socket, err := zmq.NewSocket(zmq.PUB)
	if err != nil {
		return nil, err
	}

	err = socket.Bind(address)
	if err != nil {
		return nil, err
	}

	return &ZmqConnection{
		socket: socket,
	}, nil
}

func (conn *ZmqConnection) Close() {
	conn.socket.Close()
}

func (conn *ZmqConnection) Start() {

	for {
		message, err := conn.socket.Recv(0)
		if err != nil {
			log.Println("Error receiving message:", err)
			continue
		}

		// Handle the received message here
		fmt.Println("Received message:", message)
	}
}

func (conn *ZmqConnection) GetNextMessage() {
	// get message at top of queue
}
