package node

import (
	"testing"

	"github.com/oleglegun/blockchain-btc/internal/cryptography"
	"github.com/oleglegun/blockchain-btc/internal/genproto"
	"github.com/oleglegun/blockchain-btc/internal/random"
	"github.com/oleglegun/blockchain-btc/internal/types"
	"github.com/stretchr/testify/require"
)

const (
	genesisBlockTx0Hash = "bc88af88ffccbc54dbf64bef0b865568c974844352a8b989c1ebcd914defd27c"
)

func TestNewChain(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore(), NewMemoryTxStore(), NewMemoryUTXOStore())
	require.NotNil(t, chain)
	require.Equal(t, 0, chain.Height())

	_, err := chain.GetBlockByHeight(0)
	require.Nil(t, err)
}

func TestAddBlock(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore(), NewMemoryTxStore(), NewMemoryUTXOStore())

	privKey := cryptography.NewPrivateKey()

	for i := 1; i < 10; i++ {
		block, err := createRandomSignedBlock(chain, privKey)
		require.Nil(t, err)

		blockHash := types.HashBlockBytes(block)
		height := i

		require.Nil(t, chain.AddBlock(block))

		fetchedBlock, err := chain.GetBlockByHash(blockHash)
		require.Nil(t, err)
		require.Equal(t, block, fetchedBlock)

		fetchedBlockByHeight, err := chain.GetBlockByHeight(height)
		require.Nil(t, err)
		require.Equal(t, block, fetchedBlockByHeight)
	}
}

func TestChainHeight(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore(), NewMemoryTxStore(), NewMemoryUTXOStore())

	privKey := cryptography.NewPrivateKey()

	for i := 1; i < 10; i++ {
		block, err := createRandomSignedBlock(chain, privKey)
		require.Nil(t, err)

		err = chain.AddBlock(block)
		require.Nil(t, err)
		require.Equal(t, chain.Height(), i)
	}
}

func createRandomSignedBlock(chain *Chain, privKey cryptography.PrivateKey) (*genproto.Block, error) {
	block := random.RandomBlock()

	prevBlock, err := chain.GetBlockByHeight(chain.Height())
	if err != nil {
		return nil, err
	}

	block.Header.PrevHash = types.HashBlockBytes(prevBlock)

	sig := types.SignBlock(privKey, block)
	block.Signature = sig.Bytes()
	block.PublicKey = privKey.Public().Bytes()

	return block, nil
}

func TestAddBlockWithTransactions(t *testing.T) {
	var (
		senderPrivKey   = cryptography.NewPrivateKeyFromString(genesisBlockSeed)
		receiverAddress = cryptography.NewPrivateKey().Public().Address().Bytes()
		chain           = NewChain(NewMemoryBlockStore(), NewMemoryTxStore(), NewMemoryUTXOStore())
	)
	block, err := createRandomSignedBlock(chain, senderPrivKey)
	require.Nil(t, err)

	genesisTx, err := chain.txStore.Get(genesisBlockTx0Hash)
	require.Nil(t, err)

	inputs := []*genproto.TxInput{
		{
			PrevTxHash:     types.HashTransactionBytes(genesisTx),
			PrevTxOutIndex: 0,
			PublicKey:      senderPrivKey.Public().Bytes(),
		},
	}

	outputs := []*genproto.TxOutput{
		{
			Amount:  100,
			Address: receiverAddress,
		},
		{
			Amount:  genesisBlockAmount - 100,
			Address: senderPrivKey.Public().Address().Bytes(),
		},
	}

	tx := &genproto.Transaction{
		Inputs:  inputs,
		Outputs: outputs,
	}

	sig := types.CalculateTransactionSignature(senderPrivKey, tx)

	tx.Inputs[0].Signature = sig.Bytes()

	block.Transactions = append(block.Transactions, tx)

	err = chain.AddBlock(block)
	require.Nil(t, err)
}

func TestAddBlockWithInsufficientFundsTx(t *testing.T) {
	var (
		senderPrivKey   = cryptography.NewPrivateKeyFromString(genesisBlockSeed)
		receiverAddress = cryptography.NewPrivateKey().Public().Address().Bytes()
		chain           = NewChain(NewMemoryBlockStore(), NewMemoryTxStore(), NewMemoryUTXOStore())
	)
	block, err := createRandomSignedBlock(chain, senderPrivKey)
	require.Nil(t, err)

	genesisTx, err := chain.txStore.Get(genesisBlockTx0Hash)
	require.Nil(t, err)

	inputs := []*genproto.TxInput{
		{
			PrevTxHash:     types.HashTransactionBytes(genesisTx),
			PrevTxOutIndex: 0,
			PublicKey:      senderPrivKey.Public().Bytes(),
		},
	}

	outputs := []*genproto.TxOutput{
		{
			Amount:  genesisBlockAmount,
			Address: receiverAddress,
		},
		{
			Amount:  1,
			Address: senderPrivKey.Public().Address().Bytes(),
		},
	}

	tx := &genproto.Transaction{
		Inputs:  inputs,
		Outputs: outputs,
	}

	sig := types.CalculateTransactionSignature(senderPrivKey, tx)

	tx.Inputs[0].Signature = sig.Bytes()

	block.Transactions = append(block.Transactions, tx)
	err = chain.AddBlock(block)
	require.NotNil(t, err)
}
