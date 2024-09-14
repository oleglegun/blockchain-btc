package types

import (
	"crypto/sha256"

	"github.com/oleglegun/blockchain-btc/internal/cryptography"
	"github.com/oleglegun/blockchain-btc/internal/genproto"
	"google.golang.org/protobuf/proto"
)

// CalculateBlockHash computes the hash of the given block's header.
func CalcBlockHeaderHash(header *genproto.BlockHeader) []byte {
	b, err := proto.Marshal(header)
	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256(b)

	return hash[:]
}

// CalcBlockHash computes the hash of the given block's header.
// Block hash == its header hash
func CalcBlockHash(block *genproto.Block) []byte {
	return CalcBlockHeaderHash(block.Header)
}

// CalcBlockSignature signs the given block using the provided private key.
func CalcBlockSignature(privKey cryptography.PrivateKey, block *genproto.Block) cryptography.Signature {
	return privKey.Sign(CalcBlockHash(block))
}
