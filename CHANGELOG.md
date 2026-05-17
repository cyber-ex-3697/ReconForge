# Changelog

## [4.0.0] - 2024-01-15

### Added
- Complete Go rewrite from Bash
- YAML configuration system
- Worker pool with adaptive concurrency
- SQLite database for resumable scans
- API integrations (Chaos, Shodan, GitHub, Censys, SecurityTrails, VirusTotal)
- Advanced recon modules (ASN, WAF, Takeover, Cloud Buckets)
- Docker support with docker-compose
- HTML report with charts
- JSON and Markdown reports
- Structured logging
- CI/CD pipeline with GitHub Actions
- Unit tests for all modules

### Changed
- Replaced all shell parsing with JSON parsing
- Improved rate limiting and retry logic
- Better error handling and logging

### Fixed
- Memory leaks in large scans
- Rate limiting issues
- Checkpoint resume functionality
