package blc

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
	"fmt"
	"strconv"
)

// Block keeps block headers
type Block struct {
	Timestamp     int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

// Serialize serializes the block
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// HashTransactions returns a hash of the transactions in the block
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

// NewBlock creates and returns Block
func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), transactions, prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// NewGenesisBlock creates and returns genesis Block
func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

// DeserializeBlock deserializes a block
func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}

// Show Block Information
func (b *Block) ShowBlockInfo() {
	fmt.Printf("Prev. hash: %x\n", b.PrevBlockHash)
	fmt.Printf("Hash: %x\n", b.Hash)
	fmt.Println("Transaction: ")
	for _, tx := range b.Transactions {
		fmt.Printf("Transaction ID: %x\n", tx.ID)
		for _, txi := range tx.Vin {
			fmt.Printf("Vin Txid: %x\n", txi.Txid)
			fmt.Println("Vin Vout: ", txi.Vout)
			fmt.Println("Vin ScriptSig: ", txi.ScriptSig)
		}
		fmt.Println("Transaction Vout: ", tx.Vout)
	}
	pow := NewProofOfWork(b)
	fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
	fmt.Println()
}