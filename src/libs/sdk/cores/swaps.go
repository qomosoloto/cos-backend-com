package cores

import "cos-backend-com/src/common/flake"

type TokenInput struct {
	Name     string `json:"name" validate:"required"`
	Symbol   string `json:"symbol" validate:"required"`
	Decimals int    `json:"decimals" validate:"required"`
	Address  string `json:"address" validate:"required"`
}

type CreateSwapPairInput struct {
	TxId        string     `json:"txId" validate:"required"`
	StartupId   flake.ID   `json:"startupId" validate:"required"`
	PairAddress string     `json:"pairAddress" validate:"required"`
	Token0      TokenInput `json:"token0" validate:"required"`
	Token1      TokenInput `json:"token1" validate:"required"`
}

type CreateSwapMintInput struct {
	TxId string `json:"txId" validate:"required"`

	StartupId   flake.ID   `json:"startupId" validate:"required"`
	PairAddress string     `json:"pairAddress" validate:"required"`
	Token0      TokenInput `json:"token0" validate:"required"`
	Token1      TokenInput `json:"token1" validate:"required"`
}
