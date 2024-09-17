package types

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/oleglegun/blockchain-btc/internal/cryptography"
	"github.com/oleglegun/blockchain-btc/internal/genproto"
	"google.golang.org/protobuf/proto"
)

// CalculateBlockHash computes the hash of the given block's header.
func HashBlockHeader(header *genproto.BlockHeader) []byte {
	b, err := proto.Marshal(header)
	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256(b)

	return hash[:]
}

// HashBlockBytes computes the hash of the given block's header.
// Block hash == its header hash
func HashBlockBytes(block *genproto.Block) []byte {
	return HashBlockHeader(block.Header)
}

func HashBlockString(block *genproto.Block) string {
	return hex.EncodeToString(HashBlockBytes(block))
}

// SignBlock signs the given block using the provided private key.
func SignBlock(privKey cryptography.PrivateKey, block *genproto.Block) cryptography.Signature {
	hash := HashBlockBytes(block)
	signature := privKey.Sign(hash)
	block.PublicKey = privKey.Public().Bytes()
	block.Signature = signature.Bytes()
	return privKey.Sign(HashBlockBytes(block))
}

func VerifyBlock(block *genproto.Block) bool {
	if len(block.PublicKey) != cryptography.PubKeyLen {
		return false
	}
	if len(block.Signature) != cryptography.SigLen {
		return false
	}

	hash := HashBlockBytes(block)
	pubKey := cryptography.NewPublicKeyFromBytes(block.PublicKey)
	signature := cryptography.NewSignatureFromBytes(block.Signature)

	return signature.Verify(pubKey, hash)
}
