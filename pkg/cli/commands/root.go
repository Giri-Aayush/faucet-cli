package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	apiURL  string
	network string
	verbose bool
	jsonOut bool
)

// Network to API URL mapping
var networkURLs = map[string]string{
	"starknet": "https://disgusted-melodee-aayushgiri-575fc666.koyeb.app",
	"ethereum": "https://disgusted-melodee-aayushgiri-575fc666.koyeb.app",
}

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "faucet",
	Short: "Multi-Chain Testnet Faucet CLI",
	Long: `A CLI tool to request testnet tokens from multiple blockchain networks.

Supported Networks:
  • starknet  - Starknet Sepolia (STRK, ETH)
  • ethereum  - Ethereum Sepolia (ETH)

Commands:
  request <ADDRESS> [flags]  Request testnet tokens
  quota                      Check YOUR remaining quota
  limits                     Show detailed rate limit rules
  status <ADDRESS>           Check request status
  info                       View faucet information

Examples:
  # Starknet Sepolia
  faucet request 0xYOUR_ADDRESS --network starknet              # Request STRK
  faucet request 0xYOUR_ADDRESS --network starknet --token ETH  # Request ETH
  faucet request 0xYOUR_ADDRESS --network starknet --both       # Request both

  # Ethereum Sepolia
  faucet request 0xYOUR_ADDRESS --network ethereum              # Request ETH

Rate Limits (per IP):
  • 5 requests per day
  • 1 hour cooldown per token type
  • After 5th request: 24-hour cooldown

Security:
  • Proof of Work challenge (prevents bot abuse)
  • CAPTCHA verification (human check)

Need help? Visit: https://github.com/Giri-Aayush/faucet-cli`,
	Version: "1.0.18",
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Global flags - network has no default, must be specified
	rootCmd.PersistentFlags().StringVar(&network, "network", "", "Blockchain network (required: starknet, ethereum)")
	rootCmd.PersistentFlags().StringVar(&apiURL, "api-url", "", "Override faucet API URL (optional)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().BoolVar(&jsonOut, "json", false, "Output in JSON format")

	// Add subcommands
	rootCmd.AddCommand(requestCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(limitsCmd)
	rootCmd.AddCommand(quotaCmd)
}

// ValidateNetwork checks if network is specified and valid
func ValidateNetwork() error {
	if network == "" {
		return fmt.Errorf(`--network flag is required

Please specify a network:
  --network starknet    Starknet Sepolia (STRK, ETH)
  --network ethereum    Ethereum Sepolia (ETH)

Example:
  faucet request 0xYOUR_ADDRESS --network starknet`)
	}

	validNetworks := []string{"starknet", "ethereum"}
	for _, valid := range validNetworks {
		if network == valid {
			return nil
		}
	}

	return fmt.Errorf(`invalid network: %s

Supported networks:
  starknet    Starknet Sepolia (STRK, ETH)
  ethereum    Ethereum Sepolia (ETH)`, network)
}

// GetAPIURL returns the API URL for the selected network
func GetAPIURL() string {
	// If explicit URL provided, use it
	if apiURL != "" {
		return apiURL
	}

	// Otherwise, look up by network
	if url, ok := networkURLs[network]; ok {
		return url
	}

	// Default to starknet
	return networkURLs["starknet"]
}

// GetNetwork returns the selected network
func GetNetwork() string {
	return network
}
