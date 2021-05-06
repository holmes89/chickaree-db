package core

import (
	"bytes"
	"io"
	"strconv"
)

type Request struct {
	msgCount int
	args     []Arg
	command  string
}

type Response struct {
	rtype   RESPType
	length  int
	content []byte
}

func (res Response) SetContent(c []byte) {
	res.content = c
	res.length = len(c)
}
func (res Response) Encode() []byte {

	buf := new(bytes.Buffer)
	buf.WriteRune(rune(res.rtype))

	length := strconv.Itoa(res.length)

	switch res.rtype {
	case SimpleString:
		buf.Write(res.content)
	case BulkStrings:
		buf.WriteString(length)
		buf.Write(TerminationSeq)
		buf.Write(res.content)
	case Integers:
		buf.Write(res.content)
	case Errors:
		buf.Write(res.content)
	}
	buf.Write(TerminationSeq)
	return buf.Bytes()
}

type RESPType rune

var (
	SimpleString RESPType = '+'
	Errors       RESPType = '-'
	Integers     RESPType = ':'
	BulkStrings  RESPType = '$'
	Arrays       RESPType = '*'
)

var TerminationSeq = []byte{'\r', '\n'}

var EmptyStringResp = Response{
	rtype:   BulkStrings,
	content: TerminationSeq,
	length:  0,
}
var NilStringResp = Response{
	rtype:  BulkStrings,
	length: -1,
}

var OkResp = Response{
	rtype:   SimpleString,
	content: []byte("OK"),
}

func ErrResponse(err error) Response {
	return Response{
		rtype:   Errors,
		content: []byte(err.Error()),
	}
}

func NewRequest(r io.Reader) (req Request, err error) {

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
		req.args = append(req.args, b[:len(b)-2])
	}
	req.command = string(req.args[0])
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
