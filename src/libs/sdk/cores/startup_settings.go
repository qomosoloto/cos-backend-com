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
	AssignedProposers []string `json:"assigned_proposers"`

	VoterType                   int   `json:"voterType"`
	VoterTokenLimit             int64    `json:"voterTokenLimit"`
	AssignedVoters				[]string `json:"assigned_voters"`
	//VoteAssignAddrs            []string `json:"voteAssignAddrs"`
	//VoteSupportPercent         int      `json:"voteSupportPercent"`
	//VoteMinApprovalPercent     int      `json:"voteMinApprovalPercent"`
	//VoteMinDurationHours       int64    `json:"voteMinDurationHours"`
	//VoteMaxDurationHours       int64    `json:"voteMaxDurationHours"`
	ProposerType               int      `json:"proposerType"`
	ProposerTokenLimit         int64    `json:"proposerTokenLimit"`
	ProposalSupporters           int      `json:"proposalSupporters"`
	ProposalMinApprovalPercent int      `json:"proposalMinApprovalPercent"`
	ProposalMinDuration        int      `json:"proposalMinDuration"`
	ProposalMaxDuration        int      `json:"proposalMaxDuration"`
}

type StartupSettingRevisionsResult struct {
	TokenName                  string         `json:"tokenName" db:"token_name"`     // token_name
	TokenSymbol                string         `json:"tokenSymbol" db:"token_symbol"` // token_symbol
	TokenAddr                  *string        `json:"tokenAddr" db:"token_addr"`     // token_addr
	WalletAddrs                types.JSONText `json:"walletAddrs" db:"wallet_addrs"`
	//Type                       string         `json:"type" db:"type"`
	AssignedProposers			[]string		`json:"assignedProposers" db:"assigned_proposers"`
	AssignedVoters				[]string		`json:"assignedVoters" db:"assigned_voters"`
	VoterType					int				`json:"voterType" db:"voter_type"`
	VoterTokenLimit             *flake.ID      `json:"voterTokenLimit" db:"voter_token_limit"`                  // vote_token_limit
	//VoteAssignAddrs            []string       `json:"voteAssignAddrs" db:"vote_assign_addrs"`                // vote_assign_addrs
	//VoteSupportPercent         int            `json:"voteSupportPercent" db:"vote_support_percent"`          // vote_support_percent
	//VoteMinApprovalPercent     int            `json:"voteMinApprovalPercent" db:"vote_min_approval_percent"` // vote_min_approval_percent
	//VoteMinDurationHours       flake.ID       `json:"voteMinDurationHours" db:"vote_min_duration_hours"`     // vote_min_duration_hours
	//VoteMaxDurationHours       flake.ID       `json:"voteMaxDurationHours" db:"vote_max_duration_hours"`     // vote_max_duration_hours
	ProposerType               int            `json:"proposerType" db:"proposer_type"`
	ProposerTokenLimit         flake.ID       `json:"proposerTokenLimit" db:"proposer_token_limit"`
	ProposalSupporters           int            `json:"proposalSupporters" db:"proposal_supporters"`
	ProposalMinApprovalPercent int            `json:"proposalMinApprovalPercent" db:"proposal_min_approval_percent"`
	ProposalMinDuration        int            `json:"proposalMinDuration" db:"proposal_min_duration"`
	ProposalMaxDuration        int            `json:"proposalMaxDuration" db:"proposal_max_duration"`
}
