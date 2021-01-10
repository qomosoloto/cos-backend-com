package discomodels

import (
	"context"
	"cos-backend-com/src/common/dbconn"
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/util"
	"cos-backend-com/src/libs/models"
	"cos-backend-com/src/libs/models/ethmodels"
	"cos-backend-com/src/libs/sdk/cores"
	ethSdk "cos-backend-com/src/libs/sdk/eth"

	"github.com/jmoiron/sqlx"
)

var DiscoInvestors = &discoInvestors{
	Connector: models.DefaultConnector,
}

type discoInvestors struct {
	dbconn.Connector
}

func (c *discoInvestors) CreateDiscoInvestor(ctx context.Context, startupId, uid flake.ID, input *cores.CreateDiscoInvestorInput) (err error) {
	stmt := `
		WITH get_disco_cte AS (
		    SELECT d.id FROM discos d
		    WHERE d.startup_id = ${startupId}
		)
		INSERT INTO disco_Investors(id, disco_id, uid, eth_count)
		SELECT ${id}, gdc.id, ${uid}, ${ethCount} FROM get_disco_cte gdc;
	`
	query, args := util.PgMapQueryV2(stmt, map[string]interface{}{
		"{id}":                   input.Id,
		"{startupId}":            startupId,
		"{uid}":                  uid,
		"{ethCount}":             input.EthCount,
		"{discoStateInprogress}": cores.DiscoStateInProgress,
	})
	return c.Invoke(ctx, func(db *sqlx.Tx) error {
		newCtx := dbconn.WithDB(ctx, db)
		createTransactionsInput := ethSdk.CreateTransactionsInput{
			TxId:     input.TxId,
			Source:   ethSdk.TransactionSourceDiscoInvestor,
			SourceId: input.Id,
		}

		if err := ethmodels.Transactions.Create(newCtx, &createTransactionsInput); err != nil {
			return err
		}
		_, err = db.ExecContext(newCtx, query, args...)
		return err
	})
}

func (c *discoInvestors) ListDiscoInvestor(ctx context.Context, startupId flake.ID, input *cores.ListDiscoInvestorsInput, outputs interface{}) (totalEth int64, total int, err error) {
	countStmt := `
		SELECT COUNT(*)
		FROM disco_investors di
		INNER JOIN users u ON di.uid = u.id
		INNER JOIN discos d ON d.id=di.disco_id
		WHERE d.startup_id=${startupId}
	`

	query, args := util.PgMapQueryV2(countStmt, map[string]interface{}{
		"{startupId}": startupId,
	})

	if err = c.Invoke(ctx, func(db dbconn.Q) error {
		return db.GetContext(ctx, &total, query, args...)
	}); err != nil {
		return 0, 0, err
	}

	if total > 0 {
		totalEthStmt := `
			SELECT sum(di.eth_count)
			FROM disco_investors di
			INNER JOIN users u ON di.uid = u.id
			INNER JOIN discos d ON d.id=di.disco_id
			WHERE d.startup_id=${startupId}
		`
		query, args = util.PgMapQueryV2(totalEthStmt, map[string]interface{}{
			"{startupId}": startupId,
		})
		if err = c.Invoke(ctx, func(db dbconn.Q) error {
			return db.GetContext(ctx, &totalEth, query, args...)
		}); err != nil {
			return 0, 0, err
		}

		stmt := `
			WITH res AS (
				SELECT di.*
				FROM disco_investors di
				INNER JOIN users u ON di.uid = u.id
				INNER JOIN discos d ON d.id = di.disco_id
				WHERE d.startup_id = ${startupId}
				LIMIT ${limit} OFFSET ${offset}
			)SELECT JSON_AGG(r.*) FROM res r;
		`
		query, args = util.PgMapQueryV2(stmt, map[string]interface{}{
			"{startupId}": startupId,
			"{limit}":     input.GetLimit(),
			"{offset}":    input.Offset,
		})
		if err = c.Invoke(ctx, func(db dbconn.Q) error {
			return db.GetContext(ctx, &util.PgJsonScanWrap{outputs}, query, args...)
		}); err != nil {
			return 0, 0, err
		}
	}
	return
}
