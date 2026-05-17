package database

import (
    "database/sql"
    "time"
)

type Scan struct {
    ID        int64
    Target    string
    Phase     string
    Status    string
    Progress  int
    StartedAt time.Time
    UpdatedAt time.Time
    Checkpoint string
}

func (db *Database) CreateScanRecord(target, phase, status string) (int64, error) {
    result, err := db.conn.Exec(
        "INSERT INTO scans (target, phase, status) VALUES (?, ?, ?)",
        target, phase, status,
    )
    if err != nil {
        return 0, err
    }
    return result.LastInsertId()
}

func (db *Database) GetScan(id int64) (*Scan, error) {
    row := db.conn.QueryRow(
        "SELECT id, target, phase, status, progress, started_at, updated_at, checkpoint FROM scans WHERE id = ?",
        id,
    )
    
    var scan Scan
    err := row.Scan(&scan.ID, &scan.Target, &scan.Phase, &scan.Status, &scan.Progress, &scan.StartedAt, &scan.UpdatedAt, &scan.Checkpoint)
    if err != nil {
        return nil, err
    }
    return &scan, nil
}

func (db *Database) UpdateScanStatus(id int64, phase, status string, progress int) error {
    _, err := db.conn.Exec(
        "UPDATE scans SET phase = ?, status = ?, progress = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
        phase, status, progress, id,
    )
    return err
}

func (db *Database) UpdateScanCheckpoint(id int64, checkpoint string) error {
    _, err := db.conn.Exec(
        "UPDATE scans SET checkpoint = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
        checkpoint, id,
    )
    return err
}

func (db *Database) GetAllScans() ([]Scan, error) {
    rows, err := db.conn.Query("SELECT id, target, phase, status, progress, started_at, updated_at FROM scans ORDER BY id DESC")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var scans []Scan
    for rows.Next() {
        var s Scan
        if err := rows.Scan(&s.ID, &s.Target, &s.Phase, &s.Status, &s.Progress, &s.StartedAt, &s.UpdatedAt); err != nil {
            return nil, err
        }
        scans = append(scans, s)
    }
    return scans, nil
}

func (db *Database) DeleteScan(id int64) error {
    _, err := db.conn.Exec("DELETE FROM scans WHERE id = ?", id)
    return err
}
