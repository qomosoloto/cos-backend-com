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
			ex.status,
			(SELECT count(*) FROM startups_follows_rel sfr WHERE s.id = sfr.startup_id) AS follow_count
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

func (c *exchanges) ListExchanges(ctx context.Context, input *coresSdk.ListExchangesInput, outputs interface{}) (total int, err error) {
	filterStmt := ``
	var keyword string
	if input.Keyword != "" {
		keyword = "%" + util.PgEscapeLike(input.Keyword) + "%"
		filterStmt += `AND s.name ILIKE ${keyword}`
	}

	stmt := ` 
		WITH
			exchanges_cte AS (
				SELECT
					ex.id,
					ex.tx_id,
					json_build_object('id',s.id,'name',s.name,'logo',sr.logo,'token_symbol',ssr.token_symbol) startup,
					ex.price,
					ex.newest_pooled_tokens2 liquidities,
					ex.status
				FROM exchanges ex
					INNER JOIN startups s ON s.id = ex.startup_id
					INNER JOIN startup_revisions sr ON s.current_revision_id = sr.id
					INNER JOIN startup_settings ss ON s.id = ss.startup_id
					INNER JOIN startup_setting_revisions ssr ON ss.current_revision_id = ssr.id
				WHERE 1=1` + filterStmt + `
					ORDER BY ex.created_at DESC
					LIMIT ${limit} OFFSET ${offset}
			),
			volumes_24hrs_cte AS (
				SELECT et.exchange_id, SUM(et.token_amount2) volumes_24hrs
				FROM exchanges_cte ec
				LEFT JOIN exchange_transactions et ON et.exchange_id = ec.id
					WHERE et.status = ${exTxStatusCompleted}
   						  AND (et.type = ${swap1for2} OR et.type = ${swap2for1})
						  AND et.occured_at BETWEEN (SELECT CURRENT_TIMESTAMP - interval '24 h') AND CURRENT_TIMESTAMP
					GROUP BY et.exchange_id),
			exchange_tx_rels_cte AS (
				SELECT et.exchange_id, to_char(et.occured_at, 'yyyy-mm-dd') AS occured_day, AVG(et.total_value) AS avg_price
				FROM exchanges_cte ec
				LEFT JOIN exchange_transactions et ON et.exchange_id = ec.id
					GROUP BY et.exchange_id, to_char(et.occured_at, 'yyyy-mm-dd')
					ORDER BY to_char(et.occured_at, 'yyyy-mm-dd')
					LIMIT 12
			),
			exchange_tx_rels_group_cte AS (
				SELECT etrc.exchange_id, COALESCE(json_agg(etrc), '[]'::json) price_changes
				FROM exchange_tx_rels_cte etrc
				GROUP BY etrc.exchange_id
			),
			res AS (
				SELECT ec.*, v24c.volumes_24hrs, COALESCE(etrgc.price_changes, '[]'::json) price_changes
				FROM exchanges_cte ec
				LEFT JOIN exchange_tx_rels_group_cte etrgc ON ec.id = etrgc.exchange_id
				LEFT JOIN volumes_24hrs_cte v24c ON ec.id = v24c.exchange_id
			)
			SELECT COALESCE(json_agg(r.*), '[]'::json) FROM res r;
	`

	countStmt := `
		SELECT count(*)
		FROM exchanges ex
			INNER JOIN startups s ON s.id = ex.startup_id
			INNER JOIN startup_revisions sr ON s.current_revision_id = sr.id
			INNER JOIN startup_settings ss ON s.id = ss.startup_id
			INNER JOIN startup_setting_revisions ssr ON ss.current_revision_id = ssr.id
		WHERE 1=1` + filterStmt

	query, args := util.PgMapQuery(countStmt, map[string]interface{}{
		"{keyword}": keyword,
	})

	if err = c.Invoke(ctx, func(db *sqlx.Tx) (er error) {
		return db.GetContext(ctx, &total, query, args...)
	}); err != nil {
		return
	}
	query, args = util.PgMapQuery(stmt, map[string]interface{}{
		"{keyword}":             keyword,
		"{offset}":              input.Offset,
		"{limit}":               input.GetLimit(),
		"{exTxStatusCompleted}": coresSdk.ExchangeTxStatusCompleted,
		"{swap1for2}":           coresSdk.ExchangeTxTypeSwap1for2,
		"{swap2for1}":           coresSdk.ExchangeTxTypeSwap2for1,
	})
	return total, c.Invoke(ctx, func(db *sqlx.Tx) (er error) {
		return db.GetContext(ctx, &util.PgJsonScanWrap{outputs}, query, args...)
	})
}

func (c *exchanges) CreateExchangeTx(ctx context.Context, input *coresSdk.CreateExchangeTxInput, output *coresSdk.CreateExchangeTxResult) (err error) {
	stmt := `
		INSERT INTO exchange_transactions(tx_id, exchange_id, account, type, token_amount1, token_amount2, status,
										  price_per_token1, price_per_token2)
		VALUES (${txId}, ${exchangeId}, ${account}, ${type}, ${tokenAmount1}, ${tokenAmount2}, ${status},
				${pricePerToken1}, ${pricePerToken2})
		RETURNING id, status;
	`
	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{txId}":           input.TxId,
		"{exchangeId}":     input.ExchangeId,
		"{account}":        input.Account,
		"{type}":           input.Type,
		"{tokenAmount1}":   input.TokenAmount1,
		"{tokenAmount2}":   input.TokenAmount2,
		"{status}":         input.Status,
		"{pricePerToken1}": input.PricePerToken1,
		"{pricePerToken2}": input.PricePerToken2,
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

func (c *exchanges) GetExchangeAllStatsTotal(ctx context.Context, output *coresSdk.ExchangeAllStatsTotalResult) (err error) {
	stmt := `
		WITH
			swap_48hrs_rows AS
				(SELECT * 
				FROM exchange_transactions
				WHERE (type = ${swap1for2} OR type = ${swap2for1})
					  AND status = ${exTxStatusCompleted}
					  AND occured_at BETWEEN (SELECT CURRENT_TIMESTAMP - interval '48 h') AND	CURRENT_TIMESTAMP),
			swap_24hrs_rows AS
				(SELECT * 
				FROM swap_48hrs_rows
				WHERE occured_at BETWEEN (SELECT CURRENT_TIMESTAMP - interval '24 h') AND	CURRENT_TIMESTAMP),
			volumes_24hrs_cte AS
				(SELECT COALESCE(SUM(token_amount2), 0) AS volumes_24hrs
				FROM swap_24hrs_rows),
			volumes_48hrs_cte AS
				(SELECT COALESCE(SUM(token_amount2), 0) AS volumes_48hrs
				FROM swap_48hrs_rows),
			exchange_stats_cte AS
				(SELECT COALESCE(SUM(newest_pooled_tokens2), 0) AS liquidities, COALESCE(SUM(last_pooled_tokens2), 0) AS liquidities_last
				FROM exchanges)
		
			SELECT volumes_24hrs AS volumes_24hrs,
				   CASE (volumes_48hrs-volumes_24hrs) WHEN 0 THEN 0 ELSE (2*volumes_24hrs-volumes_48hrs)/(volumes_48hrs-volumes_24hrs) END AS volumes_24hrs_rate,
				   liquidities,
				   CASE liquidities_last WHEN 0 THEN 0 ELSE (liquidities-liquidities_last)/liquidities_last END AS liquidities_rate
			FROM volumes_24hrs_cte, volumes_48hrs_cte, exchange_stats_cte
		`

	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{exTxStatusCompleted}": coresSdk.ExchangeTxStatusCompleted,
		"{swap1for2}":           coresSdk.ExchangeTxTypeSwap1for2,
		"{swap2for1}":           coresSdk.ExchangeTxTypeSwap2for1,
	})

	err = c.Invoke(ctx, func(db dbconn.Q) error {
		return db.GetContext(ctx, output, query, args...)
	})
	return
}

func (c *exchanges) GetExchangeOneStatsTotal(ctx context.Context, input *coresSdk.ExchangeOneStatsInput, output *coresSdk.ExchangeOneStatsTotalResult) (err error) {
	stmt := `
		WITH
			transaction_48hrs_rows AS
			(SELECT * 
			FROM exchange_transactions
			WHERE exchange_id = ${id}
				  AND status = ${exTxStatusCompleted}
				  AND occured_at BETWEEN (SELECT CURRENT_TIMESTAMP - interval '48 h') AND	CURRENT_TIMESTAMP),
			transaction_24hrs_rows AS
			(SELECT * 
			FROM transaction_48hrs_rows
			WHERE occured_at BETWEEN (SELECT CURRENT_TIMESTAMP - interval '24 h') AND	CURRENT_TIMESTAMP),
			swap_48hrs_rows AS
			(SELECT * 
			FROM transaction_48hrs_rows
			WHERE type = ${swap1for2} OR type = ${swap2for1}),
			swap_24hrs_rows AS
			(SELECT * 
			FROM swap_48hrs_rows
			WHERE occured_at BETWEEN (SELECT CURRENT_TIMESTAMP - interval '24 h') AND	CURRENT_TIMESTAMP),
			volumes_24hrs_cte AS
			(SELECT COALESCE(SUM(token_amount2), 0) AS volumes_24hrs
			FROM swap_24hrs_rows),
			volumes_48hrs_cte AS
			(SELECT COALESCE(SUM(token_amount2), 0) AS volumes_48hrs
			FROM swap_48hrs_rows),
			transactions_24hrs_cte AS
			(SELECT COALESCE(COUNT(*), 0) AS transactions_24hrs
			FROM transaction_24hrs_rows),
			transactions_48hrs_cte AS
			(SELECT COALESCE(COUNT(*), 0)::NUMERIC AS transactions_48hrs
			FROM transaction_48hrs_rows),
			liquidities_cte AS
			(SELECT COALESCE(newest_pooled_tokens2, 0) AS liquidities, COALESCE(last_pooled_tokens2, 0) AS liquidities_last
			FROM exchanges
			WHERE id = ${id})
			
			SELECT volumes_24hrs AS volumes_24hrs,
						 CASE (volumes_48hrs-volumes_24hrs) WHEN 0 THEN 0 ELSE (2*volumes_24hrs-volumes_48hrs)/(volumes_48hrs-volumes_24hrs) END AS volumes_24hrs_rate,
						 liquidities,
						 CASE liquidities_last WHEN 0 THEN 0 ELSE (liquidities-liquidities_last)/liquidities_last END AS liquidities_rate,
						 transactions_24hrs,
						 CASE (transactions_48hrs-transactions_24hrs) WHEN 0 THEN 0 ELSE (2*transactions_24hrs-transactions_48hrs)/(transactions_48hrs-transactions_24hrs) END AS transactions_24hrs_rate
			FROM volumes_24hrs_cte, volumes_48hrs_cte, liquidities_cte, transactions_24hrs_cte, transactions_48hrs_cte
		`

	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{id}":                  input.Id,
		"{exTxStatusCompleted}": coresSdk.ExchangeTxStatusCompleted,
		"{swap1for2}":           coresSdk.ExchangeTxTypeSwap1for2,
		"{swap2for1}":           coresSdk.ExchangeTxTypeSwap2for1,
	})

	err = c.Invoke(ctx, func(db dbconn.Q) error {
		return db.GetContext(ctx, output, query, args...)
	})
	return
}
