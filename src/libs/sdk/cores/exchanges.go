package cores

import (
	"cos-backend-com/src/common/flake"
)

type ExchangeStatus int

const (
	ExchangeStatusPending   ExchangeStatus = iota
	ExchangeStatusCompleted ExchangeStatus = iota
	ExchangeStatusUndone    ExchangeStatus = iota
)

type ExchangeTxType int

const (
	ExchangeTxTypeAddLiquidity    ExchangeTxType = 1
	ExchangeTxTypeRemoveLiquidity ExchangeTxType = 2
	ExchangeTxTypeSwap            ExchangeTxType = 3
)

type ExchangeTxStatus int

const (
	ExchangeTxStatusPending   ExchangeTxStatus = iota
	ExchangeTxStatusCompleted ExchangeTxStatus = iota
	ExchangeTxStatusUndone    ExchangeTxStatus = iota
)

type CreateExchangeInput struct {
	TxId          string         `json:"txId" validate:"required"`
	StartupId     flake.ID       `json:"startupId" validate:"required"`
	PairName      string         `json:"pairName"`
	PairAddress   string         `json:"pairAddress"`
	TokenName1    string         `json:"tokenName1" validate:"required"`
	TokenSymbol1  string         `json:"tokenSymbol1" validate:"required"`
	TokenAddress1 string         `json:"tokenAddress1"`
	TokenName2    string         `json:"tokenName2" validate:"required"`
	TokenSymbol2  string         `json:"tokenSymbol2" validate:"required"`
	TokenAddress2 string         `json:"tokenAddress2"`
	Status        ExchangeStatus `json:"status"`
}

type CreateExchangeResult struct {
	Id     flake.ID       `json:"id" db:"id"`
	Status ExchangeStatus `json:"status" db:"status"`
}

type GetExchangeInput struct {
	Id        flake.ID `json:"id"`
	StartupId flake.ID `json:"startupId"`
}

type ExchangeResult struct {
	Id      flake.ID `json:"id" db:"id"`
	TxId    string   `json:"txId" db:"tx_id"`
	Startup struct {
		Id          flake.ID `json:"id" db:"id"`
		Name        string   `json:"name" db:"name"`
		Logo        string   `json:"logo" db:"logo"`
		TokenName   string   `json:"tokenName" db:"token_name"`
		TokenSymbol string   `json:"tokenSymbol" db:"token_symbol"`
		Mission     string   `json:"mission" db:"mission"`
	} `json:"startup" db:"startup"`
	PairName    string         `json:"pairName" db:"pair_name"`
	PairAddress string         `json:"pairAddress" db:"pair_address"`
	Status      ExchangeStatus `json:"status" db:"status"`
}

type CreateExchangeTxInput struct {
	TxId         string           `json:"txId" validate:"required"`
	ExchangeId   flake.ID         `json:"exchangeId" validate:"required"`
	Account      string           `json:"account" validate:"required"`
	Type         ExchangeTxType   `json:"type" validate:"required"`
	TokenAmount1 float32          `json:"tokenAmount1"`
	TokenAmount2 float32          `json:"tokenAmount2"`
	Status       ExchangeTxStatus `json:"status"`
}

type CreateExchangeTxResult struct {
	Id     flake.ID         `json:"id" db:"id"`
	Status ExchangeTxStatus `json:"status" db:"status"`
}
