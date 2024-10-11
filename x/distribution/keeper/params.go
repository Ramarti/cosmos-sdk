package keeper

import (
	"context"

	"cosmossdk.io/math"
)

// GetUbiPool returns the current distribution ubi pool.
func (k Keeper) GetUbiPool(ctx context.Context) (math.LegacyDec, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return math.LegacyDec{}, err
	}

	return params.UbiPool, nil
}

// GetWithdrawAddrEnabled returns the current distribution withdraw address
// enabled parameter.
func (k Keeper) GetWithdrawAddrEnabled(ctx context.Context) (enabled bool, err error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return false, err
	}

	return params.WithdrawAddrEnabled, nil
}
