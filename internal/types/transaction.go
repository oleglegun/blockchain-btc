package types

import (
	"crypto/sha256"

	"github.com/oleglegun/blockchain-btc/internal/cryptography"
	"github.com/oleglegun/blockchain-btc/internal/genproto"
	"google.golang.org/protobuf/proto"
)

func CalculateTransactionHash(tx *genproto.Transaction) []byte {
	b, err := proto.Marshal(tx)
	if err != nil {
		panic(err)
	}

	hash := sha256.Sum256(b)
	return hash[:]
}

func CalculateTransactionSignature(privKey cryptography.PrivateKey, tx *genproto.Transaction) cryptography.Signature {
	return privKey.Sign(CalculateTransactionHash(tx))
}

func VerifyTransaction(tx *genproto.Transaction) bool {
	var isValid = true

	inputSignatures := make([][]byte, len(tx.Inputs))

	// Remove signatures
	for idx, input := range tx.Inputs {
		inputSignatures[idx] = input.Signature
		input.Signature = nil
	}

	for idx, input := range tx.Inputs {
		sig := cryptography.NewSignatureFromBytes(inputSignatures[idx])
		pubKey := cryptography.NewPublicKeyFromBytes(input.PublicKey)

		if !sig.Verify(pubKey, CalculateTransactionHash(tx)) {
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
