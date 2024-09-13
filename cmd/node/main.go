package main

import (
	"context"
	"log"

	"github.com/oleglegun/blockchain-btc/internal/genproto"
	"github.com/oleglegun/blockchain-btc/internal/node"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	makeNode(":3000", []string{})
	makeNode(":3001", []string{"localhost:3000"})
	makeNode(":3002", []string{"localhost:3000", "localhost:3001"})
}

/*-----------------------------------------------------------------------------
 *  Temp testing functions
 *----------------------------------------------------------------------------*/

func makeNode(listenAddr string, bootstrapNodes []string) error {
	nodeServer := node.NewNode(listenAddr)
	go func() {
		log.Fatal(nodeServer.Start())
	}()
	return nodeServer.BootstrapNetwork(bootstrapNodes)
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
