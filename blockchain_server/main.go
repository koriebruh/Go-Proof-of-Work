package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/koriebruh/block"
	"github.com/koriebruh/dto"
	"github.com/koriebruh/utils"
	"github.com/koriebruh/wallet"
	"log"
	"net/http"
	"strconv"
)

func init() {
	log.SetPrefix("[blockchain_server]: ")
}

func main() {

	// parse command line arguments so for run we do =>  go run main.go -port 8080
	port := flag.Int("port", 8080, "Port to run the blockchain server on")
	flag.Parse()
	bcs := NewBlockchainServer(SetPort(uint16(*port)))
	bcs.Run()
}

var cache = make(map[string]*block.Blockchain)

type BlockchainServer struct {
	port uint16
}

type Opts func(*BlockchainServer)

func SetPort(port uint16) Opts {
	return func(b *BlockchainServer) {
		b.port = port
	}
}

func NewBlockchainServer(ops ...Opts) *BlockchainServer {
	bs := &BlockchainServer{
		port: 5000,
	}

	for _, op := range ops {
		op(bs)
	}
	return bs
}

func (bcs *BlockchainServer) Run() {
	http.HandleFunc("/", HelloWorld)
	http.HandleFunc("/blockchain", bcs.GetChainHandler)
	http.HandleFunc("/wallet", bcs.MakeWalletHandler)
	http.HandleFunc("/transaction", bcs.CreateTransactionHandler)

	fmt.Printf("starting blockchain server on port %d\n", bcs.port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", bcs.port), nil); err != nil {
		panic(err)
	}
}

func (bcs *BlockchainServer) GetBlockchain() *block.Blockchain {
	bc, ok := cache["blockchain"]
	if !ok {
		minerWallet := wallet.NewWallet()
		bc = block.NewBlockchain(minerWallet.BlockchainAddr(), bcs.port)
		cache["blockchain"] = bc
		fmt.Printf("private key %v\n", minerWallet.PrivateKeyStr())
		fmt.Printf("public key %v\n", minerWallet.PublicKeyStr())
		fmt.Printf("blockchain address %s\n", minerWallet.BlockchainAddr())
	}

	return bc
}

func (bcs *BlockchainServer) GetChainHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case
		http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		bc := bcs.GetBlockchain()
		m, _ := bc.MarshalJSON()
		fmt.Fprint(w, string(m[:]))
		w.WriteHeader(http.StatusOK)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, Blockchain!")
}

func (bcs *BlockchainServer) MakeWalletHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		w.Header().Set("Content-Type", "application/json")
		wallet := wallet.NewWallet()
		json, _ := wallet.MarshalJSON()
		fmt.Fprint(w, string(json[:]))
		w.WriteHeader(http.StatusCreated)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (bcs *BlockchainServer) CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		var t dto.TransactionRequest
		if err := decoder.Decode(&t); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		t.Validate()

		// do validation  struct
		// do validation amount enough or no
		// send transaction to pool
		// etc

		pubKey := utils.PublicKeyFromString(*t.SenderPublicKey)
		privKey := utils.PrivateKeyFromString(*t.SenderPrivateKey, pubKey)
		value, err := strconv.ParseFloat(*t.Value, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "Invalid value", http.StatusBadRequest)
			return
		}
		value32 := float32(value)

		fmt.Printf("sender private key %s\n", pubKey)
		fmt.Printf("sender public key %s\n", privKey)
		fmt.Printf("sender value key %v\n", value32)

		transaction := wallet.NewTransaction(privKey, pubKey, *t.SenderBlockchainAddress, *t.ReceiverBlockchainAddress, value32)
		signature := transaction.GenerateSignature()
		signatureStr := signature.String()

		_ = &block.TransactionRequest{
			SenderBlockchainAddr:    *t.SenderBlockchainAddress,
			RecipientBlockchainAddr: *t.ReceiverBlockchainAddress,
			SenderPublicKey:         *t.SenderPublicKey,
			Value:                   value32,
			Signature:               signatureStr,
		}

		fmt.Fprint(w, string("success create transaction"))
		w.WriteHeader(http.StatusCreated)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}
