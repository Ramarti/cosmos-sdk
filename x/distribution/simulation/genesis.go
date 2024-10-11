package simulation

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// Simulation parameter constants
const (
	UbiPool         = "ubi_pool"
	WithdrawEnabled = "withdraw_enabled"
)

// GenUbiPool randomized UbiPool parameter.
func GenUbiPool(r *rand.Rand) math.LegacyDec {
	return math.LegacyNewDecWithPrec(1, 2).Add(math.LegacyNewDecWithPrec(int64(r.Intn(30)), 2))
}

// GenWithdrawEnabled returns a randomized WithdrawEnabled parameter.
func GenWithdrawEnabled(r *rand.Rand) bool {
	return r.Int63n(101) <= 95 // 95% chance of withdraws being enabled
}

// RandomizedGenState generates a random GenesisState for distribution
func RandomizedGenState(simState *module.SimulationState) {
	var ubiPool math.LegacyDec
	simState.AppParams.GetOrGenerate(UbiPool, &ubiPool, simState.Rand, func(r *rand.Rand) { ubiPool = GenUbiPool(r) })

	var withdrawEnabled bool
	simState.AppParams.GetOrGenerate(WithdrawEnabled, &withdrawEnabled, simState.Rand, func(r *rand.Rand) { withdrawEnabled = GenWithdrawEnabled(r) })

	distrGenesis := types.GenesisState{
		FeePool: types.InitialFeePool(),
		Params: types.Params{
			UbiPool:             ubiPool,
			WithdrawAddrEnabled: withdrawEnabled,
		},
	}

	bz, err := json.MarshalIndent(&distrGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated distribution parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&distrGenesis)
}
