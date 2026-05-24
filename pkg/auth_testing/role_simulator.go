package auth_testing

import (
    "fmt"
    "net/http"
    "net/http/cookiejar"
    "strings"
    "time"
)

// RoleSimulator simulates different user roles
type RoleSimulator struct {
    roles       map[string]*Role
    sessions    map[string]*Session
    client      *http.Client
}

// Role represents a user role
type Role struct {
    Name        string
    Permissions []string
    Description string
}

// NewRoleSimulator creates a new role simulator
func NewRoleSimulator() *RoleSimulator {
    jar, _ := cookiejar.New(nil)
    return &RoleSimulator{
        roles:    make(map[string]*Role),
        sessions: make(map[string]*Session),
        client: &http.Client{
            Jar:     jar,
            Timeout: 30 * time.Second,
        },
    }
}

// AddRole adds a role definition
func (rs *RoleSimulator) AddRole(role *Role) {
    rs.roles[role.Name] = role
}

// AddSession adds a session for a role
func (rs *RoleSimulator) AddSession(session *Session) {
    rs.sessions[session.ID] = session
}

// SimulateRole simulates a specific role
func (rs *RoleSimulator) SimulateRole(roleName string, url string) (map[string]interface{}, error) {
    role, ok := rs.roles[roleName]
    if !ok {
        return nil, fmt.Errorf("role not found: %s", roleName)
    }
    
    // Find a session for this role
    var session *Session
    for _, s := range rs.sessions {
        if s.UserRole == roleName {
            session = s
            break
        }
    }
    
    if session == nil {
        return nil, fmt.Errorf("no session found for role: %s", roleName)
    }
    
    // Make request with session
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }
    
    for _, cookie := range session.Cookies {
        req.AddCookie(cookie)
    }
    
    for k, v := range session.Headers {
        req.Header.Set(k, v)
    }
    
    resp, err := rs.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    result := map[string]interface{}{
        "role":        roleName,
        "url":         url,
        "status_code": resp.StatusCode,
        "permissions": role.Permissions,
    }
    
    return result, nil
}

// CompareRoles compares access between two roles
func (rs *RoleSimulator) CompareRoles(role1, role2, url string) (map[string]interface{}, error) {
    result1, err := rs.SimulateRole(role1, url)
    if err != nil {
        return nil, err
    }
    
    result2, err := rs.SimulateRole(role2, url)
    if err != nil {
        return nil, err
    }
    
    comparison := map[string]interface{}{
        "role1":      result1,
        "role2":      result2,
        "same_access": result1["status_code"] == result2["status_code"],
        "role1_can_access": result1["status_code"] == 200,
        "role2_can_access": result2["status_code"] == 200,
    }
    
    return comparison, nil
}

// SimulatePrivilegeEscalation tests for privilege escalation
func (rs *RoleSimulator) SimulatePrivilegeEscalation(lowRole, highRole, endpoint string) (bool, error) {
    // Try to access high-privilege endpoint with low-privilege role
    result, err := rs.SimulateRole(lowRole, endpoint)
    if err != nil {
        return false, err
    }
    
    statusCode := result["status_code"].(int)
    
    // If low role can access high-privilege endpoint, it's an escalation
    if statusCode == 200 {
        return true, nil
    }
    
    return false, nil
}

// GetAllRoles returns all defined roles
func (rs *RoleSimulator) GetAllRoles() []string {
    roles := make([]string, 0, len(rs.roles))
    for name := range rs.roles {
        roles = append(roles, name)
    }
    return roles
}

// CreateRoleFromSession creates a role from an existing session
func (rs *RoleSimulator) CreateRoleFromSession(sessionID string, roleName string) error {
    session, ok := rs.sessions[sessionID]
    if !ok {
        return fmt.Errorf("session not found: %s", sessionID)
    }
    
    session.UserRole = roleName
    
    rs.AddRole(&Role{
        Name:        roleName,
        Permissions: []string{"user"},
        Description: "Role created from session",
    })
    
    return nil
}
