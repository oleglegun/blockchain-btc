package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/oleglegun/blockchain-btc/internal/cryptography"
	"github.com/oleglegun/blockchain-btc/internal/genproto"
	"github.com/oleglegun/blockchain-btc/internal/node"
	"github.com/oleglegun/blockchain-btc/internal/random"
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

	makeTransaction(":3000")
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

	sender1PrivKey := cryptography.GeneratePrivateKey()
	sender2PrivKey := cryptography.GeneratePrivateKey()
	receiverPrivKey := cryptography.GeneratePrivateKey()

	txIn1 := &genproto.TxInput{
		PrevTxHash:     random.Random32ByteHash(),
		PrevTxOutIndex: 0,
		PublicKey:      sender1PrivKey.Public().Bytes(),
		// Signature will be set after constructing transaction
	}

	txIn2 := &genproto.TxInput{
		PrevTxHash:     random.Random32ByteHash(),
		PrevTxOutIndex: 1,
		PublicKey:      sender2PrivKey.Public().Bytes(),
		// Signature will be set after constructing transaction
	}

	txOut1 := &genproto.TxOutput{
		Amount:  9,
		Address: receiverPrivKey.Public().Address().Bytes(),
	}
	txOut2 := &genproto.TxOutput{
		Amount:  1,
		Address: sender1PrivKey.Public().Address().Bytes(),
	}

	tx := &genproto.Transaction{
		Version: 1,
		Inputs:  []*genproto.TxInput{txIn1, txIn2},
		Outputs: []*genproto.TxOutput{txOut1, txOut2},
	}

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
