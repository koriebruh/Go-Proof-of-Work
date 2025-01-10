package main

import (
	"fmt"
	"github.com/koriebruh/block"
	"github.com/koriebruh/wallet"
)

func main() {

	walletSatoshi := wallet.NewWallet()
	walletA := wallet.NewWallet()
	walletB := wallet.NewWallet()

	t := wallet.NewTransaction(walletA.PrivateKey(), walletA.PublicKey(), walletA.BlockchainAddr(), walletB.BlockchainAddr(), 1.0)

	bc := block.NewBlockchain(walletSatoshi.BlockchainAddr())
	isAdd := bc.AddTransaction(walletA.BlockchainAddr(), walletB.BlockchainAddr(), 1.0, walletA.PublicKey(), t.GenerateSignature())
	fmt.Println("is add ?", isAdd)

}
