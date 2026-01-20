package chains

import (
	"context"
	"math/big"
)

// Chain defines the interface that all blockchain implementations must satisfy.
// Each chain (Starknet, Ethereum, Arbitrum, etc.) will implement this interface
// with their specific blockchain logic.
type Chain interface {
	// TransferTokens transfers tokens to a recipient address.
	// Returns the transaction hash on success.
	TransferTokens(ctx context.Context, recipient string, token string, amount *big.Int) (string, error)

	// GetBalance returns the balance of a token for a given address.
	GetBalance(ctx context.Context, address string, token string) (*big.Int, error)

	// WaitForTransaction waits for a transaction to be confirmed.
	WaitForTransaction(ctx context.Context, txHash string) error

	// ValidateAddress validates if an address is valid for this chain.
	ValidateAddress(address string) error

	// NormalizeAddress normalizes an address to its canonical form.
	NormalizeAddress(address string) string

	// GetSupportedTokens returns the list of tokens supported by this chain.
	GetSupportedTokens() []string

	// ValidateToken checks if a token is supported by this chain.
	ValidateToken(token string) error

	// GetExplorerURL returns the block explorer URL for a transaction.
	GetExplorerURL(txHash string) string

	// GetChainName returns the name of the chain (e.g., "starknet", "ethereum").
	GetChainName() string

	// GetNetworkName returns the network name (e.g., "sepolia", "mainnet").
	GetNetworkName() string
}

// ChainConfig holds common configuration for all chains.
// Each chain implementation can embed this and add chain-specific fields.
type ChainConfig struct {
	// Network is the network name (e.g., "sepolia", "mainnet", "goerli")
	Network string

	// RPCURL is the RPC endpoint URL
	RPCURL string

	// FaucetPrivateKey is the private key of the faucet wallet
	FaucetPrivateKey string

	// FaucetAddress is the address of the faucet wallet
	FaucetAddress string

	// DripAmounts maps token symbols to their drip amounts (as strings for precision)
	DripAmounts map[string]string
}

// AmountToWei converts a float amount to wei (10^18).
// This is a common utility used by most chains.
func AmountToWei(amount float64) *big.Int {
	weiPerToken := new(big.Float).SetInt(new(big.Int).Exp(
		big.NewInt(10),
		big.NewInt(18),
		nil,
	))

	amountFloat := new(big.Float).Mul(
		big.NewFloat(amount),
		weiPerToken,
	)

	amountInt, _ := amountFloat.Int(nil)
	return amountInt
}

// WeiToAmount converts wei to a float amount.
// This is a common utility used by most chains.
func WeiToAmount(wei *big.Int) float64 {
	weiPerToken := new(big.Float).SetInt(new(big.Int).Exp(
		big.NewInt(10),
		big.NewInt(18),
		nil,
	))

	weiFloat := new(big.Float).SetInt(wei)
	amount := new(big.Float).Quo(weiFloat, weiPerToken)

	result, _ := amount.Float64()
	return result
}
