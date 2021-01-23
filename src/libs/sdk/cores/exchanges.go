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
	ExchangeTxTypeSwap1for2       ExchangeTxType = 3
	ExchangeTxTypeSwap2for1       ExchangeTxType = 4
)

type ExchangeTxStatus int

const (
	ExchangeTxStatusPending   ExchangeTxStatus = iota
	ExchangeTxStatusCompleted ExchangeTxStatus = iota
	ExchangeTxStatusUndone    ExchangeTxStatus = iota
)

type ListExchangesOrderBy string

const (
	ListExchangesOrderByTime         ListExchangesOrderBy = "createdAt"
	ListExchangesOrderByName         ListExchangesOrderBy = "startup"
	ListExchangesOrderByLiquidities  ListExchangesOrderBy = "liquidities"
	ListExchangesOrderByVolumes24Hrs ListExchangesOrderBy = "volumes24Hrs"
)

type CreateExchangeInput struct {
	TxId          string         `json:"txId" validate:"required"`
	StartupId     flake.ID       `json:"startupId" validate:"required"`
	PairName      string         `json:"pairName"`
	PairAddress   string         `json:"pairAddress"`
	TokenName1    string         `json:"tokenName1" validate:"required"`
	TokenSymbol1  string         `json:"tokenSymbol1" validate:"required"`
	TokenAddress1 string         `json:"tokenAddress1"`
	TokenDivider1 int            `json:"tokenDivider1"`
	TokenName2    string         `json:"tokenName2" validate:"required"`
	TokenSymbol2  string         `json:"tokenSymbol2" validate:"required"`
	TokenAddress2 string         `json:"tokenAddress2"`
	TokenDivider2 int            `json:"tokenDivider2"`
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
	PairName      string         `json:"pairName" db:"pair_name"`
	PairAddress   string         `json:"pairAddress" db:"pair_address"`
	Status        ExchangeStatus `json:"status" db:"status"`
	FollowCount   int            `json:"followCount" db:"follow_count"`
	TokenDivider1 int            `json:"tokenDivider1" db:"token_divider1"`
	TokenDivider2 int            `json:"tokenDivider2" db:"token_divider2"`
	TokenSymbol1  string         `json:"tokenSymbol1" db:"token_symbol1"`
	TokenSymbol2  string         `json:"tokenSymbol2" db:"token_symbol2"`
}

type ListExchangesInput struct {
	pagination.ListRequest
	Keyword string                `param:"keyword"`
	OrderBy *ListExchangesOrderBy `param:"orderBy"`
	IsDesc  bool                  `param:"isDesc"`
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
		PriceChanges []struct {
			OccuredDay string  `json:"occuredDay" db:"occured_day"`
			EndPrice   float64 `json:"endPrice" db:"end_price"`
		} `json:"priceChanges" db:"price_changes"`
	} `json:"result"`
}

type CreateExchangeTxInput struct {
	TxId           string           `json:"txId" validate:"required"`
	StartupId      flake.ID         `json:"startupId"`
	ExchangeId     flake.ID         `json:"exchangeId" validate:"required"`
	Sender         string           `json:"sender" validate:"required"`
	To             string           `json:"to"`
	Type           ExchangeTxType   `json:"type" validate:"required"`
	Name           string           `json:"name"`
	TotalValue     float64          `json:"totalValue"`
	TokenAmount1   float64          `json:"tokenAmount1"`
	TokenAmount2   float64          `json:"tokenAmount2"`
	Amount0        string           `json:"amount0"`
	Amount1        string           `json:"amount1"`
	Fee            float64          `json:"fee"`
	PricePerToken1 float64          `json:"pricePerToken1"`
	PricePerToken2 float64          `json:"pricePerToken2"`
	Status         ExchangeTxStatus `json:"status"`
	OccuredAt      string           `json:"occuredAt"`
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
	Sender         string           `json:"sender" db:"sender"`
	To             string           `json:"to" db:"receiver"`
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

type ExchangeAllStatsTotalResult struct {
	Volumes24Hrs     float64 `json:"volumes24Hrs" db:"volumes_24hrs"`
	Volumes24HrsRate float64 `json:"volumes24HrsRate" db:"volumes_24hrs_rate"`
	Liquidities      float64 `json:"liquidities" db:"liquidities"`
	LiquiditiesRate  float64 `json:"liquiditiesRate" db:"liquidities_rate"`
}

type ExchangeOneStatsInput struct {
	Id flake.ID `json:"id"`
}

type ExchangeOneStatsTotalResult struct {
	Volumes24Hrs          float64 `json:"volumes24Hrs" db:"volumes_24hrs"`
	Volumes24HrsRate      float64 `json:"volumes24HrsRate" db:"volumes_24hrs_rate"`
	Liquidities           float64 `json:"liquidities" db:"liquidities"`
	LiquiditiesRate       float64 `json:"liquiditiesRate" db:"liquidities_rate"`
	Transactions24Hrs     int     `json:"transactions24Hrs" db:"transactions_24hrs"`
	Transactions24HrsRate float64 `json:"transactions24HrsRate" db:"transactions_24hrs_rate"`
}

type ExchangeOneStatsPriceChangeResult struct {
	TokenSymbol1    string  `json:"tokenSymbol1" db:"token_symbol1"`
	TokenSymbol2    string  `json:"tokenSymbol2" db:"token_symbol2"`
	PricePerToken1  float64 `json:"pricePerToken1" db:"price_per_token1"`
	PricePerToken2  float64 `json:"pricePerToken2" db:"price_per_token2"`
	PriceChangeRate float64 `json:"priceChangeRate" db:"price_change_rate"`
	PriceChanges    []struct {
		OccuredDay string  `json:"occuredDay" db:"occured_day"`
		EndPrice   float64 `json:"endPrice" db:"end_price"`
	} `json:"priceChanges" db:"price_changes"`
}

type ExchangeBalanceInput struct {
	StartupId           flake.ID `json:"startupId" validate:"required"`
	NewestDay           string   `json:"newestDay"`
	NewestPooledTokens1 float64  `json:"newestPooledTokens1"`
	NewestPooledTokens2 float64  `json:"newestPooledTokens2"`
	LastDay             string   `json:"lastDay"`
	LastPooledTokens1   float64  `json:"lastPooledTokens1"`
	LastPooledTokens2   float64  `json:"lastPooledTokens2"`
	Reserve0            string   `json:"reserve0" validate:"required"`
	Reserve1            string   `json:"reserve1" validate:"required"`
}

type ExchangeBalanceResult struct {
	TokenDivider1       int     `json:"tokenDivider1" db:"token_divider1"`
	TokenDivider2       int     `json:"tokenDivider2" db:"token_divider2"`
	NewestDay           string  `json:"newestDay" db:"newest_day"`
	NewestPooledTokens1 float64 `json:"newestPooledTokens1" db:"newest_pooled_tokens1"`
	NewestPooledTokens2 float64 `json:"newestPooledTokens2" db:"newest_pooled_tokens2"`
	LastDay             string  `json:"lastDay" db:"last_day"`
	LastPooledTokens1   float64 `json:"lastPooledTokens1" db:"last_pooled_tokens1"`
	LastPooledTokens2   float64 `json:"lastPooledTokens2" db:"last_pooled_tokens2"`
}
