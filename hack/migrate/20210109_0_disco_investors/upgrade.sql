-- auto-generated definition
CREATE TABLE disco_investors
(
    id         BIGINT                   DEFAULT id_generator()    NOT NULL
        CONSTRAINT disco_investors_id_pk
            PRIMARY KEY,
    disco_id   BIGINT                                             NOT NULL,
    uid        BIGINT                                             NOT NULL,
    eth_count  BIGINT                                             NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);
