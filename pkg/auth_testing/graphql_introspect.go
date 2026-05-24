package auth_testing

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strings"
    "time"
)

// GraphQLIntrospector performs GraphQL introspection
type GraphQLIntrospector struct {
    client      *http.Client
    endpoint    string
    schema      map[string]interface{}
    types       []GraphQLType
    queries     []GraphQLField
    mutations   []GraphQLField
}

// GraphQLType represents a GraphQL type
type GraphQLType struct {
    Name        string   `json:"name"`
    Kind        string   `json:"kind"`
    Description string   `json:"description"`
    Fields      []GraphQLField `json:"fields"`
}

// GraphQLField represents a GraphQL field
type GraphQLField struct {
    Name        string `json:"name"`
    Type        string `json:"type"`
    Description string `json:"description"`
    Args        []GraphQLArg `json:"args"`
}

// GraphQLArg represents a GraphQL argument
type GraphQLArg struct {
    Name string `json:"name"`
    Type string `json:"type"`
}

// NewGraphQLIntrospector creates a new GraphQL introspector
func NewGraphQLIntrospector(endpoint string) *GraphQLIntrospector {
    return &GraphQLIntrospector{
        client: &http.Client{
            Timeout: 30 * time.Second,
        },
        endpoint:  endpoint,
        types:     make([]GraphQLType, 0),
        queries:   make([]GraphQLField, 0),
        mutations: make([]GraphQLField, 0),
    }
}

// Introspect performs GraphQL introspection
func (gi *GraphQLIntrospector) Introspect() error {
    // Full introspection query
    query := `
    query {
        __schema {
            types {
                name
                kind
                description
                fields {
                    name
                    type {
                        name
                        kind
                    }
                    args {
                        name
                        type {
                            name
                            kind
                        }
                    }
                }
            }
            queryType {
                fields {
                    name
                    type {
                        name
                    }
                    args {
                        name
                        type {
                            name
                        }
                    }
                }
            }
            mutationType {
                fields {
                    name
                    type {
                        name
                    }
                    args {
                        name
                        type {
                            name
                        }
                    }
                }
            }
        }
    }`
    
    requestBody := map[string]interface{}{
        "query": query,
    }
    
    jsonBody, err := json.Marshal(requestBody)
    if err != nil {
        return err
    }
    
    resp, err := gi.client.Post(gi.endpoint, "application/json", bytes.NewReader(jsonBody))
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    
    return gi.parseResponse(body)
}

// parseResponse parses the introspection response
func (gi *GraphQLIntrospector) parseResponse(data []byte) error {
    var result map[string]interface{}
    if err := json.Unmarshal(data, &result); err != nil {
        return err
    }
    
    if data, ok := result["data"]; ok {
        if dataMap, ok := data.(map[string]interface{}); ok {
            if schema, ok := dataMap["__schema"].(map[string]interface{}); ok {
                gi.schema = schema
                gi.parseTypes(schema)
                gi.parseQueries(schema)
                gi.parseMutations(schema)
            }
        }
    }
    
    return nil
}

// parseTypes extracts types from schema
func (gi *GraphQLIntrospector) parseTypes(schema map[string]interface{}) {
    if types, ok := schema["types"].([]interface{}); ok {
        for _, t := range types {
            if typeMap, ok := t.(map[string]interface{}); ok {
                typeName, _ := typeMap["name"].(string)
                
                // Skip internal types
                if strings.HasPrefix(typeName, "__") {
                    continue
                }
                
                graphQLType := GraphQLType{
                    Name:        typeName,
                    Kind:        gi.getString(typeMap, "kind"),
                    Description: gi.getString(typeMap, "description"),
                }
                
                // Parse fields
                if fields, ok := typeMap["fields"].([]interface{}); ok {
                    for _, f := range fields {
                        if fieldMap, ok := f.(map[string]interface{}); ok {
                            field := GraphQLField{
                                Name:        gi.getString(fieldMap, "name"),
                                Description: gi.getString(fieldMap, "description"),
                            }
                            
                            // Parse type
                            if typeInfo, ok := fieldMap["type"].(map[string]interface{}); ok {
                                field.Type = gi.getString(typeInfo, "name")
                            }
                            
                            // Parse args
                            if args, ok := fieldMap["args"].([]interface{}); ok {
                                for _, a := range args {
                                    if argMap, ok := a.(map[string]interface{}); ok {
                                        arg := GraphQLArg{
                                            Name: gi.getString(argMap, "name"),
                                        }
                                        if argType, ok := argMap["type"].(map[string]interface{}); ok {
                                            arg.Type = gi.getString(argType, "name")
                                        }
                                        field.Args = append(field.Args, arg)
                                    }
                                }
                            }
                            
                            graphQLType.Fields = append(graphQLType.Fields, field)
                        }
                    }
                }
                
                gi.types = append(gi.types, graphQLType)
            }
        }
    }
}

// parseQueries extracts query fields
func (gi *GraphQLIntrospector) parseQueries(schema map[string]interface{}) {
    if queryType, ok := schema["queryType"].(map[string]interface{}); ok {
        if fields, ok := queryType["fields"].([]interface{}); ok {
            for _, f := range fields {
                if fieldMap, ok := f.(map[string]interface{}); ok {
                    field := GraphQLField{
                        Name:        gi.getString(fieldMap, "name"),
                        Description: gi.getString(fieldMap, "description"),
                    }
                    
                    if typeInfo, ok := fieldMap["type"].(map[string]interface{}); ok {
                        field.Type = gi.getString(typeInfo, "name")
                    }
                    
                    gi.queries = append(gi.queries, field)
                }
            }
        }
    }
}

// parseMutations extracts mutation fields
func (gi *GraphQLIntrospector) parseMutations(schema map[string]interface{}) {
    if mutationType, ok := schema["mutationType"].(map[string]interface{}); ok {
        if fields, ok := mutationType["fields"].([]interface{}); ok {
            for _, f := range fields {
                if fieldMap, ok := f.(map[string]interface{}); ok {
                    field := GraphQLField{
                        Name:        gi.getString(fieldMap, "name"),
                        Description: gi.getString(fieldMap, "description"),
                    }
                    
                    if typeInfo, ok := fieldMap["type"].(map[string]interface{}); ok {
                        field.Type = gi.getString(typeInfo, "name")
                    }
                    
                    gi.mutations = append(gi.mutations, field)
                }
            }
        }
    }
}

// getString safely gets a string from map
func (gi *GraphQLIntrospector) getString(m map[string]interface{}, key string) string {
    if val, ok := m[key]; ok {
        if str, ok := val.(string); ok {
            return str
        }
    }
    return ""
}

// GetTypes returns all discovered types
func (gi *GraphQLIntrospector) GetTypes() []GraphQLType {
    return gi.types
}

// GetQueries returns all query fields
func (gi *GraphQLIntrospector) GetQueries() []GraphQLField {
    return gi.queries
}

// GetMutations returns all mutation fields
func (gi *GraphQLIntrospector) GetMutations() []GraphQLField {
    return gi.mutations
}

// GenerateIntrospectionQuery generates an introspection query
func (gi *GraphQLIntrospector) GenerateIntrospectionQuery() string {
    return `{
        __schema {
            types {
                name
                kind
                description
                fields {
                    name
                    type {
                        name
                        kind
                    }
                }
            }
        }
    }`
}

// ExecuteCustomQuery executes a custom GraphQL query
func (gi *GraphQLIntrospector) ExecuteCustomQuery(query string) ([]byte, error) {
    requestBody := map[string]interface{}{
        "query": query,
    }
    
    jsonBody, err := json.Marshal(requestBody)
    if err != nil {
        return nil, err
    }
    
    resp, err := gi.client.Post(gi.endpoint, "application/json", bytes.NewReader(jsonBody))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    return io.ReadAll(resp.Body)
}
