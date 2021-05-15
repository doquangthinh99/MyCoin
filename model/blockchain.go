package model

import(
	"os"
	"fmt"
	"log"
	"bytes"
	"errors"
	"github.com/boltdb/bolt"
)

const dbFile = "BlockChain.db"
const blocksBucket = "blocks"
const subsidy = 10

type BlockChain struct {
	tip []byte
	db  *bolt.DB
}

type BlockChainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func CreateFristBlock() (*Block) {
	coinbase := CreateBaseTransaction()
	return NewBlock(coinbase, []byte{})
}

func NewBlockChain() *BlockChain {
	if dbExists() == false {
		fmt.Println("No existing BlockChain found. Create one first.")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	bc := BlockChain{tip, db}

	return &bc
}

func CreateBlockChain(address string) *BlockChain {
	if dbExists() {
		fmt.Println("BlockChain already exists.")
		os.Exit(1)
	}

	var tip []byte

	fristBlock := CreateFristBlock()

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Panic(err)
		}

		err = b.Put(fristBlock.Hash, fristBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), fristBlock.Hash)
		if err != nil {
			log.Panic(err)
		}
		tip = fristBlock.Hash

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	bc := BlockChain{tip, db}

	return &bc
}

func (bc *BlockChain) Iterator() *BlockChainIterator {
	bci := &BlockChainIterator{bc.tip, bc.db}

	return bci
}

func (i *BlockChainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	i.currentHash = block.PrevBlockHash
	return block
}

func (bc *BlockChain) FindTransaction(ID []byte) (Transaction, error) {
	bci := bc.Iterator()

	for {
		block := bci.Next()

		tx := block.transaction
		if bytes.Compare(tx.Sender, ID) == 0 {
			return tx, nil
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return Transaction{}, errors.New("Transaction is not found")
}




func (bc *BlockChain) MineBlock(block *Block,miner []byte)(int){
	
	pow := NewProofOfWork(block)
	count,_ := pow.Run()
	block.transaction.Status = "Confirm"
	block.transaction.Miner = miner
	block.Count = count
	return count
}


func (bc *BlockChain) FindUnConfirmTransactions() []Transaction {
	var unconfirmTXs []Transaction
	bci := bc.Iterator()
	for {
		block := bci.Next()
		if block.transaction.Status == "Uncofirm"{
			unconfirmTXs = append(unconfirmTXs,block.transaction);
		}
		

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return unconfirmTXs
}

func (bc *BlockChain) FindConfirmTransactions() []Transaction {
	var confirmTXs []Transaction
	bci := bc.Iterator()
	for {
		block := bci.Next()
		if block.transaction.Status == "Cofirm"{
			confirmTXs = append(confirmTXs,block.transaction);
		}
		

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return confirmTXs
}

func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}