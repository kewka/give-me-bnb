package blockchain

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type Account struct {
	key *ecdsa.PrivateKey
}

func (acc Account) public() common.Address {
	return crypto.PubkeyToAddress(acc.key.PublicKey)
}

func (acc Account) PublicHex() string {
	return acc.public().Hex()
}

func (acc Account) PrivateHex() string {
	return hexutil.Encode(crypto.FromECDSA(acc.key))
}

func newAccount(key *ecdsa.PrivateKey, err error) (*Account, error) {
	if err != nil {
		return nil, err
	}
	return &Account{key: key}, nil
}

func ReadAccount(key string) (*Account, error) {
	return newAccount(crypto.HexToECDSA(key))
}

func GenerateAccount() (*Account, error) {
	return newAccount(crypto.GenerateKey())
}
