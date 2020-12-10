CREATE TABLE exchanges (
    id bigint DEFAULT comunion.id_generator() NOT NULL
        CONSTRAINT exchanges_id_pk
            PRIMARY KEY,
    tx_id text NOT NULL,
    startup_id bigint NOT NULL,
    startup_name bigint NOT NULL,
    pair_name text,
    pair_address text,
    token1_name text NOT NULL,
    token1_symbol text NOT NULL,
    token1_address text,
    token2_name text NOT NULL,
    token2_symbol text NOT NULL,
    token2_address text,
    liquidities float DEFAULT 0 NOT NULL,
    liquidities_rate float DEFAULT 0 NOT NULL,
    volumes float DEFAULT 0 NOT NULL,
    volumes_rate float DEFAULT 0 NOT NULL,
    transactions integer DEFAULT 0 NOT NULL,
    transactions_rate float DEFAULT 0 NOT NULL,
    fees float DEFAULT 0 NOT NULL,
    fees_rate float DEFAULT 0 NOT NULL,
    pooled_tokens1 float DEFAULT 0 NOT NULL,
    pooled_tokens2 float DEFAULT 0 NOT NULL,
    status integer DEFAULT 0 NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE UNIQUE INDEX exchanges_startup_id ON comunion.exchanges USING btree (startup_id);
CREATE UNIQUE INDEX exchanges_pair_address ON comunion.exchanges USING btree (pair_address);

CREATE TABLE exchange_transactions (
    id bigint DEFAULT comunion.id_generator() NOT NULL
        CONSTRAINT exchange_transactions_id_pk
            PRIMARY KEY,
    exchange_id bigint NOT NULL,
    tx_id text NOT NULL,
    account text NOT NULL,
    type integer NOT NULL,
    name text NOT NULL,
    total_value float,
    token1_amount float NOT NULL,
    token2_amount float NOT NULL,
    fee float DEFAULT 0 NOT NULL,
    status integer DEFAULT 0 NOT NULL,
    occured_at timestamp with time zone NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);