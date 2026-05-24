-- =============================================================================
-- ReconForge - Session Recorder Schema
-- Version: 1.0.0
-- Description: Tables for storing request/response pairs for replay
-- =============================================================================

-- -----------------------------------------------------------------------------
-- Sessions table - Stores browser/session information
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    user_agent TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_used DATETIME DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    metadata TEXT
);

-- Indexes for sessions table
CREATE INDEX IF NOT EXISTS idx_sessions_name ON sessions(name);
CREATE INDEX IF NOT EXISTS idx_sessions_is_active ON sessions(is_active);

-- -----------------------------------------------------------------------------
-- Requests table - Stores captured HTTP requests
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS requests (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL,
    scan_id INTEGER,
    timestamp DATETIME NOT NULL,
    method TEXT NOT NULL,
    url TEXT NOT NULL,
    path TEXT,
    query_string TEXT,
    headers TEXT,
    body TEXT,
    body_size INTEGER,
    duration INTEGER,
    metadata TEXT,
    FOREIGN KEY(session_id) REFERENCES sessions(id) ON DELETE CASCADE
);

-- Indexes for requests table
CREATE INDEX IF NOT EXISTS idx_requests_session_id ON requests(session_id);
CREATE INDEX IF NOT EXISTS idx_requests_timestamp ON requests(timestamp);
CREATE INDEX IF NOT EXISTS idx_requests_url ON requests(url);
CREATE INDEX IF NOT EXISTS idx_requests_method ON requests(method);

-- -----------------------------------------------------------------------------
-- Responses table - Stores HTTP responses for requests
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS responses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    request_id TEXT NOT NULL,
    status_code INTEGER NOT NULL,
    status_text TEXT,
    headers TEXT,
    body TEXT,
    body_size INTEGER,
    content_type TEXT,
    content_encoding TEXT,
    response_time REAL,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    hash TEXT,
    FOREIGN KEY(request_id) REFERENCES requests(id) ON DELETE CASCADE
);

-- Indexes for responses table
CREATE INDEX IF NOT EXISTS idx_responses_request_id ON responses(request_id);
CREATE INDEX IF NOT EXISTS idx_responses_status_code ON responses(status_code);
CREATE INDEX IF NOT EXISTS idx_responses_hash ON responses(hash);

-- -----------------------------------------------------------------------------
-- Replay Tasks table - Tracks replay operations
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS replay_tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    request_id TEXT NOT NULL,
    status TEXT DEFAULT 'pending',
    modifications TEXT,
    executed_at DATETIME,
    result_status INTEGER,
    result_body TEXT,
    error_message TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(request_id) REFERENCES requests(id) ON DELETE CASCADE
);

-- Indexes for replay_tasks table
CREATE INDEX IF NOT EXISTS idx_replay_tasks_request_id ON replay_tasks(request_id);
CREATE INDEX IF NOT EXISTS idx_replay_tasks_status ON replay_tasks(status);

-- -----------------------------------------------------------------------------
-- Exploitation Workflows table - Tracks exploitation sequences
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS exploitation_workflows (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    target TEXT NOT NULL,
    vulnerability TEXT,
    status TEXT DEFAULT 'pending',
    steps TEXT,
    result TEXT,
    started_at DATETIME,
    completed_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for exploitation_workflows table
CREATE INDEX IF NOT EXISTS idx_workflows_target ON exploitation_workflows(target);
CREATE INDEX IF NOT EXISTS idx_workflows_status ON exploitation_workflows(status);

-- -----------------------------------------------------------------------------
-- Exploitation Steps table - Individual steps in workflows
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS exploitation_steps (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    workflow_id TEXT NOT NULL,
    step_order INTEGER NOT NULL,
    step_type TEXT NOT NULL,
    request_id TEXT,
    modifications TEXT,
    expected_result TEXT,
    actual_result TEXT,
    status TEXT DEFAULT 'pending',
    executed_at DATETIME,
    error_message TEXT,
    FOREIGN KEY(workflow_id) REFERENCES exploitation_workflows(id) ON DELETE CASCADE,
    FOREIGN KEY(request_id) REFERENCES requests(id)
);

-- Indexes for exploitation_steps table
CREATE INDEX IF NOT EXISTS idx_steps_workflow_id ON exploitation_steps(workflow_id);
CREATE INDEX IF NOT EXISTS idx_steps_status ON exploitation_steps(status);

-- -----------------------------------------------------------------------------
-- Comparison Results table - Stores response comparison results
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS comparison_results (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    request_id TEXT NOT NULL,
    comparison_type TEXT NOT NULL,
    similarity_score REAL,
    differences TEXT,
    analyzed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(request_id) REFERENCES requests(id) ON DELETE CASCADE
);

-- Indexes for comparison_results table
CREATE INDEX IF NOT EXISTS idx_comparison_request_id ON comparison_results(request_id);

-- -----------------------------------------------------------------------------
-- Rate Limit Events table - Tracks rate limiting incidents
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS rate_limit_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    target TEXT NOT NULL,
    status_code INTEGER,
    retry_after INTEGER,
    occurred_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    resolved_at DATETIME
);

-- Indexes for rate_limit_events table
CREATE INDEX IF NOT EXISTS idx_ratelimit_target ON rate_limit_events(target);

-- -----------------------------------------------------------------------------
-- Triggers for automatic timestamp updates
-- -----------------------------------------------------------------------------
CREATE TRIGGER IF NOT EXISTS update_sessions_timestamp 
AFTER UPDATE ON sessions
BEGIN
    UPDATE sessions SET last_used = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

-- -----------------------------------------------------------------------------
-- Compression settings for PostgreSQL (if used)
-- -----------------------------------------------------------------------------
-- ALTER TABLE requests SET (toast_tuple_target = 8160);
-- ALTER TABLE responses SET (toast_tuple_target = 8160);
