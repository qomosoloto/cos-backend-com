--The nonce for login comunion is:659639759
--signature:0x207d0118dc4d143c6bf726b8026c0d5cabe2f24e648af9794923db27f9a8e0337ef956c2da176c879e93ae0e73d827cbe1d5f24c63b5918fa041d9cfd2d390461b
--public key 0x18fbdc8ed9018ae125f501f07d735faa0552c9d8

/*
    创建成功，设置成功
*/
-- INSERT INTO comunion.users (id, avatar, public_key, nonce, public_secret, private_secret, created_at, updated_at, is_hunter)
-- VALUES (comunion.fake_id('u-1') ,'https://file.comunion.io/comunion-dev/avatar/default-avatar1.png', '0x18fbdc8ed9018ae125f501f07d735faa0552c9d8', '659639759',
--         '4stvrdi66tq3nj5yt06zbcqz8alr5tvx', 'cnbqyiz8ujdi6wjp5nbr', '2020-05-31 06:57:29.439556',
--         '2020-05-31 06:57:29.439556', FALSE);
INSERT INTO comunion.startups (id, name, uid, current_revision_id, confirming_revision_id, created_at, updated_at)
VALUES (comunion.fake_id('s-1'), 'wujiu2020', comunion.fake_id('u-1'), comunion.fake_id('sr-1'), comunion.fake_id('sr-1'),
        '2020-06-13 15:56:02.466872',
        '2020-06-13 15:56:02.466872');
 INSERT INTO comunion.startup_revisions (id, startup_id, name, mission, logo, description_addr, category_id, created_at,
                                        updated_at)
VALUES (comunion.fake_id('sr-1'), comunion.fake_id('s-1'), 'wujiu2020', '成就梦想', '', 'http://baidu.com', comunion.fake_id('c-2'),
        '2020-06-13 15:56:02.466872', '2020-06-13 15:56:02.466872');
INSERT INTO comunion.transactions (id, tx_id, block_addr, source, source_id, retry_time, created_at, updated_at, state)
VALUES (comunion.fake_id('t-1'), '0xd0818eed0cf7b2e2098ae545033f26dce75a71388bcb2a94fa10532e148dd393', '0x3ba71ba6bd5df31af2c5905da27f34ceb1f6691a7c87847cbc1dd74f6af3f6e1', 'startup',
        comunion.fake_id('sr-1'), 0, '2020-06-13 15:56:02.476482', '2020-06-13 15:56:02.476482', 2);
INSERT INTO comunion.startup_settings (id, startup_id, current_revision_id, confirming_revision_id, created_at,
                                       updated_at)
VALUES (comunion.fake_id('ss-1'), comunion.fake_id('s-1'), comunion.fake_id('ssr-1'), comunion.fake_id('ssr-1'),
        '2020-06-15 16:16:44.783255',
        '2020-06-15 16:16:44.779211');
INSERT INTO comunion.startup_setting_revisions (id, startup_setting_id, token_name, token_symbol, token_addr,
                                                wallet_addrs, voter_type, voter_token_limit, assigned_voters,
                                                assigned_proposers,
                                                proposer_type, proposer_token_limit, proposal_supporters,
                                                proposal_min_approval_percent, proposal_min_duration,
                                                proposal_max_duration, created_at,
                                                updated_at)
VALUES (comunion.fake_id('ssr-1'), comunion.fake_id('ss-1'), 'wujiu', 'wujiu', '0xd0818eed0cf7b2e2098ae545033f26dce75a7138', '[
  {
    "addr": "0x18fbdc8ed9018ae125f501f07d735faa0552c9d8",
    "name": "wujiu2020"
  }
]', 1, 0, '{}', '{}', 1, 48, 48, 1, 10, 50,'2020-06-15 16:16:44.779211', '2020-06-15 16:16:44.779211');
INSERT INTO comunion.transactions (id, tx_id, block_addr, source, source_id, retry_time, created_at, updated_at, state)
VALUES (comunion.fake_id('tss-1'), '0xfb6acf18f56ef414e2681f4a99c6bc912eb1e649701f664c40354dc44fe04606', '0x3ba71ba6bd5df31af2c5905da27f34ceb1f6691a7c87847cbc1dd74f6af3f6e1',
        'startupSetting', comunion.fake_id('ssr-1'), 0, '2020-06-15 16:16:44.793364', '2020-06-15 16:16:44.793364', 2);
/*
    创建失败
*/
INSERT INTO comunion.startups (id, name, uid, current_revision_id, confirming_revision_id, created_at, updated_at)
VALUES (comunion.fake_id('s-2'), 'wujiu2021', comunion.fake_id('u-1'), null , comunion.fake_id('sr-2'),
        '2020-06-13 15:56:02.466872',
        '2020-06-13 15:56:02.466872');
INSERT INTO comunion.startup_revisions (id, startup_id, name, mission, logo, description_addr, category_id, created_at,
                                        updated_at)
VALUES (comunion.fake_id('sr-2'), comunion.fake_id('s-2'), 'wujiu2021', '成就梦想', '', 'http://baidu.com', comunion.fake_id('c-2'),
        '2020-06-13 15:56:02.466872', '2020-06-13 15:56:02.466872');
INSERT INTO comunion.transactions (id, tx_id, block_addr, source, source_id, retry_time, created_at, updated_at, state)
VALUES (comunion.fake_id('t-2'), '0xd0818eed0cf7b2e2098ae545033f26dce75a71388bcb2a94fa10532e148dd393', null , 'startup',
        comunion.fake_id('sr-2'), 0, '2020-06-13 15:56:02.476482', '2020-06-13 15:56:02.476482', 3);
/*
    创建失败，设置失败
*/
INSERT INTO comunion.startups (id, name, uid, current_revision_id, confirming_revision_id, created_at, updated_at)
VALUES (comunion.fake_id('s-3'), 'wujiu2022', comunion.fake_id('u-1'), comunion.fake_id('sr-3'), comunion.fake_id('sr-3'),
        '2020-06-13 15:56:02.466872',
        '2020-06-13 15:56:02.466872');
 INSERT INTO comunion.startup_revisions (id, startup_id, name, mission, logo, description_addr, category_id, created_at,
                                        updated_at)
VALUES (comunion.fake_id('sr-3'), comunion.fake_id('s-3'), 'wujiu2022', '成就梦想', '', 'http://baidu.com', comunion.fake_id('c-2'),
        '2020-06-13 15:56:02.466872', '2020-06-13 15:56:02.466872');
INSERT INTO comunion.transactions (id, tx_id, block_addr, source, source_id, retry_time, created_at, updated_at, state)
VALUES (comunion.fake_id('t-3'), '0xd0818eed0cf7b2e2098ae545033f26dce75a71388bcb2a94fa10532e148dd391', '0x3ba71ba6bd5df31af2c5905da27f34ceb1f6691a7c87847cbc1dd74f6af3f6e1', 'startup',
        comunion.fake_id('sr-3'), 0, '2020-06-13 15:56:02.476482', '2020-06-13 15:56:02.476482', 2);
INSERT INTO comunion.startup_settings (id, startup_id, current_revision_id, confirming_revision_id, created_at,
                                       updated_at)
VALUES (comunion.fake_id('ss-3'), comunion.fake_id('s-3'), null, comunion.fake_id('ssr-3'),
        '2020-06-15 16:16:44.783255',
        '2020-06-15 16:16:44.779211');
INSERT INTO comunion.startup_setting_revisions (id, startup_setting_id, token_name, token_symbol, token_addr,
                                                wallet_addrs, voter_type, voter_token_limit, assigned_voters,
                                                assigned_proposers,
                                                proposer_type, proposer_token_limit, proposal_supporters,
                                                proposal_min_approval_percent, proposal_min_duration,
                                                proposal_max_duration, created_at,
                                                updated_at)
VALUES (comunion.fake_id('ssr-3'), comunion.fake_id('ss-3'), 'wujiu', 'wujiu', '0xd0818eed0cf7b2e2098ae545033f26dce75a7139', '[
  {
    "addr": "0x18fbdc8ed9018ae125f501f07d735faa0552c9d8",
    "name": "wujiu2020"
  }
]', 1, 30, '{}', '{}', 1, 48, 48, 1, 50, 100,  '2020-06-15 16:16:44.779211', '2020-06-15 16:16:44.779211');
INSERT INTO comunion.transactions (id, tx_id, block_addr, source, source_id, retry_time, created_at, updated_at, state)
VALUES (comunion.fake_id('tss-3'), '0xfb6acf18f56ef414e2681f4a99c6bc912eb1e649701f664c40354dc44fe04606', null,
        'startupSetting', comunion.fake_id('ssr-3'), 0, '2020-06-15 16:16:44.793364', '2020-06-15 16:16:44.793364', 3);

INSERT INTO comunion.categories (id, name, code, source, created_at, updated_at, deleted) VALUES (comunion.fake_id('c-1'), 'Non-profit', 'Non-profit', 'startup', '2020-05-27 16:12:50.412243', '2020-05-27 16:12:50.412243', false);
INSERT INTO comunion.categories (id, name, code, source, created_at, updated_at, deleted) VALUES (comunion.fake_id('c-2'), 'Business', 'Business', 'startup', '2020-05-27 16:12:50.412243', '2020-05-27 16:12:50.412243', false);
INSERT INTO comunion.categories (id, name, code, source, created_at, updated_at, deleted) VALUES (comunion.fake_id('c-3'), 'Education', 'Education', 'startup', '2020-05-27 16:12:50.412243', '2020-05-27 16:12:50.412243', false);
