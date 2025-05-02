package test

import (
	"github.com/koriebruh/wallet"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMakeWallet(t *testing.T) {
	korieWallet := wallet.NewWallet()

	assert.Equal(t, len(korieWallet.PrivateKeyStr()), 64)
	assert.Equal(t, len(korieWallet.PublicKeyStr()), 128)
}
