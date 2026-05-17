#!/bin/bash

# =============================================================================
# ReconForge - Uninstaller Script
# =============================================================================

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}"
echo "╔══════════════════════════════════════════════════════════════════╗"
echo "║              ReconForge - Uninstaller                            ║"
echo "╚══════════════════════════════════════════════════════════════════╝"
echo -e "${NC}"

echo -e "${YELLOW}[!] This will remove ReconForge and all its components.${NC}"
echo -ne "${YELLOW}[?] Are you sure? (y/N): ${NC}"
read -r confirm

if [[ "$confirm" != "y" && "$confirm" != "Y" ]]; then
    echo -e "${GREEN}[✓] Uninstall cancelled${NC}"
    exit 0
fi

# Remove binary
echo -e "\n${BLUE}[*] Removing ReconForge binary...${NC}"
if [[ -f "/usr/local/bin/reconforge" ]]; then
    sudo rm -f /usr/local/bin/reconforge
    echo -e "${GREEN}[✓] Removed from /usr/local/bin${NC}"
fi
if [[ -f "$HOME/.local/bin/reconforge" ]]; then
    rm -f $HOME/.local/bin/reconforge
    echo -e "${GREEN}[✓] Removed from ~/.local/bin${NC}"
fi
if [[ -f "./reconforge" ]]; then
    rm -f ./reconforge
    echo -e "${GREEN}[✓] Removed from current directory${NC}"
fi

# Remove configuration
echo -e "\n${BLUE}[*] Removing configuration...${NC}"
if [[ -d "$HOME/.config/reconforge" ]]; then
    rm -rf $HOME/.config/reconforge
    echo -e "${GREEN}[✓] Removed configuration${NC}"
fi

# Remove output directories
echo -e "\n${BLUE}[*] Removing scan results...${NC}"
rm -rf ./recon_* 2>/dev/null || true
echo -e "${GREEN}[✓] Removed scan results${NC}"

# Remove logs
echo -e "\n${BLUE}[*] Removing logs...${NC}"
rm -rf ./logs 2>/dev/null || true

# Ask to remove Go tools
echo -e "\n${YELLOW}[?] Remove Go tools (subfinder, httpx, nuclei, etc.)? (y/N): ${NC}"
read -r remove_tools

if [[ "$remove_tools" == "y" || "$remove_tools" == "Y" ]]; then
    echo -e "${BLUE}[*] Removing Go tools...${NC}"
    rm -f /usr/local/bin/subfinder 2>/dev/null || true
    rm -f /usr/local/bin/httpx 2>/dev/null || true
    rm -f /usr/local/bin/nuclei 2>/dev/null || true
    rm -f /usr/local/bin/assetfinder 2>/dev/null || true
    rm -f /usr/local/bin/gau 2>/dev/null || true
    rm -f /usr/local/bin/katana 2>/dev/null || true
    rm -f /usr/local/bin/dnsx 2>/dev/null || true
    rm -f /usr/local/bin/findomain 2>/dev/null || true
    rm -f /usr/local/bin/naabu 2>/dev/null || true
    rm -f /usr/local/bin/gowitness 2>/dev/null || true
    echo -e "${GREEN}[✓] Go tools removed${NC}"
fi

echo -e "\n${GREEN}╔══════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║              ReconForge UNINSTALLED SUCCESSFULLY!                 ║${NC}"
echo -e "${GREEN}╚══════════════════════════════════════════════════════════════════╝${NC}"
