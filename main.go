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

func main() {
	var port string
	flag.StringVar(&port, "port", "6379", "port to listen on")

	PORT := ":" + port
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	fmt.Printf("listening on port %s\n", PORT)

	c, err := l.Accept()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	for {

		req, err := Read(c)
		if err != nil {
			fmt.Printf("err: %v", err)
		}
		fmt.Printf("%+v\n", req)
		switch req.command {
		case "ping":
			c.Write([]byte("+PONG\r\n"))
		default:
			fmt.Printf("unknown: %v\n", req.command)
			c.Write([]byte(fmt.Sprintf("-ERR unknown command '%s'\r\n", req.command)))
		}
	}
}

type lexer struct {
	*scanner.Scanner
	buf bytes.Buffer
}

// New returns new lexer
func Read(r io.Reader) (Request, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	req := Request{}

	for scanner.Scan() {
		if req.argCount == 0 {
			req.argCount, _ = strconv.Atoi(scanner.Text()[1:])
			continue
		}
		if scanner.Text()[0] == '$' {
			continue
		}
		if req.command == "" {
			req.command = strings.ToLower(scanner.Text())
		} else {
			req.args = append(req.args, scanner.Text())
		}

		if req.command != "" && len(req.args) == (req.argCount-1) { // don't count command
			break
		}
	}

	return req, nil
}

type Request struct {
	argCount int
	args     []string
	command  string
}
