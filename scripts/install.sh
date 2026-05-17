#!/bin/bash

# =============================================================================
# ReconForge - One-Click Installer
# =============================================================================

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}"
echo "╔══════════════════════════════════════════════════════════════════╗"
echo "║              ReconForge - One-Click Installer                    ║"
echo "╚══════════════════════════════════════════════════════════════════╝"
echo -e "${NC}"

# Check if running as root
if [[ $EUID -ne 0 ]]; then
    echo -e "${YELLOW}[!] Not running as root. Some features may be limited.${NC}"
    echo -e "${YELLOW}[!] Run with sudo for complete installation.${NC}"
    HAS_ROOT=false
else
    echo -e "${GREEN}[✓] Running with root privileges${NC}"
    HAS_ROOT=true
fi

# Detect OS
if [[ -f /etc/os-release ]]; then
    . /etc/os-release
    OS=$NAME
    VER=$VERSION_ID
    echo -e "${GREEN}[✓] Detected OS: $OS $VER${NC}"
else
    OS="Unknown"
    echo -e "${YELLOW}[!] Could not detect OS${NC}"
fi

# Update system
echo -e "\n${BLUE}[*] Updating package lists...${NC}"
if [[ "$HAS_ROOT" == true ]]; then
    apt update -y 2>/dev/null || true
    echo -e "${GREEN}[✓] System updated${NC}"
fi

# Install system dependencies
echo -e "\n${BLUE}[*] Installing system dependencies...${NC}"
if [[ "$HAS_ROOT" == true ]]; then
    apt install -y curl wget git golang-go nmap jq parallel unzip 2>/dev/null || true
    echo -e "${GREEN}[✓] System dependencies installed${NC}"
fi

# Install Go tools
echo -e "\n${BLUE}[*] Installing Go tools...${NC}"

install_go_tool() {
    local tool=$1
    local url=$2
    if ! command -v $tool &> /dev/null; then
        echo -e "  Installing $tool..."
        go install $url@latest 2>/dev/null || true
    else
        echo -e "  $tool already installed"
    fi
}

install_go_tool "subfinder" "github.com/projectdiscovery/subfinder/v2/cmd/subfinder"
install_go_tool "httpx" "github.com/projectdiscovery/httpx/cmd/httpx"
install_go_tool "nuclei" "github.com/projectdiscovery/nuclei/v3/cmd/nuclei"
install_go_tool "assetfinder" "github.com/tomnomnom/assetfinder"
install_go_tool "gau" "github.com/lc/gau/v2/cmd/gau"
install_go_tool "katana" "github.com/projectdiscovery/katana/cmd/katana"
install_go_tool "dnsx" "github.com/projectdiscovery/dnsx/cmd/dnsx"

# Copy to /usr/local/bin
if [[ "$HAS_ROOT" == true ]]; then
    cp ~/go/bin/* /usr/local/bin/ 2>/dev/null || true
fi

# Install binary tools
echo -e "\n${BLUE}[*] Installing binary tools...${NC}"

# Findomain
if ! command -v findomain &> /dev/null; then
    echo -e "  Installing findomain..."
    wget -q https://github.com/findomain/findomain/releases/latest/download/findomain-linux.zip -O /tmp/findomain.zip
    unzip -q /tmp/findomain.zip -d /tmp/
    chmod +x /tmp/findomain
    if [[ "$HAS_ROOT" == true ]]; then
        mv /tmp/findomain /usr/local/bin/
    else
        mkdir -p ~/.local/bin
        mv /tmp/findomain ~/.local/bin/
    fi
    rm /tmp/findomain.zip
fi

# Naabu
if ! command -v naabu &> /dev/null; then
    echo -e "  Installing naabu..."
    wget -q https://github.com/projectdiscovery/naabu/releases/latest/download/naabu_linux_amd64.zip -O /tmp/naabu.zip
    unzip -q /tmp/naabu.zip -d /tmp/
    chmod +x /tmp/naabu
    if [[ "$HAS_ROOT" == true ]]; then
        mv /tmp/naabu /usr/local/bin/
    else
        mv /tmp/naabu ~/.local/bin/
    fi
    rm /tmp/naabu.zip
fi

# Gowitness
if ! command -v gowitness &> /dev/null; then
    echo -e "  Installing gowitness..."
    wget -q https://github.com/sensepost/gowitness/releases/latest/download/gowitness-linux-amd64 -O /tmp/gowitness
    chmod +x /tmp/gowitness
    if [[ "$HAS_ROOT" == true ]]; then
        mv /tmp/gowitness /usr/local/bin/
    else
        mv /tmp/gowitness ~/.local/bin/
    fi
fi

# Update nuclei templates
if command -v nuclei &> /dev/null; then
    echo -e "\n${BLUE}[*] Updating nuclei templates...${NC}"
    nuclei -update-templates 2>/dev/null || true
fi

# Build ReconForge
echo -e "\n${BLUE}[*] Building ReconForge...${NC}"
go build -o reconforge cmd/reconforge/main.go 2>/dev/null || true

if [[ -f "./reconforge" ]]; then
    if [[ "$HAS_ROOT" == true ]]; then
        cp reconforge /usr/local/bin/
    else
        cp reconforge ~/.local/bin/
    fi
    echo -e "${GREEN}[✓] ReconForge built successfully${NC}"
else
    echo -e "${RED}[!] Build failed${NC}"
fi

# Setup PATH for non-root
if [[ "$HAS_ROOT" != true ]]; then
    mkdir -p ~/.local/bin
    export PATH="$PATH:~/.local/bin"
    
    if ! grep -q ".local/bin" ~/.bashrc 2>/dev/null; then
        echo 'export PATH="$PATH:$HOME/.local/bin"' >> ~/.bashrc
        echo -e "${GREEN}[✓] Added ~/.local/bin to PATH${NC}"
    fi
fi

echo -e "\n${GREEN}╔══════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║              ReconForge INSTALLED SUCCESSFULLY!                   ║${NC}"
echo -e "${GREEN}╚══════════════════════════════════════════════════════════════════╝${NC}"
echo -e "\n${YELLOW}[*] To use ReconForge:${NC}"
echo -e "  reconforge -t example.com"
echo -e "  reconforge -t example.com --deep"
echo -e "\n${YELLOW}[!] If command not found, run: source ~/.bashrc${NC}"
