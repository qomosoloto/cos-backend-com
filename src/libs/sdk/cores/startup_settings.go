package cores

import (
	"cos-backend-com/src/common/flake"

	"github.com/jmoiron/sqlx/types"
)

type StartupSettingResult struct {
	Id flake.ID `json:"id" db:"id"` // id
}

type UpdateStartupSettingInput struct {
	TxId        string `json:"txId"`
	TokenName   string `json:"tokenName"`
	TokenSymbol string `json:"tokenSymbol"`
	TokenAddr   string `json:"tokenAddr"`
	WalletAddrs []struct {
		Name string `json:"name"`
		Addr string `json:"addr"`
	} `json:"walletAddrs"`
	VoteType                   string   `json:"voteType"`
	VoteTokenLimit             int64    `json:"voteTokenLimit"`
	VoteAssignAddrs            []string `json:"voteAssignAddrs"`
	VoteSupportPercent         int      `json:"voteSupportPercent"`
	VoteMinApprovalPercent     int      `json:"voteMinApprovalPercent"`
	VoteMinDurationHours       int64    `json:"voteMinDurationHours"`
	VoteMaxDurationHours       int64    `json:"voteMaxDurationHours"`
	ProposerType               int      `json:"proposerType"`
	ProposerTokenLimit         int64    `json:"proposerTokenLimit"`
	ProposalSupportPercent     int      `json:"proposalSupportPercent"`
	ProposalMinApprovalPercent int      `json:"proposalMinApprovalPercent"`
	ProposalMinDuration        int      `json:"proposalMinDuration"`
	ProposalMaxDuration        int      `json:"proposalMaxDuration"`
}

type StartupSettingRevisionsResult struct {
	TokenName                  string         `json:"tokenName" db:"token_name"`     // token_name
	TokenSymbol                string         `json:"tokenSymbol" db:"token_symbol"` // token_symbol
	TokenAddr                  *string        `json:"tokenAddr" db:"token_addr"`     // token_addr
	WalletAddrs                types.JSONText `json:"walletAddrs" db:"wallet_addrs"`
	Type                       string         `json:"type" db:"type"`                                        // type
	VoteTokenLimit             *flake.ID      `json:"voteTokenLimit" db:"vote_token_limit"`                  // vote_token_limit
	VoteAssignAddrs            []string       `json:"voteAssignAddrs" db:"vote_assign_addrs"`                // vote_assign_addrs
	VoteSupportPercent         int            `json:"voteSupportPercent" db:"vote_support_percent"`          // vote_support_percent
	VoteMinApprovalPercent     int            `json:"voteMinApprovalPercent" db:"vote_min_approval_percent"` // vote_min_approval_percent
	VoteMinDurationHours       flake.ID       `json:"voteMinDurationHours" db:"vote_min_duration_hours"`     // vote_min_duration_hours
	VoteMaxDurationHours       flake.ID       `json:"voteMaxDurationHours" db:"vote_max_duration_hours"`     // vote_max_duration_hours
	ProposerType               int            `json:"proposerType" db:"proposer_type"`
	ProposerTokenLimit         flake.ID       `json:"proposerTokenLimit" db:"proposer_token_limit"`
	ProposalSupportPercent     int            `json:"proposalSupportPercent" db:"proposal_support_percent"`
	ProposalMinApprovalPercent int            `json:"proposalMinApprovalPercent" db:"proposal_min_approval_percent"`
	ProposalMinDuration        int            `json:"proposalMinDuration" db:"proposal_min_duration"`
	ProposalMaxDuration        int            `json:"proposalMaxDuration" db:"proposal_max_duration"`
}
