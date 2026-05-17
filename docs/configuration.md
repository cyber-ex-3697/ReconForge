# ReconForge - Configuration Guide

## Config File Location

- Default: `./config.yaml`
- User config: `~/.config/reconforge/config.yaml`

## Configuration Structure

```yaml
version: "4.0.0"

scan:
  profile: standard
  threads: 50
  deep_mode: false
  timeout: 30
  retries: 3
  rate_limit: 10

rate:
  default: 50
  aggressive: 100
  stealth: 10

api:
  chaos_key: ""
  github_token: ""
  shodan_key: ""
  censys_key: ""
  securitytrails_key: ""
  virustotal_key: ""

tools:
  subfinder_version: "2.14.0"
  nuclei_version: "3.1.0"
  httpx_version: "1.3.8"
  wordlists:
    - "/usr/share/wordlists/subdomains.txt"
  resolvers:
    - "1.1.1.1"
    - "8.8.8.8"


Scan Profiles

Quick Profile

scan:
  profile: quick
  threads: 100
  deep_mode: false
  timeout: 10

Standard Profile (Default)

scan:
  profile: standard
  threads: 50
  deep_mode: false
  timeout: 30

Full Profile (Deep Scan)

scan:
  profile: full
  threads: 200
  deep_mode: true
  timeout: 60
  retries: 5

Rate Limiting

Default (Balanced)

rate:
  default: 50

Aggressive (Fast)

rate:
  aggressive: 100

Stealth (Avoid Detection)

rate:
  stealth: 10

API Configuration

Chaos API

api:
  chaos_key: "your-chaos-key"

Get key from: https://cloud.projectdiscovery.io/

GitHub Token

api:
  github_token: "ghp_your_token"

Create token: https://github.com/settings/tokens

Shodan API

api:
  shodan_key: "your-shodan-key"

Get key from: https://account.shodan.io/

Wordlists

Custom Wordlist

tools:
  wordlists:
    - "/path/to/custom/wordlist.txt"
    - "/usr/share/wordlists/subdomains.txt"

Common Wordlist Locations

/usr/share/wordlists/subdomains.txt

/usr/share/seclists/Discovery/DNS/subdomains-top1million-5000.txt

/usr/share/wordlists/dirb/common.txt

Resolvers

Custom DNS Resolvers

tools:
  resolvers:
    - "1.1.1.1"
    - "8.8.8.8"
    - "9.9.9.9"
    - "208.67.222.222"

Environment Variables

# Override config values
export RECONFORGE_THREADS=100
export RECONFORGE_DEEP=true
export RECONFORGE_OUTPUT=/custom/path

# API keys as environment variables
export CHAOS_KEY="your-key"
export GITHUB_TOKEN="your-token"

Profile Presets

Bug Bounty Profile

scan:
  threads: 200
  deep_mode: true
  timeout: 60
  retries: 5
rate:
  default: 30

Corporate Assessment Profile

scan:
  threads: 50
  deep_mode: true
  timeout: 30
rate:
  default: 20

Quick Assessment Profile

scan:
  threads: 100
  deep_mode: false
  timeout: 10
rate:
  default: 50


