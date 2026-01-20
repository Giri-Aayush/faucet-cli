package starknet

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
)

// Client implements the chains.Chain interface for Starknet.
type Client struct {
	account     *account.Account
	provider    *rpc.Provider
	config      *Config
	tokenAddrs  map[string]*felt.Felt
}

// NewClient creates a new Starknet chain client.
func NewClient(cfg *Config) (*Client, error) {
	ctx := context.Background()

	// Initialize RPC provider
	provider, err := rpc.NewProvider(ctx, cfg.RPCURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create provider: %w", err)
	}

	// Parse private key
	privKeyBI, ok := new(big.Int).SetString(cfg.FaucetPrivateKey, 0)
	if !ok {
		return nil, fmt.Errorf("invalid private key format")
	}

	// Setup keystore
	ks := account.NewMemKeystore()
	ks.Put(cfg.FaucetAddress, privKeyBI)

	// Parse account address
	accAddress, err := utils.HexToFelt(cfg.FaucetAddress)
	if err != nil {
		return nil, fmt.Errorf("invalid account address: %w", err)
	}

	// Create account (Cairo 2 - latest version)
	accnt, err := account.NewAccount(provider, accAddress, cfg.FaucetAddress, ks, 2)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	// Parse token addresses from config.json
	tokenAddrs := make(map[string]*felt.Felt)

	ethAddrStr := cfg.GetTokenAddress("ETH")
	if ethAddrStr != "" {
		ethAddr, err := utils.HexToFelt(ethAddrStr)
		if err != nil {
			return nil, fmt.Errorf("invalid ETH token address: %w", err)
		}
		tokenAddrs["ETH"] = ethAddr
	}

	strkAddrStr := cfg.GetTokenAddress("STRK")
	if strkAddrStr != "" {
		strkAddr, err := utils.HexToFelt(strkAddrStr)
		if err != nil {
			return nil, fmt.Errorf("invalid STRK token address: %w", err)
		}
		tokenAddrs["STRK"] = strkAddr
	}

	return &Client{
		account:    accnt,
		provider:   provider,
		config:     cfg,
		tokenAddrs: tokenAddrs,
	}, nil
}

// TransferTokens transfers tokens to a recipient.
func (c *Client) TransferTokens(
	ctx context.Context,
	recipient string,
	token string,
	amount *big.Int,
) (string, error) {
	// Parse recipient address
	recipientFelt, err := utils.HexToFelt(recipient)
	if err != nil {
		return "", fmt.Errorf("invalid recipient address: %w", err)
	}

	// Get token address
	tokenAddress, ok := c.tokenAddrs[token]
	if !ok {
		return "", fmt.Errorf("unsupported token: %s", token)
	}

	// Convert amount to Cairo uint256 format (low, high)
	low := new(big.Int).And(amount, new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 128), big.NewInt(1)))
	high := new(big.Int).Rsh(amount, 128)

	lowFelt := new(felt.Felt).SetBigInt(low)
	highFelt := new(felt.Felt).SetBigInt(high)

	// Build transfer call
	call := rpc.InvokeFunctionCall{
		ContractAddress: tokenAddress,
		FunctionName:    "transfer",
		CallData: []*felt.Felt{
			recipientFelt,
			lowFelt,
			highFelt,
		},
	}

	// Build and send invoke transaction
	tx, err := c.account.BuildAndSendInvokeTxn(ctx, []rpc.InvokeFunctionCall{call}, nil)
	if err != nil {
		return "", fmt.Errorf("transaction failed: %w", err)
	}

	return tx.Hash.String(), nil
}

// GetBalance gets the token balance of an address.
func (c *Client) GetBalance(ctx context.Context, address string, token string) (*big.Int, error) {
	// Parse address
	addrFelt, err := utils.HexToFelt(address)
	if err != nil {
		return nil, fmt.Errorf("invalid address: %w", err)
	}

	// Get token address
	tokenAddress, ok := c.tokenAddrs[token]
	if !ok {
		return nil, fmt.Errorf("unsupported token: %s", token)
	}

	// Call balanceOf
	balanceSelector := utils.GetSelectorFromNameFelt("balanceOf")

	result, err := c.provider.Call(ctx, rpc.FunctionCall{
		ContractAddress:    tokenAddress,
		EntryPointSelector: balanceSelector,
		Calldata:           []*felt.Felt{addrFelt},
	}, rpc.BlockID{Tag: "latest"})

	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	if len(result) < 2 {
		return nil, fmt.Errorf("unexpected balance result length")
	}

	// Convert from uint256 (low, high) to big.Int
	low := result[0].BigInt(big.NewInt(0))
	high := result[1].BigInt(big.NewInt(0))

	balance := new(big.Int).Add(
		low,
		new(big.Int).Lsh(high, 128),
	)

	return balance, nil
}

// WaitForTransaction waits for a transaction to be accepted.
func (c *Client) WaitForTransaction(ctx context.Context, txHash string) error {
	txHashFelt, err := utils.HexToFelt(txHash)
	if err != nil {
		return fmt.Errorf("invalid tx hash: %w", err)
	}

	// Poll for transaction receipt
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			// Check transaction receipt
			receipt, err := c.provider.TransactionReceipt(ctx, txHashFelt)
			if err != nil {
				continue
			}

			// Check if transaction is accepted
			if receipt != nil {
				return nil
			}
		}
	}
}

// ValidateAddress validates a Starknet address format.
func (c *Client) ValidateAddress(address string) error {
	return ValidateAddress(address)
}

// NormalizeAddress normalizes a Starknet address to 66 characters.
func (c *Client) NormalizeAddress(address string) string {
	return NormalizeAddress(address)
}

// GetSupportedTokens returns the list of tokens supported by Starknet.
func (c *Client) GetSupportedTokens() []string {
	return []string{"ETH", "STRK"}
}

// ValidateToken checks if a token is supported.
func (c *Client) ValidateToken(token string) error {
	return ValidateToken(token)
}

// GetExplorerURL returns the block explorer URL for a transaction.
func (c *Client) GetExplorerURL(txHash string) string {
	if c.config.ExplorerURL != "" {
		return c.config.ExplorerURL + txHash
	}
	// Fallback to default
	if c.config.Network == "mainnet" {
		return fmt.Sprintf("https://voyager.online/tx/%s", txHash)
	}
	return fmt.Sprintf("https://sepolia.voyager.online/tx/%s", txHash)
}

// GetChainName returns the chain name.
func (c *Client) GetChainName() string {
	return "starknet"
}

// GetNetworkName returns the network name.
func (c *Client) GetNetworkName() string {
	return c.config.Network
}

// GetConfig returns the chain configuration.
func (c *Client) GetConfig() *Config {
	return c.config
}
