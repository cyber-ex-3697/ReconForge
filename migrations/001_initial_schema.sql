-- =============================================================================
-- ReconForge - Initial Database Schema
-- Version: 1.0.0
-- Description: Core tables for scan management, subdomains, and vulnerabilities
-- =============================================================================

-- -----------------------------------------------------------------------------
-- Scans table - Tracks all scan sessions
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS scans (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    target TEXT NOT NULL,
    phase TEXT NOT NULL,
    status TEXT NOT NULL,
    progress INTEGER DEFAULT 0,
    started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    completed_at DATETIME,
    checkpoint TEXT,
    metadata TEXT,
    error_count INTEGER DEFAULT 0
);

-- Indexes for scans table
CREATE INDEX IF NOT EXISTS idx_scans_target ON scans(target);
CREATE INDEX IF NOT EXISTS idx_scans_status ON scans(status);
CREATE INDEX IF NOT EXISTS idx_scans_started_at ON scans(started_at);

-- -----------------------------------------------------------------------------
-- Subdomains table - Stores discovered subdomains
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS subdomains (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    scan_id INTEGER NOT NULL,
    subdomain TEXT NOT NULL,
    resolved BOOLEAN DEFAULT FALSE,
    ip_address TEXT,
    cname TEXT,
    source TEXT,
    discovered_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_seen DATETIME DEFAULT CURRENT_TIMESTAMP,
    confidence REAL DEFAULT 1.0,
    FOREIGN KEY(scan_id) REFERENCES scans(id) ON DELETE CASCADE
);

-- Indexes for subdomains table
CREATE INDEX IF NOT EXISTS idx_subdomains_scan_id ON subdomains(scan_id);
CREATE INDEX IF NOT EXISTS idx_subdomains_subdomain ON subdomains(subdomain);
CREATE INDEX IF NOT EXISTS idx_subdomains_resolved ON subdomains(resolved);

-- -----------------------------------------------------------------------------
-- Live Hosts table - Stores live/responding hosts
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS live_hosts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    scan_id INTEGER NOT NULL,
    url TEXT NOT NULL,
    status_code INTEGER,
    title TEXT,
    content_length INTEGER,
    response_time REAL,
    server_header TEXT,
    technologies TEXT,
    screenshot_path TEXT,
    discovered_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_checked DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(scan_id) REFERENCES scans(id) ON DELETE CASCADE
);

-- Indexes for live_hosts table
CREATE INDEX IF NOT EXISTS idx_live_hosts_scan_id ON live_hosts(scan_id);
CREATE INDEX IF NOT EXISTS idx_live_hosts_url ON live_hosts(url);
CREATE INDEX IF NOT EXISTS idx_live_hosts_status_code ON live_hosts(status_code);

-- -----------------------------------------------------------------------------
-- URLs table - Stores discovered URLs
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS urls (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    scan_id INTEGER NOT NULL,
    url TEXT NOT NULL,
    method TEXT DEFAULT 'GET',
    response_status INTEGER,
    content_type TEXT,
    content_length INTEGER,
    parameters TEXT,
    discovered_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    is_api BOOLEAN DEFAULT FALSE,
    is_js BOOLEAN DEFAULT FALSE,
    FOREIGN KEY(scan_id) REFERENCES scans(id) ON DELETE CASCADE
);

-- Indexes for urls table
CREATE INDEX IF NOT EXISTS idx_urls_scan_id ON urls(scan_id);
CREATE INDEX IF NOT EXISTS idx_urls_url ON urls(url);
CREATE INDEX IF NOT EXISTS idx_urls_is_api ON urls(is_api);

-- -----------------------------------------------------------------------------
-- Vulnerabilities table - Stores discovered vulnerabilities
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS vulnerabilities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    scan_id INTEGER NOT NULL,
    url TEXT NOT NULL,
    template_id TEXT NOT NULL,
    severity TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    matched_at TEXT,
    cvss REAL,
    evidence TEXT,
    remediation TEXT,
    is_false_positive BOOLEAN DEFAULT FALSE,
    confidence REAL DEFAULT 1.0,
    discovered_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    verified_at DATETIME,
    FOREIGN KEY(scan_id) REFERENCES scans(id) ON DELETE CASCADE
);

-- Indexes for vulnerabilities table
CREATE INDEX IF NOT EXISTS idx_vulnerabilities_scan_id ON vulnerabilities(scan_id);
CREATE INDEX IF NOT EXISTS idx_vulnerabilities_severity ON vulnerabilities(severity);
CREATE INDEX IF NOT EXISTS idx_vulnerabilities_template_id ON vulnerabilities(template_id);

-- -----------------------------------------------------------------------------
-- Ports table - Stores open ports discovered during scanning
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS ports (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    scan_id INTEGER NOT NULL,
    host TEXT NOT NULL,
    port INTEGER NOT NULL,
    protocol TEXT DEFAULT 'tcp',
    service TEXT,
    version TEXT,
    discovered_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(scan_id) REFERENCES scans(id) ON DELETE CASCADE
);

-- Indexes for ports table
CREATE INDEX IF NOT EXISTS idx_ports_scan_id ON ports(scan_id);
CREATE INDEX IF NOT EXISTS idx_ports_host ON ports(host);

-- -----------------------------------------------------------------------------
-- Technologies table - Stores detected technologies
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS technologies (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    scan_id INTEGER NOT NULL,
    url TEXT NOT NULL,
    tech_name TEXT NOT NULL,
    version TEXT,
    category TEXT,
    confidence REAL DEFAULT 1.0,
    detected_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(scan_id) REFERENCES scans(id) ON DELETE CASCADE
);

-- Indexes for technologies table
CREATE INDEX IF NOT EXISTS idx_technologies_scan_id ON technologies(scan_id);
CREATE INDEX IF NOT EXISTS idx_technologies_tech_name ON technologies(tech_name);

-- -----------------------------------------------------------------------------
-- Takeovers table - Stores subdomain takeover findings
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS takeovers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    scan_id INTEGER NOT NULL,
    subdomain TEXT NOT NULL,
    cname TEXT,
    service TEXT,
    status TEXT,
    validated BOOLEAN DEFAULT FALSE,
    discovered_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(scan_id) REFERENCES scans(id) ON DELETE CASCADE
);

-- Indexes for takeovers table
CREATE INDEX IF NOT EXISTS idx_takeovers_scan_id ON takeovers(scan_id);
CREATE INDEX IF NOT EXISTS idx_takeovers_status ON takeovers(status);

-- -----------------------------------------------------------------------------
-- Scan Results Summary Table - Denormalized for quick reporting
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS scan_summary (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    scan_id INTEGER NOT NULL,
    subdomain_count INTEGER DEFAULT 0,
    live_host_count INTEGER DEFAULT 0,
    url_count INTEGER DEFAULT 0,
    vulnerability_count INTEGER DEFAULT 0,
    critical_count INTEGER DEFAULT 0,
    high_count INTEGER DEFAULT 0,
    medium_count INTEGER DEFAULT 0,
    low_count INTEGER DEFAULT 0,
    port_count INTEGER DEFAULT 0,
    technology_count INTEGER DEFAULT 0,
    takeover_count INTEGER DEFAULT 0,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(scan_id) REFERENCES scans(id) ON DELETE CASCADE,
    UNIQUE(scan_id)
);

-- Indexes for scan_summary table
CREATE INDEX IF NOT EXISTS idx_scan_summary_scan_id ON scan_summary(scan_id);

-- -----------------------------------------------------------------------------
-- Triggers for automatically updating timestamps
-- -----------------------------------------------------------------------------
CREATE TRIGGER IF NOT EXISTS update_scans_timestamp 
AFTER UPDATE ON scans
BEGIN
    UPDATE scans SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_scan_summary_timestamp 
AFTER UPDATE ON scan_summary
BEGIN
    UPDATE scan_summary SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
