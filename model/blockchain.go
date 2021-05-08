package model

import(
	"os"
	"encoding/hex"
	"fmt"
	"log"
	"bytes"
	"errors"
	"crypto/ecdsa"
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

func CreateFristBlock(coinbase *Transaction) (*Block) {
	return NewBlock([]*Transaction{coinbase}, []byte{})
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

	cbtx := NewCoinBaseTX(address, "DQT")
	fristBlock := CreateFristBlock(cbtx)

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
	fmt.Printf("TestNew: %d\n\n", len(block.transactions))
	return block
}

func (bc *BlockChain) FindTransaction(ID []byte) (Transaction, error) {
	bci := bc.Iterator()

	for {
		block := bci.Next()

		tx := block.transaction{
		if bytes.Compare(tx.Sender, ID) == 0 {
			return *tx, nil
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return Transaction{}, errors.New("Transaction is not found")
}




func (bc *BlockChain) MineBlock(transactions []*Transaction) *Block {
	var lastHash []byte

	for _, tx := range transactions {
		if bc.VerifyTransaction(tx) != true {
			log.Panic("ERROR: Invalid transaction")
		}
	}

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(transactions, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}

		bc.tip = newBlock.Hash

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return newBlock
}


func (bc *BlockChain) FindUnConfirmTransactions() []Transaction {
	var unconfirmTXs []Transaction

	for {
		block := bci.Next()
		if block.transaction.Status == "Uncofirm"{
			unconfirmTXs = unconfirmTXs.append(unconfirmTXs,block.transaction);
		}
		

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return unconfirmTXs
}

func (bc *BlockChain) FindConfirmTransactions() []Transaction {
	var confirmTXs []Transaction

	for {
		block := bci.Next()
		if block.transaction.Status == "Cofirm"{
			unconfirmTXs = unconfirmTXs.append(unconfirmTXs,block.transaction);
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