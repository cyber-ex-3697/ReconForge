package report

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
)

type ScreenshotCapturer struct {
    outputDir string
    timeout   int
}

func NewScreenshotCapturer(outputDir string, timeout int) *ScreenshotCapturer {
    return &ScreenshotCapturer{
        outputDir: outputDir,
        timeout:   timeout,
    }
}

func (s *ScreenshotCapturer) Capture(url string) (string, error) {
    outputFile := filepath.Join(s.outputDir, sanitizeFilename(url)+".png")
    
    cmd := exec.Command("gowitness", "single", "-u", url, "--destination", s.outputDir, "--timeout", fmt.Sprintf("%d", s.timeout))
    err := cmd.Run()
    if err != nil {
        return "", err
    }
    
    return outputFile, nil
}

func (s *ScreenshotCapturer) CaptureBatch(urls []string) ([]string, error) {
    var files []string
    for _, url := range urls {
        file, err := s.Capture(url)
        if err == nil {
            files = append(files, file)
        }
    }
    return files, nil
}

func (s *ScreenshotCapturer) CaptureFromFile(urlFile string) error {
    cmd := exec.Command("gowitness", "file", "-f", urlFile, "--destination", s.outputDir, "--timeout", fmt.Sprintf("%d", s.timeout))
    return cmd.Run()
}

func sanitizeFilename(url string) string {
    // Remove special characters for filename
    return strings.ReplaceAll(strings.ReplaceAll(url, "https://", ""), "http://", "")
}
