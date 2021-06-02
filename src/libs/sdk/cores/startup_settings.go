package cores

import (
	"cos-backend-com/src/common/flake"

	"github.com/jmoiron/sqlx/types"
)

type StartupSettingResult struct {
	Id flake.ID `json:"id" db:"id"` // id
}

type UpdateStartupSettingInput struct {
	TxId        string `json:"txId" validate:"required"`
	TokenName   string `json:"tokenName" validate:"required"`
	TokenSymbol string `json:"tokenSymbol" validate:"required"`
	TokenAddr   string `json:"tokenAddr" validate:"required"`
	WalletAddrs []struct {
		Name string `json:"name" validate:"required"`
		Addr string `json:"addr" validate:"required"`
	} `json:"walletAddrs"`
	AssignedProposers []string `json:"assignedProposers"`

	VoterType                  int      `json:"voterType" validate:"required"`
	VoterTokenLimit            int64    `json:"voterTokenLimit"`
	AssignedVoters             []string `json:"assignedVoters"`
	ProposerType               int      `json:"proposerType" validate:"required"`
	ProposerTokenLimit         int64    `json:"proposerTokenLimit"`
	ProposalSupporters         int      `json:"proposalSupporters" validate:"required"`
	ProposalMinApprovalPercent int      `json:"proposalMinApprovalPercent" validate:"required"`
	ProposalMinDuration        int      `json:"proposalMinDuration" validate:"required"`
	ProposalMaxDuration        int      `json:"proposalMaxDuration" validate:"required"`
}

type StartupSettingRevisionsResult struct {
	TokenName                  string         `json:"tokenName" db:"token_name"`     // token_name
	TokenSymbol                string         `json:"tokenSymbol" db:"token_symbol"` // token_symbol
	TokenAddr                  *string        `json:"tokenAddr" db:"token_addr"`     // token_addr
	WalletAddrs                types.JSONText `json:"walletAddrs" db:"wallet_addrs"`
	AssignedProposers          []string       `json:"assignedProposers" db:"assigned_proposers"`
	AssignedVoters             []string       `json:"assignedVoters" db:"assigned_voters"`
	VoterType                  int            `json:"voterType" db:"voter_type"`
	VoterTokenLimit            *flake.ID      `json:"voterTokenLimit" db:"voter_token_limit"` // vote_token_limit
	ProposerType               int            `json:"proposerType" db:"proposer_type"`
	ProposerTokenLimit         flake.ID       `json:"proposerTokenLimit" db:"proposer_token_limit"`
	ProposalSupporters         int            `json:"proposalSupporters" db:"proposal_supporters"`
	ProposalMinApprovalPercent int            `json:"proposalMinApprovalPercent" db:"proposal_min_approval_percent"`
	ProposalMinDuration        int            `json:"proposalMinDuration" db:"proposal_min_duration"`
	ProposalMaxDuration        int            `json:"proposalMaxDuration" db:"proposal_max_duration"`
}
