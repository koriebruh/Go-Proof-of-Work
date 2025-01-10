package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/koriebruh/utils"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	privateKey     *ecdsa.PrivateKey
	publicKey      *ecdsa.PublicKey
	blockchainAddr string
}

func NewWallet() *Wallet {
	// 1. Create ECDSA privateKey (32 bytes) publicKey (64 bytes)
	w := new(Wallet)
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	w.privateKey = privateKey
	w.publicKey = &w.privateKey.PublicKey
	// 2. Perform SHA-256 hashing on the publicKey (32 bytes)
	h2 := sha256.New()
	h2.Write(w.publicKey.X.Bytes())
	h2.Write(w.publicKey.Y.Bytes())
	digest2 := h2.Sum(nil)
	// 3. Perform RIPEMD-160 hashing on the result of SHA-265 (20 bytes)
	h3 := ripemd160.New()
	h3.Write(digest2)
	digest3 := h3.Sum(nil)
	// 4. add version byte in front of RIPEMD-160 hash (0x00 for Main Network)
	vd4 := make([]byte, 21)
	vd4[0] = 0x00
	copy(vd4[1:], digest3[:]) // merge
	// 5. Perform SHA-256 hash on the extended RIPEMD-160 result
	h5 := sha256.New()
	h5.Write(vd4)
	digest5 := h5.Sum(nil)
	// 6. Perform SHA-256 hash on the result of the previous SHA-256
	h6 := sha256.New()
	h6.Write(digest5)
	digest6 := h6.Sum(nil)
	// 7. Take the first 4 byte of the second SHA-256 hash for checksum
	chsum := digest6[:4]
	// 8. Add the checksum bytes from 7 at the end of extended RIPEMD-160 hash from 4 (25 bytes)
	dc8 := make([]byte, 25)
	copy(dc8[:21], vd4[:4])
	copy(dc8[:21], chsum[:])
	// 9. Convert the result from a byte string into base58
	addr := base58.Encode(dc8)
	w.blockchainAddr = addr
	return w
}

func (w *Wallet) PrivateKey() *ecdsa.PrivateKey {
	return w.privateKey
}

func (w *Wallet) PrivateKeyStr() string {
	return fmt.Sprintf("%x", w.privateKey.D.Bytes())
}

func (w *Wallet) PublicKey() *ecdsa.PublicKey {
	return w.publicKey
}

func (w *Wallet) PublicKeyStr() string {
	return fmt.Sprintf("%x%x", w.publicKey.X.Bytes(), w.publicKey.Y.Bytes())
}

// BlockchainAddr this is next usage for address in the blockchain
func (w *Wallet) BlockchainAddr() string {
	return w.blockchainAddr
}

type Transaction struct {
	senderPrivateKey        *ecdsa.PrivateKey
	senderPublicKey         *ecdsa.PublicKey
	senderBlockchainAddr    string
	recipientBlockchainAddr string
	value                   float32
}

func NewTransaction(senderPrivateKey *ecdsa.PrivateKey, senderPublicKey *ecdsa.PublicKey, senderBlockchainAddr string, recipientBlockchainAddr string, value float32) *Transaction {
	return &Transaction{senderPrivateKey, senderPublicKey, senderBlockchainAddr, recipientBlockchainAddr, value}
}

func (t *Transaction) GenerateSignature() *utils.Signature {
	m, _ := json.Marshal(t)
	hash := sha256.Sum256([]byte(m))
	r, s, _ := ecdsa.Sign(rand.Reader, t.senderPrivateKey, hash[:])
	return &utils.Signature{
		R: r,
		S: s,
	}
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

// SUMMARY, digitalSignature Generate
// 1. take tx data convert into json
// 2. hash tx data
// 3. sign the data use PrivateKey with ECDSA and return R,S
// 4. save R, S
// R means random value Create a different signature each time, even if the signed data is the same.
// S means signatureValue Binding the transaction data with the sender's private key, making the signature unique for that data.
