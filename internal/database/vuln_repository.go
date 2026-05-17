package database

type Vulnerability struct {
    ID         int64
    ScanID     int64
    URL        string
    TemplateID string
    Severity   string
    Name       string
    MatchedAt  string
    CVSS       float64
}

func (db *Database) AddVulnerability(vuln *Vulnerability) error {
    _, err := db.conn.Exec(
        "INSERT INTO vulnerabilities (scan_id, url, template_id, severity, name, matched_at, cvss) VALUES (?, ?, ?, ?, ?, ?, ?)",
        vuln.ScanID, vuln.URL, vuln.TemplateID, vuln.Severity, vuln.Name, vuln.MatchedAt, vuln.CVSS,
    )
    return err
}

func (db *Database) AddVulnerabilitiesBatch(scanID int64, vulns []Vulnerability) error {
    tx, err := db.conn.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    stmt, err := tx.Prepare("INSERT INTO vulnerabilities (scan_id, url, template_id, severity, name, matched_at, cvss) VALUES (?, ?, ?, ?, ?, ?, ?)")
    if err != nil {
        return err
    }
    defer stmt.Close()
    
    for _, v := range vulns {
        if _, err := stmt.Exec(scanID, v.URL, v.TemplateID, v.Severity, v.Name, v.MatchedAt, v.CVSS); err != nil {
            return err
        }
    }
    
    return tx.Commit()
}

func (db *Database) GetVulnerabilities(scanID int64) ([]Vulnerability, error) {
    rows, err := db.conn.Query("SELECT id, url, template_id, severity, name, matched_at, cvss FROM vulnerabilities WHERE scan_id = ? ORDER BY severity DESC", scanID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var vulns []Vulnerability
    for rows.Next() {
        var v Vulnerability
        if err := rows.Scan(&v.ID, &v.URL, &v.TemplateID, &v.Severity, &v.Name, &v.MatchedAt, &v.CVSS); err != nil {
            return nil, err
        }
        vulns = append(vulns, v)
    }
    return vulns, nil
}

func (db *Database) GetVulnerabilitiesBySeverity(scanID int64, severity string) ([]Vulnerability, error) {
    rows, err := db.conn.Query(
        "SELECT id, url, template_id, name, matched_at, cvss FROM vulnerabilities WHERE scan_id = ? AND severity = ?",
        scanID, severity,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var vulns []Vulnerability
    for rows.Next() {
        var v Vulnerability
        if err := rows.Scan(&v.ID, &v.URL, &v.TemplateID, &v.Name, &v.MatchedAt, &v.CVSS); err != nil {
            return nil, err
        }
        v.Severity = severity
        vulns = append(vulns, v)
    }
    return vulns, nil
}
