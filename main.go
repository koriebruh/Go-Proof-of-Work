package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

// when print something will be print tamp stamp first
func init() {
	log.SetPrefix("Blockchain: ")
}

type Block struct {
	nonce       int
	prevHash    [32]byte
	timeStamp   int64
	transaction []string
}

func NewBlock(prevHash [32]byte, nonce int) *Block {
	b := new(Block)
	b.timeStamp = time.Now().UnixNano()
	b.prevHash = prevHash
	b.nonce = nonce
	return b
}

func (b *Block) Print() {
	// print data in 1 block
	fmt.Printf("timeStamp  	: %d\n", b.timeStamp)
	fmt.Printf("prevHash   	: %x\n", b.prevHash) //x mean hexadecimal
	fmt.Printf("nonce		: %d\n", b.nonce)
	fmt.Printf("transaction	: %s\n", b.transaction)
}

// Hash the block
func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	fmt.Println(string(m))
	return sha256.Sum256([]byte(m))
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		TimeStamp   int64    `json:"time_stamp"`
		PrevHash    [32]byte `json:"prev_hash"`
		Nonce       int      `json:"nonce"`
		Transaction []string `json:"transaction"`
	}{
		TimeStamp:   b.timeStamp,
		PrevHash:    b.prevHash,
		Nonce:       b.nonce,
		Transaction: b.transaction,
	})
}

type Blockchain struct {
	transactionPool []*string
	chain           []*Block
}

// CreateBlock add block and insert into chain
func (bc *Blockchain) CreateBlock(prevHash [32]byte, nonce int) *Block {
	b := NewBlock(prevHash, nonce)
	bc.chain = append(bc.chain, b)
	return b
}

func NewBlockchain() *Blockchain {
	bc := new(Blockchain)
	var genesisHash [32]byte
	bc.CreateBlock(genesisHash, 0)
	return bc
}

func (bc *Blockchain) Print() {
	// print every blocks
	for i, block := range bc.chain {
		fmt.Printf("%s <= chain %d => %s\n", strings.Repeat("-", 32), i,
			strings.Repeat("-", 32))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 60))
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func main() {
	blockChain := NewBlockchain()
	prevHash := blockChain.LastBlock().Hash()

	blockChain.CreateBlock(prevHash, 832)
	blockChain.CreateBlock(prevHash, 1256)

	blockChain.Print()

}
