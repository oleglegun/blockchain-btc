package node

import (
	"testing"

	"github.com/oleglegun/blockchain-btc/internal/cryptography"
	"github.com/oleglegun/blockchain-btc/internal/genproto"
	"github.com/oleglegun/blockchain-btc/internal/random"
	"github.com/oleglegun/blockchain-btc/internal/types"
	"github.com/stretchr/testify/require"
)

func TestNewChain(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore())
	require.NotNil(t, chain)
	require.Equal(t, 0, chain.Height())

	_, err := chain.GetBlockByHeight(0)
	require.Nil(t, err)
}

func TestAddBlock(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore())

	privKey := cryptography.GeneratePrivateKey()

	for i := 1; i < 10; i++ {
		block, err := createSignedBlock(chain, privKey)
		require.Nil(t, err)

		blockHash := types.HashBlock(block)
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
	chain := NewChain(NewMemoryBlockStore())

	privKey := cryptography.GeneratePrivateKey()

	for i := 1; i < 10; i++ {
		block, err := createSignedBlock(chain, privKey)
		require.Nil(t, err)

		err = chain.AddBlock(block)
		require.Nil(t, err)
		require.Equal(t, chain.Height(), i)
	}
}

func createSignedBlock(chain *Chain, privKey cryptography.PrivateKey) (*genproto.Block, error) {
	block := random.RandomBlock()

	prevBlock, err := chain.GetBlockByHeight(chain.Height())
	if err != nil {
		return nil, err
	}

	block.Header.PrevHash = types.HashBlock(prevBlock)

	sig := types.SignBlock(privKey, block)
	block.Signature = sig.Bytes()
	block.PublicKey = privKey.Public().Bytes()

	return block, nil
}
