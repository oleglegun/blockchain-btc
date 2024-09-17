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

type Chain struct {
	blockStore   BlockStorer
	blockHeaders *BlockHeaderList
}

func NewChain(bs BlockStorer) *Chain {
	chain := &Chain{
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
	return c.blockStore.Put(block)
}

func (c *Chain) ValidateBlock(block *genproto.Block) error {
	if !types.VerifyBlock(block) {
		return fmt.Errorf("block with hash %s has an invalid signature", hex.EncodeToString(types.HashBlock(block)))
	}

	currentBlock, err := c.GetBlockByHeight(c.Height())
	if err != nil {
		return err
	}

	currentHash := types.HashBlock(currentBlock)

	nextHash := hex.EncodeToString(types.HashBlock(block))

	if !bytes.Equal(block.Header.PrevHash, currentHash) {
		return fmt.Errorf("block with hash %s is not a successor of the current block", nextHash)
	}

	return nil
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
	privKey := cryptography.NewPrivateKey()
	block := &genproto.Block{
		Header: &genproto.BlockHeader{
			Version:   1,
			Height:    0,
			Timestamp: time.Now().Unix(),
		},
	}

	types.SignBlock(privKey, block)

	return block
}
