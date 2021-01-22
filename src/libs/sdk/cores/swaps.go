package cores

import "cos-backend-com/src/common/flake"

type TokenInput struct {
	Name     string `json:"name" validate:"required"`
	Symbol   string `json:"symbol" validate:"required"`
	Decimals int    `json:"decimals" validate:"required"`
	Address  string `json:"address" validate:"required"`
}

type SwapPairCreatedInput struct {
	TxId        string     `json:"txId" validate:"required"`
	StartupId   flake.ID   `json:"startupId" validate:"required"`
	PairAddress string     `json:"pairAddress" validate:"required"`
	Token0      TokenInput `json:"token0" validate:"required"`
	Token1      TokenInput `json:"token1" validate:"required"`
}

type SwapMintInput struct {
	TxId      string   `json:"txId" validate:"required"`
	StartupId flake.ID `json:"startupId" validate:"required"`
	Sender    string   `json:"sender" validate:"required"`
	Amount0   string   `json:"amount0" validate:"required"`
	Amount1   string   `json:"amount1" validate:"required"`
	OccuredAt string   `json:"occuredAt" validate:"required"`
}

type SwapBurnInput struct {
	TxId      string   `json:"txId" validate:"required"`
	StartupId flake.ID `json:"startupId" validate:"required"`
	Sender    string   `json:"sender" validate:"required"`
	Amount0   string   `json:"amount0" validate:"required"`
	Amount1   string   `json:"amount1" validate:"required"`
	To        string   `json:"to" validate:"required"`
}

type SwapSwapInput struct {
	TxId       string   `json:"txId" validate:"required"`
	StartupId  flake.ID `json:"startupId" validate:"required"`
	Sender     string   `json:"sender" validate:"required"`
	Amount0In  string   `json:"amount0In" validate:"required"`
	Amount1In  string   `json:"amount1In" validate:"required"`
	Amount0Out string   `json:"amount0Out" validate:"required"`
	Amount1Out string   `json:"amount1Out" validate:"required"`
	To         string   `json:"to" validate:"required"`
}

type SwapSyncInput struct {
	StartupId flake.ID `json:"startupId" validate:"required"`
	Reserve0  string   `json:"reserve0" validate:"required"`
	Reserve1  string   `json:"reserve1" validate:"required"`
	OccuredAt string   `json:"occuredAt" validate:"required"`
}
