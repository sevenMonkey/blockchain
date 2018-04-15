package main

import (
	"strconv"
	"bytes"
	"crypto/sha256"
	"time"
)

// block keeps block header
type Block struct {
	// The timestamp when block created
	Timestamp int64
	// The actual valuable information containing in the block
	Data []byte
	// The hash of the previous block
	PreBlockHash []byte
	// hash are block header
	Hash []byte
	Nonce int
}

func (b *Block)SetHash()  {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PreBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}

// NewBlock creates and returns Block
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}
