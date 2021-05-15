package redis

import (
	"bufio"
	"fmt"
	"net"

	"github.com/holmes89/chickaree-db/pkg/core"
)

type Client struct {
	// incoming chan string
	outgoing chan []byte
	reader   *bufio.Reader
	writer   *bufio.Writer
	conn     net.Conn
	repo     core.Repository
}

func (client *Client) Read() {
	for {
		req, err := NewRequest(client.reader)
		if err == nil {
			client.outgoing <- client.Handle(req)
		} else {
			break
		}

	}

	client.conn.Close()
	fmt.Println("client disconnected")
	client = nil
}

func (client *Client) Write() {
	for data := range client.outgoing {
		client.writer.Write(data)
		client.writer.Flush()
	}
}

func (client *Client) Listen() {
	go client.Read()
	go client.Write()
}

func NewClient(connection net.Conn, repo core.Repository) *Client {
	if connection == nil {
		panic("no connection")
	}
	writer := bufio.NewWriter(connection)
	reader := bufio.NewReader(connection)

	client := &Client{
		outgoing: make(chan []byte),
		conn:     connection,
		reader:   reader,
		writer:   writer,
		repo:     repo,
	}
	client.Listen()

	return client
}
