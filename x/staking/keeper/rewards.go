package keeper

import (
	"context"
	"fmt"
	"strconv"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

func ParseTokenTypeNormalized(tokenTypeStr string) (types.TokenType, error) {
	tokenTypeNum, err := strconv.Atoi(tokenTypeStr)
	if err != nil {
		return types.TokenType_LOCKED, fmt.Errorf("invalid token type %s", tokenTypeStr)
	}

	tokenType := types.TokenType(tokenTypeNum)
	switch tokenType {
	case types.TokenType_LOCKED, types.TokenType_UNLOCKED:
		// do nothing
	default:
		return types.TokenType_LOCKED, fmt.Errorf("unsupported token type %s", tokenTypeStr)
	}

	return tokenType, nil
}

func (k Keeper) GetTokenTypeRewardsMultiplier(ctx context.Context, tokenType types.TokenType) (math.LegacyDec, error) {
	fmt.Printf("fxxxxxxxxk %v %s\n", int32(tokenType), tokenType)
	tokenTypes, err := k.GetTokenTypes(ctx)
	if err != nil {
		return math.LegacyDec{}, err
	}
	fmt.Printf("fccccck %+v\n", tokenTypes.TokenTypeInfoMap)
	tokenTypeInfo, ok := tokenTypes.TokenTypeInfoMap[int32(tokenType)]
	if !ok {
		return math.LegacyDec{}, fmt.Errorf("invalid token type %v", tokenType)
	}
	return tokenTypeInfo.RewardsMultiplier, nil
}

func ParsePeriodTypeNormalized(periodTypeStr string) (types.PeriodType, error) {
	periodTypeNum, err := strconv.Atoi(periodTypeStr)
	if err != nil {
		return types.PeriodType_FLEXIBLE, fmt.Errorf("invalid period type %s", periodTypeStr)
	}

	periodType := types.PeriodType(periodTypeNum)
	switch periodType {
	case types.PeriodType_FLEXIBLE, types.PeriodType_THREE_MONTHS, types.PeriodType_ONE_YEAR:
		// do nothing
	default:
		return types.PeriodType_FLEXIBLE, fmt.Errorf("unsupported period type %s", periodTypeStr)
	}

	return periodType, nil
}

func (k Keeper) GetPeriodRewardsMultiplier(ctx context.Context, periodType types.PeriodType) (math.LegacyDec, error) {
	ps, err := k.GetPeriods(ctx)
	if err != nil {
		return math.LegacyZeroDec(), err
	}
	period, ok := ps.PeriodMap[int32(periodType)]
	if !ok {
		return math.LegacyDec{}, fmt.Errorf("invalid period type %v", periodType)
	}
	return period.RewardsMultiplier, nil
}
