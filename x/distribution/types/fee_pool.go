package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// zero fee pool
func InitialFeePool() FeePool {
	return FeePool{
		UbiPool: sdk.DecCoins{},
	}
}

// ValidateGenesis validates the fee pool for a genesis state
func (f FeePool) ValidateGenesis() error {
	if f.UbiPool.IsAnyNegative() {
		return fmt.Errorf("negative UbiPool in distribution fee pool, is %v",
			f.UbiPool)
	}

	return nil
}
