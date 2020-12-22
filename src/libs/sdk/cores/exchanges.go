package cores

import "cos-backend-com/src/common/flake"

type ExchangeStatus int

const (
	ExchangeStatusPending   ExchangeStatus = iota
	ExchangeStatusCompleted ExchangeStatus = iota
	ExchangeStatusUndone    ExchangeStatus = iota
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
