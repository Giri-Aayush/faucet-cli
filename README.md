# Faucet CLI

A multi-chain testnet faucet for blockchain developers. Request tokens on Starknet Sepolia and Ethereum Sepolia directly from your terminal.

[![npm version](https://img.shields.io/npm/v/faucet-cli)](https://www.npmjs.com/package/faucet-cli)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Downloads](https://img.shields.io/npm/dm/faucet-cli)](https://www.npmjs.com/package/faucet-cli)

## Installation

```bash
npm install -g faucet-cli
```

## Quick Start

```bash
# Starknet Sepolia - request STRK
faucet request 0xYOUR_STARKNET_ADDRESS --network starknet

# Starknet Sepolia - request ETH
faucet request 0xYOUR_STARKNET_ADDRESS --network starknet --token ETH

# Ethereum Sepolia - request ETH
faucet request 0xYOUR_ETH_ADDRESS --network ethereum
```

## Supported Networks

| Network | Tokens | Amount per Request | Cooldown |
|---------|--------|-------------------|----------|
| Starknet Sepolia | STRK | 2 STRK | 24 hours |
| Starknet Sepolia | ETH | 0.001 ETH | 24 hours |
| Ethereum Sepolia | ETH | 0.001 ETH | 24 hours |

## Commands

### `faucet request <address> --network <network>`

Request testnet tokens for the specified address.

```bash
# Request STRK (default for Starknet)
faucet request 0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7 --network starknet

# Request ETH on Starknet
faucet request 0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7 --network starknet --token ETH

# Request both STRK and ETH
faucet request 0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7 --network starknet --both

# Request ETH on Ethereum Sepolia
faucet request 0x742d35Cc6634C0532925a3b844Bc9e7595f8fE21 --network ethereum
```

**Options:**
| Flag | Description |
|------|-------------|
| `--network <network>` | Network to use: `starknet` or `ethereum` (required) |
| `--token <type>` | Token to request: `STRK` or `ETH` (default: `STRK` for Starknet) |
| `--both` | Request both STRK and ETH (Starknet only) |
| `--json` | Output response in JSON format |
| `--verbose, -v` | Enable detailed logging |

### `faucet status <address> --network <network>`

Check cooldown status for an address.

```bash
faucet status 0xYOUR_ADDRESS --network starknet
faucet status 0xYOUR_ADDRESS --network ethereum
```

### `faucet info --network <network>`

Display faucet configuration and current balance.

```bash
faucet info --network starknet
faucet info --network ethereum
```

### `faucet quota --network <network>`

Check your remaining request quota.

```bash
faucet quota --network starknet
```

### `faucet limits --network <network>`

Show detailed rate limit rules.

```bash
faucet limits --network starknet
```

## How It Works

1. **Request Initiated**: You submit your wallet address
2. **Human Verification**: Answer a simple math question (prevents basic bots)
3. **Proof of Work**: Your machine solves a computational challenge (SHA-256 based, difficulty 4)
4. **Rate Check**: Server validates you haven't exceeded limits
5. **Token Transfer**: Tokens are sent to your address
6. **Confirmation**: Transaction hash with explorer link is returned

The entire process takes approximately 30 seconds.

## Rate Limits

Requests are throttled to prevent abuse:

| Limit Type | Requests per Hour | Requests per Day |
|------------|------------------|------------------|
| Per IP | 10 | 20 |
| Per Address | 2 | 5 |

## Security

The faucet implements defense-in-depth:

- **Proof of Work**: SHA-256 hashcash-style challenge (difficulty 4, ~65k iterations average)
- **CAPTCHA**: Interactive arithmetic verification
- **Rate Limiting**: Redis-backed per-IP and per-address throttling
- **Challenge TTL**: PoW challenges expire after 5 minutes
- **Balance Protection**: Faucet pauses at 5% remaining balance
- **Address Validation**: EIP-55 checksum for Ethereum, format validation for Starknet

## Architecture

```
┌─────────────┐     HTTPS      ┌─────────────┐     RPC      ┌──────────────┐
│  Faucet CLI │ ─────────────► │   Backend   │ ──────────► │  Blockchain  │
│   (Go)      │                │  (Go/Koyeb) │              │   Network    │
└─────────────┘                └──────┬──────┘              └──────────────┘
                                      │
                                      ▼
                               ┌─────────────┐
                               │    Redis    │
                               │ (Rate Limit)│
                               └─────────────┘
```

**Tech Stack:**
- CLI: Go 1.23+ with Cobra
- Backend: Go HTTP server on Koyeb
- Starknet: [starknet.go](https://github.com/NethermindEth/starknet.go) v0.17.0
- Ethereum: [go-ethereum](https://github.com/ethereum/go-ethereum)
- Rate Limiting: Redis
- Explorers: Voyager (Starknet), Etherscan (Ethereum)

## Platform Support

| Platform | Architecture | Binary |
|----------|-------------|--------|
| Linux | x64 | `faucet-linux-amd64` |
| Linux | ARM64 | `faucet-linux-arm64` |
| macOS | Intel | `faucet-macos-amd64` |
| macOS | Apple Silicon | `faucet-macos-arm64` |
| Windows | x64 | `faucet-windows-amd64.exe` |

The npm package automatically downloads the correct binary for your platform.

## Development

```bash
# Clone the repository
git clone https://github.com/Giri-Aayush/faucet-cli.git
cd faucet-cli

# Build the CLI
go build -o faucet ./cmd/cli

# Run locally
./faucet --help
```

## Contributing

Contributions welcome. Please open an issue first to discuss proposed changes.

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/improvement`)
3. Commit changes (`git commit -am 'Add feature'`)
4. Push to branch (`git push origin feature/improvement`)
5. Open a Pull Request

## License

MIT License - see [LICENSE](LICENSE) for details.

## Links

- [npm Package](https://www.npmjs.com/package/faucet-cli)
- [GitHub Repository](https://github.com/Giri-Aayush/faucet-cli)
- [Report Issues](https://github.com/Giri-Aayush/faucet-cli/issues)

---

Made by [Aayush Giri](https://github.com/Giri-Aayush)
