package cryptography

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePrivateKey(t *testing.T) {
	privKey := GeneratePrivateKey()

	assert.Equal(t, len(privKey.Bytes()), PrivKeyLen)
	pubKey := privKey.Public()

	assert.Equal(t, len(pubKey.Bytes()), PubKeyLen)
}

func TestPrivateKeySign(t *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.Public()

	differentPrivKey := GeneratePrivateKey()
	differentPubKey := differentPrivKey.Public()

	msg := []byte("message")
	differentMsg := []byte("different message")
	emptyMsg := make([]byte, 0)

	t.Run("Valid message verification", func(t *testing.T) {
		sig := privKey.Sign(msg)
		if !sig.Verify(pubKey, msg) {
			t.Errorf("Signature verification failed for valid message")
		}
	})

	t.Run("Different message verification", func(t *testing.T) {
		sig := privKey.Sign(msg)
		if sig.Verify(pubKey, differentMsg) {
			t.Errorf("Signature verification should fail for a different message")
		}
	})

	t.Run("Empty message verification", func(t *testing.T) {
		emptySig := privKey.Sign(emptyMsg)
		if !emptySig.Verify(pubKey, emptyMsg) {
			t.Errorf("Signature verification failed for empty message")
		}
	})

	t.Run("Different keys verification", func(t *testing.T) {
		sig := privKey.Sign(msg)
		if sig.Verify(differentPubKey, msg) {
			t.Errorf("Signature verification should fail for a different public key")
		}
	})
}

func TestNewPrivateKeyFromString(t *testing.T) {
	var (
		seed               = "852d9b8e11b181bcf81aad70689521c469a2a5d830a8cbe1df1a497a91c93c84"
		privKey            = GeneratePrivateKeyFromString(seed)
		addressStrExpected = "a6461be4eac9ff331cfa7709f657ab1094064007"
		address            = privKey.Public().Address()
	)

	assert.Equal(t, PrivKeyLen, len(privKey.Bytes()))

	assert.Equal(t, addressStrExpected, address.String())
}

func TestPublicKeyToAddress(t *testing.T) {
	seed := make([]byte, SeedLen)

	privKey := GeneratePrivateKeyFromSeed(seed)
	pubKey := privKey.Public()

	address := pubKey.Address()

	assert.Equal(t, len(address.Bytes()), AddressLen)
}
