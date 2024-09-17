package node

import (
	"fmt"
	"sync"

	"github.com/oleglegun/blockchain-btc/internal/genproto"
	"github.com/oleglegun/blockchain-btc/internal/types"
)

//-----------------------------------------------------------------------------
//  TxStorer
//-----------------------------------------------------------------------------

type TxStore interface {
	Get(hash string) (*genproto.Transaction, error)
	Put(*genproto.Transaction) error
}

type MemoryTxStore struct {
	sync.RWMutex
	txMap map[string]*genproto.Transaction
}

func NewMemoryTxStore() *MemoryTxStore {
	return &MemoryTxStore{
		txMap: make(map[string]*genproto.Transaction),
	}
}

func (s *MemoryTxStore) Get(hash string) (*genproto.Transaction, error) {
	s.RLock()
	defer s.RUnlock()
	tx, ok := s.txMap[hash]
	if !ok {
		return nil, fmt.Errorf("transaction [%s] is not found", hash)
	}

	return tx, nil
}

func (s *MemoryTxStore) Put(tx *genproto.Transaction) error {
	hash := types.HashTransactionString(tx)

	s.Lock()
	defer s.Unlock()
	s.txMap[hash] = tx

	return nil
}

//-----------------------------------------------------------------------------
//  UTXOStorer
//-----------------------------------------------------------------------------

type UTXOStore interface {
	Get(hash string) (*UTXO, error)
	Put(utxo *UTXO) error
}

type MemoryUTXOStore struct {
	sync.RWMutex
	utxoMap map[string]*UTXO
}

func NewMemoryUTXOStore() *MemoryUTXOStore {
	return &MemoryUTXOStore{
		utxoMap: make(map[string]*UTXO),
	}
}

func (s *MemoryUTXOStore) Get(hash string) (*UTXO, error) {
	s.RLock()
	defer s.RUnlock()

	utxo, ok := s.utxoMap[hash]
	if !ok {
		return nil, fmt.Errorf("utxo [%s] is not found", hash)
	}

	return utxo, nil
}

func (s *MemoryUTXOStore) Put(utxo *UTXO) error {
	s.Lock()
	defer s.Unlock()

	key := getUTXOKey(utxo.Hash, utxo.OutIndex)
	s.utxoMap[key] = utxo

	return nil
}

func getUTXOKey(hash string, outIndex int) string {
	return fmt.Sprintf("%s:%d", hash, outIndex)
}

//-----------------------------------------------------------------------------
//  BlockStorer
//-----------------------------------------------------------------------------

type BlockStore interface {
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
	hash := types.HashBlockString(block)
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
