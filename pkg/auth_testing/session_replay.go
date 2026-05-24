package auth_testing

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/http/cookiejar"
    "strings"
    "time"
)

// SessionReplay replays HTTP requests with session management
type SessionReplay struct {
    client      *http.Client
    sessions    map[string]*Session
    history     []RequestResponse
    replayMode  string // "single", "multi", "sequential"
}

// RequestResponse stores a request-response pair
type RequestResponse struct {
    ID         string                 `json:"id"`
    Timestamp  time.Time              `json:"timestamp"`
    Request    RequestInfo            `json:"request"`
    Response   ResponseInfo           `json:"response"`
    SessionID  string                 `json:"session_id"`
}

// RequestInfo contains request details
type RequestInfo struct {
    Method  string              `json:"method"`
    URL     string              `json:"url"`
    Headers map[string]string   `json:"headers"`
    Body    string              `json:"body"`
}

// ResponseInfo contains response details
type ResponseInfo struct {
    StatusCode int                 `json:"status_code"`
    Headers    map[string]string   `json:"headers"`
    Body       string              `json:"body"`
    Length     int                 `json:"length"`
}

// NewSessionReplay creates a new session replay manager
func NewSessionReplay() *SessionReplay {
    jar, _ := cookiejar.New(nil)
    return &SessionReplay{
        client: &http.Client{
            Jar:     jar,
            Timeout: 30 * time.Second,
        },
        sessions: make(map[string]*Session),
        history:  make([]RequestResponse, 0),
    }
}

// AddSession adds a session to the replay manager
func (sr *SessionReplay) AddSession(session *Session) {
    sr.sessions[session.ID] = session
}

// RecordRequest records a request-response pair
func (sr *SessionReplay) RecordRequest(req *http.Request, resp *http.Response, sessionID string) (*RequestResponse, error) {
    // Read request body
    reqBody := ""
    if req.Body != nil {
        bodyBytes, err := io.ReadAll(req.Body)
        if err != nil {
            return nil, err
        }
        reqBody = string(bodyBytes)
        req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
    }
    
    // Read response body
    respBody := ""
    if resp.Body != nil {
        bodyBytes, err := io.ReadAll(resp.Body)
        if err != nil {
            return nil, err
        }
        respBody = string(bodyBytes)
        resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))
    }
    
    // Extract headers
    reqHeaders := make(map[string]string)
    for k, v := range req.Header {
        if len(v) > 0 {
            reqHeaders[k] = v[0]
        }
    }
    
    respHeaders := make(map[string]string)
    for k, v := range resp.Header {
        if len(v) > 0 {
            respHeaders[k] = v[0]
        }
    }
    
    rr := RequestResponse{
        ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
        Timestamp: time.Now(),
        Request: RequestInfo{
            Method:  req.Method,
            URL:     req.URL.String(),
            Headers: reqHeaders,
            Body:    reqBody,
        },
        Response: ResponseInfo{
            StatusCode: resp.StatusCode,
            Headers:    respHeaders,
            Body:       respBody,
            Length:     len(respBody),
        },
        SessionID: sessionID,
    }
    
    sr.history = append(sr.history, rr)
    return &rr, nil
}

// ReplayRequest replays a previously recorded request
func (sr *SessionReplay) ReplayRequest(requestID string, sessionID string) (*RequestResponse, error) {
    // Find the request in history
    var original *RequestResponse
    for _, rr := range sr.history {
        if rr.ID == requestID {
            original = &rr
            break
        }
    }
    
    if original == nil {
        return nil, fmt.Errorf("request not found: %s", requestID)
    }
    
    // Get the session
    session, ok := sr.sessions[sessionID]
    if !ok {
        return nil, fmt.Errorf("session not found: %s", sessionID)
    }
    
    // Create new request
    bodyReader := strings.NewReader(original.Request.Body)
    req, err := http.NewRequest(original.Request.Method, original.Request.URL, bodyReader)
    if err != nil {
        return nil, err
    }
    
    // Add headers
    for k, v := range original.Request.Headers {
        req.Header.Set(k, v)
    }
    
    // Add session cookies
    for _, cookie := range session.Cookies {
        req.AddCookie(cookie)
    }
    
    // Execute request
    resp, err := sr.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    // Record the new request-response
    return sr.RecordRequest(req, resp, sessionID)
}

// ReplayWithDifferentRole replays request with a different user role
func (sr *SessionReplay) ReplayWithDifferentRole(requestID string, targetRole string) (*RequestResponse, error) {
    // Find a session with the target role
    var targetSession *Session
    for _, session := range sr.sessions {
        if session.UserRole == targetRole {
            targetSession = session
            break
        }
    }
    
    if targetSession == nil {
        return nil, fmt.Errorf("no session found with role: %s", targetRole)
    }
    
    return sr.ReplayRequest(requestID, targetSession.ID)
}

// ReplayAllWithDifferentRoles replays all requests with different roles
func (sr *SessionReplay) ReplayAllWithDifferentRoles(roles []string) map[string][]RequestResponse {
    results := make(map[string][]RequestResponse)
    
    for _, role := range roles {
        var resultsForRole []RequestResponse
        for _, rr := range sr.history {
            replayed, err := sr.ReplayWithDifferentRole(rr.ID, role)
            if err == nil {
                resultsForRole = append(resultsForRole, *replayed)
            }
        }
        results[role] = resultsForRole
    }
    
    return results
}

// GetHistory returns all recorded requests
func (sr *SessionReplay) GetHistory() []RequestResponse {
    return sr.history
}

// ClearHistory clears the request history
func (sr *SessionReplay) ClearHistory() {
    sr.history = make([]RequestResponse, 0)
}

// ExportHistory exports history to JSON
func (sr *SessionReplay) ExportHistory() ([]byte, error) {
    return json.MarshalIndent(sr.history, "", "  ")
}
