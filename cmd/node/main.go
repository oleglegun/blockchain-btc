package main

import (
	"context"
	"log"
	"time"

	"github.com/oleglegun/blockchain-btc/internal/genproto"
	"github.com/oleglegun/blockchain-btc/internal/node"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	makeNode(":3001", []string{})
	time.Sleep(time.Second)
	makeNode(":3002", []string{})
	time.Sleep(time.Second)
	makeNode(":3003", []string{"localhost:3001"})
	time.Sleep(time.Second)
	makeNode(":3004", []string{"localhost:3001", "localhost:3002"})

	time.Sleep(10 * time.Second)
}

/*-----------------------------------------------------------------------------
 *  Temp testing functions
 *----------------------------------------------------------------------------*/

func makeNode(listenAddr string, bootstrapNodes []string) error {
	nodeServer := node.NewNode(listenAddr)
	go func() {
		log.Fatal(nodeServer.Start(bootstrapNodes))
	}()

	return nil
}

func makeTransaction(addr string) {
	clientConn, err := grpc.NewClient("dns:///localhost"+addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	defer clientConn.Close()

	tx := &genproto.Transaction{}

	_, err = genproto.NewNodeClient(clientConn).HandleTransaction(context.Background(), tx)
	if err != nil {
		log.Fatal(err)
	}
}

func makeHandshake(addr string) {
	clientConn, err := grpc.NewClient("dns:///localhost"+addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	defer clientConn.Close()

	nodeInfo := &genproto.NodeInfo{
		Version:    "2.0",
		Height:     11,
		ListenAddr: "localhost:3000",
	}

	_, err = genproto.NewNodeClient(clientConn).Handshake(context.Background(), nodeInfo)
	if err != nil {
		log.Fatal(err)
	}
}
