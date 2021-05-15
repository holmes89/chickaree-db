package redis

import (
	"fmt"
	"strings"

	"github.com/holmes89/chickaree-db/pkg/core"
)

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
	entry := core.Entry{
		Type:  "primative",
		Key:   args[0],
		Value: args[1],
	}
	err := c.repo.Set(entry)
	if err != nil {
		return ErrResponse(err)
	}
	return OkResp
}

func (c *Client) get(args []Arg) Response {

	r, err := c.repo.Get(args[0])

	if err == core.ErrNotFound {
		return NilStringResp
	}

	if err != nil {
		return ErrResponse(err)
	}
	res := Response{
		rtype:   BulkStrings,
		content: r.Value,
		length:  len(r.Value),
	}
	return res
}

func (c *Client) hSet(args []Arg) Response {
	// var count int
	// if (len(args) < 3) || (len(args[1:])%2 != 0) {
	// 	return ErrResponse(fmt.Errorf("invalid arg count %d", len(args)))
	// }
	// err := r.db.Update(func(tx *bolt.Tx) error {
	// 	b, err := tx.CreateBucketIfNotExists(args[0])
	// 	if err != nil {
	// 		return err
	// 	}
	// 	for i := 1; i < len(args); i += 2 {
	// 		if err := b.Put(args[i], args[i+1]); err != nil {
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
