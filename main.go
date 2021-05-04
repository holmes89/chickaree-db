package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"strconv"
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
		fmt.Printf("%+v", req)
		// command := strings.TrimSpace(string(line))
		// switch command {
		// case "*1":
		// 	c.Write([]byte("+PONG\r\n"))
		// default:
		// 	fmt.Printf("unknown: %v\n", command)
		// 	c.Write([]byte(fmt.Sprintf("-ERR unknown command '%s'\r\n", command)))
		// }
	}
}

type lexer struct {
	*scanner.Scanner
	buf bytes.Buffer
}

// New returns new lexer
func Read(r io.Reader) (Request, error) {
	var s scanner.Scanner
	s.Init(r)
	s.Mode &^= scanner.ScanChars | scanner.ScanRawStrings
	l := &lexer{
		Scanner: &s,
	}
	return l.readForm()
}

func (l *lexer) readForm() (Request, error) {
	req := Request{}
	test, err := l.readUntilNextLine()
	if err != nil {
		return req, err
	}
	fmt.Println(test)
	test, err = l.readUntilNextLine()
	if err != nil {
		return req, err
	}
	fmt.Println(test)
	cs, err := l.readUntilNextLine()
	if err != nil {
		return req, err
	}
	i, err := strconv.Atoi(cs[1:])
	if err != nil {
		return req, fmt.Errorf("invalid token %s", cs)
	}
	req.argCount = i

	_, err = l.readUntilNextLine() // next line is info telling us the next line size
	if err != nil {
		return req, err
	}

	req.command, err = l.readUntilNextLine()
	if err != nil {
		return req, err
	}

	for j := i - 1; j > 0; j-- {
		_, err = l.readUntilNextLine() // next line is info telling us the next line size
		if err != nil {
			return req, err
		}

		arg, err := l.readUntilNextLine()
		if err != nil {
			return req, err
		}
		req.args = append(req.args, arg)
	}
	return req, nil
}

func (l *lexer) readUntilNextLine() (string, error) {
	for {
		r := l.Next()
		if r == scanner.EOF {
			return "", errors.New("invalid command")
		}
		if r == '\r' && l.Peek() == '\n' {
			l.Next()
			break
		}
		l.buf.WriteRune(r)
	}
	return l.String(), nil
}

type Request struct {
	argCount int
	args     []string
	command  string
}
