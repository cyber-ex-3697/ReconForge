# ReconForge - Installation Guide

## System Requirements

| Requirement | Minimum | Recommended |
|-------------|---------|-------------|
| OS | Ubuntu 20.04+, Kali Linux, Debian 11+ | Ubuntu 22.04+ |
| RAM | 4GB | 8GB+ |
| CPU | 2 cores | 4+ cores |
| Disk | 10GB | 20GB+ |
| Go | 1.21+ | 1.21+ |

## Quick Installation (One-Liner)

```bash
curl -sSL https://raw.githubusercontent.com/reconforge/reconforge/main/scripts/install.sh | sudo bash

Manual Installation

Step 1: Install Go

# Download Go
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

Step 2: Clone Repository

git clone https://github.com/cyber-ex-3697/reconforge.git
cd reconforge

Step 3: Install Dependencies

# Run installer script
chmod +x scripts/install.sh
./scripts/install.sh

Step 4: Build

go build -o reconforge cmd/reconforge/main.go

Step 5: Verify Installation

./reconforge -version

Docker Installation

# Pull image
docker pull reconforge/reconforge:latest

# Run container
docker run --rm reconforge/reconforge:latest -t example.com

From Source

git clone https://github.com/reconforge/reconforge.git
cd reconforge
make build
sudo make install

Post-Installation

echo 'export PATH="$PATH:$HOME/go/bin"' >> ~/.bashrc
source ~/.bashrc

Configure API Keys (Optional)

# Edit config file
nano ~/.config/reconforge/config.yaml

# Add your API keys
api:
  chaos_key: "your-key-here"
  github_token: "your-token-here"

Troubleshooting

"command not found"

# Add to PATH
export PATH="$PATH:$HOME/go/bin"

"missing dependencies"

# Run installer again
./scripts/install.sh

Build fails

# Clean and rebuild
go clean -modcache
go mod tidy
go build -o reconforge cmd/reconforge/main.go

Next Steps

Read Usage Guide

Configure API Integrations

Learn Plugin Development
