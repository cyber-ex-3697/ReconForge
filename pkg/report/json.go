package report

import (
    "encoding/json"
    "os"
)

type JSONExporter struct {
    outputFile string
}

func NewJSONExporter(outputFile string) *JSONExporter {
    return &JSONExporter{
        outputFile: outputFile,
    }
}

func (j *JSONExporter) Export(data *ReportData) error {
    jsonData, err := json.MarshalIndent(data, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(j.outputFile, jsonData, 0644)
}

func (j *JSONExporter) ExportMinified(data *ReportData) error {
    jsonData, err := json.Marshal(data)
    if err != nil {
        return err
    }
    return os.WriteFile(j.outputFile, jsonData, 0644)
}
