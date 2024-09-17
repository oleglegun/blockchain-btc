package types

import (
	"testing"

	"github.com/oleglegun/blockchain-btc/internal/cryptography"
	"github.com/oleglegun/blockchain-btc/internal/genproto"
	"github.com/oleglegun/blockchain-btc/internal/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashBlock(t *testing.T) {
	block := random.RandomBlock()
	hash := HashBlockBytes(block)
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
	assert.True(t, sig.Verify(pubKey, HashBlockBytes(block)))

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

func TestCalculateRootHash(t *testing.T) {
	block := random.RandomBlock()
	tx := &genproto.Transaction{
		Version: 1,
	}

	block.Transactions = append(block.Transactions, tx)

	initialRootHash := block.Header.RootHash

	rootHash, err := CalculateRootHash(block)
	assert.Nil(t, err)

	require.NotEqual(t, initialRootHash, rootHash)

}

func TestVerifyRootHash(t *testing.T) {
	privKey := cryptography.NewPrivateKey()

	block := random.RandomBlock()
	tx := &genproto.Transaction{
		Version: 1,
	}

	block.Transactions = append(block.Transactions, tx)

	_ = SignBlock(privKey, block)

	assert.True(t, VerifyRootHash(block))

	block.Header.RootHash = []byte("invalid")
	assert.False(t, VerifyRootHash(block))
}
