package main

import (
	"context"
	"flag"
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

func main() {
	nodeCount := flag.Int("nodeCount", 3, "Number of nodes in the network")
	flag.Parse()

	log.Printf("Running blockchain with %d nodes", *nodeCount)

	// Create and start the specified number of nodes
	for i := 1; i <= *nodeCount; i++ {
		port := 3000 + i
		listenAddr := fmt.Sprintf(":%d", port)
		bootstrapNodes := make([]string, 0, *nodeCount)

		if i == 1 {
			// The first node is a validator and does not have any bootstrap nodes
			makeNode(listenAddr, true, bootstrapNodes)
		} else {
			// Subsequent nodes are not validators and bootstrap from the previous node
			// Nodes will discover each other through the nodes gossip protocol
			makeNode(listenAddr, false, []string{fmt.Sprintf("localhost:%d", port-1)})
		}

		// Sleep for a second to allow the node to start
		time.Sleep(time.Second)
	}

	// Continuously make transactions to the node running on port 3002
	// that will be shared with the network
	for {
		// In this demo transactions are not verified and are always accepted by the validator loop
		// But the system itself is capable of verifying transactions
		go makeTransaction(":3002")
		time.Sleep(time.Millisecond * 1000)
	}
}

/*-----------------------------------------------------------------------------
 *  Temp testing functions
 *----------------------------------------------------------------------------*/

func makeNode(listenAddr string, isValidator bool, bootstrapNodes []string) error {
	nodeConfig := node.NodeConfig{
		Version:    "1",
		ListenAddr: listenAddr,
	}

	if isValidator {
		privKey := cryptography.NewPrivateKey()
		nodeConfig.PrivateKey = &privKey
	}

	chain := node.NewChain(node.NewMemoryBlockStore(), node.NewMemoryTxStore(), node.NewMemoryUTXOStore())

	nodeServer := node.NewNode(nodeConfig, chain)

	go func() {
		log.Fatal(nodeServer.Start(bootstrapNodes))
	}()

	return nil
}

var clientConnCache = make(map[string]*grpc.ClientConn)

func makeTransaction(addr string) {
	clientConn, err := getClientConn(addr)
	if err != nil {
		log.Fatal(err)
	}

	sender1PrivKey := cryptography.NewPrivateKey()
	sender2PrivKey := cryptography.NewPrivateKey()
	receiverPrivKey := cryptography.NewPrivateKey()

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

func getClientConn(addr string) (*grpc.ClientConn, error) {
	if conn, exists := clientConnCache[addr]; exists {
		return conn, nil
	}

	conn, err := grpc.NewClient("dns:///localhost"+addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	clientConnCache[addr] = conn
	return conn, nil
}
