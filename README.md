# 🔍 ReconForge - Enterprise Reconnaissance Framework

[![Version](https://img.shields.io/badge/version-4.0.0-blue.svg)](https://github.com/cyber-ex-3697/ReconForge)
[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Kali](https://img.shields.io/badge/Kali-Linux-557C94?logo=kalilinux&logoColor=white)](https://www.kali.org/)
[![Ubuntu](https://img.shields.io/badge/Ubuntu-22.04+-E95420?logo=ubuntu&logoColor=white)](https://ubuntu.com/)

## 🚀 What is ReconForge?

**ReconForge** is an **Enterprise-Grade Automated Reconnaissance Framework** written in Go. It automates the entire security assessment workflow including subdomain enumeration, live host detection, URL discovery, vulnerability scanning, port scanning, screenshot capture, and advanced recon modules.

## ✨ Features

| Feature | Description |
|---------|-------------|
| 🔍 **Subdomain Enumeration** | 6+ passive sources + active brute force |
| 🌐 **Live Host Detection** | HTTP/HTTPS probing with technology fingerprinting |
| 🔗 **URL Discovery** | Historical URLs + active crawling + JS extraction |
| ⚠️ **Vulnerability Scanning** | Nuclei with 5000+ templates |
| 🚪 **Port Scanning** | Naabu + Nmap integration |
| 📸 **Screenshots** | Gowitness for visual recon |
| 🔄 **Resumable Scans** | SQLite database for checkpoints |
| 🐳 **Docker Support** | Containerized deployment |
| 📊 **Professional Reports** | HTML, JSON, Markdown with charts |
| 🔌 **Plugin Architecture** | Extend functionality with custom plugins |

## 🎯 Advanced Recon Modules

- ✅ ASN Enumeration (amass intel)
- ✅ Subdomain Takeover Detection (subzy)
- ✅ Cloud Bucket Detection (AWS/GCP/Azure)
- ✅ WAF Detection (wafw00f)
- ✅ Technology Fingerprinting
- ✅ Favicon Hashing
- ✅ JS Endpoint Extraction
- ✅ API Endpoint Detection

## 📦 Installation

### One-Liner Install

```bash
curl -sSL https://raw.githubusercontent.com/cyber-ex-3697/ReconForge/main/scripts/install.sh | sudo bash


Manual Installation

# Clone repository
git clone https://github.com/cyber-ex-3697/ReconForge.git
cd ReconForge

# Install dependencies
chmod +x scripts/install.sh
sudo ./scripts/install.sh

# Build
go build -o reconforge cmd/reconforge/main.go

# Verify installation
./reconforge -version


🚀 Quick Start

Basic Scan

./reconforge -t example.com

Deep Scan (All Features)

./reconforge -t example.com --deep

Fast Scan with Custom Threads

./reconforge -t example.com -T 200

Custom Output Directory

./reconforge -t example.com -o my_scan_results


📊 Command Line Options

Usage: reconforge -t TARGET [OPTIONS]

Options:
  -t, --target     Target domain (required)
  -T, --threads    Number of threads (default: 50)
  -d, --deep       Deep scan mode (enables all features)
  -o, --output     Custom output directory
  -j, --json       JSON output only
  -c, --config     Config file path (default: config.yaml)
  --version        Show version
  -h, --help       Show help


🐳 Docker Usage

# Pull image
docker pull cyberex3697/reconforge:latest

# Run basic scan
docker run --rm cyberex3697/reconforge:latest -t example.com

# Run deep scan with volume mount
docker run --rm -v $(pwd)/output:/app/output cyberex3697/reconforge:latest -t example.com --deep


🔧 Configuration

Edit config.yaml to customize:

scan:
  threads: 50
  deep_mode: false
  timeout: 30
  retries: 3
  rate_limit: 10

api:
  chaos_key: "your-key"
  github_token: "your-token"
  shodan_key: "your-key"



⚠️ Disclaimer

This tool is for authorized security testing only. 
Unauthorized scanning of systems is illegal. Use at your own risk.

📝 License

This project is licensed under the MIT License - see LICENSE file for details.

📞 Contact

Author: UMAR RUMAN (CYBER EX STUDY)

Instagram: @CYBER_EX_3697

YouTube: CyberEX3697

GitHub: cyber-ex-3697

Buy Me A Coffee [ To my bank account ]

IBAN : PK42TMFB0000000097301736

⭐ Star History

If you find this tool useful, please give it a star! ⭐

