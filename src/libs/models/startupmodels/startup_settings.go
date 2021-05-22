package startupmodels

import (
	"context"
	"cos-backend-com/src/common/dbconn"
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/types"
	"cos-backend-com/src/common/util"
	"cos-backend-com/src/libs/models"
	"cos-backend-com/src/libs/models/ethmodels"
	coresSdk "cos-backend-com/src/libs/sdk/cores"
	ethSdk "cos-backend-com/src/libs/sdk/eth"

	"github.com/jmoiron/sqlx"
)

var StartupSettings = &startupSettings{
	Connector: models.DefaultConnector,
}

// startUps represents controller for 'startups'.
type startupSettings struct {
	dbconn.Connector
}

func (c *startupSettings) UpsertWithRevision(ctx context.Context, startupId flake.ID, input *coresSdk.UpdateStartupSettingInput, startupSettingId *flake.ID) (err error) {
	stmt := `
		UPDATE startup_settings SET
		(
		    confirming_revision_id, updated_at
		) = (
		    ${confirmingRevisionId}, CURRENT_TIMESTAMP
		) WHERE id = ${id};
	`
	return c.Invoke(ctx, func(db *sqlx.Tx) error {
		newCtx := dbconn.WithDB(ctx, db)
		if er := c.Upsert(newCtx, startupId, startupSettingId); er != nil {
			return er
		}

		var startupSettingsRevisionId flake.ID
		if er := c.CreateRevision(newCtx, *startupSettingId, input, &startupSettingsRevisionId); er != nil {
			return er
		}

		query, args := util.PgMapQueryV2(stmt, map[string]interface{}{
			"{id}":                   *startupSettingId,
			"{confirmingRevisionId}": startupSettingsRevisionId,
		})

		return c.Invoke(newCtx, func(db dbconn.Q) (er error) {
			_, er = db.ExecContext(newCtx, query, args...)
			return er
		})
	})
}

func (c *startupSettings) Upsert(ctx context.Context, startupId flake.ID, output interface{}) (err error) {
	stmt := `
		INSERT INTO startup_settings(startup_id)
		VALUES (${startupId})
		ON CONFLICT(startup_id) DO UPDATE SET updated_at = current_timestamp RETURNING id;
	`

	query, args := util.PgMapQueryV2(stmt, map[string]interface{}{
		"{startupId}": startupId,
	})

	return c.Invoke(ctx, func(db *sqlx.DB) (er error) {
		return db.GetContext(ctx, output, query, args...)
	})
}

func (c *startupSettings) CreateRevision(ctx context.Context, startupSettingId flake.ID, input *coresSdk.UpdateStartupSettingInput, revisionId *flake.ID) (err error) {
	//stmt := `
	//	INSERT INTO startup_setting_revisions(startup_setting_id, token_name, token_symbol, token_addr, wallet_addrs, type, vote_token_limit, vote_assign_addrs, vote_support_percent, vote_min_approval_percent, vote_min_duration_hours, vote_max_duration_hours, proposer_type, proposer_token_limit, proposal_supports, proposal_min_approval_percent, proposal_min_duration, proposal_max_duration)
	//	VALUES (${startupSettingId},${tokenName},${tokenSymbol},${tokenAddr},${walletAddrs},${type},${voteTokenLimit},ARRAY[${voteAssignAddrs}],${voteSupportPercent},${voteMinApprovalPercent},${voteMinDurationHours},${voteMaxDurationHours},${proposerType},${proposerTokenLimit},${proposalSupports},${proposalMinApprovalPercent},${proposalMinDuration},${proposalMaxDuration})
	//	RETURNING id;
	//`
	stmt := `
		INSERT INTO startup_setting_revisions(startup_setting_id, token_name, token_symbol, token_addr, wallet_addrs, type, vote_token_limit, vote_support_percent, vote_min_approval_percent, vote_min_duration_hours, vote_max_duration_hours, proposer_type, proposer_token_limit, proposal_supports, proposal_min_approval_percent, proposal_min_duration, proposal_max_duration)
		VALUES (${startupSettingId},${tokenName},${tokenSymbol},${tokenAddr},${walletAddrs},${type},${voteTokenLimit},${voteSupportPercent},${voteMinApprovalPercent},${voteMinDurationHours},${voteMaxDurationHours},${proposerType},${proposerTokenLimit},${proposalSupports},${proposalMinApprovalPercent},${proposalMinDuration},${proposalMaxDuration})
		RETURNING id;
	`

	query, args := util.PgMapQueryV2(stmt, map[string]interface{}{
		"{startupSettingId}":           startupSettingId,
		"{tokenName}":                  input.TokenName,
		"{tokenSymbol}":                input.TokenSymbol,
		"{tokenAddr}":                  input.TokenAddr,
		"{walletAddrs}":                types.JSONAny{input.WalletAddrs},
		"{type}":                       input.VoteType,
		"{voteTokenLimit}":             input.VoteTokenLimit,
		"{voteAssignAddrs}":            input.VoteAssignAddrs,
		"{voteSupportPercent}":         input.VoteSupportPercent,
		"{voteMinApprovalPercent}":     input.VoteMinApprovalPercent,
		"{voteMaxDurationHours}":       input.VoteMaxDurationHours,
		"{voteMinDurationHours}":       input.VoteMinDurationHours,
		"{proposerType}":               input.ProposerType,
		"{proposerTokenLimit}":         input.ProposerTokenLimit,
		"{proposalSupports}":           input.ProposalSupports,
		"{proposalMinApprovalPercent}": input.ProposalMinApprovalPercent,
		"{proposalMinDuration}":        input.ProposalMinDuration,
		"{proposalMaxDuration}":        input.ProposalMaxDuration,
	})
	return c.Invoke(ctx, func(db *sqlx.Tx) error {
		newCtx := dbconn.WithDB(ctx, db)
		if er := db.GetContext(newCtx, revisionId, query, args...); er != nil {
			return er
		}
		createTransactionsInput := ethSdk.CreateTransactionsInput{
			TxId:     input.TxId,
			Source:   ethSdk.TransactionSourceStartupSetting,
			SourceId: *revisionId,
		}

		return ethmodels.Transactions.Create(newCtx, &createTransactionsInput)
	})
}

func (c *startupSettings) Restore(ctx context.Context, uid, id flake.ID) (err error) {
	stmt := `
		WITH get_startup_id_cte AS(
		    SELECT id FROM startups WHERE id = ${id} AND uid = ${uid}
		)
		UPDATE startup_settings ss
		SET (confirming_revision_id,updated_at)= (current_revision_id,current_timestamp)
		FROM get_startup_id_cte gsic
		WHERE ss.startup_id = gsic.id;
	`

	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{id}":  id,
		"{uid}": uid,
	})

	return c.Invoke(ctx, func(db dbconn.Q) (er error) {
		_, er = db.ExecContext(ctx, query, args...)
		return er
	})
}
