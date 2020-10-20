ALTER TABLE bounties ADD is_closed BOOLEAN DEFAULT FALSE;

alter table comunion.bounties
	add serial_no SERIAL;