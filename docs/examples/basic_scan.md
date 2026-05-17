# Basic Scan Example

## Command

```bash
./reconforge -t example.com

Expected Output

╔══════════════════════════════════════════════════════════════════╗
║                    RECONFORGE v4.0                               ║
║              Enterprise Reconnaissance Framework                 ║
╚══════════════════════════════════════════════════════════════════╝

  Target: example.com
  Mode: Standard Scan
  Threads: 50
  Output: recon_example_com_20250115_120000

[INFO] PHASE 1: Subdomain Enumeration
  Running subfinder... done
  Running assetfinder... done
  Running findomain... done
[✓] Found 15 subdomains

[INFO] PHASE 2: Live Host Detection
  Running httpx... done
[✓] Found 8 live hosts

[INFO] PHASE 3: URL Discovery
  Running gau... done
[✓] Found 1250 URLs

[✓] SCAN COMPLETED SUCCESSFULLY
  Results saved in: recon_example_com_20250115_120000
  Duration: 2m30s

Results

View Live Hosts

cat recon_example_com_*/live_hosts.txt

Output:

https://example.com
https://www.example.com
https://api.example.com
https://cdn.example.com

View Subdomains

cat recon_example_com_*/subdomains.txt

View HTML Report

firefox recon_example_com_*/report.html


What Gets Scanned

Phase		What

1		All subdomains
2		HTTP/HTTPS live hosts
3		Historical URLs


Time Estimate

Target Size			Estimated Time

Small (<100 subdomains)		1-2 minutes
Medium (100-500 subdomains)	3-5 minutes
Large (>500 subdomains)		5-10 minutes

Next Steps
Try Deep Scan for thorough analysis

Configure API Keys for better results

Check Configuration Guide for tuning
