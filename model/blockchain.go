package models

type BlockChain struct {
	Blocks []*Block
}

func (bc *BlockChain)AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}

func CreateFristBlock() (*Block) {
	return NewBlock("Frist Block", []byte{})
}

func NewBlockChain() *BlockChain {
	return &BlockChain{[]*Block{CreateFristBlock()}}
}