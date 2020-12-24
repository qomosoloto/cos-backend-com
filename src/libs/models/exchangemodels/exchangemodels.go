package exchangemodels

import (
	"context"
	"cos-backend-com/src/common/dbconn"
	"cos-backend-com/src/common/util"
	"cos-backend-com/src/libs/models"
	"cos-backend-com/src/libs/models/ethmodels"
	coresSdk "cos-backend-com/src/libs/sdk/cores"
	ethSdk "cos-backend-com/src/libs/sdk/eth"
	"github.com/jmoiron/sqlx"
)

var Exchanges = &exchanges{
	Connector: models.DefaultConnector,
}

type exchanges struct {
	dbconn.Connector
}

func (c *exchanges) CreateExchange(ctx context.Context, input *coresSdk.CreateExchangeInput, output *coresSdk.CreateExchangeResult) (err error) {
	stmt := `
		INSERT INTO exchanges(tx_id, startup_id, token_name1, token_symbol1, token_address1, token_name2, token_symbol2, status)
		VALUES (${txId}, ${startupId}, ${tokenName1}, ${tokenSymbol1}, ${tokenAddress1}, ${tokenName2}, ${tokenSymbol2}, ${status})
		RETURNING id, status;
	`
	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{txId}":          input.TxId,
		"{startupId}":     input.StartupId,
		"{tokenName1}":    input.TokenName1,
		"{tokenSymbol1}":  input.TokenSymbol1,
		"{tokenAddress1}": input.TokenAddress1,
		"{tokenName2}":    input.TokenName2,
		"{tokenSymbol2}":  input.TokenSymbol2,
		"{status}":        input.Status,
	})

	return c.Invoke(ctx, func(db *sqlx.Tx) error {
		newCtx := dbconn.WithDB(ctx, db)
		if er := db.GetContext(newCtx, output, query, args...); er != nil {
			return er
		}
		createTransactionsInput := ethSdk.CreateTransactionsInput{
			TxId:     input.TxId,
			Source:   ethSdk.TransactionSourceExchange,
			SourceId: output.Id,
		}

		return ethmodels.Transactions.Create(newCtx, &createTransactionsInput)
	})
}

func (c *exchanges) GetExchange(ctx context.Context, input *coresSdk.GetExchangeInput, output *coresSdk.ExchangeResult) (err error) {
	where := ""
	if input.Id != 0 {
		where += `ex.id = ${id}`
	} else if input.StartupId != 0 {
		where += `ex.startup_id = ${startupId}`
	} else {
		where += "1 = 2"
	}
	stmt := `
	WITH res AS (
	    SELECT 
			ex.id,
			ex.tx_id,
			json_build_object('id',s.id,'name',s.name,'logo',sr.logo,'mission',sr.mission,'token_name',ssr.token_name,
							  'token_symbol',ssr.token_symbol) startup,
			ex.pair_name,
			ex.pair_address,
			ex.status
	    FROM exchanges ex
			INNER JOIN startups s ON s.id = ex.startup_id
			INNER JOIN startup_revisions sr ON s.current_revision_id = sr.id
			INNER JOIN startup_settings ss ON s.id = ss.startup_id
			INNER JOIN startup_setting_revisions ssr ON ss.current_revision_id = ssr.id
	    WHERE ` + where + `
	)
	SELECT row_to_json(res.*) FROM res
	`
	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{id}":        input.Id,
		"{startupId}": input.StartupId,
	})

	return c.Invoke(ctx, func(db dbconn.Q) (er error) {
		return db.GetContext(ctx, &util.PgJsonScanWrap{output}, query, args...)
	})
}

func (c *exchanges) CreateExchangeTx(ctx context.Context, input *coresSdk.CreateExchangeTxInput, output *coresSdk.CreateExchangeTxResult) (err error) {
	stmt := `
		INSERT INTO exchange_transactions(tx_id, exchange_id, account, type, token_amount1, token_amount2, status)
		VALUES (${txId}, ${exchangeId}, ${account}, ${type}, ${tokenAmount1}, ${tokenAmount2}, ${status})
		RETURNING id, status;
	`
	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{txId}":         input.TxId,
		"{exchangeId}":   input.ExchangeId,
		"{account}":      input.Account,
		"{type}":         input.Type,
		"{tokenAmount1}": input.TokenAmount1,
		"{tokenAmount2}": input.TokenAmount2,
		"{status}":       input.Status,
	})

	return c.Invoke(ctx, func(db *sqlx.Tx) error {
		newCtx := dbconn.WithDB(ctx, db)
		if er := db.GetContext(newCtx, output, query, args...); er != nil {
			return er
		}
		createTransactionsInput := ethSdk.CreateTransactionsInput{
			TxId:     input.TxId,
			Source:   ethSdk.TransactionSourceExchangeTx,
			SourceId: output.Id,
		}

		return ethmodels.Transactions.Create(newCtx, &createTransactionsInput)
	})
}

func (c *exchanges) GetExchangeTx(ctx context.Context, input *coresSdk.GetExchangeTxInput, output *coresSdk.ExchangeTxResult) (err error) {
	where := ""
	if input.Id != 0 {
		where += `et.id = ${id}`
	} else if input.TxId != "" {
		where += `et.tx_id = ${txId}`
	} else {
		where += "1 = 2"
	}
	stmt := `
		SELECT et.*
		FROM exchange_transactions et
		WHERE ` + where

	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{id}":   input.Id,
		"{txId}": input.TxId,
	})

	err = c.Invoke(ctx, func(db dbconn.Q) error {
		return db.GetContext(ctx, output, query, args...)
	})
	return
}
