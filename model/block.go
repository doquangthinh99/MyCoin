package model

import(
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Count			int
}

func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
	hash := sha256.Sum256(headers)

	b.Hash = hash[:]
}

func NewBlock(data string, prevBlockHash []byte) (*Block) {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{},0}
	pow := NewProofOfWork(block)

	count,hash := pow.Run()

	block.Hash = hash[:]
	block.Count = count

	return block
}