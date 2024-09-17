package types

import (
	"testing"

	"github.com/oleglegun/blockchain-btc/internal/cryptography"
	"github.com/oleglegun/blockchain-btc/internal/genproto"
	"github.com/oleglegun/blockchain-btc/internal/random"
	"github.com/stretchr/testify/assert"
)

func TestCalculateTransactionHash(t *testing.T) {
	sender1PrivKey := cryptography.NewPrivateKey()
	sender1Addr := sender1PrivKey.Public().Address().Bytes()

	sender2PrivKey := cryptography.NewPrivateKey()
	_ = sender2PrivKey.Public().Address().Bytes()

	receiverPrivKey := cryptography.NewPrivateKey()
	receiverAddr := receiverPrivKey.Public().Address().Bytes()

	// Initial balance: 10

	txIn1 := &genproto.TxInput{
		PrevTxHash:     random.Random32ByteHash(),
		PrevTxOutIndex: 0,
		PublicKey:      sender1PrivKey.Public().Bytes(),
		// Signature will be set after constructing transaction
	}

	txIn2 := &genproto.TxInput{
		PrevTxHash:     random.Random32ByteHash(),
		PrevTxOutIndex: 1,
		PublicKey:      sender2PrivKey.Public().Bytes(),
		// Signature will be set after constructing transaction
	}

	txOut1 := &genproto.TxOutput{
		Amount:  9,
		Address: receiverAddr,
	}
	txOut2 := &genproto.TxOutput{
		Amount:  1,
		Address: sender1Addr,
	}

	tx := &genproto.Transaction{
		Version: 1,
		Inputs:  []*genproto.TxInput{txIn1, txIn2},
		Outputs: []*genproto.TxOutput{txOut1, txOut2},
	}

	sig1 := CalculateTransactionSignature(sender1PrivKey, tx)
	sig2 := CalculateTransactionSignature(sender2PrivKey, tx)
	txIn1.Signature = sig1.Bytes()
	txIn2.Signature = sig2.Bytes()

	assert.True(t, VerifyTransaction(tx))
}
