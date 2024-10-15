package keeper

import (
	"context"

	abci "github.com/cometbft/cometbft/abci/types"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// AllocateTokens performs reward and fee distribution to all validators based
// on the F1 fee distribution specification.
func (k Keeper) AllocateTokens(ctx context.Context, bondedVotes []abci.VoteInfo) error {
	// fetch and clear the collected fees for distribution, since this is
	// called in BeginBlock, collected fees will be from the previous block
	// (and distributed to the previous proposer)
	feeCollector := k.authKeeper.GetModuleAccount(ctx, k.feeCollectorName)
	feesCollectedInt := k.bankKeeper.GetAllBalances(ctx, feeCollector.GetAddress())
	feesCollected := sdk.NewDecCoinsFromCoins(feesCollectedInt...)

	// transfer collected fees to the distribution module account
	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, k.feeCollectorName, types.ModuleName, feesCollectedInt)
	if err != nil {
		return err
	}

	// temporary workaround to keep CanWithdrawInvariant happy
	// general discussions here: https://github.com/cosmos/cosmos-sdk/issues/2906#issuecomment-441867634
	feePool, err := k.FeePool.Get(ctx)
	if err != nil {
		return err
	}

	validators := make([]stakingtypes.ValidatorI, len(bondedVotes))
	var (
		totalPreviousPower int64
		totalRewardsShares = math.LegacyZeroDec()
	)
	for i := range bondedVotes {
		totalPreviousPower += bondedVotes[i].Validator.Power

		val, err := k.stakingKeeper.ValidatorByConsAddr(ctx, bondedVotes[i].Validator.Address)
		if err != nil {
			return err
		}

		validators[i] = val
		totalRewardsShares = totalRewardsShares.Add(val.GetDelegatorRewardsShares())
	}

	if totalPreviousPower == 0 {
		feePool.Ubi = feePool.Ubi.Add(feesCollected...)
		return k.FeePool.Set(ctx, feePool)
	}

	// calculate fraction allocated to validators
	remaining := feesCollected
	ubi, err := k.GetUbi(ctx)
	if err != nil {
		return err
	}

	voteMultiplier := math.LegacyOneDec().Sub(ubi)
	feeMultiplier := feesCollected.MulDecTruncate(voteMultiplier)

	// allocate tokens proportionally to voting power
	//
	// TODO: Consider parallelizing later
	//
	// Ref: https://github.com/cosmos/cosmos-sdk/pull/3099#discussion_r246276376
	for i := range validators {
		// TODO: Consider micro-slashing for missing votes.
		//
		// Ref: https://github.com/cosmos/cosmos-sdk/issues/2525#issuecomment-430838701
		powerFraction := validators[i].GetDelegatorRewardsShares().QuoTruncate(totalRewardsShares)
		reward := feeMultiplier.MulDecTruncate(powerFraction)

		err = k.AllocateTokensToValidator(ctx, validators[i], reward)
		if err != nil {
			return err
		}

		remaining = remaining.Sub(reward)
	}

	// allocate community funding
	feePool.Ubi = feePool.Ubi.Add(remaining...)
	return k.FeePool.Set(ctx, feePool)
}

// AllocateTokensToValidator allocate tokens to a particular validator,
// splitting according to commission.
func (k Keeper) AllocateTokensToValidator(ctx context.Context, val stakingtypes.ValidatorI, tokens sdk.DecCoins) error {
	// split tokens between validator and delegators according to commission
	commission := tokens.MulDec(val.GetCommission())
	shared := tokens.Sub(commission)

	valBz, err := k.stakingKeeper.ValidatorAddressCodec().StringToBytes(val.GetOperator())
	if err != nil {
		return err
	}

	// update current commission
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCommission,
			sdk.NewAttribute(sdk.AttributeKeyAmount, commission.String()),
			sdk.NewAttribute(types.AttributeKeyValidator, val.GetOperator()),
		),
	)
	currentCommission, err := k.GetValidatorAccumulatedCommission(ctx, valBz)
	if err != nil {
		return err
	}

	currentCommission.Commission = currentCommission.Commission.Add(commission...)
	err = k.SetValidatorAccumulatedCommission(ctx, valBz, currentCommission)
	if err != nil {
		return err
	}

	// update current rewards
	currentRewards, err := k.GetValidatorCurrentRewards(ctx, valBz)
	if err != nil {
		return err
	}

	currentRewards.Rewards = currentRewards.Rewards.Add(shared...)
	err = k.SetValidatorCurrentRewards(ctx, valBz, currentRewards)
	if err != nil {
		return err
	}

	// update outstanding rewards
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRewards,
			sdk.NewAttribute(sdk.AttributeKeyAmount, tokens.String()),
			sdk.NewAttribute(types.AttributeKeyValidator, val.GetOperator()),
		),
	)

	outstanding, err := k.GetValidatorOutstandingRewards(ctx, valBz)
	if err != nil {
		return err
	}

	outstanding.Rewards = outstanding.Rewards.Add(tokens...)
	return k.SetValidatorOutstandingRewards(ctx, valBz, outstanding)
}
