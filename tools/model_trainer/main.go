package main

import (
    "bufio"
    "encoding/csv"
    "encoding/json"
    "flag"
    "fmt"
    "os"
    "strconv"
    "strings"
    "time"
)

// ModelConfig holds training configuration
type ModelConfig struct {
    ModelType    string   `json:"model_type"`
    Features     []string `json:"features"`
    Target       string   `json:"target"`
    TestSize     float64  `json:"test_size"`
    RandomState  int      `json:"random_state"`
    NEstimators  int      `json:"n_estimators"`
    MaxDepth     int      `json:"max_depth"`
}

// TrainingData represents training data sample
type TrainingData struct {
    Features map[string]float64
    Label    float64
}

func main() {
    var dataPath string
    var configPath string
    var outputPath string
    var modelType string
    
    flag.StringVar(&dataPath, "data", "", "Training data CSV file")
    flag.StringVar(&configPath, "config", "", "Model configuration JSON")
    flag.StringVar(&outputPath, "output", "model.pkl", "Output model file")
    flag.StringVar(&modelType, "type", "priority", "Model type (priority, vuln)")
    flag.Parse()
    
    if dataPath == "" {
        fmt.Println("Usage: model_trainer -data <training.csv> -type <priority|vuln>")
        fmt.Println("\nOptions:")
        fmt.Println("  -data     Training data CSV file")
        fmt.Println("  -config   Model configuration JSON")
        fmt.Println("  -output   Output model file")
        fmt.Println("  -type     Model type (priority, vuln)")
        os.Exit(1)
    }
    
    fmt.Println("=== ReconForge Model Trainer ===\n")
    
    // Load training data
    fmt.Printf("[*] Loading training data: %s\n", dataPath)
    data, err := loadTrainingData(dataPath)
    if err != nil {
        fmt.Printf("[!] Error loading data: %v\n", err)
        os.Exit(1)
    }
    fmt.Printf("[✓] Loaded %d samples\n", len(data))
    
    // Load or create config
    config := loadConfig(configPath, modelType)
    
    // Train model
    fmt.Printf("[*] Training %s model...\n", modelType)
    model, err := trainModel(data, config)
    if err != nil {
        fmt.Printf("[!] Training failed: %v\n", err)
        os.Exit(1)
    }
    
    // Save model
    fmt.Printf("[*] Saving model to: %s\n", outputPath)
    if err := saveModel(model, outputPath, config); err != nil {
        fmt.Printf("[!] Save failed: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Println("\n[✓] Training complete!")
}

func loadTrainingData(path string) ([]TrainingData, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    
    reader := csv.NewReader(file)
    records, err := reader.ReadAll()
    if err != nil {
        return nil, err
    }
    
    if len(records) < 2 {
        return nil, fmt.Errorf("no data rows found")
    }
    
    headers := records[0]
    var data []TrainingData
    
    for i := 1; i < len(records); i++ {
        sample := TrainingData{
            Features: make(map[string]float64),
        }
        
        for j, value := range records[i] {
            if j >= len(headers) {
                continue
            }
            
            floatVal, err := strconv.ParseFloat(value, 64)
            if err != nil {
                continue
            }
            
            if headers[j] == "label" || headers[j] == "target" {
                sample.Label = floatVal
            } else {
                sample.Features[headers[j]] = floatVal
            }
        }
        
        data = append(data, sample)
    }
    
    return data, nil
}

func loadConfig(path string, modelType string) *ModelConfig {
    config := &ModelConfig{
        ModelType:   modelType,
        TestSize:    0.2,
        RandomState: 42,
        NEstimators: 100,
        MaxDepth:    10,
    }
    
    if path != "" {
        data, err := os.ReadFile(path)
        if err == nil {
            json.Unmarshal(data, config)
        }
    }
    
    // Set default features based on model type
    if len(config.Features) == 0 {
        if modelType == "priority" {
            config.Features = []string{
                "entropy", "tech_count", "has_api", "has_admin",
                "has_backup", "asn_reputation", "response_time",
            }
            config.Target = "priority_score"
        } else {
            config.Features = []string{
                "entropy", "tech_vulnerable", "has_admin", "has_api",
                "param_count", "known_cve_count", "update_frequency",
            }
            config.Target = "vulnerable"
        }
    }
    
    return config
}

func trainModel(data []TrainingData, config *ModelConfig) (interface{}, error) {
    // This is a simplified training simulation
    // In production, actual ML model would be trained here
    
    fmt.Printf("  Features: %d\n", len(config.Features))
    fmt.Printf("  Estimators: %d\n", config.NEstimators)
    fmt.Printf("  Max Depth: %d\n", config.MaxDepth)
    
    // Simulate training time
    time.Sleep(2 * time.Second)
    
    // Return a mock model
    model := map[string]interface{}{
        "type":       config.ModelType,
        "features":   config.Features,
        "trained_at": time.Now(),
        "version":    "1.0.0",
    }
    
    return model, nil
}

func saveModel(model interface{}, path string, config *ModelConfig) error {
    output := map[string]interface{}{
        "model":      model,
        "config":     config,
        "exported_at": time.Now(),
    }
    
    data, err := json.MarshalIndent(output, "", "  ")
    if err != nil {
        return err
    }
    
    return os.WriteFile(path, data, 0644)
}
