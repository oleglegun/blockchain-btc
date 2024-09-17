package cryptography

import (
	"bytes"

	"github.com/cbergoon/merkletree"
)

type TxHash struct {
	hash []byte
}

func NewTxHash(hash []byte) *TxHash {
	return &TxHash{hash: hash}
}

func (h TxHash) CalculateHash() ([]byte, error) {
	return h.hash, nil
}

func (h TxHash) Equals(other merkletree.Content) (bool, error) {
	equals := bytes.Equal(h.hash, other.(TxHash).hash)
	return equals, nil
}
