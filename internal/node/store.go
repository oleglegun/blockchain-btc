package node

import (
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/oleglegun/blockchain-btc/internal/genproto"
	"github.com/oleglegun/blockchain-btc/internal/types"
)

type BlockStorer interface {
	Put(*genproto.Block) error
	Get(hash string) (*genproto.Block, error)
}

type MemoryBlockStore struct {
	sync.RWMutex
	blocks map[string]*genproto.Block
}

func NewMemoryBlockStore() *MemoryBlockStore {
	return &MemoryBlockStore{
		blocks: make(map[string]*genproto.Block),
	}
}

func (s *MemoryBlockStore) Put(block *genproto.Block) error {
	hash := hex.EncodeToString(types.HashBlock(block))
	s.Lock()
	defer s.Unlock()
	s.blocks[hash] = block

	return nil
}

func (s *MemoryBlockStore) Get(hash string) (*genproto.Block, error) {
	s.RLock()
	defer s.RUnlock()
	block, ok := s.blocks[hash]
	if !ok {
		return nil, fmt.Errorf("block [%s] is not found", hash)
	}

	return block, nil
}
