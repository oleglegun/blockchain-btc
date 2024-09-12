package node

import (
	"context"
	"log"

	"github.com/oleglegun/blockchain-btc/internal/genproto"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/emptypb"
)

const nodeVersion = "1.0"

type Node struct {
	version string
	genproto.UnimplementedNodeServer
}

func NewNode() *Node {
	return &Node{
		version: nodeVersion,
	}
}

func (n *Node) HandleTransaction(ctx context.Context, tx *genproto.Transaction) (*emptypb.Empty, error) {
	peer, ok := peer.FromContext(ctx)
	if !ok {
		log.Fatal("cannot get peer from the context")
	}
	log.Println("received tx from:", peer.Addr)
	_ = peer

	return &emptypb.Empty{}, nil
}

func (n *Node) Handshake(ctx context.Context, peerNodeInfo *genproto.NodeInfo) (*genproto.NodeInfo, error) {
	thisNodeInfo := &genproto.NodeInfo{
		Version: n.version,
		Height:  10,
	}

	p, ok := peer.FromContext(ctx)
	if !ok {
		log.Fatal("cannot get peer from the context")
	}

	log.Printf("received node info from peer %s: %+v\n", p.Addr, peerNodeInfo)
	return thisNodeInfo, nil
}
