#!/bin/bash

# =============================================================================
# ReconForge - Self-Updater Script
# =============================================================================

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}"
echo "╔══════════════════════════════════════════════════════════════════╗"
echo "║              ReconForge - Self-Updater                           ║"
echo "╚══════════════════════════════════════════════════════════════════╝"
echo -e "${NC}"

# Check current version
CURRENT_VERSION=$(./reconforge -version 2>/dev/null | grep -oP 'v\K[0-9.]+' || echo "0.0.0")
echo -e "${GREEN}[✓] Current version: v$CURRENT_VERSION${NC}"

# Check for updates (GitHub API)
echo -e "\n${BLUE}[*] Checking for updates...${NC}"
LATEST_VERSION=$(curl -s https://api.github.com/repos/reconforge/reconforge/releases/latest | grep -oP '"tag_name": "\Kv[0-9.]+' || echo "v$CURRENT_VERSION")

if [[ "$LATEST_VERSION" == "v$CURRENT_VERSION" ]]; then
    echo -e "${GREEN}[✓] Already up to date!${NC}"
    exit 0
fi

echo -e "${YELLOW}[!] New version available: $LATEST_VERSION${NC}"
echo -ne "${YELLOW}[?] Download and install? (y/N): ${NC}"
read -r confirm

if [[ "$confirm" != "y" && "$confirm" != "Y" ]]; then
    echo -e "${GREEN}[✓] Update cancelled${NC}"
    exit 0
fi

# Backup current version
echo -e "\n${BLUE}[*] Backing up current version...${NC}"
if [[ -f "./reconforge" ]]; then
    cp ./reconforge ./reconforge.bak
    echo -e "${GREEN}[✓] Backup created: reconforge.bak${NC}"
fi

# Download latest version
echo -e "\n${BLUE}[*] Downloading latest version...${NC}"
wget -q https://github.com/reconforge/reconforge/releases/latest/download/reconforge-linux-amd64 -O ./reconforge.new
chmod +x ./reconforge.new

# Verify download
if [[ -f "./reconforge.new" ]]; then
    mv ./reconforge.new ./reconforge
    echo -e "${GREEN}[✓] Updated successfully!${NC}"
    
    # Copy to system path
    if [[ -f "/usr/local/bin/reconforge" ]]; then
        sudo cp ./reconforge /usr/local/bin/
        echo -e "${GREEN}[✓] Updated in /usr/local/bin${NC}"
    fi
    if [[ -f "$HOME/.local/bin/reconforge" ]]; then
        cp ./reconforge $HOME/.local/bin/
        echo -e "${GREEN}[✓] Updated in ~/.local/bin${NC}"
    fi
else
    echo -e "${RED}[!] Download failed. Restoring backup...${NC}"
    if [[ -f "./reconforge.bak" ]]; then
        mv ./reconforge.bak ./reconforge
    fi
    exit 1
fi

# Update nuclei templates
if command -v nuclei &> /dev/null; then
    echo -e "\n${BLUE}[*] Updating nuclei templates...${NC}"
    nuclei -update-templates 2>/dev/null || true
fi

NEW_VERSION=$(./reconforge -version 2>/dev/null | grep -oP 'v\K[0-9.]+' || echo "unknown")
echo -e "\n${GREEN}╔══════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║              ReconForge UPDATED TO v$NEW_VERSION!                    ║${NC}"
echo -e "${GREEN}╚══════════════════════════════════════════════════════════════════╝${NC}"
