package proposalmodels

import (
	"context"
	"cos-backend-com/src/common/dbconn"
	"cos-backend-com/src/common/util"
	"cos-backend-com/src/libs/models"
	coresSdk "cos-backend-com/src/libs/sdk/cores"
	"github.com/jmoiron/sqlx"
)

var Proposals = &proposals{
	Connector: models.DefaultConnector,
}

type proposals struct {
	dbconn.Connector
}

func (c *proposals) CreateProposal(ctx context.Context, input *coresSdk.CreateProposalInput, output *coresSdk.CreateProposalResult) (err error) {
	stmt := `
		INSERT INTO proposals(tx_id, startup_id, wallet_addr, contract_addr, status, title, type, user_id, contact, description,
							  voter_type, support_percentage, minimum_approval_percentage, duration, has_payment, payment_addr,
							  payment_type, payment_months, payment_date, payment_amount, total_payment_amount)
		VALUES (${txId}, ${startupId}, ${walletAddr}, ${contractAddr}, ${status}, ${title}, ${type}, ${userId}, ${contact}, ${description},
				${voterType}, ${supportPercentage}, ${minimumApprovalPercentage}, ${duration}, ${hasPayment}, ${paymentAddr},
				${paymentType}, ${paymentMonths}, ${paymentDate}, ${paymentAmount}, ${totalPaymentAmount})
		RETURNING id, status;
	`
	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{txId}":                      input.TxId,
		"{startupId}":                 input.StartupId,
		"{walletAddr}":                input.WalletAddr,
		"{contractAddr}":              input.ContractAddr,
		"{status}":                    input.Status,
		"{title}":                     input.Title,
		"{type}":                      input.Type,
		"{userId}":                    input.UserId,
		"{contact}":                   input.Contact,
		"{description}":               input.Description,
		"{voterType}":                 input.VoterType,
		"{supportPercentage}":         input.SupportPercentage,
		"{minimumApprovalPercentage}": input.MinimumApprovalPercentage,
		"{duration}":                  input.Duration,
		"{hasPayment}":                input.HasPayment,
		"{paymentAddr}":               input.PaymentAddr,
		"{paymentType}":               input.PaymentType,
		"{paymentMonths}":             input.PaymentMonths,
		"{paymentDate}":               input.PaymentDate,
		"{paymentAmount}":             input.PaymentAmount,
		"{totalPaymentAmount}":        input.TotalPaymentAmount,
	})

	return c.Invoke(ctx, func(db *sqlx.Tx) error {
		return db.GetContext(ctx, output, query, args...)
	})
}

func (c *proposals) GetProposal(ctx context.Context, input *coresSdk.GetProposalInput, output *coresSdk.ProposalResult) (err error) {
	where := ""
	if input.Id != 0 {
		where += `pr.id = ${id}`
	} else if input.TxId != "" {
		where += `pr.tx_id = ${txId}`
	} else {
		where += "1 = 2"
	}
	stmt := `
	WITH res AS (
	    SELECT 
			pr.id,
			pr.tx_id,
			json_build_object('id',s.id,'name',s.name,'logo',sr.logo,'token_symbol',ssr.token_symbol) startup,
			json_build_object('id',pr.user_id,'name','') comer,
			pr.wallet_addr,
			pr.contract_addr
	    FROM proposals pr
			INNER JOIN startups s ON s.id = pr.startup_id
			INNER JOIN startup_revisions sr ON s.current_revision_id = sr.id
			INNER JOIN startup_settings ss ON s.id = ss.startup_id
			INNER JOIN startup_setting_revisions ssr ON ss.current_revision_id = ssr.id
	    WHERE ` + where + `
	)
	SELECT row_to_json(res.*) FROM res
	`
	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{id}":   input.Id,
		"{txId}": input.TxId,
	})

	return c.Invoke(ctx, func(db dbconn.Q) (er error) {
		return db.GetContext(ctx, &util.PgJsonScanWrap{output}, query, args...)
	})
}
