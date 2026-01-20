package utils

import (
	starknet "github.com/Giri-Aayush/starknet-faucet/chains/starknet-sepolia"
)

// ValidateStarknetAddress validates a Starknet address format.
// This is a convenience wrapper around the chain-specific validator
// for backward compatibility.
func ValidateStarknetAddress(address string) error {
	return starknet.ValidateAddress(address)
}

// NormalizeStarknetAddress normalizes a Starknet address to 66 characters.
// This is a convenience wrapper around the chain-specific normalizer
// for backward compatibility.
func NormalizeStarknetAddress(address string) string {
	return starknet.NormalizeAddress(address)
}

// ValidateToken validates a token type for Starknet.
// This is a convenience wrapper around the chain-specific validator
// for backward compatibility.
func ValidateToken(token string) error {
	return starknet.ValidateToken(token)
}
