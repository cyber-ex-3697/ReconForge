package database

import (
    "database/sql"
    "fmt"
)

type Migration struct {
    Version int
    Up      string
    Down    string
}

var migrations = []Migration{
    {
        Version: 1,
        Up: `CREATE TABLE IF NOT EXISTS scans (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            target TEXT NOT NULL,
            phase TEXT NOT NULL,
            status TEXT NOT NULL,
            progress INTEGER DEFAULT 0,
            started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            checkpoint TEXT
        );`,
        Down: `DROP TABLE IF EXISTS scans;`,
    },
    {
        Version: 2,
        Up: `CREATE TABLE IF NOT EXISTS subdomains (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            scan_id INTEGER,
            subdomain TEXT NOT NULL,
            resolved BOOLEAN DEFAULT FALSE,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY(scan_id) REFERENCES scans(id) ON DELETE CASCADE
        );`,
        Down: `DROP TABLE IF EXISTS subdomains;`,
    },
    {
        Version: 3,
        Up: `CREATE TABLE IF NOT EXISTS vulnerabilities (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            scan_id INTEGER,
            url TEXT,
            template_id TEXT,
            severity TEXT,
            name TEXT,
            matched_at TEXT,
            cvss REAL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY(scan_id) REFERENCES scans(id) ON DELETE CASCADE
        );`,
        Down: `DROP TABLE IF EXISTS vulnerabilities;`,
    },
    {
        Version: 4,
        Up: `CREATE INDEX IF NOT EXISTS idx_subdomains_scan_id ON subdomains(scan_id);
              CREATE INDEX IF NOT EXISTS idx_vulnerabilities_scan_id ON vulnerabilities(scan_id);
              CREATE INDEX IF NOT EXISTS idx_vulnerabilities_severity ON vulnerabilities(severity);`,
        Down: `DROP INDEX IF EXISTS idx_subdomains_scan_id;
               DROP INDEX IF EXISTS idx_vulnerabilities_scan_id;
               DROP INDEX IF EXISTS idx_vulnerabilities_severity;`,
    },
}

func (db *Database) RunMigrations() error {
    // Create migrations table if not exists
    _, err := db.conn.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
        version INTEGER PRIMARY KEY,
        applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`)
    if err != nil {
        return err
    }

    // Get current version
    var currentVersion int
    row := db.conn.QueryRow(`SELECT COALESCE(MAX(version), 0) FROM schema_migrations`)
    row.Scan(&currentVersion)

    // Run pending migrations
    for _, m := range migrations {
        if m.Version > currentVersion {
            fmt.Printf("Running migration %d...\n", m.Version)
            if _, err := db.conn.Exec(m.Up); err != nil {
                return fmt.Errorf("migration %d failed: %v", m.Version, err)
            }
            if _, err := db.conn.Exec(`INSERT INTO schema_migrations (version) VALUES (?)`, m.Version); err != nil {
                return err
            }
        }
    }
    return nil
}

func (db *Database) Rollback() error {
    var currentVersion int
    row := db.conn.QueryRow(`SELECT COALESCE(MAX(version), 0) FROM schema_migrations`)
    row.Scan(&currentVersion)

    for i := len(migrations) - 1; i >= 0; i-- {
        if migrations[i].Version == currentVersion {
            fmt.Printf("Rolling back migration %d...\n", currentVersion)
            if _, err := db.conn.Exec(migrations[i].Down); err != nil {
                return err
            }
            if _, err := db.conn.Exec(`DELETE FROM schema_migrations WHERE version = ?`, currentVersion); err != nil {
                return err
            }
            break
        }
    }
    return nil
}
