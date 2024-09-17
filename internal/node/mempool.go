package node

import (
	"sync"
	"time"

	"github.com/oleglegun/blockchain-btc/internal/genproto"
	"github.com/oleglegun/blockchain-btc/internal/types"
)

type Mempool struct {
	sync.RWMutex
	txMap map[string]*genproto.Transaction
	// txTimestampMap contains all transactions' timestamps (including cleared)
	txTimestampMap map[string]time.Time
}

func NewMempool() *Mempool {
	return &Mempool{
		txMap:          make(map[string]*genproto.Transaction),
		txTimestampMap: make(map[string]time.Time),
	}
}

func (p *Mempool) Size() int {
	p.RLock()
	defer p.RUnlock()
	return len(p.txMap)
}

func (p *Mempool) Clear() []*genproto.Transaction {
	p.Lock()
	defer p.Unlock()

	txList := make([]*genproto.Transaction, len(p.txMap))
	idx := 0
	for _, tx := range p.txMap {
		txList[idx] = tx
		idx++
	}
	p.txMap = make(map[string]*genproto.Transaction)

	return txList
}

func (p *Mempool) ClearProcessed(threshold time.Duration) []string {
	p.Lock()
	defer p.Unlock()

	now := time.Now()

	clearedTxHashList := make([]string, 0)

	for hash, timestamp := range p.txTimestampMap {
		if now.Sub(timestamp) > threshold {
			delete(p.txMap, hash)
			delete(p.txTimestampMap, hash)
			clearedTxHashList = append(clearedTxHashList, hash)
		}
	}

	return clearedTxHashList
}

// Has checks if the PROCESSED given transaction is present in the mempool.
// It returns true if the transaction is present, false otherwise.
func (p *Mempool) Has(tx *genproto.Transaction) bool {
	hash := types.HashTransactionString(tx)

	p.RLock()
	defer p.RUnlock()

	_, ok := p.txTimestampMap[hash]
	return ok
}

func (p *Mempool) Add(tx *genproto.Transaction) bool {
	hash := types.HashTransactionString(tx)

	p.Lock()
	defer p.Unlock()

	if _, exists := p.txTimestampMap[hash]; exists {
		return false
	}

	p.txMap[hash] = tx
	p.txTimestampMap[hash] = time.Now()
	return true
}
