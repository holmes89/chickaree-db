package core

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Client struct {
	// incoming chan string
	outgoing chan []byte
	reader   *bufio.Reader
	writer   *bufio.Writer
	conn     net.Conn
	repo     Repository
}

func (client *Client) Read() {
	for {
		req, err := NewRequest(client.reader)
		if err == nil {
			switch strings.ToLower(req.command) {
			case "command":
				client.outgoing <- OkResp.Encode()
			case "ping":
				client.outgoing <- Response{
					rtype:   SimpleString,
					content: []byte("PONG"),
				}.Encode()
			case "set":
				client.outgoing <- client.repo.Set(req.args).Encode()
			case "hset":
				client.outgoing <- client.repo.HSet(req.args).Encode()
			case "hget":
				client.outgoing <- client.repo.HGet(req.args).Encode()
			case "hexists":
				client.outgoing <- client.repo.HExists(req.args).Encode()
			case "hgetall":
				client.outgoing <- client.repo.HGetAll(req.args).Encode()
			case "get":
				client.outgoing <- client.repo.Get(req.args).Encode()
			case "del":
				client.outgoing <- client.repo.Del(req.args).Encode()
			default:
				err := fmt.Errorf("unknown command '%s'", req.command)
				client.outgoing <- ErrResponse(err).Encode()
			}
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

func NewClient(connection net.Conn, repo Repository) *Client {
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
