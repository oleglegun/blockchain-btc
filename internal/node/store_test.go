package node

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/oleglegun/blockchain-btc/internal/random"
	"github.com/oleglegun/blockchain-btc/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestNewMemoryBlockStore(t *testing.T) {
	store := NewMemoryBlockStore()

	assert.NotNil(t, store)
	assert.NotNil(t, store.blocks)
}

func TestMemoryBlockStore_Put(t *testing.T) {
	store := NewMemoryBlockStore()
	block := random.RandomBlock()

	err := store.Put(block)
	assert.Nil(t, err)

	hash := types.HashBlockString(block)
	storedBlock, exists := store.blocks[hash]
	assert.True(t, exists)
	assert.Equal(t, block, storedBlock)
}

func TestMemoryBlockStore_Get(t *testing.T) {
	store := NewMemoryBlockStore()
	block := random.RandomBlock()

	hash := types.HashBlockString(block)
	store.blocks[hash] = block

	retrievedBlock, err := store.Get(hash)
	assert.Nil(t, err)
	assert.Equal(t, block, retrievedBlock)
}

func TestMemoryBlockStore_Get_NonExistent(t *testing.T) {
	store := NewMemoryBlockStore()
	nonExistentHash := hex.EncodeToString(random.Random32ByteHash())

	block, err := store.Get(nonExistentHash)
	assert.Nil(t, block)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("block [%s] is not found", nonExistentHash), err.Error())
}
