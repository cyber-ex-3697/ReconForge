# ReconForge - Usage Guide

## Command Line Options

Usage: reconforge -t TARGET [OPTIONS]

Options:
-t, --target Target domain (required)
-T, --threads Number of threads (default: 50)
-d, --deep Deep scan mode
-o, --output Custom output directory
-j, --json JSON output only
-c, --config Config file path (default: config.yaml)
--version Show version
-h, --help Show help

## Basic Commands

### Standard Scan

```bash
./reconforge -t example.com

Deep Scan (All Features)

./reconforge -t example.com --deep

Fast Scan (More Threads)

./reconforge -t example.com -T 200

Custom Output Directory

./reconforge -t example.com -o my_scan_results

JSON Output Only

./reconforge -t example.com --json

Custom Config File

./reconforge -t example.com -c myconfig.yaml

Scan Profiles
Quick Scan (Fast)

./reconforge -t example.com -T 100

Standard Scan (Balanced)

./reconforge -t example.com -T 50

Stealth Scan (Slow, Avoid Detection)

./reconforge -t example.com -T 10

Viewing Reports

# HTML report (best for viewing)
firefox recon_*/report.html

# JSON report (for automation)
cat recon_*/report.json | jq .

# Text results
cat recon_*/live_hosts.txt

Resume Interrupted Scan

# Scan automatically saves checkpoints
./reconforge -r recon_*/checkpoint.json

Exit Codes


Code	Meaning
0	Success
1	General error
2	Target unreachable
3	Missing dependencies



