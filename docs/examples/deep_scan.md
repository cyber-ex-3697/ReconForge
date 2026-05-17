# Deep Scan Example

## Command

```bash
./reconforge -t example.com --deep


Expected Output

╔══════════════════════════════════════════════════════════════════╗
║                    RECONFORGE v4.0                               ║
║              Enterprise Reconnaissance Framework                 ║
╚══════════════════════════════════════════════════════════════════╝

  Target: example.com
  Mode: Deep Scan
  Threads: 50
  Output: recon_example_com_20250115_120000

[INFO] PHASE 1: Subdomain Enumeration
  Running subfinder... done
  Running assetfinder... done
  Running findomain... done
  Running chaos API... done
[✓] Found 25 subdomains

[INFO] PHASE 2: Live Host Detection
  Running httpx... done
[✓] Found 12 live hosts

[INFO] PHASE 3: URL Discovery
  Running gau... done
  Running katana crawler... done
[✓] Found 3250 URLs

[INFO] PHASE 4: Vulnerability Assessment
  Running nuclei... done
[✓] Found 5 potential vulnerabilities

[INFO] PHASE 5: Port Scanning
  Running naabu... done
[✓] Found open ports on 3 hosts

[INFO] PHASE 6: Screenshot Capture
  Running gowitness... done
[✓] Screenshots captured: 12

[✓] SCAN COMPLETED SUCCESSFULLY
  Results saved in: recon_example_com_20250115_120000
  Duration: 15m30s


Additional Features in Deep Scan

Feature	Description

Port Scanning		Scans for open ports (top 1000)
Screenshots		Captures website screenshots
Advanced Recon		ASN, takeover detection
API Integration		Chaos, Shodan, GitHub
Deeper Crawling		Katana crawler with depth 5


Results

View Vulnerabilities

cat recon_example_com_*/vulnerabilities.txt

View Open Ports

cat recon_example_com_*/open_ports.txt

View Screenshots

ls recon_example_com_*/screenshots/

View Full HTML Report with Charts

firefox recon_example_com_*/report.html


Time Estimate

Target Size		Estimated Time

Small			5-10 minutes
Medium			15-30 minutes
Large			30-60 minutes


Resource Usage

Resource		Usage
CPU			2-4 cores
RAM			2-4 GB
Disk			100-500 


When to Use Deep Scan


Bug Bounty - Comprehensive coverage needed

Security Assessment - Thorough analysis required

Unknown Target - Don't know what to expect

Compliance - Need complete report


When Not to Use


Quick checks

Rate-limited targets

Resource-constrained environment

Large targets with timeout concerns


# Reduce threads for rate-limited targets
./reconforge -t example.com --deep -T 20

# Increase threads for faster scan
./reconforge -t example.com --deep -T 100Optimization Tips

