package types

import (
	"testing"

	"github.com/oleglegun/blockchain-btc/internal/cryptography"
	"github.com/oleglegun/blockchain-btc/internal/random"
	"github.com/stretchr/testify/assert"
)

func TestHashBlock(t *testing.T) {
	block := random.RandomBlock()
	hash := HashBlock(block)
	assert.Equal(t, 32, len(hash))
}

func TestSignBlock(t *testing.T) {
	var (
		block   = random.RandomBlock()
		privKey = cryptography.NewPrivateKey()
		pubKey  = privKey.Public()
	)

	sig := SignBlock(privKey, block)

	assert.Equal(t, 64, len(sig.Bytes()))
	assert.True(t, sig.Verify(pubKey, HashBlock(block)))

	assert.Equal(t, pubKey.Bytes(), block.PublicKey)
	assert.Equal(t, sig.Bytes(), block.Signature)

	assert.True(t, VerifyBlock(block))

	invalidPrivKey := cryptography.NewPrivateKey()
	block.PublicKey = invalidPrivKey.Public().Bytes()

	assert.False(t, VerifyBlock(block))

}

func TestVerifyBlock(t *testing.T) {
	var (
		block   = random.RandomBlock()
		privKey = cryptography.NewPrivateKey()
		pubKey  = privKey.Public()
	)

	sig := SignBlock(privKey, block)
	block.Signature = sig.Bytes()
	block.PublicKey = pubKey.Bytes()

	assert.True(t, VerifyBlock(block))

	block.PublicKey = cryptography.NewPrivateKey().Public().Bytes()
	assert.False(t, VerifyBlock(block))
}
