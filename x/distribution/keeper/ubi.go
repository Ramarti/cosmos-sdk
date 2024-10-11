package keeper

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
)

func (k *Keeper) GetUbiPoolBalanceByDenom(ctx context.Context, denom string) (math.LegacyDec, error) {
	feePool, err := k.FeePool.Get(ctx)
	if err != nil {
		return math.LegacyDec{}, err
	}

	return feePool.UbiPool.AmountOf(denom), nil
}

func (k *Keeper) WithdrawUbiPooByDenomToModule(ctx context.Context, denom string, recipientModule string) error {
	feePool, err := k.FeePool.Get(ctx)
	if err != nil {
		return err
	}

	amt := feePool.UbiPool.AmountOf(denom).TruncateInt()
	coins := sdk.NewCoins(sdk.NewCoin(denom, amt))

	// NOTE the ubi pool isn't a module account, however its coins
	// are held in the distribution module account. Thus the ubi pool
	// must be reduced separately from the SendCoinsFromModuleToModule call
	newPool, negative := feePool.UbiPool.SafeSub(sdk.NewDecCoinsFromCoins(coins...))
	if negative {
		return types.ErrBadDistribution
	}
	feePool.UbiPool = newPool

	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, recipientModule, coins); err != nil {
		return err
	}

	return k.FeePool.Set(ctx, feePool)
}
