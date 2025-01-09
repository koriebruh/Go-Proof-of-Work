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

const (
	MINING_DIFFICULTY int = 3
	MINING_SENDER         = "coinbase"
	MINING_REWARD         = 1.0
)

type Block struct {
	timeStamp    int64
	nonce        int
	prevHash     [32]byte
	transactions []*Transaction
}

func NewBlock(prevHash [32]byte, nonce int, transactions []*Transaction) *Block {
	b := new(Block)
	b.timeStamp = time.Now().UnixNano()
	b.prevHash = prevHash
	b.nonce = nonce
	b.transactions = transactions
	return b
}

func (b *Block) Print() {
	// print data in 1 block
	fmt.Printf("timeStamp  	: %d\n", b.timeStamp)
	fmt.Printf("prevHash   	: %x\n", b.prevHash) //x mean hexadecimal
	fmt.Printf("nonce		: %d\n", b.nonce)
	fmt.Printf("transaction	:\n")
	for _, transaction := range b.transactions {
		transaction.Print()
	}
	fmt.Printf("\n")
}

// Hash the block
func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	return sha256.Sum256([]byte(m))
}

func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		TimeStamp    int64          `json:"time_stamp"`
		PrevHash     [32]byte       `json:"prev_hash"`
		Nonce        int            `json:"nonce"`
		Transactions []*Transaction `json:"transactions"`
	}{
		TimeStamp:    b.timeStamp,
		PrevHash:     b.prevHash,
		Nonce:        b.nonce,
		Transactions: b.transactions,
	})
}

type Blockchain struct {
	transactionPool []*Transaction
	chain           []*Block
	blockchainAddr  string
}

// CreateBlock add block and insert into chain
func (bc *Blockchain) CreateBlock(prevHash [32]byte, nonce int) *Block {
	b := NewBlock(prevHash, nonce, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{} // clean transaction pool
	return b
}

func NewBlockchain(blockchainAddr string) *Blockchain {
	bc := new(Blockchain)
	b := Block{}
	bc.CreateBlock(b.prevHash, 0)
	bc.blockchainAddr = blockchainAddr
	return bc
}

func (bc *Blockchain) Print() {
	// print every blocks
	for i, block := range bc.chain {
		fmt.Printf("%s <= chain %d => %s\n", strings.Repeat("<=", 16), i,
			strings.Repeat("=>", 16))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 60))
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) AddTransaction(sender string, recipient string, value float32) {
	transaction := NewTransaction(sender, recipient, value)
	bc.transactionPool = append(bc.transactionPool, transaction)
}

func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, t := range bc.transactionPool {
		transactions = append(transactions,
			NewTransaction(t.senderBlockchainAddr, t.recipientBlockchainAddr, t.value))
	}
	return transactions
}

func (bc *Blockchain) ValidProof(nonce int, prevHash [32]byte, transaction []*Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty) // how many 0
	guessBlock := Block{0, nonce, prevHash, transaction}
	guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())
	return guessHashStr[:difficulty] == zeros // ex diff = 3 will true if  000xxxx
}

func (bc *Blockchain) ProofOfWork() int {
	transactionPool := bc.CopyTransactionPool()
	prevHash := bc.LastBlock().Hash()
	nonce := 0
	for {
		if bc.ValidProof(nonce, prevHash, transactionPool, MINING_DIFFICULTY) {
			break
		}
		nonce += 1
	}
	return nonce
}

// Mining  func for mine and gift reward miner
func (bc *Blockchain) Mining() bool {
	bc.AddTransaction(MINING_SENDER, bc.blockchainAddr, MINING_REWARD)
	nonce := bc.ProofOfWork()
	prevHash := bc.LastBlock().Hash()
	bc.CreateBlock(prevHash, nonce)
	log.Println("action : mining ==> success")
	return true
}

// CalculateTotalAmount to get total amount
func (bc *Blockchain) CalculateTotalAmount(blockchainAddr string) float32 {
	var totalAmount float32 // to catch value
	for _, b := range bc.chain {
		for _, t := range b.transactions {
			value := t.value
			// if addr as recipient got +
			if blockchainAddr == t.recipientBlockchainAddr {
				totalAmount += value
			}
			// if addr as sender got -
			if blockchainAddr == t.senderBlockchainAddr {
				totalAmount -= value
			}
		}
	}
	return totalAmount
}

type Transaction struct {
	senderBlockchainAddr    string
	recipientBlockchainAddr string
	value                   float32
}

func NewTransaction(sender string, recipient string, value float32) *Transaction {
	return &Transaction{sender, recipient, value}
}

func (t *Transaction) Print() {
	fmt.Printf("        %s\n", strings.Repeat("-", 44))
	fmt.Printf("		sender_blockchain_address : %s\n", t.senderBlockchainAddr)
	fmt.Printf("		sender_recipient_address  : %s\n", t.recipientBlockchainAddr)
	fmt.Printf("		value				      : %.1f\n", t.value)
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		SenderBlockchainAddr    string  `json:"sender_blockchain_addr"`
		RecipientBlockchainAddr string  `json:"recipient_blockchain_addr"`
		Value                   float32 `json:"value"`
	}{
		SenderBlockchainAddr:    t.senderBlockchainAddr,
		RecipientBlockchainAddr: t.recipientBlockchainAddr,
		Value:                   t.value,
	})
}

func main() {
	blockChain := NewBlockchain("JF92")

	blockChain.AddTransaction("Allah", "Jamal", 99999999999999999999909999999)
	blockChain.Mining()

	blockChain.AddTransaction("jamal", "mom", 24000000000000)
	blockChain.AddTransaction("jamal", "dad", 20000000000000)
	blockChain.Mining()

	blockChain.AddTransaction("jamal", "sis", 210000000)
	blockChain.AddTransaction("C", "D", 2)
	//blockChain.AddTransaction("D", "C", 2)
	blockChain.Mining()
	blockChain.Print()

	fmt.Println("JAMAL DUID", blockChain.CalculateTotalAmount("Jamal"))
	fmt.Println("JF92 DUID", blockChain.CalculateTotalAmount("JF92"))
	fmt.Println("C DUID", blockChain.CalculateTotalAmount("C"))
	fmt.Println("D DUID", blockChain.CalculateTotalAmount("D"))
	//fmt.Println("D DUID", blockChain.CalculateTotalAmount("D"))

}

// NOTE FOR SOLVE :
// - PrevHash
// - MinerAddr
// - Add validate totalAmount cant negative

// How actually Tx in blockchain
// example :
// Rei's as sender make tx where inside (sender, recipient, value) This tx signed with digital signature generate from
// privateKey Rei's. The resulting transaction (including the signature) is broadcast to the blockchain network.
// Nodes validation the digital signature rei's usage publicKey rei's to confirm Rei is the real sender.
// Afterward, nodes check if Rei has sufficient funds to complete the transaction.
