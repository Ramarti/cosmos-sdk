package types

import (
	fmt "fmt"
	"strconv"
	"time"

	"cosmossdk.io/math"
)

const (
	FlexiblePeriodDelegationID = "0"
)

var StakingPeriods = map[PeriodType]Period{
	PeriodType_FLEXIBLE: {
		PeriodType:        PeriodType_FLEXIBLE,
		Duration:          time.Duration(0),
		RewardsMultiplier: math.LegacyNewDec(1),
	},
	PeriodType_THREE_MONTHS: {
		PeriodType:        PeriodType_THREE_MONTHS,
		Duration:          time.Hour * 24 * 30 * 3,
		RewardsMultiplier: math.LegacyNewDecWithPrec(1051, 3), // 1.051
	},
	PeriodType_ONE_YEAR: {
		PeriodType:        PeriodType_ONE_YEAR,
		Duration:          time.Hour * 24 * 365,
		RewardsMultiplier: math.LegacyNewDecWithPrec(116, 2), // 1.16
	},
	PeriodType_EIGHTEEN_MONTHS: {
		PeriodType:        PeriodType_EIGHTEEN_MONTHS,
		Duration:          time.Hour * 24 * 30 * 18,
		RewardsMultiplier: math.LegacyNewDecWithPrec(134, 2), // 1.34
	},
}

func GetStakingPeriodInfo(periodType PeriodType) (Period, bool) {
	period, ok := StakingPeriods[periodType]
	return period, ok
}

func ParsePeriodTypeNormalized(periodTypeStr string) (PeriodType, error) {
	periodTypeNum, err := strconv.Atoi(periodTypeStr)
	if err != nil {
		return PeriodType_FLEXIBLE, fmt.Errorf("invalid period type %s", periodTypeStr)
	}

	periodType := PeriodType(periodTypeNum)
	switch periodType {
	case PeriodType_FLEXIBLE, PeriodType_THREE_MONTHS, PeriodType_ONE_YEAR, PeriodType_EIGHTEEN_MONTHS:
		// do nothing
	default:
		return PeriodType_FLEXIBLE, fmt.Errorf("unsupported period type %s", periodTypeStr)
	}

	return periodType, nil
}
