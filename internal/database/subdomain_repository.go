package database

type Subdomain struct {
    ID        int64
    ScanID    int64
    Subdomain string
    Resolved  bool
}

func (db *Database) AddSubdomain(scanID int64, subdomain string, resolved bool) error {
    _, err := db.conn.Exec(
        "INSERT INTO subdomains (scan_id, subdomain, resolved) VALUES (?, ?, ?)",
        scanID, subdomain, resolved,
    )
    return err
}

func (db *Database) AddSubdomainsBatch(scanID int64, subdomains []string) error {
    tx, err := db.conn.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    stmt, err := tx.Prepare("INSERT INTO subdomains (scan_id, subdomain, resolved) VALUES (?, ?, 0)")
    if err != nil {
        return err
    }
    defer stmt.Close()
    
    for _, sub := range subdomains {
        if _, err := stmt.Exec(scanID, sub); err != nil {
            return err
        }
    }
    
    return tx.Commit()
}

func (db *Database) GetSubdomains(scanID int64) ([]string, error) {
    rows, err := db.conn.Query("SELECT subdomain FROM subdomains WHERE scan_id = ?", scanID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var subdomains []string
    for rows.Next() {
        var sub string
        if err := rows.Scan(&sub); err != nil {
            return nil, err
        }
        subdomains = append(subdomains, sub)
    }
    return subdomains, nil
}

func (db *Database) MarkSubdomainResolved(scanID int64, subdomain string) error {
    _, err := db.conn.Exec(
        "UPDATE subdomains SET resolved = 1 WHERE scan_id = ? AND subdomain = ?",
        scanID, subdomain,
    )
    return err
}
