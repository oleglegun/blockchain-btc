package types

import (
	"testing"

	"github.com/oleglegun/blockchain-gg/internal/cryptography"
	"github.com/oleglegun/blockchain-gg/internal/random"
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

	sig := SignBlock(privKey, block)

	assert.Equal(t, 64, len(sig.Bytes()))
	assert.True(t, sig.Verify(pubKey, CalculateBlockHash(block)))
}
