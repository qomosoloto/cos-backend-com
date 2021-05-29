package cores

import (
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/validate"
	"gopkg.in/go-playground/validator.v9"
)

type ProposalStatus int

func init() {
	validate.Default.RegisterValidation("special-proposal-states", func(fl validator.FieldLevel) bool {
		v := fl.Field().Int()
		return v == 4 || v == 5 || v == 6
	})
}

type CreateProposalInput struct {
	TxId                      string         `json:"txId" validate:"required"`
	StartupId                 string         `json:"startupId" validate:"required"`
	WalletAddr                string         `json:"walletAddr" validate:"required"`
	ContractAddr              string         `json:"contractAddr" validate:"required"`
	Status                    ProposalStatus `json:"status" validate:"required"`
	Title                     string         `json:"title" validate:"required"`
	Type                      int            `json:"type" validate:"required"`
	UserId                    flake.ID       `json:"userId"`
	Contact                   string         `json:"contact" validate:"required"`
	Description               string         `json:"description" validate:"required"`
	VoterType                 int            `json:"voterType" validate:"required"`
	SupportPercentage         int            `json:"supportPercentage" validate:"required"`
	MinimumApprovalPercentage int            `json:"minimumApprovalPercentage" validate:"required"`
	Duration                  int            `json:"duration" validate:"required"`
	HasPayment                int            `json:"hasPayment" validate:"required"`
	PaymentAddr               string         `json:"PaymentAddr"`
	PaymentType               int            `json:"paymentType"`
	PaymentMonths             int            `json:"paymentMonths"`
	PaymentDate               string         `json:"paymentDate"`
	PaymentAmount             float64        `json:"paymentAmount"`
	TotalPaymentAmount        float64        `json:"totalPaymentAmount"`
}

type CreateProposalResult struct {
	Id     flake.ID       `json:"id" db:"id"`
	Status ProposalStatus `json:"status" db:"status"`
}

type UpdateProposalStatusInput struct {
	Id     flake.ID       `json:"id"`
	Status ProposalStatus `json:"status" validate:"special-proposal-states"`
}

type UpdateProposalStatusResult struct {
	Id flake.ID `json:"id"`
}

type VoteProposalInput struct {
	Id         flake.ID `json:"id"`
	TxId       string   `json:"txId"`
	Amount     float32  `json:"amount" validate:"required"`
	IsApproved bool     `json:"isApproved" validate:"required"`
	WalletAddr string   `json:"walletAddr" validate:"required"`
	CreatedAt  string   `json:"createdAt" validate:"required"`
}

type VoteProposalResult struct {
}
