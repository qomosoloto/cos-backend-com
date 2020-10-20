package cores

import (
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/pagination"
	"cos-backend-com/src/libs/sdk/eth"
	"time"

	"github.com/jmoiron/sqlx/types"
)

type BountyStatus int

const (
	BountyStatusOpen       BountyStatus = iota
	BountyStatusInProgress BountyStatus = iota
	BountyStatusClosed     BountyStatus = iota
)

type BountyType string

const (
	BountyTypeContest     BountyType = "contest"
	BountyTypeCooperative BountyType = "cooperative"
)

func (b BountyType) Validate() bool {
	switch b {
	case BountyTypeContest, BountyTypeCooperative:
		return true
	}
	return false
}

type BountyHunterRelStatus int

const (
	UndertakeBountyStatusNull      BountyHunterRelStatus = iota
	UndertakeBountyStatusStartWork BountyHunterRelStatus = iota
	UndertakeBountyStatusSubmitted BountyHunterRelStatus = iota
	UndertakeBountyStatusPaid      BountyHunterRelStatus = iota
	UndertakeBountyStatusQuited    BountyHunterRelStatus = iota
	UndertakeBountyStatusRejected  BountyHunterRelStatus = iota
)

type CreateBountyInput struct {
	Id                  flake.ID   `json:"id" validate:"required"`
	Title               string     `json:"title" validate:"required"`
	TxId                string     `json:"txId" validate:"required"`
	Type                BountyType `json:"type" validate:"func=self.Validate"`
	Keywords            []string   `json:"keywords"`
	ContactEmail        string     `json:"contactEmail" validate:"required"`
	Intro               string     `json:"intro" validate:"required"`
	DescriptionAddr     string     `json:"descriptionAddr"`
	DescriptionFileAddr string     `json:"descriptionFileAddr"`
	Duration            int        `json:"duration" validate:"required"`
	Payments            []struct {
		Token string  `json:"token"`
		Value float64 `json:"value"`
	} `json:"payments" validate:"required"`
}

type BountyOutput struct {
	Id       flake.ID `json:"id" db:"id"`
	SerialNo int      `json:"serialNo" db:"serial_no"`
	Startup  struct {
		Id   flake.ID `json:"id" db:"id"`
		Name string   `json:"name" db:"name"`
		Logo string   `json:"logo" db:"logo"`
	} `json:"startup" db:"startup"`
	CreatedBy struct {
		Id       flake.ID `json:"id" db:"id"`
		Name     string   `json:"name" db:"name"`
		IsHunter bool     `json:"isHunter" db:"is_hunter"`
	} `json:"createdBy" db:"created_by"`
	Title               string         `json:"title" db:"title"`
	Type                string         `json:"type" db:"type"`
	Keywords            []string       `json:"keywords" db:"keywords"`
	Intro               string         `json:"intro" db:"intro"`
	ContactEmail        string         `json:"contactEmail" db:"contact_email"`
	DescriptionAddr     string         `json:"descriptionAddr" db:"description_addr"`
	DescriptionFileAddr string         `json:"descriptionFileAddr" db:"description_file_addr"`
	Duration            int            `json:"duration" db:"duration"`
	Payments            types.JSONText `json:"payments" db:"payments"`
	Hunters             []struct {
		UserId           flake.ID              `json:"userId" db:"user_id"`
		Name             string                `json:"name" db:"name"`                // name
		Status           BountyHunterRelStatus `json:"status" db:"status"`            // status
		StartedAt        *time.Time            `json:"startedAt" db:"started_at"`     // started_at
		SubmittedAt      *time.Time            `json:"submittedAt" db:"submitted_at"` // submitted_at
		QuitedAt         *time.Time            `json:"quitedAt" db:"quited_at"`       // quited_at
		PaidAt           *time.Time            `json:"paidAt" db:"paid_at"`           // paid_at
		RejectedAt       *time.Time            `json:"rejectedAt" db:"rejected_at"`   // rejected_at
		PaidTokens       types.JSONText        `json:"paidTokens" db:"paid_tokens"`   // paid_tokens
		TransactionState eth.TransactionState  `json:"transactionState" db:"transaction_state"`
	} `json:"hunters" db:"hunters"`
	Status           BountyStatus         `json:"status" db:"status"`
	CreatedAt        time.Time            `json:"createdAt" db:"created_at"`
	ExpiredAt        time.Time            `json:"expiredAt" db:"expired_at"`
	BlockAddr        string               `json:"blockAddr" db:"block_addr"`
	TransactionState eth.TransactionState `json:"transactionState" db:"transaction_state"`
}

type ListBountiesInput struct {
	Keyword string `param:"keyword"`
	pagination.ListRequest
}

type ListBountiesResult struct {
	pagination.ListResult
	Result []BountyOutput `json:"result"`
}

type CreateUndertakeBountyInput struct {
	BountyId flake.ID `json:"bountyId" validate:"required"`
	UserId   flake.ID `json:"userId" validate:"required"`
	TxId     string   `json:"txId" validate:"required"`
}

type UpdateUndertakeBountyInput struct {
	BountyId flake.ID              `json:"bountyId" validate:"required"`
	UserId   flake.ID              `json:"userId" validate:"required"`
	Status   BountyHunterRelStatus `json:"status"`
}

type UndertakeBountyResult struct {
	Id       flake.ID              `json:"id" db:"id"`
	BountyId flake.ID              `json:"bountyId" db:"bounty_id"`
	Status   BountyHunterRelStatus `json:"status" db:"status"`
}
