package startupmodels

import (
	"context"
	"cos-backend-com/src/common/dbconn"
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/util"
	"cos-backend-com/src/libs/models"
	"cos-backend-com/src/libs/models/ethmodels"
	coresSdk "cos-backend-com/src/libs/sdk/cores"
	ethSdk "cos-backend-com/src/libs/sdk/eth"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var Startups = &startups{
	Connector: models.DefaultConnector,
}

type startups struct {
	dbconn.Connector
}

func (c *startups) List(ctx context.Context, input *coresSdk.ListStartupsInput, outputs interface{}) (total int, err error) {
	filterStmt := ``
	if input.CategoryId != 0 {
		filterStmt += `AND sr.category_id = ${categoryId}`
	}
	var keyword string
	if input.Keyword != "" {
		keyword = "%" + util.PgEscapeLike(input.Keyword) + "%"
		filterStmt += `AND s.name ILIKE ${keyword}`
	}

	stmt := ` 
		WITH res AS(
			SELECT
				s.id,
				sr.name,
				sr.logo,
				sr.mission,
				sr.description_addr,
				c AS category,
				s.created_at,
				(SELECT count(*) FROM bounties b WHERE s.id = b.startup_id) AS bounty_count,
				(SELECT count(*) FROM startups_follows_rel sfr WHERE s.id = sfr.startup_id) AS follow_count
			FROM startups s
				INNER JOIN startup_revisions sr ON s.current_revision_id = sr.id
				INNER JOIN categories c ON c.id = sr.category_id
				INNER JOIN startup_settings ss ON s.id = ss.startup_id
				INNER JOIN startup_setting_revisions ssr ON ss.current_revision_id = ssr.id
			WHERE 1=1` + filterStmt + `
				ORDER BY s.created_at DESC
				LIMIT ${limit} OFFSET ${offset}
		)
		SELECT COALESCE(json_agg(r.*), '[]'::json) FROM res r;
	`

	countStmt := `
		SELECT count(*)
		FROM startups s
			INNER JOIN startup_revisions sr ON s.current_revision_id = sr.id
			INNER JOIN categories c ON c.id = sr.category_id
			INNER JOIN startup_settings ss ON s.id = ss.startup_id
			INNER JOIN startup_setting_revisions ssr ON ss.current_revision_id = ssr.id
		WHERE 1=1` + filterStmt

	query, args := util.PgMapQuery(countStmt, map[string]interface{}{
		"{categoryId}": input.CategoryId,
		"{isIRO}":      input.IsIRO,
		"{keyword}":    keyword,
	})

	if err = c.Invoke(ctx, func(db *sqlx.Tx) (er error) {
		return db.GetContext(ctx, &total, query, args...)
	}); err != nil {
		return
	}
	query, args = util.PgMapQuery(stmt, map[string]interface{}{
		"{categoryId}": input.CategoryId,
		"{isIRO}":      input.IsIRO,
		"{keyword}":    keyword,
		"{offset}":     input.Offset,
		"{limit}":      input.GetLimit(),
	})
	return total, c.Invoke(ctx, func(db *sqlx.Tx) (er error) {
		return db.GetContext(ctx, &util.PgJsonScanWrap{outputs}, query, args...)
	})
}

func (c *startups) ListMe(ctx context.Context, uid flake.ID, input *coresSdk.ListStartupsInput, outputs interface{}) (total int, err error) {
	filterStmt := ``
	if input.IsInBlock != nil {
		if *input.IsInBlock {
			filterStmt += `AND t1.state = ${transactionStateSuccess}`
		}
		if *input.IsInBlock {
			filterStmt += `AND t1.state != ${transactionStateSuccess}`
		}
	}
	stmt := `
		WITH res AS(
			SELECT
				s.id,
				sr.name,
				sr.logo,
				sr.mission,
				sr.description_addr,
				c AS category,
			    t1.state,
			    t2.state AS setting_state
			FROM startups s
			    INNER JOIN startup_revisions sr ON s.confirming_revision_id = sr.id
			    INNER JOIN transactions t1 ON t1.source_id = sr.id AND t1.source = ${sourceStartup}
			    INNER JOIN categories c ON c.id = sr.category_id
			    LEFT JOIN startup_settings ss ON s.id = ss.startup_id
			    LEFT JOIN startup_setting_revisions ssr ON ss.confirming_revision_id = ssr.id
			    LEFT JOIN transactions t2 ON t2.source_id = ssr.id AND t2.source = ${sourceStartupSetting}
			WHERE s.uid = ${uid}
				` + filterStmt + `
			ORDER BY s.created_at DESC
			LIMIT ${limit} OFFSET ${offset}
		)
		SELECT COALESCE(json_agg(r.*), '[]'::json) FROM res r;
	`

	countStmt := `
		SELECT COUNT(*)
		FROM startups s
		    INNER JOIN startup_revisions sr ON s.confirming_revision_id = sr.id
		    INNER JOIN transactions t1 ON t1.source_id = sr.id AND t1.source = ${sourceStartup}
		    INNER JOIN categories c ON c.id = sr.category_id
		WHERE s.uid = ${uid}
			` + filterStmt + `;
	`

	query, args := util.PgMapQuery(countStmt, map[string]interface{}{
		"{uid}":                     uid,
		"{categoryId}":              input.CategoryId,
		"{isIRO}":                   input.IsIRO,
		"{sourceStartup}":           ethSdk.TransactionSourceStartup,
		"{sourceStartupSetting}":    ethSdk.TransactionSourceStartupSetting,
		"{transactionStateSuccess}": ethSdk.TransactionStateSuccess,
	})

	if err = c.Invoke(ctx, func(db *sqlx.Tx) (er error) {
		return db.GetContext(ctx, &total, query, args...)
	}); err != nil {
		return
	}
	query, args = util.PgMapQuery(stmt, map[string]interface{}{
		"{uid}":                     uid,
		"{categoryId}":              input.CategoryId,
		"{isIRO}":                   input.IsIRO,
		"{sourceStartup}":           ethSdk.TransactionSourceStartup,
		"{sourceStartupSetting}":    ethSdk.TransactionSourceStartupSetting,
		"{transactionStateSuccess}": ethSdk.TransactionStateSuccess,
		"{offset}":                  input.Offset,
		"{limit}":                   input.GetLimit(),
	})
	return total, c.Invoke(ctx, func(db *sqlx.Tx) (er error) {
		return db.GetContext(ctx, &util.PgJsonScanWrap{outputs}, query, args...)
	})
}

func (c *startups) ListMeFollowed(ctx context.Context, uid flake.ID, input *coresSdk.ListStartupsInput, outputs interface{}) (total int, err error) {
	stmt := `
		WITH res AS(
			SELECT
				s.id,
				sr.name,
				sr.logo,
				sr.mission,
				sr.description_addr,
				c AS category,
			    t1.state,
			    t2.state AS setting_state
			FROM startups s
			    INNER JOIN startup_revisions sr ON s.confirming_revision_id = sr.id
			    INNER JOIN transactions t1 ON t1.source_id = sr.id AND t1.source = ${sourceStartup}
			    INNER JOIN categories c ON c.id = sr.category_id
				INNER JOIN startups_follows_rel sfr ON s.id = sfr.startup_id AND sfr.user_id = ${uid}
			    LEFT JOIN startup_settings ss ON s.id = ss.startup_id
			    LEFT JOIN startup_setting_revisions ssr ON ss.confirming_revision_id = ssr.id
			    LEFT JOIN transactions t2 ON t2.source_id = ssr.id AND t2.source = ${sourceStartupSetting}
			WHERE 1=1
			ORDER BY s.created_at DESC
			LIMIT ${limit} OFFSET ${offset}
		)
		SELECT COALESCE(json_agg(r.*), '[]'::json) FROM res r;
	`

	countStmt := `
		SELECT COUNT(*)
		FROM startups s
		    INNER JOIN startup_revisions sr ON s.confirming_revision_id = sr.id
		    INNER JOIN transactions t1 ON t1.source_id = sr.id AND t1.source = ${sourceStartup}
		    INNER JOIN categories c ON c.id = sr.category_id
			INNER JOIN startups_follows_rel sfr ON s.id = sfr.startup_id AND sfr.user_id = ${uid}
	`

	query, args := util.PgMapQuery(countStmt, map[string]interface{}{
		"{uid}":                  uid,
		"{categoryId}":           input.CategoryId,
		"{isIRO}":                input.IsIRO,
		"{sourceStartup}":        ethSdk.TransactionSourceStartup,
		"{sourceStartupSetting}": ethSdk.TransactionSourceStartupSetting,
	})

	if err = c.Invoke(ctx, func(db *sqlx.Tx) (er error) {
		return db.GetContext(ctx, &total, query, args...)
	}); err != nil {
		return
	}
	query, args = util.PgMapQuery(stmt, map[string]interface{}{
		"{uid}":                  uid,
		"{categoryId}":           input.CategoryId,
		"{isIRO}":                input.IsIRO,
		"{sourceStartup}":        ethSdk.TransactionSourceStartup,
		"{sourceStartupSetting}": ethSdk.TransactionSourceStartupSetting,
		"{offset}":               input.Offset,
		"{limit}":                input.GetLimit(),
	})
	return total, c.Invoke(ctx, func(db *sqlx.Tx) (er error) {
		return db.GetContext(ctx, &util.PgJsonScanWrap{outputs}, query, args...)
	})
}

func (c *startups) Get(ctx context.Context, id flake.ID, output interface{}) (err error) {
	stmt := `
	WITH res AS (
	    SELECT 
			s.id,
			sr.name,
			sr.logo,
			sr.mission,
			sr.description_addr,
			c   AS category,
			ssr AS settings,
			t1  AS transaction,
			s.created_at,
			(SELECT count(*) FROM startups_follows_rel sfr WHERE s.id = sfr.startup_id) AS follow_count
	    FROM startups s
			INNER JOIN startup_revisions sr ON s.current_revision_id = sr.id
			INNER JOIN transactions t1 ON t1.source_id = sr.id AND t1.source = ${sourceStartup} AND t1.state = ${stateSuccess}
			INNER JOIN categories c ON c.id = sr.category_id
			INNER JOIN startup_settings ss ON s.id = ss.startup_id
			INNER JOIN startup_setting_revisions ssr ON ss.current_revision_id = ssr.id
			INNER JOIN transactions t2 ON t2.source_id = ssr.id AND t2.source = ${sourceStartupSetting} AND t2.state = ${stateSuccess}
	    WHERE s.id = ${id}
	)
	SELECT row_to_json(res.*) FROM res
	`
	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{id}":                   id,
		"{sourceStartup}":        ethSdk.TransactionSourceStartup,
		"{sourceStartupSetting}": ethSdk.TransactionSourceStartupSetting,
		"{stateSuccess}":         ethSdk.TransactionStateSuccess,
	})

	return c.Invoke(ctx, func(db dbconn.Q) (er error) {
		return db.GetContext(ctx, &util.PgJsonScanWrap{output}, query, args...)
	})
}

func (c *startups) GetMe(ctx context.Context, uid, id flake.ID, output interface{}) (err error) {
	stmt := `
	WITH res AS (
	    SELECT 
			s.id,
			sr.name,
			sr.logo,
			sr.mission,
			sr.description_addr,
			c   AS category,
			ssr AS settings,
			t1  AS transaction,
			s.created_at,
			(SELECT count(*) FROM startups_follows_rel sfr WHERE s.id = sfr.startup_id) AS follow_count
	    FROM startups s
			LEFT JOIN startup_revisions sr ON s.confirming_revision_id = sr.id
			LEFT JOIN transactions t1 ON t1.source_id = sr.id AND t1.source = ${sourceStartup}
			LEFT JOIN categories c ON c.id = sr.category_id
			LEFT JOIN startup_settings ss ON s.id = ss.startup_id
			LEFT JOIN startup_setting_revisions ssr ON ss.confirming_revision_id = ssr.id
	    WHERE s.uid = ${uid}
			AND s.id = ${id}
	)
	SELECT row_to_json(res.*) FROM res
	`
	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{uid}":                  uid,
		"{id}":                   id,
		"{sourceStartup}":        ethSdk.TransactionSourceStartup,
		"{sourceStartupSetting}": ethSdk.TransactionSourceStartupSetting,
		"{stateSuccess}":         ethSdk.TransactionStateSuccess,
	})

	return c.Invoke(ctx, func(db dbconn.Q) (er error) {
		return db.GetContext(ctx, &util.PgJsonScanWrap{output}, query, args...)
	})
}

func (c *startups) CreateWithRevision(ctx context.Context, uid flake.ID, input *coresSdk.CreateStartupInput, startupId *flake.ID) (err error) {
	stmt := `
		UPDATE startups SET
		(
		    confirming_revision_id, updated_at
		) = (
		    ${confirmingRevisionId}, CURRENT_TIMESTAMP
		) WHERE id = ${id};
	`
	return c.Invoke(ctx, func(db *sqlx.Tx) error {
		newCtx := dbconn.WithDB(ctx, db)
		if er := c.Create(newCtx, uid, input, startupId); er != nil {
			return er
		}

		var startupRevisionId flake.ID
		if er := c.CreateRevision(newCtx, *startupId, &input.CreateStartupRevisionInput, &startupRevisionId); er != nil {
			return er
		}

		query, args := util.PgMapQuery(stmt, map[string]interface{}{
			"{id}":                   *startupId,
			"{confirmingRevisionId}": startupRevisionId,
		})

		return c.Invoke(newCtx, func(db dbconn.Q) (er error) {
			_, er = db.ExecContext(newCtx, query, args...)
			return er
		})
	})
}

func (c *startups) Create(ctx context.Context, uid flake.ID, input *coresSdk.CreateStartupInput, output interface{}) (err error) {
	stmt := `
		INSERT INTO startups(id, name, uid)
		VALUES (${id}, ${name},${uid}) RETURNING id;
	`

	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{uid}":  uid,
		"{id}":   input.Id,
		"{name}": input.Name,
	})

	return c.Invoke(ctx, func(db dbconn.Q) (er error) {
		return db.GetContext(ctx, output, query, args...)
	})
}

func (c *startups) CreateRevision(ctx context.Context, startupId flake.ID, input *coresSdk.CreateStartupRevisionInput, revisionId *flake.ID) (err error) {
	stmt := `
		INSERT INTO startup_revisions(startup_id, name, mission, logo, description_addr, category_id)
		VALUES (${startupId}, ${name}, ${mission}, ${logo}, ${descriptionAddr}, ${categoryId}) RETURNING id;
	`

	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{startupId}":       startupId,
		"{name}":            input.Name,
		"{mission}":         input.Mission,
		"{logo}":            input.Logo,
		"{descriptionAddr}": input.DescriptionAddr,
		"{categoryId}":      input.CategoryId,
	})
	return c.Invoke(ctx, func(db *sqlx.Tx) error {
		newCtx := dbconn.WithDB(ctx, db)
		if er := db.GetContext(newCtx, revisionId, query, args...); er != nil {
			return er
		}
		createTransactionsInput := ethSdk.CreateTransactionsInput{
			TxId:     input.TxId,
			Source:   ethSdk.TransactionSourceStartup,
			SourceId: *revisionId,
		}

		return ethmodels.Transactions.Create(newCtx, &createTransactionsInput)
	})
}

func (c *startups) UpdateWithRevision(ctx context.Context, uid, id flake.ID, input *coresSdk.UpdateStartupInput) (err error) {
	stmt := `
		UPDATE startups SET
		(
		    name, confirming_revision_id, updated_at
		) = (
		    ${name}, ${confirmingRevisionId}, CURRENT_TIMESTAMP
		) WHERE id = ${id} AND uid = ${uid};
	`
	return c.Invoke(ctx, func(db *sqlx.Tx) error {
		newCtx := dbconn.WithDB(ctx, db)
		var startupRevisionId flake.ID
		if er := c.CreateRevision(newCtx, id, &input.CreateStartupRevisionInput, &startupRevisionId); er != nil {
			return er
		}

		query, args := util.PgMapQuery(stmt, map[string]interface{}{
			"{id}":                   id,
			"{uid}":                  uid,
			"{name}":                 input.Name,
			"{confirmingRevisionId}": startupRevisionId,
		})

		return c.Invoke(newCtx, func(db dbconn.Q) (er error) {
			_, er = db.ExecContext(newCtx, query, args...)
			return er
		})
	})
}

func (c *startups) Update(ctx context.Context, uid, id flake.ID, input *coresSdk.UpdateStartupInput, output interface{}) (err error) {
	stmt := `
	UPDATE startups SET (
	    name, updated_at
	) = (
	    ${name}, current_timestamp
	)
	WHERE id = ${id}
	AND uid = ${uid}
	RETURNING id`

	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{id}":   id,
		"{uid}":  uid,
		"{name}": input.Name,
	})

	return c.Invoke(ctx, func(db dbconn.Q) (er error) {
		return db.GetContext(ctx, &output, query, args...)
	})
}

func (c *startups) Restore(ctx context.Context, uid, id flake.ID) (err error) {
	stmt := `
		UPDATE startups s
		SET (confirming_revision_id,updated_at)= (current_revision_id,current_timestamp)
		WHERE s.uid = ${uid}
		  AND s.id = ${id};
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

func (c *startups) NextId(ctx context.Context) (netxtId flake.ID, err error) {
	err = c.Invoke(ctx, func(db dbconn.Q) error {
		return db.GetContext(ctx, &netxtId, `SELECT id_generator()`)
	})
	return
}

func (c *startups) HasFollowed(ctx context.Context, uid, startupId flake.ID, output interface{}) (err error) {
	stmt := "SELECT EXISTS(SELECT 1 FROM startups_follows_rel WHERE user_id = ${uid} AND startup_id= ${startupId}) AS has_followed"

	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{uid}":       uid,
		"{startupId}": startupId,
	})
	err = c.Invoke(ctx, func(db dbconn.Q) error {
		return db.GetContext(ctx, output, query, args...)
	})
	fmt.Println(output)
	return
}

func (c *startups) GetToken(ctx context.Context, startupId flake.ID, output interface{}) (err error) {
	stmt := `
		SELECT ssr.token_name, ssr.token_symbol
		FROM startup_setting_revisions ssr
			LEFT JOIN startup_settings ss ON ssr."id" = ss.confirming_revision_id
		WHERE ss.startup_id = ${startupId}
	`
	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{startupId}": startupId,
	})

	err = c.Invoke(ctx, func(db dbconn.Q) error {
		return db.GetContext(ctx, output, query, args...)
	})
	return
}

func (c *startups) IsTokenAddrBinding(ctx context.Context, input *coresSdk.IsTokenAddrBindingInput, output interface{}) (err error) {
	stmt := `
		SELECT case WHEN res.amount > 1 then (SELECT id from startup_setting_revisions where token_addr=${tokenAddr} limit 1) else 0 END
		FROM 
		(
			SELECT count(*) as amount from startup_setting_revisions WHERE token_addr=${tokenAddr}

		) as res
	`
	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{tokenAddr}": input.TokenAddr,
	})

	err = c.Invoke(ctx, func(db dbconn.Q) error {
		return db.GetContext(ctx, output, query, args...)
	})
	return
}

func (c *startups) GetId(ctx context.Context, tokenAddr string) (startupId flake.ID, err error) {
	tokenWhere := `token_addr = '` + tokenAddr + `'`
	stmt := `
		SELECT ss.startup_id
		FROM startup_setting_revisions ssr
			LEFT JOIN startup_settings ss ON ss.id = ssr.startup_setting_id
		WHERE ` + tokenWhere + ` AND ss.confirming_revision_id = ssr.id
		ORDER BY ssr.created_at
		LIMIT 1
	`
	query, args := util.PgMapQuery(stmt, map[string]interface{}{})

	err = c.Invoke(ctx, func(db dbconn.Q) error {
		return db.GetContext(ctx, &startupId, query, args...)
	})
	return
}
