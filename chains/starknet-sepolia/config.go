package starknet

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Giri-Aayush/starknet-faucet/internal/config"
	"github.com/joho/godotenv"
)

// Config holds Starknet-specific configuration.
// Distribution settings come from local config.json, secrets from .env
type Config struct {
	// Network is the Starknet network (sepolia, mainnet)
	Network string

	// RPCURL is the Starknet RPC endpoint URL (from .env)
	RPCURL string

	// FaucetPrivateKey is the private key of the faucet wallet (from .env)
	FaucetPrivateKey string

	// FaucetAddress is the address of the faucet wallet (from .env)
	FaucetAddress string

	// Token configuration (from local config.json)
	Tokens map[string]config.TokenConfig

	// MinBalanceProtectPct stops distributing when balance drops to this percentage
	MinBalanceProtectPct int

	// ExplorerURL for transaction links
	ExplorerURL string
}

// getChainDir returns the directory where this chain's config.json is located
func getChainDir() string {
	// First, try relative to current working directory
	if _, err := os.Stat("chains/starknet-sepolia/config.json"); err == nil {
		return "chains/starknet-sepolia"
	}

	// Try relative to executable
	if execPath, err := os.Executable(); err == nil {
		execDir := filepath.Dir(execPath)
		chainDir := filepath.Join(execDir, "chains", "starknet-sepolia")
		if _, err := os.Stat(filepath.Join(chainDir, "config.json")); err == nil {
			return chainDir
		}
	}

	// Fallback: try using runtime.Caller for development
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		return filepath.Dir(filename)
	}

	// Final fallback
	return "chains/starknet-sepolia"
}

// LoadConfig loads Starknet configuration from local config.json and .env (secrets)
func LoadConfig() (*Config, error) {
	// Load .env for secrets
	_ = godotenv.Load()

	// Load chain-specific config from this directory's config.json
	chainDir := getChainDir()
	chainConfig, err := config.LoadChainConfig(chainDir)
	if err != nil {
		return nil, fmt.Errorf("failed to load starknet config.json: %w", err)
	}

	// Load secrets from environment
	rpcURL := os.Getenv("STARKNET_RPC_URL")
	if rpcURL == "" {
		return nil, fmt.Errorf("STARKNET_RPC_URL is required in .env")
	}

	privateKey := os.Getenv("STARKNET_PRIVATE_KEY")
	if privateKey == "" {
		return nil, fmt.Errorf("STARKNET_PRIVATE_KEY is required in .env")
	}

	address := os.Getenv("STARKNET_ADDRESS")
	if address == "" {
		return nil, fmt.Errorf("STARKNET_ADDRESS is required in .env")
	}

	// Get network/chain ID
	network := "sepolia"
	if chainConfig.ChainID != nil {
		if s, ok := chainConfig.ChainID.(string); ok {
			network = s
		}
	}

	cfg := &Config{
		Network:              network,
		RPCURL:               rpcURL,
		FaucetPrivateKey:     privateKey,
		FaucetAddress:        address,
		Tokens:               chainConfig.Tokens,
		MinBalanceProtectPct: chainConfig.MinBalanceProtectPct,
		ExplorerURL:          chainConfig.ExplorerURL,
	}

	return cfg, nil
}

// GetDripAmount returns the drip amount for a given token
func (c *Config) GetDripAmount(token string) string {
	if tc, ok := c.Tokens[token]; ok {
		return tc.DripAmount
	}
	return "0"
}

// GetTokenAddress returns the contract address for a token
func (c *Config) GetTokenAddress(token string) string {
	if tc, ok := c.Tokens[token]; ok {
		return tc.ContractAddress
	}
	return ""
}

// GetMaxTokensPerHour returns the max hourly distribution limit for a token
func (c *Config) GetMaxTokensPerHour(token string) float64 {
	if tc, ok := c.Tokens[token]; ok {
		return tc.MaxPerHour
	}
	return 0
}

// GetMaxTokensPerDay returns the max daily distribution limit for a token
func (c *Config) GetMaxTokensPerDay(token string) float64 {
	if tc, ok := c.Tokens[token]; ok {
		return tc.MaxPerDay
	}
	return 0
}

// GetMinBalanceProtectPct returns the minimum balance protection percentage
func (c *Config) GetMinBalanceProtectPct() int {
	return c.MinBalanceProtectPct
}

// GetFaucetAddress returns the faucet wallet address
func (c *Config) GetFaucetAddress() string {
	return c.FaucetAddress
}

// GetExplorerURL returns the block explorer URL for transactions
func (c *Config) GetExplorerURL() string {
	return c.ExplorerURL
}
