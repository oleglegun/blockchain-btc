package node

import (
	"testing"

	"github.com/oleglegun/blockchain-btc/internal/random"
	"github.com/oleglegun/blockchain-btc/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestAddBlock(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore())

	for i := 0; i < 10; i++ {
		block := random.RandomBlock()
		blockHash := types.CalcBlockHash(block)
		height := i

		assert.Nil(t, chain.AddBlock(block))

		fetchedBlock, err := chain.GetBlockByHash(blockHash)
		assert.Nil(t, err)
		assert.Equal(t, block, fetchedBlock)

		fetchedBlockByHeight, err := chain.GetBlockByHeight(height)
		assert.Nil(t, err)
		assert.Equal(t, block, fetchedBlockByHeight)
	}
}

func TestChainHeight(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore())

	for i := 0; i < 10; i++ {
		block := random.RandomBlock()
		err := chain.AddBlock(block)
		assert.Nil(t, err)
		assert.Equal(t, chain.Height(), i)
	}
}
