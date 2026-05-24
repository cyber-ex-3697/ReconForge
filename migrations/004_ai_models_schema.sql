-- =============================================================================
-- ReconForge - AI Models Schema
-- Version: 1.0.0
-- Description: Tables for AI model storage, training data, and predictions
-- =============================================================================

-- -----------------------------------------------------------------------------
-- ML Models table - Stores trained machine learning models
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS ml_models (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    version TEXT NOT NULL,
    model_type TEXT NOT NULL,
    description TEXT,
    model_path TEXT,
    parameters TEXT,
    accuracy REAL,
    precision REAL,
    recall REAL,
    f1_score REAL,
    trained_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT FALSE,
    metadata TEXT,
    UNIQUE(name, version)
);

-- Indexes for ml_models table
CREATE INDEX IF NOT EXISTS idx_ml_models_name ON ml_models(name);
CREATE INDEX IF NOT EXISTS idx_ml_models_is_active ON ml_models(is_active);

-- -----------------------------------------------------------------------------
-- Training Data table - Stores data used for model training
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS training_data (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    model_id INTEGER,
    feature_set TEXT NOT NULL,
    label TEXT,
    features TEXT,
    weight REAL DEFAULT 1.0,
    source TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(model_id) REFERENCES ml_models(id) ON DELETE SET NULL
);

-- Indexes for training_data table
CREATE INDEX IF NOT EXISTS idx_training_model_id ON training_data(model_id);
CREATE INDEX IF NOT EXISTS idx_training_label ON training_data(label);

-- -----------------------------------------------------------------------------
-- Predictions table - Stores prediction results
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS predictions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    model_id INTEGER,
    target TEXT NOT NULL,
    features TEXT,
    predicted_label TEXT,
    confidence REAL,
    probabilities TEXT,
    prediction_time REAL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(model_id) REFERENCES ml_models(id) ON DELETE SET NULL
);

-- Indexes for predictions table
CREATE INDEX IF NOT EXISTS idx_predictions_model_id ON predictions(model_id);
CREATE INDEX IF NOT EXISTS idx_predictions_target ON predictions(target);
CREATE INDEX IF NOT EXISTS idx_predictions_created_at ON predictions(created_at);

-- -----------------------------------------------------------------------------
-- Feature Importance table - Stores feature importance scores
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS feature_importance (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    model_id INTEGER NOT NULL,
    feature_name TEXT NOT NULL,
    importance_score REAL NOT NULL,
    rank INTEGER,
    FOREIGN KEY(model_id) REFERENCES ml_models(id) ON DELETE CASCADE
);

-- Indexes for feature_importance table
CREATE INDEX IF NOT EXISTS idx_feature_model_id ON feature_importance(model_id);

-- -----------------------------------------------------------------------------
-- Model Performance Metrics table - Tracks model performance over time
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS model_metrics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    model_id INTEGER NOT NULL,
    accuracy REAL,
    precision REAL,
    recall REAL,
    f1_score REAL,
    confusion_matrix TEXT,
    evaluation_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(model_id) REFERENCES ml_models(id) ON DELETE CASCADE
);

-- Indexes for model_metrics table
CREATE INDEX IF NOT EXISTS idx_metrics_model_id ON model_metrics(model_id);

-- -----------------------------------------------------------------------------
-- Subdomain Priority Scores table - Stores ML-based priority scores
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS priority_scores (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    scan_id INTEGER NOT NULL,
    subdomain TEXT NOT NULL,
    entropy_score REAL,
    tech_score REAL,
    asn_score REAL,
    historical_score REAL,
    behavior_score REAL,
    final_score REAL,
    priority_rank INTEGER,
    calculated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(scan_id, subdomain)
);

-- Indexes for priority_scores table
CREATE INDEX IF NOT EXISTS idx_priority_scan_id ON priority_scores(scan_id);
CREATE INDEX IF NOT EXISTS idx_priority_final_score ON priority_scores(final_score);
CREATE INDEX IF NOT EXISTS idx_priority_rank ON priority_scores(priority_rank);

-- -----------------------------------------------------------------------------
-- Vulnerability Probability table - Stores ML-based vulnerability predictions
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS vuln_probability (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    scan_id INTEGER NOT NULL,
    subdomain TEXT NOT NULL,
    vulnerability_type TEXT NOT NULL,
    probability REAL NOT NULL,
    confidence REAL,
    factors TEXT,
    predicted_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(scan_id, subdomain, vulnerability_type)
);

-- Indexes for vuln_probability table
CREATE INDEX IF NOT EXISTS idx_vuln_prob_scan_id ON vuln_probability(scan_id);
CREATE INDEX IF NOT EXISTS idx_vuln_prob_subdomain ON vuln_probability(subdomain);
CREATE INDEX IF NOT EXISTS idx_vuln_probability ON vuln_probability(probability);

-- -----------------------------------------------------------------------------
-- Feedback Loop table - Stores user feedback for model improvement
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS feedback (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    prediction_id INTEGER,
    is_correct BOOLEAN,
    actual_label TEXT,
    user_rating INTEGER,
    comments TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(prediction_id) REFERENCES predictions(id) ON DELETE SET NULL
);

-- Indexes for feedback table
CREATE INDEX IF NOT EXISTS idx_feedback_prediction_id ON feedback(prediction_id);

-- -----------------------------------------------------------------------------
-- Feature Store table - Centralized feature storage
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS feature_store (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    entity_type TEXT NOT NULL,
    entity_id TEXT NOT NULL,
    feature_name TEXT NOT NULL,
    feature_value REAL,
    feature_string TEXT,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(entity_type, entity_id, feature_name)
);

-- Indexes for feature_store table
CREATE INDEX IF NOT EXISTS idx_feature_store_entity ON feature_store(entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_feature_store_name ON feature_store(feature_name);

-- -----------------------------------------------------------------------------
-- Model Version Control table
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS model_versions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    model_id INTEGER NOT NULL,
    version TEXT NOT NULL,
    changelog TEXT,
    released_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    is_deprecated BOOLEAN DEFAULT FALSE,
    FOREIGN KEY(model_id) REFERENCES ml_models(id) ON DELETE CASCADE
);

-- Indexes for model_versions table
CREATE INDEX IF NOT EXISTS idx_model_versions_model_id ON model_versions(model_id);
CREATE INDEX IF NOT EXISTS idx_model_versions_version ON model_versions(version);

-- -----------------------------------------------------------------------------
-- Triggers for updating feature timestamps
-- -----------------------------------------------------------------------------
CREATE TRIGGER IF NOT EXISTS update_feature_store_timestamp 
AFTER UPDATE ON feature_store
BEGIN
    UPDATE feature_store SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
