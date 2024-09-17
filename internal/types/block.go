package types

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"

	"github.com/cbergoon/merkletree"
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
	hashRoot, err := CalculateRootHash(block)
	if err != nil {
		panic(err)
	}
	block.Header.RootHash = hashRoot
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
	if !VerifyRootHash(block) {
		return false
	}

	hash := HashBlockBytes(block)
	pubKey := cryptography.NewPublicKeyFromBytes(block.PublicKey)
	signature := cryptography.NewSignatureFromBytes(block.Signature)

	return signature.Verify(pubKey, hash)
}

func VerifyRootHash(b *genproto.Block) bool {
	merkleRoot, err := CalculateRootHash(b)
	if err != nil {
		return false
	}

	return bytes.Equal(merkleRoot, b.Header.RootHash)

}

// CalculateRootHash updates the root hash of the given block by computing the Merkle tree
// of the block's transactions and setting the root hash in the block's header.
//
// If the Merkle tree is invalid, an error is returned.
func CalculateRootHash(b *genproto.Block) ([]byte, error) {
	list := make([]merkletree.Content, len(b.Transactions))
	for idx, tx := range b.Transactions {
		list[idx] = cryptography.NewTxHash(HashTransactionBytes(tx))
	}

	tree, err := merkletree.NewTree(list)
	if err != nil {
		return nil, err
	}

	return tree.MerkleRoot(), nil
}
