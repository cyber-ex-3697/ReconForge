package distributed

import (
    "time"
    "encoding/json"
)

// TaskType defines the type of scan task
type TaskType string

const (
    TaskSubdomain  TaskType = "subdomain"
    TaskLiveHost   TaskType = "livehost"
    TaskURL        TaskType = "url"
    TaskVulnerability TaskType = "vulnerability"
    TaskPort       TaskType = "portscan"
    TaskScreenshot TaskType = "screenshot"
    TaskTakeover   TaskType = "takeover"
    TaskJS         TaskType = "js_analysis"
    TaskAuth       TaskType = "auth_testing"
)

// TaskStatus defines task state
type TaskStatus string

const (
    StatusPending   TaskStatus = "pending"
    StatusRunning   TaskStatus = "running"
    StatusCompleted TaskStatus = "completed"
    StatusFailed    TaskStatus = "failed"
    StatusRetry     TaskStatus = "retry"
)

// Task represents a unit of work
type Task struct {
    ID          string            `json:"id"`
    Type        TaskType          `json:"type"`
    Target      string            `json:"target"`
    Payload     map[string]interface{} `json:"payload"`
    Status      TaskStatus        `json:"status"`
    Priority    int               `json:"priority"`
    RetryCount  int               `json:"retry_count"`
    MaxRetries  int               `json:"max_retries"`
    CreatedAt   time.Time         `json:"created_at"`
    UpdatedAt   time.Time         `json:"updated_at"`
    WorkerID    string            `json:"worker_id"`
    Result      *TaskResult       `json:"result"`
}

// TaskResult contains task execution result
type TaskResult struct {
    Success     bool                   `json:"success"`
    Data        interface{}            `json:"data"`
    Error       string                 `json:"error"`
    Duration    time.Duration          `json:"duration"`
    Metadata    map[string]interface{} `json:"metadata"`
}

// NewTask creates a new task
func NewTask(id string, taskType TaskType, target string) *Task {
    return &Task{
        ID:         id,
        Type:       taskType,
        Target:     target,
        Status:     StatusPending,
        Priority:   0,
        MaxRetries: 3,
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
        Payload:    make(map[string]interface{}),
    }
}

// SetPayload sets task payload
func (t *Task) SetPayload(key string, value interface{}) {
    t.Payload[key] = value
}

// GetPayload retrieves payload value
func (t *Task) GetPayload(key string) interface{} {
    return t.Payload[key]
}

// ToJSON serializes task to JSON
func (t *Task) ToJSON() ([]byte, error) {
    return json.Marshal(t)
}

// FromJSON deserializes task from JSON
func (t *Task) FromJSON(data []byte) error {
    return json.Unmarshal(data, t)
}

// CanRetry checks if task can be retried
func (t *Task) CanRetry() bool {
    return t.RetryCount < t.MaxRetries
}

// IncrementRetry increases retry counter
func (t *Task) IncrementRetry() {
    t.RetryCount++
    t.Status = StatusRetry
    t.UpdatedAt = time.Now()
}
