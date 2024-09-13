package node

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/oleglegun/blockchain-btc/internal/genproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	nodeVersion = "1.0"
)

type Node struct {
	genproto.UnimplementedNodeServer

	version    string
	listenAddr string
	log        *log.Logger

	mu    sync.RWMutex
	peers map[genproto.NodeClient]*genproto.NodeInfo
}

func NewNode(listenAddr string) *Node {
	return &Node{
		peers:      make(map[genproto.NodeClient]*genproto.NodeInfo),
		version:    nodeVersion,
		listenAddr: listenAddr,
		log:        log.New(os.Stderr, fmt.Sprintf("[ NODE %s ] ", listenAddr), log.LstdFlags|log.Lmsgprefix),
	}
}

func (n *Node) Start() error {
	grpcServer := grpc.NewServer()
	tpcListener, err := net.Listen("tcp", n.listenAddr)
	if err != nil {
		return err
	}

	genproto.RegisterNodeServer(grpcServer, n)

	n.log.Println("running...")
	return grpcServer.Serve(tpcListener)
}

// bootstrapNetwork connects the current node to the specified list of listen socket addresses (ip:port).
// It creates a new client connection to each peer node and performs a handshake to exchange node information.
func (n *Node) BootstrapNetwork(listenSocketAddrs []string) error {
	for _, listenSocketAddr := range listenSocketAddrs {
		peerClient, err := newNodeClient(listenSocketAddr)
		if err != nil {
			return err
		}

		peerNodeInfo, err := peerClient.Handshake(context.TODO(), n.getNodeInfo())
		if err != nil {
			n.log.Println("handshake error:", err)
			continue
		}

		n.addPeer(peerClient, peerNodeInfo)
	}

	return nil
}

//-----------------------------------------------------------------------------
//  GRPC Service methods
//-----------------------------------------------------------------------------

// Handshake is called when a new peer node connects to the current node.
// It exchanges node information with the peer node and adds the peer to the list of connected peers.
//
// If there is an error creating the client connection to the peer node, the function will return an error.
func (n *Node) Handshake(ctx context.Context, peerNodeInfo *genproto.NodeInfo) (*genproto.NodeInfo, error) {
	thisNodeInfo := n.getNodeInfo()

	if thisNodeInfo.Version != peerNodeInfo.Version {
		n.log.Println("incompatible node versions")
		return nil, fmt.Errorf("incompatible node versions")
	}

	peerClient, err := newNodeClient(peerNodeInfo.ListenAddr)
	if err != nil {
		return nil, err
	}

	n.addPeer(peerClient, peerNodeInfo)

	return thisNodeInfo, nil
}

func (n *Node) HandleTransaction(ctx context.Context, tx *genproto.Transaction) (*emptypb.Empty, error) {
	peer, ok := peer.FromContext(ctx)
	if !ok {
		n.log.Fatal("cannot get peer from the context")
	}
	n.log.Println("received tx from:", peer.Addr)
	_ = peer

	return &emptypb.Empty{}, nil
}

//-----------------------------------------------------------------------------
//  Other methods
//-----------------------------------------------------------------------------

func (n *Node) addPeer(peerClient genproto.NodeClient, peerNodeInfo *genproto.NodeInfo) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.log.Printf("new peer connected: %s, height: %d", peerNodeInfo.ListenAddr, peerNodeInfo.Height)
	n.peers[peerClient] = peerNodeInfo
}

func (n *Node) removePeer(c genproto.NodeClient) {
	n.mu.Lock()
	defer n.mu.Unlock()

	delete(n.peers, c)
}

func (n *Node) getNodeInfo() *genproto.NodeInfo {
	return &genproto.NodeInfo{
		Version:    n.version,
		Height:     10,
		ListenAddr: n.listenAddr,
	}
}

func newNodeClient(listenSocketAddr string) (genproto.NodeClient, error) {
	clientConn, err := grpc.NewClient("dns:///"+listenSocketAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return genproto.NewNodeClient(clientConn), nil

}
