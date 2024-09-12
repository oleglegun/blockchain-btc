package cryptography

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
)

const (
	privKeyLen = 64
	pubKeyLen  = 32
	seedLen    = 32
	addressLen = 20
	sigLen     = 64
)

func GeneratePrivateKey() PrivateKey {
	seed := make([]byte, seedLen)
	_, err := io.ReadFull(rand.Reader, seed)
	if err != nil {
		log.Fatal(err)
	}
	return PrivateKey{key: ed25519.NewKeyFromSeed(seed)}
}

func GeneratePrivateKeyFromSeed(seed []byte) PrivateKey {
	if len(seed) != seedLen {
		log.Fatalf("seed length should be %d bytes", seedLen)
	}

	return PrivateKey{key: ed25519.NewKeyFromSeed(seed)}
}

func GeneratePrivateKeyFromString(str string) PrivateKey {
	seed, err := hex.DecodeString(str)
	if err != nil {
		log.Fatal(err)
	}
	return GeneratePrivateKeyFromSeed(seed)
}

/*-----------------------------------------------------------------------------
 *  PrivateKey
 *----------------------------------------------------------------------------*/

type PrivateKey struct {
	key ed25519.PrivateKey
}

func (p PrivateKey) Bytes() []byte {
	return p.key
}

func (p PrivateKey) Sign(msg []byte) Signature {
	return Signature{
		value: ed25519.Sign(p.key, msg),
	}
}

func (p PrivateKey) Public() PublicKey {
	b := make([]byte, pubKeyLen)
	copy(b, p.Bytes()[32:])

	return PublicKey{
		key: b,
	}
}

/*-----------------------------------------------------------------------------
 *  PublicKey
 *----------------------------------------------------------------------------*/

type PublicKey struct {
	key ed25519.PublicKey
}

func NewPublicKeyFromBytes(b []byte) PublicKey {
	if len(b) != pubKeyLen {
		log.Fatal("invalid public key length")
	}
	return PublicKey{
		key: ed25519.PublicKey(b),
	}
}

func (p PublicKey) Bytes() []byte {
	return p.key
}

func (p PublicKey) Address() Address {
	return Address{
		// Return the last 20 bytes of the public key
		value: p.key[len(p.key)-addressLen:],
	}
}

/*-----------------------------------------------------------------------------
 *  Signature
 *----------------------------------------------------------------------------*/

// Change to pointer if needed (for len and cap)
type Signature struct {
	value []byte
}

func NewSignatureFromBytes(b []byte) Signature {
	if len(b) != sigLen {
		log.Fatal("signature length is incorrect")
	}
	return Signature{
		value: b,
	}
}

func (s Signature) Bytes() []byte {
	return s.value
}

func (s Signature) Verify(pubKey PublicKey, msg []byte) bool {
	return ed25519.Verify(pubKey.key, msg, s.value)
}

/*-----------------------------------------------------------------------------
 *  Address
 *----------------------------------------------------------------------------*/

type Address struct {
	value []byte
}

func (a Address) String() string {
	return hex.EncodeToString(a.value)
}

func (a Address) Bytes() []byte {
	return a.value
}
