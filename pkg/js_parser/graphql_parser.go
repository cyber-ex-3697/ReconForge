package js_parser

import (
    "encoding/json"
    "regexp"
    "strings"
)

// GraphQLOperation represents a GraphQL operation
type GraphQLOperation struct {
    Name      string            `json:"name"`
    Type      string            `json:"type"` // query, mutation, subscription
    Fields    []string          `json:"fields"`
    Arguments map[string]string `json:"arguments"`
    Raw       string            `json:"raw"`
}

// GraphQLParser parses GraphQL queries and schemas
type GraphQLParser struct {
    operations []GraphQLOperation
    schema     map[string]interface{}
}

// NewGraphQLParser creates a new GraphQL parser
func NewGraphQLParser() *GraphQLParser {
    return &GraphQLParser{
        operations: make([]GraphQLOperation, 0),
        schema:     make(map[string]interface{}),
    }
}

// ParseQuery parses a GraphQL query
func (p *GraphQLParser) ParseQuery(query string) []GraphQLOperation {
    p.operations = make([]GraphQLOperation, 0)
    
    // Find query definitions
    queryPattern := regexp.MustCompile(`(query|mutation|subscription)\s+(\w+)\s*\(([^)]*)\)\s*{([^}]+)}`)
    matches := queryPattern.FindAllStringSubmatch(query, -1)
    
    for _, match := range matches {
        if len(match) >= 5 {
            op := GraphQLOperation{
                Type:      match[1],
                Name:      match[2],
                Arguments: p.parseArguments(match[3]),
                Fields:    p.parseFields(match[4]),
                Raw:       match[0],
            }
            p.operations = append(p.operations, op)
        }
    }
    
    // Also find inline queries
    inlinePattern := regexp.MustCompile(`{([a-zA-Z0-9_\s{},]+)}`)
    matches = inlinePattern.FindAllStringSubmatch(query, -1)
    for _, match := range matches {
        if len(match) > 1 {
            op := GraphQLOperation{
                Type:   "query",
                Name:   "inline",
                Fields: p.parseFields(match[1]),
                Raw:    match[0],
            }
            p.operations = append(p.operations, op)
        }
    }
    
    return p.operations
}

// parseArguments parses argument string
func (p *GraphQLParser) parseArguments(argStr string) map[string]string {
    args := make(map[string]string)
    argStr = strings.TrimSpace(argStr)
    
    if argStr == "" {
        return args
    }
    
    parts := strings.Split(argStr, ",")
    for _, part := range parts {
        part = strings.TrimSpace(part)
        if strings.Contains(part, ":") {
            kv := strings.SplitN(part, ":", 2)
            if len(kv) == 2 {
                args[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
            }
        } else if strings.Contains(part, "=") {
            kv := strings.SplitN(part, "=", 2)
            if len(kv) == 2 {
                args[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
            }
        }
    }
    
    return args
}

// parseFields parses field selection set
func (p *GraphQLParser) parseFields(fieldStr string) []string {
    fields := make([]string, 0)
    fieldStr = strings.TrimSpace(fieldStr)
    
    // Simple field extraction
    for _, field := range strings.Fields(fieldStr) {
        field = strings.TrimSpace(field)
        if field != "" && !strings.Contains(field, "{") && !strings.Contains(field, "}") {
            fields = append(fields, field)
        }
    }
    
    return fields
}

// ParseIntrospection parses GraphQL introspection response
func (p *GraphQLParser) ParseIntrospection(data []byte) (map[string]interface{}, error) {
    var result map[string]interface{}
    if err := json.Unmarshal(data, &result); err != nil {
        return nil, err
    }
    
    p.schema = result
    return result, nil
}

// ExtractTypes extracts types from introspection
func (p *GraphQLParser) ExtractTypes() []string {
    types := make([]string, 0)
    
    if data, ok := p.schema["data"]; ok {
        if dataMap, ok := data.(map[string]interface{}); ok {
            if schema, ok := dataMap["__schema"].(map[string]interface{}); ok {
                if typesList, ok := schema["types"].([]interface{}); ok {
                    for _, t := range typesList {
                        if typeMap, ok := t.(map[string]interface{}); ok {
                            if name, ok := typeMap["name"].(string); ok {
                                if !strings.HasPrefix(name, "__") {
                                    types = append(types, name)
                                }
                            }
                        }
                    }
                }
            }
        }
    }
    
    return types
}

// GetOperations returns all parsed operations
func (p *GraphQLParser) GetOperations() []GraphQLOperation {
    return p.operations
}
