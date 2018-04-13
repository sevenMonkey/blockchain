package main

type Blockchain struct {
	blocks []*Block
}

func (bc *Blockchain)AddBlock(data string)  {
	prevBlock := bc.blocks[len(bc.blocks) - 1]
	NewBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, NewBlock)
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

