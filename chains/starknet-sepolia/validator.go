package starknet

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// Starknet address regex: 0x followed by up to 64 hex characters
	addressRegex = regexp.MustCompile(`^0x[0-9a-fA-F]{1,64}$`)
)

// ValidateAddress validates a Starknet address format.
func ValidateAddress(address string) error {
	if address == "" {
		return fmt.Errorf("address cannot be empty")
	}

	if !strings.HasPrefix(address, "0x") {
		return fmt.Errorf("address must start with 0x")
	}

	if !addressRegex.MatchString(address) {
		return fmt.Errorf("invalid Starknet address format")
	}

	return nil
}

// NormalizeAddress normalizes a Starknet address to 66 characters (0x + 64 hex).
func NormalizeAddress(address string) string {
	if len(address) >= 66 {
		return address
	}

	hexPart := address[2:]
	paddedHex := fmt.Sprintf("%064s", hexPart)
	paddedHex = strings.ReplaceAll(paddedHex, " ", "0")
	return "0x" + paddedHex
}

// ValidateToken validates a token type for Starknet.
func ValidateToken(token string) error {
	token = strings.ToUpper(token)
	if token != "ETH" && token != "STRK" {
		return fmt.Errorf("invalid token: must be ETH or STRK")
	}
	return nil
}
