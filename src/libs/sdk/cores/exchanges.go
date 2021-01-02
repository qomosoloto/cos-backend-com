package cores

import (
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/pagination"
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

type ListExchangesInput struct {
	Keyword string `param:"keyword"`
	pagination.ListRequest
}

type ListExchangesResult struct {
	pagination.ListResult
	Result []struct {
		Id           flake.ID           `json:"id" db:"id"`
		TxId         string             `json:"txId" db:"tx_id"`
		Startup      StartupShortResult `json:"startup" db:"startup"`
		Price        float64            `json:"price" db:"price"`
		Liquidities  float64            `json:"liquidities" db:"liquidities"`
		Volumes24Hrs float64            `json:"volumes24Hrs" db:"volumes_24hrs"`
		Status       ExchangeStatus     `json:"status" db:"status"`
	} `json:"result"`
}

type CreateExchangeTxInput struct {
	TxId         string           `json:"txId" validate:"required"`
	ExchangeId   flake.ID         `json:"exchangeId" validate:"required"`
	Account      string           `json:"account" validate:"required"`
	Type         ExchangeTxType   `json:"type" validate:"required"`
	TokenAmount1 float64          `json:"tokenAmount1"`
	TokenAmount2 float64          `json:"tokenAmount2"`
	Status       ExchangeTxStatus `json:"status"`
}

type CreateExchangeTxResult struct {
	Id     flake.ID         `json:"id" db:"id"`
	Status ExchangeTxStatus `json:"status" db:"status"`
}

type GetExchangeTxInput struct {
	Id   flake.ID `json:"id"`
	TxId string   `json:"txId"`
}

type ExchangeTxResult struct {
	Id             flake.ID         `json:"id" db:"id"`
	TxId           string           `json:"txId" db:"tx_id"`
	ExchangeId     flake.ID         `json:"exchangeId" db:"exchange_id"`
	Account        string           `json:"account" db:"account"`
	Type           ExchangeTxType   `json:"type" db:"type"`
	Name           string           `json:"name" db:"name"`
	TotalValue     float64          `json:"totalValue" db:"total_value"`
	TokenAmount1   float64          `json:"tokenAmount1" db:"token_amount1"`
	TokenAmount2   float64          `json:"tokenAmount2" db:"token_amount2"`
	Fee            float64          `json:"fee" db:"fee"`
	PricePerToken1 float64          `json:"pricePerToken1" db:"price_per_token1"`
	PricePerToken2 float64          `json:"pricePerToken2" db:"price_per_token2"`
	Status         ExchangeTxStatus `json:"status" db:"status"`
	OccuredAt      string           `json:"occuredAt" db:"occured_at"`
}
