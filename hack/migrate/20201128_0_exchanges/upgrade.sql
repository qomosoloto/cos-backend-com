CREATE TABLE exchanges (
    id bigint DEFAULT comunion.id_generator() NOT NULL
        CONSTRAINT exchanges_id_pk
            PRIMARY KEY,
    tx_id text NOT NULL,
    startup_id bigint NOT NULL,
    pair_name text,
    pair_address text,
    token_name1 text NOT NULL,
    token_symbol1 text NOT NULL,
    token_address1 text,
    token_divider1 integer DEFAULT 1 NOT NULL,
    token_name2 text NOT NULL,
    token_symbol2 text NOT NULL,
    token_address2 text,
    token_divider2 integer DEFAULT 1 NOT NULL,
    newest_day text DEFAULT '' NOT NULL,
    newest_pooled_tokens1 float DEFAULT 0 NOT NULL,
    newest_pooled_tokens2 float DEFAULT 0 NOT NULL,
    last_day text DEFAULT '' NOT NULL,
    last_pooled_tokens1 float DEFAULT 0 NOT NULL,
    last_pooled_tokens2 float DEFAULT 0 NOT NULL,
    price float DEFAULT 0 NOT NULL,
    fees float DEFAULT 0 NOT NULL,
    reserve0 text DEFAULT '' NOT NULL,
    reserve1 text DEFAULT '' NOT NULL,
    status integer DEFAULT 0 NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

COMMENT ON COLUMN comunion.exchanges.status IS '0：待确认，1：已完成，2：未完成';

CREATE UNIQUE INDEX exchanges_startup_id ON comunion.exchanges USING btree (startup_id);
CREATE UNIQUE INDEX exchanges_pair_address ON comunion.exchanges USING btree (pair_address);

CREATE TABLE exchange_transactions (
    id bigint DEFAULT comunion.id_generator() NOT NULL
        CONSTRAINT exchange_transactions_id_pk
            PRIMARY KEY,
    tx_id text NOT NULL,
    exchange_id bigint NOT NULL,
    sender text DEFAULT '' NOT NULL,
    receiver text DEFAULT '' NOT NULL,
    type integer NOT NULL,
    name text DEFAULT '' NOT NULL,
    total_value float DEFAULT 0 NOT NULL,
    token_amount1 float DEFAULT 0 NOT NULL,
    token_amount2 float DEFAULT 0 NOT NULL,
    amount0 text DEFAULT '' NOT NULL,
    amount1 text DEFAULT '' NOT NULL,
    fee float DEFAULT 0 NOT NULL,
    price_per_token1 float DEFAULT 0 NOT NULL,
    price_per_token2 float DEFAULT 0 NOT NULL,
    status integer DEFAULT 0 NOT NULL,
    occured_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

COMMENT ON COLUMN comunion.exchange_transactions.type IS '1：增加流动性，2：删除流动性，3：1兑换2，4：2兑换1';
COMMENT ON COLUMN comunion.exchange_transactions.status IS '0：待确认，1：已完成，2：未完成';
CREATE UNIQUE INDEX exchange_transactions_tx_id ON comunion.exchange_transactions USING btree (tx_id);