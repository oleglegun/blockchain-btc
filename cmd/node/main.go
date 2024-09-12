package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/oleglegun/blockchain-btc/internal/genproto"
	"github.com/oleglegun/blockchain-btc/internal/node"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const listenAddr = ":3000"

func main() {
	nodeServer := node.NewNode()

	tcpListener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	genproto.RegisterNodeServer(grpcServer, nodeServer)

	go func() {
		for {
			time.Sleep(2 * time.Second)
			makeHandshake(listenAddr)
			makeTransaction(listenAddr)
		}
	}()

	fmt.Println("node running...")
	grpcServer.Serve(tcpListener)
}

/*-----------------------------------------------------------------------------
 *  Temp testing functions
 *----------------------------------------------------------------------------*/

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
		Version: "2.0",
		Height:  11,
	}

	_, err = genproto.NewNodeClient(clientConn).Handshake(context.Background(), nodeInfo)
	if err != nil {
		log.Fatal(err)
	}
}
