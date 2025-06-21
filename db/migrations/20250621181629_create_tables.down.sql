DROP INDEX IF EXISTS idx_domains_name;
DROP INDEX IF EXISTS idx_wallets_address;
DROP INDEX IF EXISTS idx_metadata_public;
DROP INDEX IF EXISTS idx_metadata_report_id;
DROP INDEX IF EXISTS idx_reports_report_type_id;
DROP INDEX IF EXISTS idx_reports_submitter_id;

DROP TABLE IF EXISTS domains;
DROP TABLE IF EXISTS wallet_addresses;
DROP TABLE IF EXISTS case_metadata;
DROP TABLE IF EXISTS reports;
DROP TABLE IF EXISTS report_types;
DROP TABLE IF EXISTS reporters;
DROP TABLE IF EXISTS users;
