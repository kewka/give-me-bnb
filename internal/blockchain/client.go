package blockchain

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kewka/give-me-bnb/internal/blockchain/erc20"
)

const (
	GasLimit         = 21000
	TransferGasLimit = 100000
)

type Client struct {
	eclient *ethclient.Client
}

func NewClient(url string) (*Client, error) {
	eclient, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}
	return &Client{eclient: eclient}, nil
}

func (c *Client) Close() {
	c.eclient.Close()
}

func (c *Client) Wait(ctx context.Context, tx string) error {
	hash := common.HexToHash(tx)
	for {
		receipt, err := c.eclient.TransactionReceipt(ctx, hash)
		if receipt != nil {
			return nil
		}

		if err != nil && err != ethereum.NotFound {
			return err
		}

		time.Sleep(time.Second)
	}
}

func (c *Client) sendNative(
	ctx context.Context,
	from *Account,
	to *Account,
	amount *big.Int,
	gasPrice *big.Int,
) (string, error) {
	nonce, err := c.eclient.PendingNonceAt(ctx, common.HexToAddress(from.PublicHex()))
	if err != nil {
		return "", err
	}

	tx := types.NewTransaction(nonce, to.public(), amount, GasLimit, gasPrice, nil)

	chainId, err := c.eclient.NetworkID(ctx)
	if err != nil {
		return "", err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), from.key)
	if err != nil {
		return "", err
	}

	return signedTx.Hash().Hex(), c.eclient.SendTransaction(ctx, signedTx)
}

func (c *Client) WithdrawNative(ctx context.Context, from *Account, to *Account) (string, error) {
	balance, err := c.eclient.BalanceAt(ctx, from.public(), nil)
	if err != nil {
		return "", err
	}

	gasPrice, err := c.eclient.SuggestGasPrice(ctx)
	if err != nil {
		return "", err
	}

	gasAmount := new(big.Int).Mul(big.NewInt(GasLimit), gasPrice)
	amount := new(big.Int).Sub(balance, gasAmount)
	return c.sendNative(ctx, from, to, amount, gasPrice)
}

func (c *Client) WithdrawToken(ctx context.Context, from *Account, to *Account, tokenAddress string) (string, error) {
	contract, err := erc20.NewErc20(common.HexToAddress(tokenAddress), c.eclient)
	if err != nil {
		return "", err
	}

	balance, err := contract.BalanceOf(&bind.CallOpts{}, from.public())
	if err != nil {
		return "", err
	}

	chainId, err := c.eclient.ChainID(ctx)
	if err != nil {
		return "", err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(from.key, chainId)
	if err != nil {
		return "", err
	}

	auth.Context = ctx
	auth.GasLimit = TransferGasLimit
	auth.NoSend = true

	transferTx, err := contract.Transfer(auth, to.public(), balance)
	if err != nil {
		return "", err
	}

	nativeTx, err := c.sendNative(ctx, to, from, transferTx.Cost(), transferTx.GasPrice())
	if err != nil {
		return "", err
	}

	if err := c.Wait(ctx, nativeTx); err != nil {
		return "", err
	}

	return transferTx.Hash().Hex(), c.eclient.SendTransaction(ctx, transferTx)
}
