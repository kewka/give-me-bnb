package main

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	GasLimit uint64 = 21000
	// 0.99 BNB
	BnbValue = big.NewInt(0.99e18)
)

func WaitTx(ctx context.Context, client *ethclient.Client, tx string) error {
	hash := common.HexToHash(tx)
	for {
		receipt, err := client.TransactionReceipt(ctx, hash)
		if receipt != nil {
			return nil
		}

		if err != nil && err != ethereum.NotFound {
			return err
		}

		time.Sleep(time.Second)
	}
}

func SendBnb(ctx context.Context, client *ethclient.Client, from *ecdsa.PrivateKey, to string) (string, error) {
	nonce, err := client.PendingNonceAt(ctx, crypto.PubkeyToAddress(from.PublicKey))
	if err != nil {
		return "", err
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return "", err
	}

	tx := types.NewTransaction(nonce, common.HexToAddress(to), BnbValue, GasLimit, gasPrice, nil)

	chainId, err := client.NetworkID(ctx)
	if err != nil {
		return "", err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), from)
	if err != nil {
		return "", err
	}

	return signedTx.Hash().Hex(), client.SendTransaction(ctx, signedTx)
}
