-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- USERS TABLE
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT UNIQUE NOT NULL,
    name TEXT,
    contact_info TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- REPORTERS TABLE
CREATE TABLE reporters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    email TEXT,
    affiliation TEXT,
    total_reports INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

-- REPORT TYPES TABLE
CREATE TABLE report_types (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

-- REPORTS TABLE
CREATE TABLE reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    submitted_by TEXT CHECK (submitted_by IN ('user', 'reporter')) NOT NULL,
    submitter_id UUID NOT NULL,
    report_type_id INT NOT NULL REFERENCES report_types(id),
    title TEXT NOT NULL,
    description TEXT,
    is_public BOOLEAN DEFAULT TRUE,
    status TEXT DEFAULT 'Open',
    created_at TIMESTAMP DEFAULT NOW()
);

-- CASE METADATA TABLE
CREATE TABLE case_metadata (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    report_id UUID NOT NULL REFERENCES reports(id) ON DELETE CASCADE,
    key TEXT NOT NULL,
    value TEXT,
    is_public BOOLEAN DEFAULT TRUE
);

-- WALLET ADDRESSES TABLE
CREATE TABLE wallet_addresses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    report_id UUID NOT NULL REFERENCES reports(id) ON DELETE CASCADE,
    address TEXT NOT NULL,
    network TEXT
);

-- DOMAINS TABLE
CREATE TABLE domains (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    report_id UUID NOT NULL REFERENCES reports(id) ON DELETE CASCADE,
    domain_name TEXT NOT NULL
);

-- INDEXES
CREATE INDEX idx_reports_submitter_id ON reports(submitter_id);
CREATE INDEX idx_reports_report_type_id ON reports(report_type_id);
CREATE INDEX idx_metadata_report_id ON case_metadata(report_id);
CREATE INDEX idx_metadata_public ON case_metadata(report_id) WHERE is_public = TRUE;
CREATE INDEX idx_wallets_address ON wallet_addresses(address);
CREATE INDEX idx_domains_name ON domains(domain_name);
