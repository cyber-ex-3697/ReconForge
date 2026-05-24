package js_parser

import (
    "encoding/json"
    "fmt"
    "regexp"
    "strings"
)

// ASTNode represents an Abstract Syntax Tree node
type ASTNode struct {
    Type     string                 `json:"type"`
    Name     string                 `json:"name"`
    Value    interface{}            `json:"value"`
    Children []ASTNode              `json:"children"`
    Start    int                    `json:"start"`
    End      int                    `json:"end"`
    Raw      string                 `json:"raw"`
    Meta     map[string]interface{} `json:"meta"`
}

// ASTParser parses JavaScript into AST
type ASTParser struct {
    source      string
    ast         []ASTNode
    patterns    []*regexp.Regexp
}

// NewASTParser creates a new AST parser
func NewASTParser() *ASTParser {
    return &ASTParser{
        patterns: []*regexp.Regexp{
            regexp.MustCompile(`\bfetch\s*\(\s*['"]([^'"]+)['"]`),
            regexp.MustCompile(`\baxios\s*\.\s*(get|post|put|delete)\s*\(\s*['"]([^'"]+)['"]`),
            regexp.MustCompile(`\bXMLHttpRequest\s*\(\s*['"]([^'"]+)['"]`),
            regexp.MustCompile(`\bapi\s*\.\s*(get|post)\s*\(\s*['"]([^'"]+)['"]`),
            regexp.MustCompile(`/api/[a-zA-Z0-9/_-]+`),
            regexp.MustCompile(`/v\d/[a-zA-Z0-9/_-]+`),
            regexp.MustCompile(`/graphql`),
            regexp.MustCompile(`/rest/[a-zA-Z0-9/_-]+`),
        },
    }
}

// Parse parses JavaScript source code
func (p *ASTParser) Parse(source string) ([]ASTNode, error) {
    p.source = source
    p.ast = make([]ASTNode, 0)
    
    // Extract function calls
    p.extractFunctionCalls()
    
    // Extract object literals
    p.extractObjectLiterals()
    
    // Extract string literals that look like URLs
    p.extractURLLiterals()
    
    return p.ast, nil
}

// extractFunctionCalls finds function calls in source
func (p *ASTParser) extractFunctionCalls() {
    // Match function patterns
    patterns := []string{
        `(\w+)\s*\(([^)]*)\)`,
        `(\w+)\s*\.\s*(\w+)\s*\(([^)]*)\)`,
    }
    
    for _, pattern := range patterns {
        re := regexp.MustCompile(pattern)
        matches := re.FindAllStringSubmatch(p.source, -1)
        for _, match := range matches {
            node := ASTNode{
                Type: "FunctionCall",
                Name: match[1],
                Raw:  match[0],
            }
            if len(match) > 2 {
                node.Value = match[2]
            }
            p.ast = append(p.ast, node)
        }
    }
}

// extractObjectLiterals finds object definitions
func (p *ASTParser) extractObjectLiterals() {
    re := regexp.MustCompile(`(\w+)\s*:\s*['"]([^'"]+)['"]`)
    matches := re.FindAllStringSubmatch(p.source, -1)
    for _, match := range matches {
        p.ast = append(p.ast, ASTNode{
            Type:  "Property",
            Name:  match[1],
            Value: match[2],
            Raw:   match[0],
        })
    }
}

// extractURLLiterals extracts URL-like strings
func (p *ASTParser) extractURLLiterals() {
    for _, pattern := range p.patterns {
        matches := pattern.FindAllStringSubmatch(p.source, -1)
        for _, match := range matches {
            url := match[0]
            if len(match) > 1 && strings.HasPrefix(match[1], "/") {
                url = match[1]
            } else if len(match) > 2 && strings.HasPrefix(match[2], "/") {
                url = match[2]
            }
            p.ast = append(p.ast, ASTNode{
                Type:  "URLEndpoint",
                Name:  "url",
                Value: url,
                Raw:   match[0],
            })
        }
    }
}

// ExtractEndpoints extracts all URL endpoints from parsed AST
func (p *ASTParser) ExtractEndpoints() []string {
    endpoints := make(map[string]bool)
    for _, node := range p.ast {
        if node.Type == "URLEndpoint" {
            if val, ok := node.Value.(string); ok && val != "" {
                endpoints[val] = true
            }
        }
    }
    
    result := make([]string, 0, len(endpoints))
    for ep := range endpoints {
        result = append(result, ep)
    }
    return result
}

// GetAST returns the parsed AST
func (p *ASTParser) GetAST() []ASTNode {
    return p.ast
}

// ToJSON serializes AST to JSON
func (p *ASTParser) ToJSON() (string, error) {
    data, err := json.MarshalIndent(p.ast, "", "  ")
    if err != nil {
        return "", err
    }
    return string(data), nil
}
