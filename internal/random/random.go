package random

import (
	"crypto/rand"
	mathrand "math/rand/v2"
	"time"

	"github.com/oleglegun/blockchain-btc/internal/genproto"
)

func Random32ByteHash() []byte {
	hash := make([]byte, 32)
	rand.Read(hash)
	return hash
}

func Random64ByteHash() []byte {
	hash := make([]byte, 64)
	rand.Read(hash)
	return hash
}

func RandomBlock() *genproto.Block {
	blockHeader := &genproto.BlockHeader{
		Version:   1,
		Height:    int32(mathrand.IntN(1e3)),
		PrevHash:  Random32ByteHash(),
		RootHash:  Random32ByteHash(),
		Timestamp: time.Now().UnixNano(),
	}

	return &genproto.Block{
		Header: blockHeader,
		PublicKey: Random32ByteHash(),
		Signature: Random64ByteHash(),
	}

}
