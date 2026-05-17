#!/bin/bash

# =============================================================================
# ReconForge - Performance Benchmark Script
# =============================================================================

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}"
echo "╔══════════════════════════════════════════════════════════════════╗"
echo "║              ReconForge - Performance Benchmark                  ║"
echo "╚══════════════════════════════════════════════════════════════════╝"
echo -e "${NC}"

# Test target
TARGET="scanme.nmap.org"
echo -e "${GREEN}[*] Benchmark target: $TARGET${NC}"
echo ""

# Function to measure time
measure_time() {
    local cmd="$1"
    local name="$2"
    
    echo -ne "${YELLOW}[*] Running $name...${NC}"
    
    START=$(date +%s.%N)
    eval $cmd > /dev/null 2>&1
    END=$(date +%s.%N)
    DIFF=$(echo "$END - $START" | bc)
    
    echo -e " ${GREEN}${DIFF}s${NC}"
    echo "$name: ${DIFF}s" >> benchmark_results.txt
}

# Create results file
echo "ReconForge Benchmark Results - $(date)" > benchmark_results.txt
echo "Target: $TARGET" >> benchmark_results.txt
echo "" >> benchmark_results.txt

# System info
echo -e "\n${BLUE}[*] System Information:${NC}"
echo -e "  OS: $(uname -s)"
echo -e "  Kernel: $(uname -r)"
echo -e "  Architecture: $(uname -m)"
echo -e "  CPU: $(grep -c processor /proc/cpuinfo) cores"
echo -e "  RAM: $(free -h | grep Mem | awk '{print $2}')"

echo "OS: $(uname -s)" >> benchmark_results.txt
echo "Kernel: $(uname -r)" >> benchmark_results.txt
echo "Architecture: $(uname -m)" >> benchmark_results.txt
echo "CPU Cores: $(grep -c processor /proc/cpuinfo)" >> benchmark_results.txt
echo "RAM: $(free -h | grep Mem | awk '{print $2}')" >> benchmark_results.txt
echo "" >> benchmark_results.txt

# Benchmark 1: Subdomain enumeration
echo -e "\n${BLUE}[1] Subdomain Enumeration Benchmark${NC}"
measure_time "./reconforge -t $TARGET -T 50" "Subdomain Enumeration"

# Benchmark 2: Live host detection
echo -e "\n${BLUE}[2] Live Host Detection Benchmark${NC}"
# First get subdomains
./reconforge -t $TARGET -T 50 > /dev/null 2>&1
SUBDOMAINS_FILE=$(ls -t recon_* 2>/dev/null | head -1)
if [[ -n "$SUBDOMAINS_FILE" ]]; then
    measure_time "httpx -l $SUBDOMAINS_FILE/subdomains.txt -silent -threads 50" "Live Host Detection"
fi

# Benchmark 3: URL discovery
echo -e "\n${BLUE}[3] URL Discovery Benchmark${NC}"
measure_time "gau --subs $TARGET" "URL Discovery"

# Benchmark 4: Full scan
echo -e "\n${BLUE}[4] Full Scan Benchmark${NC}"
measure_time "./reconforge -t $TARGET --deep -T 50" "Full Scan"

# Memory usage
echo -e "\n${BLUE}[5] Memory Usage${NC}"
MEMORY=$(ps aux | grep reconforge | grep -v grep | awk '{sum+=$6} END {print sum/1024 " MB"}' 2>/dev/null || echo "N/A")
echo -e "  Memory used: $MEMORY"
echo "Memory Usage: $MEMORY" >> benchmark_results.txt

# Results summary
echo -e "\n${GREEN}╔══════════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║                    BENCHMARK COMPLETE!                            ║${NC}"
echo -e "${GREEN}╚══════════════════════════════════════════════════════════════════╝${NC}"
echo -e "\n${YELLOW}[*] Results saved to: benchmark_results.txt${NC}"
cat benchmark_results.txt
