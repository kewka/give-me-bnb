package faucet

import "errors"

type Currency struct {
	symbol  string
	address string
}

func (c Currency) Symbol() string {
	return c.symbol
}

func (c Currency) Address() string {
	return c.address
}

var (
	BNB  = Currency{"BNB", ""}
	BTC  = Currency{"BTC", "0x6ce8da28e2f864420840cf74474eff5fd80e65b8"}
	BUSD = Currency{"BUSD", "0xed24fc36d5ee211ea25a80239fb8c4cfd80f12ee"}
	DAI  = Currency{"DAI", "0xec5dcb5dbf4b114c9d0f65bccab49ec54f6a0867"}
	ETH  = Currency{"ETH", "0xd66c6b4f0be8ce5b39d52e0fd1344c389929b378"}
	USDC = Currency{"USDC", "0x64544969ed7ebf5f083679233325356ebe738930"}
	USDT = Currency{"USDT", "0x337610d27c682e347c9cd60bd4b3b107c9d34ddd"}
	XRP  = Currency{"XRP", "0xa83575490d7df4e2f47b7d38ef351a2722ca45b9"}
)

func NewCurrency(symbol string) (Currency, error) {
	switch symbol {
	case BNB.symbol:
		return BNB, nil
	case BTC.symbol:
		return BTC, nil
	case BUSD.symbol:
		return BUSD, nil
	case DAI.symbol:
		return DAI, nil
	case ETH.symbol:
		return ETH, nil
	case USDC.symbol:
		return USDC, nil
	case USDT.symbol:
		return USDT, nil
	case XRP.symbol:
		return XRP, nil
	default:
		return Currency{}, errors.New("invalid currency symbol")
	}
}
