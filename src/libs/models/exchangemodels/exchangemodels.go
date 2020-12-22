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
