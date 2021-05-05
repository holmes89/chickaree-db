package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/holmes89/chickaree-db/pkg/core"
)

func main() {

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

	repo := core.NewRepo("chickaree.db")
	defer repo.Close()

	fmt.Printf("listening on port %s\n", PORT)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err.Error())
		}
		_ = core.NewClient(conn, repo)
		fmt.Println("client connected")
	}
}
