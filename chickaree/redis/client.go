package redis

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/holmes89/chickaree-db/chickaree"
)

type Client struct {
	// incoming chan string
	outgoing chan []byte
	reader   *bufio.Reader
	writer   *bufio.Writer
	conn     net.Conn
	client   chickaree.ChickareeDBClient
}

func (c *Client) Read() {
	for {
		req, err := NewRequest(c.reader)
		if err == nil {
			c.outgoing <- c.Handle(req)
		} else {
			break
		}

	}

	c.conn.Close()
	fmt.Println("client disconnected")
}

func (c *Client) Write() {
	for data := range c.outgoing {
		c.writer.Write(data)
		c.writer.Flush()
	}
}

func (client *Client) Listen() {
	go client.Read()
	go client.Write()
}

func NewClient(connection net.Conn, cl chickaree.ChickareeDBClient) *Client {
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
		client:   cl,
	}
	client.Listen()

	return client
}

func (c *Client) Handle(req Request) []byte {
	switch strings.ToLower(req.Command) {
	case "command":
		return OkResp.Encode()
	case "ping":
		return Response{
			rtype:   SimpleString,
			content: []byte("PONG"),
		}.Encode()
	case "set":
		return c.set(req.Args).Encode()
	case "hset":
		return c.hSet(req.Args).Encode()
	case "hget":
		return c.hGet(req.Args).Encode()
	case "hexists":
		return c.hExists(req.Args).Encode()
	case "hgetall":
		return c.hGetAll(req.Args).Encode()
	case "get":
		return c.get(req.Args).Encode()
	case "del":
		return c.del(req.Args).Encode()
	default:
		err := fmt.Errorf("unknown command '%s'", req.Command)
		return ErrResponse(err).Encode()
	}
}

func (c *Client) set(args []Arg) Response {
	ctx := context.TODO()
	req := &chickaree.SetRequest{
		Key:   string(args[0]),
		Value: args[1],
	}
	_, err := c.client.Set(ctx, req)
	if err != nil {
		ErrResponse(err)
	}
	return OkResp
}

func (c *Client) get(args []Arg) Response {
	ctx := context.TODO()
	req := &chickaree.GetRequest{
		Key: string(args[0]),
	}
	resp, err := c.client.Get(ctx, req)
	if err != nil {
		ErrResponse(err)
	}

	if len(resp.Data) == 0 {
		return NilStringResp
	}

	return Response{
		rtype:   BulkStrings,
		length:  len(resp.Data),
		content: resp.Data,
	}
}

func (c *Client) hSet(args []Arg) Response {

	return OkResp
}

func (c *Client) hGet(args []Arg) Response {
	// var b []byte
	// r.db.View(func(tx *bolt.Tx) error {
	// 	b = tx.Bucket(args[0]).Get(args[1])
	// 	return nil
	// })
	// if len(b) == 0 {
	// 	return NilStringResp
	// }
	// res := Response{
	// 	rtype:   BulkStrings,
	// 	content: b,
	// 	length:  len(b),
	// }
	// return res
	return OkResp
}

// TODO maybe need to store types along with value
func (c *Client) hGetAll(args []Arg) ResponseArray {
	// var res ResponseArray
	// r.db.View(func(tx *bolt.Tx) error {
	// 	return tx.Bucket(args[0]).ForEach(func(k, v []byte) error {
	// 		res = append(res, Response{
	// 			rtype:   BulkStrings,
	// 			content: k,
	// 			length:  len(k),
	// 		})
	// 		res = append(res, Response{
	// 			rtype:   BulkStrings,
	// 			content: v,
	// 			length:  len(v),
	// 		})
	// 		return nil
	// 	})
	// })

	return nil

}

func (c *Client) hExists(args []Arg) Response {

	// var b []byte
	// r.db.View(func(tx *bolt.Tx) error {
	// 	b = tx.Bucket(args[0]).Get(args[1])
	// 	return nil
	// })
	// count := "1"
	// if len(b) == 0 {
	// 	count = "0"
	// }
	// res := Response{
	// 	rtype:   Integers,
	// 	content: []byte(count),
	// }
	// return res
	return OkResp
}

func (c *Client) del(args []Arg) Response {
	// var count int
	// err := r.db.Update(func(tx *bolt.Tx) error {
	// 	b := tx.Bucket(defaultBucket)

	// 	for _, a := range args {
	// 		if b.Get(a) == nil {
	// 			continue
	// 		}
	// 		if err := b.Delete(a); err != nil {
	// 			return err
	// 		}
	// 		count++
	// 	}
	// 	return nil
	// })
	// if err != nil {
	// 	ErrResponse(err)
	// }
	// c := fmt.Sprintf("%d", count)
	// return Response{
	// 	rtype:   Integers,
	// 	content: []byte(c),
	// }
	return OkResp
}
