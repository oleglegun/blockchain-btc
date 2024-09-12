package types

import (
	"crypto/sha256"

	"github.com/oleglegun/blockchain-btc/internal/cryptography"
	"github.com/oleglegun/blockchain-btc/internal/genproto"
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

// CalculateBlockSignature signs the given block using the provided private key.
func CalculateBlockSignature(privKey cryptography.PrivateKey, block *genproto.Block) cryptography.Signature {
	return privKey.Sign(CalculateBlockHash(block))
}
