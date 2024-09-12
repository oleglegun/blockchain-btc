package types

import (
	"testing"

	"github.com/oleglegun/blockchain-btc/internal/cryptography"
	"github.com/oleglegun/blockchain-btc/internal/random"
	"github.com/stretchr/testify/assert"
)

func TestCalculateBlockHash(t *testing.T) {
	block := random.RandomBlock()
	hash := CalculateBlockHash(block)
	assert.Equal(t, 32, len(hash))
}

func TestSignBlock(t *testing.T) {
	var (
		block   = random.RandomBlock()
		privKey = cryptography.GeneratePrivateKey()
		pubKey  = privKey.Public()
	)

	sig := CalculateBlockSignature(privKey, block)

	assert.Equal(t, 64, len(sig.Bytes()))
	assert.True(t, sig.Verify(pubKey, CalculateBlockHash(block)))
}
