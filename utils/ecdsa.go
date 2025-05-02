package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"math/big"
)

type Signature struct {
	R *big.Int
	S *big.Int
}

func (s *Signature) String() string {
	return fmt.Sprintf("%064x%064x", s.R, s.S)
}

func String2BigIntTuple(s string) (big.Int, big.Int) {
	bx, _ := hex.DecodeString(s[:64])
	by, _ := hex.DecodeString(s[64:])

	var bigX, bigY big.Int
	bigX.SetBytes(bx)
	bigY.SetBytes(by)

	return bigX, bigY
}

func SignatureFromString(s string) *Signature {
	x, y := String2BigIntTuple(s)
	return &Signature{R: &x, S: &y}
}

func PublicKeyFromString(s string) *ecdsa.PublicKey {
	bx, by := String2BigIntTuple(s)
	return &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     &bx,
		Y:     &by,
	}
}

func PrivateKeyFromString(s string, pubKey *ecdsa.PublicKey) *ecdsa.PrivateKey {
	a, _ := hex.DecodeString(s[:])
	var bi big.Int
	bi.SetBytes(a)
	return &ecdsa.PrivateKey{
		PublicKey: *pubKey,
		D:         &bi,
	}
}
