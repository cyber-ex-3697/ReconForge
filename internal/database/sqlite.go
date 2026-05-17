package database

import (
    "database/sql"
    "time"
    _ "modernc.org/sqlite"
)

type Database struct {
    conn *sql.DB
    path string
}

func New(dbPath string) (*Database, error) {
    conn, err := sql.Open("sqlite", dbPath)
    if err != nil {
        return nil, err
    }
    
    conn.SetMaxOpenConns(1)
    conn.SetMaxIdleConns(1)
    
    db := &Database{
        conn: conn,
        path: dbPath,
    }
    
    if err := db.migrate(); err != nil {
        return nil, err
    }
    
    return db, nil
}

func (db *Database) migrate() error {
    queries := []string{
        `CREATE TABLE IF NOT EXISTS scans (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            target TEXT NOT NULL,
            phase TEXT NOT NULL,
            status TEXT NOT NULL,
            progress INTEGER DEFAULT 0,
            started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            checkpoint TEXT
        );`,
        `CREATE TABLE IF NOT EXISTS subdomains (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            scan_id INTEGER,
            subdomain TEXT NOT NULL,
            resolved BOOLEAN DEFAULT FALSE,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY(scan_id) REFERENCES scans(id)
        );`,
        `CREATE TABLE IF NOT EXISTS vulnerabilities (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            scan_id INTEGER,
            url TEXT,
            template_id TEXT,
            severity TEXT,
            name TEXT,
            matched_at TEXT,
            cvss REAL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY(scan_id) REFERENCES scans(id)
        );`,
    }
    
    for _, query := range queries {
        if _, err := db.conn.Exec(query); err != nil {
            return err
        }
    }
    return nil
}

func (db *Database) Close() error {
    return db.conn.Close()
}

func (db *Database) CreateScan(target, phase, status string) (int64, error) {
    result, err := db.conn.Exec(
        "INSERT INTO scans (target, phase, status) VALUES (?, ?, ?)",
        target, phase, status,
    )
    if err != nil {
        return 0, err
    }
    return result.LastInsertId()
}

func (db *Database) UpdateScan(id int64, phase, status string, progress int) error {
    _, err := db.conn.Exec(
        "UPDATE scans SET phase = ?, status = ?, progress = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
        phase, status, progress, id,
    )
    return err
}
