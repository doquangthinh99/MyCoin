package model

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"encoding/binary"
	"log"
	"strings"
	"encoding/hex"
)

var (
	maxCount = math.MaxInt64
)

const difficulty = "0"

type ProofOfWork struct {
	block  *Block
}

func NewProofOfWork(b *Block) *ProofOfWork {

	pow := &ProofOfWork{b}

	return pow
}

func (pow *ProofOfWork) prepareData(count int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.HashTransactions(),
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(count)),
		},
		[]byte{},
	)

	return data
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	count := 0

	fmt.Printf("Mining a new block")
	for count < maxCount {
		data := pow.prepareData(count)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])

		if strings.HasPrefix( binToStr(hash[:]), difficulty ) {
			break
		} else {
			count++
		}
	}
	fmt.Printf("\n\n%d\n",count)

	return count, hash[:]
}

func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func (pow *ProofOfWork) Validate() bool {
	isValid := strings.HasPrefix( binToStr(pow.block.Hash[:]), difficulty)

	return isValid
}

func binToStr( bytes []byte ) string {
	return hex.EncodeToString( bytes )
}