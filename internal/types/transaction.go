package types

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/oleglegun/blockchain-btc/internal/cryptography"
	"github.com/oleglegun/blockchain-btc/internal/genproto"
	"google.golang.org/protobuf/proto"
)

func HashTransactionBytes(tx *genproto.Transaction) []byte {
	b, err := proto.Marshal(tx)
	if err != nil {
		panic(err)
	}

	hash := sha256.Sum256(b)
	return hash[:]
}

func HashTransactionString(tx *genproto.Transaction) string {
	return hex.EncodeToString(HashTransactionBytes(tx))
}

func CalculateTransactionSignature(privKey cryptography.PrivateKey, tx *genproto.Transaction) cryptography.Signature {
	return privKey.Sign(HashTransactionBytes(tx))
}

// VerifyTransaction verifies the transaction by checking the signature of each input.
func VerifyTransaction(tx *genproto.Transaction) bool {
	var isValid = true

	inputSignatures := make([][]byte, len(tx.Inputs))

	// Remove signatures
	for idx, input := range tx.Inputs {
		if len(input.Signature) == 0 {
			panic("tx signature is empty")
		}
		inputSignatures[idx] = input.Signature
		input.Signature = nil
	}

	for idx, input := range tx.Inputs {
		sig := cryptography.NewSignatureFromBytes(inputSignatures[idx])
		pubKey := cryptography.NewPublicKeyFromBytes(input.PublicKey)

		if !sig.Verify(pubKey, HashTransactionBytes(tx)) {
			isValid = false
			break
		}
	}

	// Put back signatures
	for idx, input := range tx.Inputs {
		input.Signature = inputSignatures[idx]
	}

	return isValid
}
