package bountymodels

import (
	"context"
	"cos-backend-com/src/common/dbconn"
	"cos-backend-com/src/common/dbquery"
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/types"
	"cos-backend-com/src/common/util"
	"cos-backend-com/src/libs/apierror"
	"cos-backend-com/src/libs/models"
	"cos-backend-com/src/libs/models/ethmodels"
	coresSdk "cos-backend-com/src/libs/sdk/cores"
	ethSdk "cos-backend-com/src/libs/sdk/eth"
	"time"

	"github.com/jmoiron/sqlx"
)

var Bounties = &bounties{
	Connector: models.DefaultConnector,
}

// categories represents controller for 'categories'.
type bounties struct {
	dbconn.Connector
}

func (c *bounties) CreateBounty(ctx context.Context, startupId, uid flake.ID, input *coresSdk.CreateBountyInput) (err error) {
	stmt := `
		INSERT INTO bounties(id, startup_id, user_id, title, type, keywords, contact_email, intro, description_addr, description_file_addr,
			duration, expired_at, payments)
		VALUES (${id}, ${startupId}, ${userId}, ${title}, ${type}, ARRAY [${keywords}], ${contactEmail}, ${intro}, ${descriptionAddr},
			${descriptionFileAddr}, ${duration}, ${expiredAt}, ${payments});
	`

	query, args := util.PgMapQueryV2(stmt, map[string]interface{}{
		"{id}":                  input.Id,
		"{startupId}":           startupId,
		"{userId}":              uid,
		"{title}":               input.Title,
		"{type}":                input.Type,
		"{keywords}":            input.Keywords,
		"{contactEmail}":        input.ContactEmail,
		"{intro}":               input.Intro,
		"{descriptionAddr}":     input.DescriptionAddr,
		"{descriptionFileAddr}": input.DescriptionFileAddr,
		"{duration}":            input.Duration,
		"{expiredAt}":           time.Now().AddDate(0, 0, input.Duration),
		"{payments}":            types.JSONAny{input.Payments},
	})

	return c.Invoke(ctx, func(db *sqlx.Tx) error {
		newCtx := dbconn.WithDB(ctx, db)
		createTransactionsInput := ethSdk.CreateTransactionsInput{
			TxId:     input.TxId,
			Source:   ethSdk.TransactionSourceBounty,
			SourceId: input.Id,
		}

		if err := ethmodels.Transactions.Create(newCtx, &createTransactionsInput); err != nil {
			return err
		}
		_, err = db.ExecContext(newCtx, query, args...)
		return err
	})
}

// List Categories interface{} by input
func (c *bounties) ListBounties(ctx context.Context, startupId, uid flake.ID, isOwner bool, input *coresSdk.ListBountiesInput, outputs interface{}) (total int, err error) {
	plan := &dbquery.Plan{}
	plan.RetTotal = true

	keyword := ""
	if input.Keyword != "" {
		keyword = "%" + util.PgEscapeLike(input.Keyword) + "%"
		plan.AddCond(`AND b.title ILIKE ${keyword}`)
	}

	if startupId != flake.ID(0) {
		plan.AddCond(`AND b.startup_id = ${startupId}`)
	}

	plan.OrderBySql = ` ORDER BY is_closed ASC, created_at DESC`
	plan.LimitSql = ` LIMIT ${limit} OFFSET ${offset}`

	plan.Params = map[string]interface{}{
		"{keyword}":   keyword,
		"{uid}":       uid,
		"{startupId}": startupId,
		"{offset}":    input.Offset,
		"{limit}":     input.GetLimit(),
	}

	total, err = c.Query(ctx, uid, isOwner, outputs, plan)
	return
}

func (c *bounties) Query(ctx context.Context, uid flake.ID, isOwner bool, m interface{}, plan *dbquery.Plan) (total int, err error) {
	filterSql := `
        FROM bounties b
        INNER JOIN startups s ON s.id = b.startup_id
        INNER JOIN startup_revisions sr ON s.id = sr.startup_id
		INNER JOIN users u ON b.user_id = u.id
		LEFT JOIN hunters h ON u.id = h.user_id
	`
	joinCondition := ``

	if uid != 0 {
		joinCondition += "INNER JOIN bounties_hunters_rel bhr ON bhr.bounty_id = b.id AND bhr.uid = ${uid}"
	}

	if isOwner {
		joinCondition += "LEFT JOIN transactions t ON b.id = t.source_id AND t.source = ${source}"
	} else {
		joinCondition += "INNER JOIN transactions t ON b.id = t.source_id AND t.source = ${source}"
	}

	plan.Params["{source}"] = ethSdk.TransactionSourceBounty
	plan.Params["{undertakeBountySource}"] = ethSdk.TransactionSourceUndertakeBounty

	if plan.RetTotal {
		query :=
			`SELECT COUNT(*)
			` + filterSql + `
			` + joinCondition + `
			WHERE 1=1
			` + plan.Conditions
		query, args := util.PgMapQueryV2(query, plan.Params)

		err = c.Invoke(ctx, func(db dbconn.Q) error {
			return db.GetContext(ctx, &total, query, args...)
		})
		if err != nil {
			return
		}
	}

	query := `
	WITH bounties_cte AS (
		SELECT b.*,
			t.block_addr,
			t.state transaction_state,
			json_build_object('id',s.id,'name',s.name ,'logo' ,sr.logo) startup,
			json_build_object('id',b.user_id,'name',coalesce(h.name,u.public_key),'is_hunter',CASE WHEN h.name IS NOT NULL THEN TRUE ELSE FALSE END) created_by,
			(CASE WHEN b.status !=2 THEN 0 ELSE b.status END) order_status
		` + filterSql + `
        ` + joinCondition + `
		WHERE 1=1
		` + plan.Conditions + `
		` + plan.OrderBySql + `
		` + plan.LimitSql + `
	),bounty_hunter_rels_cte AS (
		SELECT bhr.bounty_id, bhr.uid as user_id, bhr.status, bhr.started_at, bhr.submitted_at, bhr.quited_at, bhr.rejected_at, bhr.paid_at, bhr.paid_tokens,COALESCE(h.name, u.public_key) AS name, t.state transaction_state
		FROM bounties_cte bc
		LEFT JOIN bounties_hunters_rel bhr ON bhr.bounty_id = bc.id
		LEFT JOIN users u ON bhr.uid = u.id
		LEFT JOIN hunters h ON u.id = h.user_id
		LEFT JOIN transactions t ON t.source_id = bhr.id AND t.source = ${undertakeBountySource}
	),bounty_hunter_rels_aggregate_cte AS (
		SELECT bhrc.bounty_id, COALESCE(json_agg(bhrc), '[]'::json) hunters
		FROM bounty_hunter_rels_cte bhrc
		GROUP BY bhrc.bounty_id
	),res AS (
	    SELECT bc.*, COALESCE(bhrac.hunters, '[]'::json) hunters
	    FROM bounties_cte bc
		LEFT JOIN bounty_hunter_rels_aggregate_cte bhrac ON bc.id = bhrac.bounty_id
		` + plan.OrderBySql + `
	)
	SELECT
		COALESCE(json_agg(r.*), '[]'::json)
	FROM res r;`
	query, args := util.PgMapQueryV2(query, plan.Params)

	err = c.Invoke(ctx, func(db dbconn.Q) error {
		return db.GetContext(ctx, &util.PgJsonScanWrap{m}, query, args...)
	})
	return
}

func (c *bounties) GetBounty(ctx context.Context, id flake.ID, isOwner bool, output interface{}) (err error) {
	plan := &dbquery.Plan{}
	plan.AddCond(`AND b.id = ${id}`)

	plan.Params = map[string]interface{}{
		"{id}": id,
	}
	var v struct{ Id flake.ID }
	_, err = c.Query(ctx, 0, isOwner, &util.PgJsonScanWrapValues{&[]interface{}{&v}, &[]interface{}{output}}, plan)
	if err != nil {
		return
	}
	if v.Id == 0 {
		err = apierror.ErrNotFound
		return
	}
	return
}

func (c *bounties) CreateUndertakeBounty(ctx context.Context, input *coresSdk.CreateUndertakeBountyInput, output *coresSdk.UndertakeBountyResult) (err error) {
	stmt := `
		INSERT INTO bounties_hunters_rel(bounty_id, uid, status)
		VALUES (${bountyId}, ${uid}, ${status}) RETURNING id, bounty_id, status;
	`
	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{bountyId}": input.BountyId,
		"{uid}":      input.UserId,
		"{status}":   coresSdk.UndertakeBountyStatusNull,
	})

	return c.Invoke(ctx, func(db *sqlx.Tx) error {
		newCtx := dbconn.WithDB(ctx, db)
		if er := db.GetContext(newCtx, output, query, args...); er != nil {
			return er
		}
		createTransactionsInput := ethSdk.CreateTransactionsInput{
			TxId:     input.TxId,
			Source:   ethSdk.TransactionSourceUndertakeBounty,
			SourceId: output.Id,
		}

		return ethmodels.Transactions.Create(newCtx, &createTransactionsInput)
	})
}

func (c *bounties) SubmittedUndertakeBounty(ctx context.Context, input *coresSdk.UpdateUndertakeBountyInput, output *coresSdk.UndertakeBountyResult) (err error) {
	input.Status = coresSdk.UndertakeBountyStatusSubmitted
	return c.UpdateUndertakeBounty(ctx, input, output)
}

func (c *bounties) QuitedUndertakeBounty(ctx context.Context, input *coresSdk.UpdateUndertakeBountyInput, output *coresSdk.UndertakeBountyResult) (err error) {
	input.Status = coresSdk.UndertakeBountyStatusQuited
	return c.UpdateUndertakeBounty(ctx, input, output)
}

func (c *bounties) PaidUndertakeBounty(ctx context.Context, input *coresSdk.UpdateUndertakeBountyInput, output *coresSdk.UndertakeBountyResult) (err error) {
	input.Status = coresSdk.UndertakeBountyStatusPaid
	return c.UpdateUndertakeBounty(ctx, input, output)
}

func (c *bounties) RejectedUndertakeBounty(ctx context.Context, input *coresSdk.UpdateUndertakeBountyInput, output *coresSdk.UndertakeBountyResult) (err error) {
	input.Status = coresSdk.UndertakeBountyStatusRejected
	return c.UpdateUndertakeBounty(ctx, input, output)
}

func (c *bounties) UpdateUndertakeBounty(ctx context.Context, input *coresSdk.UpdateUndertakeBountyInput, output *coresSdk.UndertakeBountyResult) (err error) {
	fields := ""
	values := ""
	if input.Status == coresSdk.UndertakeBountyStatusSubmitted {
		fields += "status, submitted_at, updated_at"
		values += "${status}, current_timestamp, current_timestamp"
	} else if input.Status == coresSdk.UndertakeBountyStatusQuited {
		fields += "status, quited_at, updated_at"
		values += "${status}, current_timestamp, current_timestamp"
	} else if input.Status == coresSdk.UndertakeBountyStatusPaid {
		fields += "status, paid_at, updated_at"
		values += "${status}, current_timestamp, current_timestamp"
	} else if input.Status == coresSdk.UndertakeBountyStatusRejected {
		fields += "status, rejected_at, updated_at"
		values += "${status}, current_timestamp, current_timestamp"
	}

	stmt := `
	UPDATE bounties_hunters_rel SET (
	    ` + fields + `
	) = (
	` + values + `
	)
	WHERE bounty_id = ${bountyId} AND uid = ${uid}
	RETURNING id, bounty_id, status;
	`
	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{bountyId}": input.BountyId,
		"{uid}":      input.UserId,
		"{status}":   input.Status,
	})

	return c.Invoke(ctx, func(db dbconn.Q) (er error) {
		return db.GetContext(ctx, output, query, args...)
	})
}

func (c *bounties) ClosedBounty(ctx context.Context) (err error) {
	stmt := `
		UPDATE bounties SET status = ${statusClosed}, is_closed = true
		WHERE status != ${statusClosed} AND expired_at<current_timestamp;
	`

	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{statusClosed}": coresSdk.BountyStatusClosed,
	})

	return c.Invoke(ctx, func(db dbconn.Q) (er error) {
		_, er = db.ExecContext(ctx, query, args...)
		return
	})
}
