package node

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"sync"
	"time"

	"github.com/oleglegun/blockchain-btc/internal/cryptography"
	"github.com/oleglegun/blockchain-btc/internal/genproto"
	"github.com/oleglegun/blockchain-btc/internal/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	nodeVersion = "1.0"
	blockTime   = time.Second * 5
)

type NodeConfig struct {
	Version    string
	ListenAddr string
	PrivateKey *cryptography.PrivateKey
}

type Node struct {
	genproto.UnimplementedNodeServer

	NodeConfig
	log *slog.Logger

	peersLock sync.RWMutex
	peers     map[string]ConnectedPeer
	mempool   *Mempool
}

type ConnectedPeer struct {
	peerClient genproto.NodeClient
	nodeInfo   *genproto.NodeInfo
}

func NewNode(config NodeConfig) *Node {
	logHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: false,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	})

	return &Node{
		NodeConfig: config,
		log:        slog.New(logHandler).With("node", config.ListenAddr),
		peers:      make(map[string]ConnectedPeer),
		mempool:    NewMempool(),
	}
}

// Start runs the node server, listening on the configured listen address and
// registering the node gRPC service. It bootstraps the node by connecting to
// the specified list of bootstrap nodes.
func (n *Node) Start(bootstrapNodes []string) error {
	grpcServer := grpc.NewServer()
	tpcListener, err := net.Listen("tcp", n.ListenAddr)
	if err != nil {
		return err
	}

	genproto.RegisterNodeServer(grpcServer, n)

	n.log.Debug("running...")

	if len(bootstrapNodes) > 0 {
		n.log.Debug("need connect to", "peers", bootstrapNodes)
		go n.bootstrapNetwork(bootstrapNodes)
	}

	if n.PrivateKey != nil {
		go n.runValidatorLoop()
	}

	return grpcServer.Serve(tpcListener)
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
		n.log.Debug("incompatible node versions")
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
		n.log.Error("cannot get peer from the context")

	}

	if n.mempool.Add(tx) {
		txHash := types.HashTransactionString(tx)
		n.log.Debug("received tx", "from", peer.Addr, "tx", txHash)

		go func() {
			if err := n.broadcast(tx); err != nil {
				n.log.Error("failed to broadcast transaction", "error", err)
			}
		}()
	}

	return &emptypb.Empty{}, nil
}

//-----------------------------------------------------------------------------
//  Other methods
//-----------------------------------------------------------------------------

// bootstrapNetwork connects the current node to the specified list of listen socket addresses (ip:port).
// It creates a new client connection to each peer node and performs a handshake to exchange node information.
func (n *Node) bootstrapNetwork(listenSocketAddrs []string) error {
	for _, listenSocketAddr := range listenSocketAddrs {
		peerClient, peerNodeInfo, err := n.dialPeerNode(listenSocketAddr)
		if err != nil {
			n.log.Error("cannot dial to peer node", "error", err)
			continue
		}

		n.addPeer(peerClient, peerNodeInfo)
	}

	return nil
}

func (n *Node) addPeer(peerClient genproto.NodeClient, peerNodeInfo *genproto.NodeInfo) {
	n.peersLock.Lock()
	n.peers[peerNodeInfo.ListenAddr] = ConnectedPeer{
		peerClient: peerClient,
		nodeInfo:   peerNodeInfo,
	}
	n.log.Debug("connected nodes", "count", len(n.peers))

	n.peersLock.Unlock()

	n.log.Debug("new peer connected", "peer", peerNodeInfo.ListenAddr)

	absentPeerList := n.getAbsentPeerList(peerNodeInfo.PeerList)

	if len(absentPeerList) > 0 {
		n.log.Debug("need connect to", "peers", absentPeerList)
		go n.bootstrapNetwork(absentPeerList)
	}
}

func (n *Node) removePeer(peerListenAddr string) {
	n.peersLock.Lock()
	defer n.peersLock.Unlock()

	// TODO: close the peerClient
	delete(n.peers, peerListenAddr)
}

func (n *Node) runValidatorLoop() {
	ticker := time.NewTicker(blockTime)
	n.log.Debug("running validation loop", "pubKey", n.PrivateKey.Public())

	for {
		<-ticker.C
		txList := n.mempool.Clear()
		n.mempool.ClearProcessed(time.Minute)
		n.log.Info("new block imminent", "txs", len(txList))
	}
}

// broadcast broadcasts message to all known peers
func (n *Node) broadcast(msg any) error {
	n.peersLock.RLock()
	defer n.peersLock.RUnlock()

	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch v := msg.(type) {
	case *genproto.Transaction:
		for _, peer := range n.peers {
			wg.Add(1)
			go func(peer ConnectedPeer) {
				defer wg.Done()
				_, err := peer.peerClient.HandleTransaction(ctx, v)
				if err != nil {
					n.log.Error("failed to broadcast transaction to peer", "peer", peer.nodeInfo.ListenAddr, "error", err)
				}
			}(peer)
		}
	default:
		n.log.Error("unsupported message type for broadcast", "type", fmt.Sprintf("%T", msg))
		return fmt.Errorf("unsupported message type: %T", msg)
	}

	wg.Wait()
	return nil
}

func (n *Node) getAbsentPeerList(peerAddrList []string) []string {
	absentPeerAddresses := []string{}

	connectedPeers := n.getPeerList()

	for _, peerAddr := range peerAddrList {
		if peerAddr == n.ListenAddr {
			// Skip own address
			continue
		}

		isConnected := false
		for _, connectedPeer := range connectedPeers {
			if peerAddr == connectedPeer {
				isConnected = true
				break
			}
		}

		if !isConnected {
			absentPeerAddresses = append(absentPeerAddresses, peerAddr)
		}
	}

	return absentPeerAddresses
}

func (n *Node) getNodeInfo() *genproto.NodeInfo {
	return &genproto.NodeInfo{
		Version:    n.Version,
		Height:     0,
		ListenAddr: n.ListenAddr,
		PeerList:   n.getPeerList(),
	}
}

func (n *Node) dialPeerNode(peerListenAddr string) (genproto.NodeClient, *genproto.NodeInfo, error) {
	peerClient, err := newNodeClient(peerListenAddr)
	if err != nil {
		return nil, nil, err
	}

	peerNodeInfo, err := peerClient.Handshake(context.TODO(), n.getNodeInfo())
	if err != nil {
		n.log.Debug("handshake error", "error", err)
		return nil, nil, err
	}

	return peerClient, peerNodeInfo, nil
}

func (n *Node) getPeerList() []string {
	n.peersLock.RLock()
	defer n.peersLock.RUnlock()

	peerList := make([]string, 0, len(n.peers))
	for k := range n.peers {
		peerList = append(peerList, k)
	}

	return peerList
}

func newNodeClient(listenSocketAddr string) (genproto.NodeClient, error) {
	clientConn, err := grpc.NewClient("dns:///"+listenSocketAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return genproto.NewNodeClient(clientConn), nil

}
