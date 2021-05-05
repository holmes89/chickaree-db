package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"text/scanner"
)

var allClients map[*Client]int

type Client struct {
	// incoming chan string
	outgoing   chan string
	reader     *bufio.Reader
	writer     *bufio.Writer
	conn       net.Conn
	connection *Client
}

func (client *Client) Read() {
	for {
		req, err := parseRequest(client.reader)
		if err == nil {
			switch strings.ToLower(req.command) {
			case "command":
				client.outgoing <- "+OK\r\n"
			case "ping":
				client.outgoing <- "+PONG\r\n"
			default:
				fmt.Printf("DONT KNOW: %s\n", req.command)
			}
			fmt.Printf("%+v\n", req)
		} else {
			break
		}

	}

	client.conn.Close()
	delete(allClients, client)
	if client.connection != nil {
		client.connection.connection = nil
	}
	client = nil
}

func (client *Client) Write() {
	for data := range client.outgoing {
		fmt.Println(data)
		client.writer.WriteString(data)
		client.writer.Flush()
	}
}

func (client *Client) Listen() {
	go client.Read()
	go client.Write()
}

func NewClient(connection net.Conn) *Client {
	if connection == nil {
		panic("no connection")
	}
	writer := bufio.NewWriter(connection)
	reader := bufio.NewReader(connection)

	client := &Client{
		outgoing: make(chan string),
		conn:     connection,
		reader:   reader,
		writer:   writer,
	}
	client.Listen()

	return client
}

func main() {

	allClients = make(map[*Client]int)

	var port string
	flag.StringVar(&port, "port", "6379", "port to listen on")

	flag.Parse()

	PORT := ":" + port
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()

	fmt.Printf("listening on port %s\n", PORT)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err.Error())
		}
		if conn == nil {
			panic("no conn")
		}
		client := NewClient(conn)
		for clientList, _ := range allClients {
			if clientList.connection == nil {
				client.connection = clientList
				clientList.connection = client
				fmt.Println("Connected")
			}
		}
		allClients[client] = 1
		fmt.Println(len(allClients))
	}
}

// func main() {
// 	var port string
// 	flag.StringVar(&port, "port", "6379", "port to listen on")

// 	PORT := ":" + port
// 	l, err := net.Listen("tcp", PORT)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer l.Close()

// 	fmt.Printf("listening on port %s\n", PORT)

// 	c, err := l.Accept()
// 	if err != nil {
// 		fmt.Printf("error: %v\n", err)
// 		return
// 	}

// 	for {
// 		req, err := Read(c)
// 		if err != nil {
// 			fmt.Printf("err: %v", err)
// 			c.Close()
// 			os.Exit(0)
// 		}
// 		fmt.Printf("%+v\n", req)
// 		switch req.command {
// 		case "ping":
// 			c.Write([]byte("+PONG\r\n"))
// 		default:
// 			fmt.Printf("unknown: %v\n", req.command)
// 			c.Write([]byte(fmt.Sprintf("-ERR unknown command '%s'\r\n", req.command)))
// 		}
// 	}
// }

type lexer struct {
	*scanner.Scanner
	buf bytes.Buffer
}

func parseRequest(r io.Reader) (req Request, err error) {

	req.msgCount, err = getSize(r)
	if err != nil {
		return req, err
	}
	for i := req.msgCount; i > 0; i-- {
		bufsize, err := getSize(r)
		if err != nil {
			return req, err
		}
		bufsize = bufsize + 2 // capture carriage return and newline
		b := make([]byte, bufsize)

		c, err := r.Read(b)
		if err != nil {
			return req, err
		}
		if c != int(bufsize) {
			return req, io.EOF
		}
		req.args = append(req.args, string(b[:len(b)-2]))
	}
	req.command = req.args[0]
	req.args = req.args[1:]

	return req, nil
}

func getSize(r io.Reader) (int, error) {
	buf := bytes.NewBuffer([]byte{})
	b := make([]byte, 1)
	creturn := false
	for {
		c, err := r.Read(b)
		if err != nil {
			return 0, err
		}
		if c == 0 {
			return 0, io.EOF
		}

		if b[0] == byte('\r') {
			creturn = true
			continue
		}

		if creturn && b[0] == byte('\n') {
			break
		}

		creturn = false
		buf.Write(b)
	}
	num, _ := strconv.Atoi(string(buf.Bytes()[1:]))
	return num, nil
}

// New returns new lexer
// TODO this should be rewritten to respect the values that are being passed by Redis itself to define byte array size and accurately read closed connections
// func Read(r io.Reader) (Request, error) {
// 	scanner := bufio.NewScanner(r)
// 	scanner.Split(bufio.ScanLines)

// 	req := Request{}

// 	for scanner.Scan() {

// 		if req.argCount == 0 {
// 			req.argCount, _ = strconv.Atoi(scanner.Text()[1:])
// 			continue
// 		}
// 		if scanner.Text()[0] == '$' {
// 			continue
// 		}
// 		if req.command == "" {
// 			req.command = strings.ToLower(scanner.Text())
// 		} else {
// 			req.args = append(req.args, scanner.Text())
// 		}

// 		if req.command != "" && len(req.args) == (req.argCount-1) { // don't count command
// 			break
// 		}
// 	}

// 	if scanner.Err() != nil {
// 		fmt.Println(scanner.Err())
// 		return req, io.EOF
// 	}

// 	return req, nil
// }

type Request struct {
	msgCount int
	args     []string
	command  string
}
