package cryptography

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
)

const (
	PrivKeyLen = 64
	PubKeyLen  = 32
	SeedLen    = 32
	AddressLen = 20
	SigLen     = 64
)

func NewPrivateKey() PrivateKey {
	seed := make([]byte, SeedLen)
	_, err := io.ReadFull(rand.Reader, seed)
	if err != nil {
		log.Fatal(err)
	}
	return PrivateKey{key: ed25519.NewKeyFromSeed(seed)}
}

func NewPrivateKeyFromBytes(b []byte) PrivateKey {
	if len(b) != PrivKeyLen {
		log.Fatalf("private key length should be %d bytes", PrivKeyLen)
	}

	return PrivateKey{key: ed25519.PrivateKey(b)}
}

func NewPrivateKeyFromSeed(seed []byte) PrivateKey {
	if len(seed) != SeedLen {
		log.Fatalf("seed length should be %d bytes", SeedLen)
	}

	return PrivateKey{key: ed25519.NewKeyFromSeed(seed)}
}

// NewPrivateKeyFromString creates a new PrivateKey from a hexadecimal string representation of the private key seed.
//
// The input string must be exactly 64 hexadecimal characters (32 bytes) long, representing the 32-byte private key seed.
// If the input string is not valid hex or the length is incorrect, the function will log a fatal error.
//
// This function is useful for deserializing a private key from a string representation, such as when reading from a configuration file or other storage.
func NewPrivateKeyFromString(str string) PrivateKey {
	seed, err := hex.DecodeString(str)
	if err != nil {
		log.Fatal(err)
	}
	return NewPrivateKeyFromSeed(seed)
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
	b := make([]byte, PubKeyLen)
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
	if len(b) != PubKeyLen {
		log.Fatal("invalid public key length")
	}
	return PublicKey{
		key: ed25519.PublicKey(b),
	}
}

func (p PublicKey) Bytes() []byte {
	return p.key
}

func (p PublicKey) String() string {
	return hex.EncodeToString(p.key)
}

func (p PublicKey) Address() Address {
	return Address{
		// Return the last 20 bytes of the public key.
		// Similar to how Ethereum addresses are derived from the last 20 bytes of the Keccak-256 hash of the public key.
		value: p.key[len(p.key)-AddressLen:],
	}
}

/*-----------------------------------------------------------------------------
 *  Signature
 *----------------------------------------------------------------------------*/

type Signature struct {
	value []byte
}

func NewSignatureFromBytes(b []byte) Signature {
	if len(b) != SigLen {
		log.Fatal("invalid signature length")
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
