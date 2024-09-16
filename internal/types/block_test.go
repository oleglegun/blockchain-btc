package types

import (
	"testing"

	"github.com/oleglegun/blockchain-btc/internal/cryptography"
	"github.com/oleglegun/blockchain-btc/internal/random"
	"github.com/stretchr/testify/assert"
)

func TestCalculateBlockHash(t *testing.T) {
	block := random.RandomBlock()
	hash := HashBlock(block)
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
	assert.True(t, sig.Verify(pubKey, HashBlock(block)))

	assert.Equal(t, pubKey.Bytes(), block.PublicKey)
	assert.Equal(t, sig.Bytes(), block.Signature)

	assert.True(t, VerifyBlock(block))

	invalidPrivKey := cryptography.GeneratePrivateKey()
	block.PublicKey = invalidPrivKey.Public().Bytes()

	assert.False(t, VerifyBlock(block))

}
