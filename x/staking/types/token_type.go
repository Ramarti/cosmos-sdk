package types

import (
	fmt "fmt"
	"strconv"

	"cosmossdk.io/math"
)

var TokenTypes = map[TokenType]TokenTypeInfo{
	TokenType_LOCKED: {
		TokenType:         TokenType_LOCKED,
		RewardsMultiplier: math.LegacyNewDecWithPrec(5, 1),
	},
	TokenType_UNLOCKED: {
		TokenType:         TokenType_UNLOCKED,
		RewardsMultiplier: math.LegacyNewDec(1),
	},
}

func GetTokenTypeInfo(tokenType TokenType) (TokenTypeInfo, bool) {
	tokenTypeInfo, ok := TokenTypes[tokenType]
	return tokenTypeInfo, ok
}

func ParseTokenTypeNormalized(tokenTypeStr string) (TokenType, error) {
	tokenTypeNum, err := strconv.Atoi(tokenTypeStr)
	if err != nil {
		return TokenType_LOCKED, fmt.Errorf("invalid token type %s", tokenTypeStr)
	}

	tokenType := TokenType(tokenTypeNum)
	switch tokenType {
	case TokenType_LOCKED, TokenType_UNLOCKED:
		// do nothing
	default:
		return TokenType_LOCKED, fmt.Errorf("unsupported token type %s", tokenTypeStr)
	}

	return tokenType, nil
}
