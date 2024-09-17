package node

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/oleglegun/blockchain-btc/internal/cryptography"
	"github.com/oleglegun/blockchain-btc/internal/genproto"
	"github.com/oleglegun/blockchain-btc/internal/types"
)

const (
	genesisBlockSeed    = "b69d0c49b336d58aca501f6a0ba60c933b5904bea73a6c288639b9c3c830627f"
	genesisBlockVersion = 1
	genesisBlockHeight  = 0
	genesisBlockAmount  = 1e6
)

type Chain struct {
	txStore      TxStore
	blockStore   BlockStore
	utxoStore    UTXOStore
	blockHeaders *BlockHeaderList
}

func NewChain(bs BlockStore, txs TxStore, utxos UTXOStore) *Chain {
	chain := &Chain{
		txStore:      txs,
		utxoStore:    utxos,
		blockStore:   bs,
		blockHeaders: NewBlockHeaderList(),
	}

	chain.addBlock(createGenesisBlock())
	return chain
}

func (c *Chain) AddBlock(block *genproto.Block) error {
	if err := c.ValidateBlock(block); err != nil {
		return err
	}

	return c.addBlock(block)
}

// We cannot validate the genesis block, that is why AddBlock is split in 2 funcs.
func (c *Chain) addBlock(block *genproto.Block) error {
	c.blockHeaders.Add(block.Header)

	for _, tx := range block.Transactions {
		if err := c.txStore.Put(tx); err != nil {
			return fmt.Errorf("failed to put transaction into store: %w", err)
		}

		hash := types.HashTransactionString(tx)

		for idx, out := range tx.Outputs {
			utxo := NewUTXO(hash, idx, out.Amount)

			if err := c.utxoStore.Put(utxo); err != nil {
				return fmt.Errorf("failed to put utxo into store: %w", err)
			}
		}
	}

	return c.blockStore.Put(block)
}

func (c *Chain) ValidateBlock(block *genproto.Block) error {
	if !types.VerifyBlock(block) {
		return fmt.Errorf("block with hash %s has an invalid signature", types.HashBlockString(block))
	}

	currentBlock, err := c.GetBlockByHeight(c.Height())
	if err != nil {
		return err
	}

	currentHash := types.HashBlockBytes(currentBlock)

	if !bytes.Equal(block.Header.PrevHash, currentHash) {
		return fmt.Errorf("block with hash %s is not a successor of the current block", types.HashBlockString(block))
	}

	for _, tx := range block.Transactions {
		if err := c.ValidateTransaction(tx); err != nil {
			return fmt.Errorf("failed to validate transaction: %w", err)
		}
	}

	return nil
}

func (c *Chain) ValidateTransaction(tx *genproto.Transaction) error {
	if !types.VerifyTransaction(tx) {
		return fmt.Errorf("transaction with hash %s is invalid", types.HashTransactionString(tx))
	}

	inputSum, err := c.sumTotalInputAmount(tx)
	if err != nil {
		return fmt.Errorf("failed to sum total input amount: %w", err)
	}

	outputSum, err := c.sumTotalOutputAmount(tx)
	if err != nil {
		return fmt.Errorf("failed to sum total output amount: %w", err)
	}

	if inputSum < outputSum {
		return fmt.Errorf("transaction with hash %s has insufficient funds", types.HashTransactionString(tx))
	}

	return nil
}

func (c *Chain) sumTotalInputAmount(tx *genproto.Transaction) (int64, error) {
	var sumInputs int64

	for _, input := range tx.Inputs {
		key := fmt.Sprintf("%s:%d", hex.EncodeToString(input.PrevTxHash), input.PrevTxOutIndex)
		utxo, err := c.utxoStore.Get(key)
		if err != nil {
			return 0, fmt.Errorf("failed to get utxo %s: %w", key, err)
		}

		if utxo.IsSpent {
			return 0, fmt.Errorf("utxo %s is already spent", key)
		}

		sumInputs += utxo.Amount
	}

	return sumInputs, nil
}

func (c *Chain) sumTotalOutputAmount(tx *genproto.Transaction) (int64, error) {
	var sumOutputs int64
	for _, output := range tx.Outputs {
		if output.Amount < 0 {
			return 0, fmt.Errorf("transaction with hash %s has negative output amount", types.HashTransactionString(tx))
		}

		sumOutputs += output.Amount
	}

	return sumOutputs, nil
}

func (c *Chain) GetBlockByHash(hash []byte) (*genproto.Block, error) {
	hashString := hex.EncodeToString(hash)
	block, err := c.blockStore.Get(hashString)
	if err != nil {
		return nil, fmt.Errorf("failed to get block by hash %s: %w", hashString, err)
	}
	return block, nil
}

func (c *Chain) GetBlockByHeight(height int) (*genproto.Block, error) {
	if height > c.Height() {
		return nil, fmt.Errorf("block with height %d doesn't exist", height)
	}

	blockHeader := c.blockHeaders.Get(height)
	hash := types.HashBlockHeader(blockHeader)
	block, err := c.GetBlockByHash(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to get block at height %d: %w", height, err)
	}
	return block, nil
}

func (c *Chain) Height() int {
	return c.blockHeaders.Height()
}

//-----------------------------------------------------------------------------
//  BlockHeaderList
//-----------------------------------------------------------------------------

type BlockHeaderList struct {
	headerList []*genproto.BlockHeader
}

func NewBlockHeaderList() *BlockHeaderList {
	return &BlockHeaderList{
		headerList: make([]*genproto.BlockHeader, 0),
	}
}

func (hs *BlockHeaderList) Add(h *genproto.BlockHeader) {
	hs.headerList = append(hs.headerList, h)
}

func (hs *BlockHeaderList) Get(height int) *genproto.BlockHeader {
	return hs.headerList[height]
}

func (hs *BlockHeaderList) Height() int {
	// blockchain always has a genesis block
	return len(hs.headerList) - 1
}

func createGenesisBlock() *genproto.Block {
	privKey := cryptography.NewPrivateKeyFromString(genesisBlockSeed)
	block := &genproto.Block{
		Header: &genproto.BlockHeader{
			Version:   genesisBlockVersion,
			Height:    genesisBlockHeight,
			Timestamp: time.Now().Unix(),
		},
	}

	tx := &genproto.Transaction{
		Version: genesisBlockVersion,
		Inputs:  []*genproto.TxInput{},
		Outputs: []*genproto.TxOutput{
			{
				Amount:  genesisBlockAmount,
				Address: privKey.Public().Address().Bytes(),
			},
		},
	}

	block.Transactions = append(block.Transactions, tx)

	types.SignBlock(privKey, block)

	return block
}
