package cores

import (
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/pagination"
	"time"
)

// DiscoInvestorsModel represents a row from 'disco_investors'.
type DiscoInvestorsModel struct {
	Id        flake.ID  `json:"id" db:"id"`                // id (PK)
	DiscoId   flake.ID  `json:"discoId" db:"disco_id"`     // disco_id
	UId       flake.ID  `json:"uid" db:"uid"`              // uid
	EthCount  int64     `json:"ethCount" db:"eth_count"`   // eth_count
	CreatedAt time.Time `json:"createdAt" db:"created_at"` // created_at
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"` // updated_at
}

// DiscoInvestorsInput represents an input for 'disco_investors'.
type CreateDiscoInvestorInput struct {
	Id       flake.ID `json:"id" validate:"required"`
	EthCount int64    `json:"ethCount" validate:"required"`
	TxId     string   `json:"txId" validate:"required"`
}

// DiscoInvestorsResult represents an output for 'disco_investors'.
type DiscoInvestorOutput struct {
	Id         flake.ID  `json:"id" db:"id"`                  // id (PK)
	DiscoId    flake.ID  `json:"discoId" db:"disco_id"`       // disco_id
	UId        flake.ID  `json:"uid" db:"uid"`                // uid
	WalletAddr string    `json:"walletAddr" db:"wallet_addr"` // wallet_addr
	EthCount   int64     `json:"ethCount" db:"eth_count"`     // eth_count
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`   // created_at
}

// ListDiscoInvestorsInput represents query params for 'disco_investors'.
type ListDiscoInvestorsInput struct {
	pagination.ListRequest
}

type ListDiscoInvestorsResult struct {
	pagination.ListResult
	Result   []DiscoInvestorOutput `json:"result"`
	TotalEth int64                 `json:"totalEth"`
}
