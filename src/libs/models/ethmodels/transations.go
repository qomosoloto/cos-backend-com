package ethmodels

import (
	"context"
	"cos-backend-com/src/common/dbconn"
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/util"
	"cos-backend-com/src/libs/models"
	"cos-backend-com/src/libs/sdk/cores"
	ethSdk "cos-backend-com/src/libs/sdk/eth"

	"github.com/jmoiron/sqlx"
)

var Transactions = &transactions{
	Connector: models.DefaultConnector,
}

type transactions struct {
	dbconn.Connector
}

func (c *transactions) List(ctx context.Context, outputs interface{}) (err error) {
	stmt := `
		WITH res AS(
			SELECT *
			FROM transactions
			WHERE state = 1
		)
		SELECT COALESCE(json_agg(r.*), '[]'::json) FROM res r;
	`

	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{state}": ethSdk.TransactionStateWaitConfirm,
	})

	return c.Invoke(ctx, func(db *sqlx.DB) (er error) {
		return db.GetContext(ctx, &util.PgJsonScanWrap{Value: outputs}, query, args...)
	})
}

func (c *transactions) Create(ctx context.Context, input *ethSdk.CreateTransactionsInput) (err error) {
	stmt := `
		INSERT INTO transactions(tx_id, source, source_id)
		VALUES (${txId}, ${source}, ${sourceId})
		ON CONFLICT(tx_id) WHERE state !=3
		DO UPDATE SET
		(
		    state, retry_time, source_id, updated_at
		)= (
		    1, 0, ${sourceId}, current_timestamp
		);
	`

	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{source}":   input.Source,
		"{sourceId}": input.SourceId,
		"{txId}":     input.TxId,
	})

	return c.Invoke(ctx, func(db *sqlx.DB) (er error) {
		_, err = db.ExecContext(ctx, query, args...)
		return err
	})
}

func (c *transactions) Update(ctx context.Context, id flake.ID, input *ethSdk.UpdateTransactionsInput) (err error) {
	stmt := `
		UPDATE transactions SET (
			block_addr, state, retry_time
		) = (
			${blockAddr}, ${state}, retry_time+1
		) WHERE id = ${id}
	`

	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{id}":        id,
		"{blockAddr}": input.BlockAddr,
		"{state}":     input.State,
	})

	return c.Invoke(ctx, func(db *sqlx.DB) (er error) {
		_, err = db.ExecContext(ctx, query, args...)
		return err
	})
}

func (c *transactions) UpdateWithConfirmSource(ctx context.Context, id, sourceId flake.ID, source ethSdk.TransactionSource, input *ethSdk.UpdateTransactionsInput) (err error) {
	return c.Invoke(ctx, func(db *sqlx.Tx) (er error) {
		newCtx := dbconn.WithDB(ctx, db)
		if er = c.Update(newCtx, id, input); er != nil {
			return er
		}
		switch source {
		case ethSdk.TransactionSourceStartup:
			if er = c.ConfirmStartup(newCtx, sourceId); er != nil {
				return er
			}
		case ethSdk.TransactionSourceStartupSetting:
			if er = c.ConfirmStartupSetting(newCtx, sourceId); er != nil {
				return er
			}
		case ethSdk.TransactionSourceUndertakeBounty:
			if er = c.ConfirmUndertakeBounty(newCtx, sourceId); er != nil {
				return er
			}
		}
		return
	})
}

func (c *transactions) ConfirmStartup(ctx context.Context, id flake.ID) (err error) {
	stmt := `
		UPDATE startups s
		SET (current_revision_id,updated_at)= (confirming_revision_id,current_timestamp)
		FROM startup_revisions sr
		WHERE sr.id = ${id}
		AND sr.startup_id = s.id;
	`

	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{id}": id,
	})

	return c.Invoke(ctx, func(db dbconn.Q) (er error) {
		_, er = db.ExecContext(ctx, query, args...)
		return er
	})
}

func (c *transactions) ConfirmStartupSetting(ctx context.Context, id flake.ID) (err error) {
	stmt := `
		UPDATE startup_settings s
		SET (current_revision_id,updated_at)= (confirming_revision_id,current_timestamp)
		FROM startup_setting_revisions sr
		WHERE sr.id = ${id}
		AND sr.startup_setting_id = s.id;
	`

	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{id}": id,
	})

	return c.Invoke(ctx, func(db dbconn.Q) (er error) {
		_, er = db.ExecContext(ctx, query, args...)
		return er
	})
}

func (c *transactions) ConfirmUndertakeBounty(ctx context.Context, id flake.ID) (err error) {
	stmt := `
	WITH update_bounties_hunter_rel_cte AS (
		UPDATE bounties_hunters_rel
			SET status = ${undertakeBountyStatusStartWork},
				started_at = current_timestamp
			WHERE id = ${sourceId}
			RETURNING bounty_id
	)
	UPDATE bounties
	SET status = ${bountyStatusInProgress}
	WHERE id = (SELECT bounty_id FROM update_bounties_hunter_rel_cte)
		AND status = 0;
	`

	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{undertakeBountyStatusStartWork}": cores.UndertakeBountyStatusStartWork,
		"{bountyStatusInProgress}":         cores.BountyStatusInProgress,
		"{sourceId}":                       id,
	})

	return c.Invoke(ctx, func(db dbconn.Q) (er error) {
		_, er = db.ExecContext(ctx, query, args...)
		return er
	})
}
