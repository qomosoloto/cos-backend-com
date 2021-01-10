package eth

import (
	"cos-backend-com/src/common/flake"
)

type TransactionState int

const (
	TransactionStateNull        TransactionState = 0
	TransactionStateWaitConfirm TransactionState = 1
	TransactionStateSuccess     TransactionState = 2
	TransactionStateFailed      TransactionState = 3
)

func (t TransactionState) Validate() bool {
	switch t {
	case TransactionStateNull, TransactionStateWaitConfirm, TransactionStateSuccess, TransactionStateFailed:
		return true
	}
	return false
}

type TransactionSource string

const (
	TransactionSourceStartup         TransactionSource = "startup"
	TransactionSourceBounty          TransactionSource = "bounty"
	TransactionSourceStartupSetting  TransactionSource = "startupSetting"
	TransactionSourceUndertakeBounty TransactionSource = "undertakeBounty"
	TransactionSourceExchange        TransactionSource = "exchange"
	TransactionSourceExchangeTx      TransactionSource = "exchangeTransaction"
	TransactionSourceDisco           TransactionSource = "disco"
	TransactionSourceDiscoInvestor   TransactionSource = "discoInvestor"
)

func (t TransactionSource) Validate() bool {
	switch t {
	case TransactionSourceStartup, TransactionSourceStartupSetting:
		return true
	}
	return false
}

type CreateTransactionsInput struct {
	TxId     string            `json:"txId"`
	Source   TransactionSource `json:"source"`
	SourceId flake.ID          `json:"sourceId"`
}

type UpdateTransactionsInput struct {
	BlockAddr string           `json:"blockAddr"`
	State     TransactionState `json:"state"`
}

type TransactionsResult struct {
	TxId      string `json:"txId" db:"tx_id"`           // tx_id
	BlockAddr string `json:"blockAddr" db:"block_addr"` // block_addr
}
