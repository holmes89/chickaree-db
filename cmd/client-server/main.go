package main

import (
	"flag"
	"log"

	"github.com/holmes89/chickaree-db/chickaree"
	"github.com/holmes89/chickaree-db/chickaree/redis"
	"google.golang.org/grpc"
)

func main() {

	var port string
	flag.StringVar(&port, "port", "6379", "port to listen on")
	flag.Parse()

	opts := []grpc.DialOption{grpc.WithInsecure()}
	conn, err := grpc.Dial(":8080", opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := chickaree.NewChickareeDBClient(conn)

	tcpServer := redis.NewTCPServer(port, client)
	defer tcpServer.Close()

	log.Println(<-tcpServer.Run())

}
