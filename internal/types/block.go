package types

import (
	"crypto/sha256"

	"github.com/oleglegun/blockchain-gg/internal/cryptography"
	"github.com/oleglegun/blockchain-gg/internal/genproto"
	"google.golang.org/protobuf/proto"
)

// CalculateBlockHash computes the hash of the given block's header.
func CalculateBlockHash(block *genproto.Block) []byte {
	b, err := proto.Marshal(block.Header)
	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256(b)

	return hash[:]
}

// SignBlock signs the given block using the provided private key.
func SignBlock(privKey cryptography.PrivateKey, block *genproto.Block) cryptography.Signature {
	return privKey.Sign(CalculateBlockHash(block))
}
