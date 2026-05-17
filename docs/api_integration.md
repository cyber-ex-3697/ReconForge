# ReconForge - API Integration Guide

## Supported APIs

| API | Purpose | Free Tier |
|-----|---------|-----------|
| Chaos | Subdomain enumeration | 1000 requests/month |
| Shodan | Host discovery | Limited |
| GitHub | Code search | 60 requests/hour |
| Censys | Certificate search | Limited |
| SecurityTrails | DNS data | Limited |
| VirusTotal | Malware check | 4 requests/minute |

## Chaos API

### Get API Key

1. Visit https://cloud.projectdiscovery.io/
2. Sign up with GitHub
3. Copy your API key

### Configure

```yaml
api:
  chaos_key: "your-chaos-key"


Usage

./reconforge -t example.com

GitHub API

Get Token

Go to https://github.com/settings/tokens

Click "Generate new token"

Select repo and security_events scopes

Copy token

Configure

api:
  github_token: "ghp_your_token_here"

Shodan API

Get API Key

Visit https://account.shodan.io/

Create free account

Copy API key

Configure

api:
  shodan_key: "your_shodan_key"

Censys API

Get Credentials

Visit https://search.censys.io/

Sign up

Get API ID and Secret

Configure

api:
  censys_key: "your_api_id"
  censys_secret: "your_api_secret"

SecurityTrails API

Get API Key

Visit https://securitytrails.com/

Sign up for free account

Get API key

Configure

api:
  securitytrails_key: "your_api_key"VirusTotal API

Get API Key

Visit https://www.virustotal.com/

Create free account

Get API key

Configure

api:
  virustotal_key: "your_api_key"

Encrypted Key Storage

Save Keys Securely

# Keys are automatically encrypted
./reconforge --save-keys

# Load encrypted keys
./reconforge --load-keys

Key File Location

~/.config/reconforge/keys.enc

Testing API Connections

# Test all APIs
./reconforge --test-apis

# Test specific API
./reconforge --test-api chaos

Rate Limiting

APIs have rate limits. Configure in config.yaml:

api:
  rate_limit: 10  # requests per second
  retry_count: 3
  timeout: 30

Troubleshooting

API Key Not Working

# Verify key is set
cat ~/.config/reconforge/config.yaml | grep api

# Test connection
curl -H "Authorization: your-key" https://api.example.com

Rate Limited

# Reduce requests per second
./reconforge -t example.com -T 10

Quota Exceeded
Some APIs have daily limits. Check your account dashboard.
