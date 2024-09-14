package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/oleglegun/blockchain-btc/internal/genproto"
	"github.com/oleglegun/blockchain-btc/internal/node"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const nodeCount = 4

func main() {
	for i := 0; i < nodeCount; i++ {
		port := 3000 + i
		listenAddr := fmt.Sprintf(":%d", port)
		bootstrapNodes := make([]string, 0, nodeCount)
		if i > 0 {
			bootstrapNodes = append(bootstrapNodes, fmt.Sprintf("localhost:%d", port-1))
		}

		makeNode(listenAddr, bootstrapNodes)
		time.Sleep(time.Second)
	}

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
