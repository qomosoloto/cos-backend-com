-- auto-generated definition
CREATE TABLE discos
(
    id                      BIGINT                   DEFAULT id_generator()    NOT NULL
        CONSTRAINT discos_id_pk
            PRIMARY KEY,
    startup_id              BIGINT                                             NOT NULL,
    wallet_addr             TEXT                                               NOT NULL,
    token_addr              TEXT                                               NOT NULL,
    description             TEXT                                               NOT NULL,
    fund_raising_started_at TIMESTAMP WITH TIME ZONE                           NOT NULL,
    fund_raising_ended_at   TIMESTAMP WITH TIME ZONE                           NOT NULL,
    investment_reward       BIGINT                                             NOT NULL,
    reward_decline_rate     INTEGER                                            NOT NULL,
    share_token             BIGINT                                             NOT NULL,
    min_fund_raising        BIGINT                                             NOT NULL,
    add_liquidity_pool      BIGINT                                             NOT NULL,
    total_deposit_token     BIGINT                                             NOT NULL,
    created_at              TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at              TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    state                   INTEGER                  DEFAULT 0                 NOT NULL
);

CREATE UNIQUE INDEX discos_startup_id_uindex
ON discos (startup_id);

