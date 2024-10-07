package keeper

import (
	"context"
	"fmt"
	"time"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

// UnbondingTime - The time duration for unbonding
func (k Keeper) UnbondingTime(ctx context.Context) (time.Duration, error) {
	params, err := k.GetParams(ctx)
	return params.UnbondingTime, err
}

// MaxValidators - Maximum number of validators
func (k Keeper) MaxValidators(ctx context.Context) (uint32, error) {
	params, err := k.GetParams(ctx)
	return params.MaxValidators, err
}

// MaxEntries - Maximum number of simultaneous unbonding
// delegations or redelegations (per pair/trio)
func (k Keeper) MaxEntries(ctx context.Context) (uint32, error) {
	params, err := k.GetParams(ctx)
	return params.MaxEntries, err
}

// HistoricalEntries = number of historical info entries
// to persist in store
func (k Keeper) HistoricalEntries(ctx context.Context) (uint32, error) {
	params, err := k.GetParams(ctx)
	return params.HistoricalEntries, err
}

// BondDenom - Bondable coin denomination
func (k Keeper) BondDenom(ctx context.Context) (string, error) {
	params, err := k.GetParams(ctx)
	return params.BondDenom, err
}

// PowerReduction - is the amount of staking tokens required for 1 unit of consensus-engine power.
// Currently, this returns a global variable that the app developer can tweak.
// TODO: we might turn this into an on-chain param:
// https://github.com/cosmos/cosmos-sdk/issues/8365
func (k Keeper) PowerReduction(ctx context.Context) math.Int {
	return sdk.DefaultPowerReduction
}

// MinCommissionRate - Minimum validator commission rate
func (k Keeper) MinCommissionRate(ctx context.Context) (math.LegacyDec, error) {
	params, err := k.GetParams(ctx)
	return params.MinCommissionRate, err
}

// MinDelegation - Minimum delegation amount
func (k Keeper) MinDelegation(ctx context.Context) (math.Int, error) {
	params, err := k.GetParams(ctx)
	return params.MinDelegation, err
}

// SetParams sets the x/staking module parameters.
// CONTRACT: This method performs no validation of the parameters.
func (k Keeper) SetParams(ctx context.Context, params types.Params) error {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := k.cdc.Marshal(&params)
	if err != nil {
		return err
	}
	return store.Set(types.ParamsKey, bz)
}

// GetParams gets the x/staking module parameters.
func (k Keeper) GetParams(ctx context.Context) (params types.Params, err error) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.ParamsKey)
	if err != nil {
		return params, err
	}

	if bz == nil {
		return params, nil
	}

	err = k.cdc.Unmarshal(bz, &params)
	return params, err
}

// SetPeriods sets the x/staking module periods.
// CONTRACT: This method performs no validation of the periods.
func (k Keeper) SetPeriods(ctx context.Context, periods types.Periods) error {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := k.cdc.Marshal(&periods)
	if err != nil {
		return err
	}
	return store.Set(types.PeriodsKey, bz)
}

// GetPeriods gets the x/staking module periods.
func (k Keeper) GetPeriods(ctx context.Context) (periods types.Periods, err error) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.PeriodsKey)
	if err != nil {
		return periods, err
	}

	if bz == nil {
		return periods, nil
	}

	err = k.cdc.Unmarshal(bz, &periods)
	return periods, err
}

// GetPeriod gets the x/staking module specific period based on period type.
func (k Keeper) GetPeriod(ctx context.Context, periodType types.PeriodType) (period types.Period, err error) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.PeriodsKey)
	if err != nil {
		return period, err
	}

	if bz == nil {
		return period, nil
	}

	var periods types.Periods
	if err := k.cdc.Unmarshal(bz, &periods); err != nil {
		return period, err
	}

	p, ok := periods.PeriodMap[int32(periodType)]
	if !ok {
		return period, fmt.Errorf("period not found for period type %v", periodType)
	}
	period = *p

	return period, nil
}

// SetTokenTypes sets the x/staking module token types.
// CONTRACT: This method performs no validation of the token types.
func (k Keeper) SetTokenTypes(ctx context.Context, tokenTypes types.TokenTypes) error {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := k.cdc.Marshal(&tokenTypes)
	if err != nil {
		return err
	}
	return store.Set(types.TokenTypesKey, bz)
}

// GetTokenTypes gets the x/staking module token types.
func (k Keeper) GetTokenTypes(ctx context.Context) (tokenTypes types.TokenTypes, err error) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.TokenTypesKey)
	if err != nil {
		return tokenTypes, err
	}

	if bz == nil {
		return tokenTypes, nil
	}

	err = k.cdc.Unmarshal(bz, &tokenTypes)
	return tokenTypes, err
}

// GetTokenType gets the x/staking module token type based on the token type.
func (k Keeper) GetTokenType(ctx context.Context, tokenType types.TokenType) (tokenTypeInfo types.TokenTypeInfo, err error) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.TokenTypesKey)
	if err != nil {
		return tokenTypeInfo, err
	}

	if bz == nil {
		return tokenTypeInfo, nil
	}

	var tokenTypes types.TokenTypes
	if err := k.cdc.Unmarshal(bz, &tokenTypes); err != nil {
		return tokenTypeInfo, err
	}

	t, ok := tokenTypes.TokenTypeInfoMap[int32(tokenType)]
	if !ok {
		return tokenTypeInfo, fmt.Errorf("invalid token type %v", tokenType)
	}
	tokenTypeInfo = *t

	return tokenTypeInfo, err
}
