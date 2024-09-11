package random

import (
	"crypto/rand"
	mathrand "math/rand/v2"
	"time"

	"github.com/oleglegun/blockchain-gg/internal/genproto"
)

func RandomHash() []byte {
	hash := make([]byte, 32)
	rand.Read(hash)
	return hash
}

func RandomBlock() *genproto.Block {
	blockHeader := &genproto.BlockHeader{
		Version:   1,
		Height:    int32(mathrand.IntN(1e3)),
		PrevHash:  RandomHash(),
		RootHash:  RandomHash(),
		Timestamp: time.Now().UnixNano(),
	}

	return &genproto.Block{
		Header: blockHeader,
	}

}
