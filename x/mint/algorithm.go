package mint

import (
	"crypto/sha256"
	"math"
	"strconv"
)

const (
	BlkTime               = 5
	AvgDaysPerMonth       = 30
	DayinSecond           = 24 * 3600
	MonthsPerYear         = 12
	AvgBlksPerMonth       = AvgDaysPerMonth * DayinSecond / BlkTime
	Period                = AvgBlksPerMonth
	MineablePeriodAsYears = 50
	DefaultMineableBlks   = MineablePeriodAsYears * 12 * AvgBlksPerMonth // 40 years mining
	TestMineableBlks      = 2 * 3600 / BlkTime                           // 2 hours mining
	TotalMineableBlks     = DefaultMineableBlks
	BlkRadianIntv         = 2.0 * math.Pi / float64(Period)

	ValidatorNumbers    = 7          // the number of validators
	ValidatorProvisions = float64(1) // 1 for each validator
	// Ignore validaotorProvisin because it is set as small as enough to be neglected.
	// if you want to set it bigger as your , you should care about this part again.
	// Including this part will decrease interoperability of the source code.
	// ValidatorTotalProvisions = ValidatorProvisions * ValidatorNumbers // 1 for each validator
	ValidatorTotalProvisions = 0

	// IssuerAmount = float64(1000000) // this is for test. 0 for production, 1000000 for test
	IssuerAmount = 0 // 0 for production

	FixedMineProvision       = float64(2925000000)
	MineTotalProvisions      = FixedMineProvision - ValidatorTotalProvisions - IssuerAmount // ~36,000,000 for 40 years
	AnnualProvisions         = MineTotalProvisions / float64(MineablePeriodAsYears)         // ~900000 per year
	AnnualProvisionAsSatoshi = int64(AnnualProvisions * sscq2satoshi)
	MonthProvisions          = AnnualProvisions / float64(MonthsPerYear) // ~75000 per month

	// this is for export case,that's,this is activated if there exporting accounts exist.
	UserProvisions = float64(325000000) // if not, this should be set as zero

	CurrentProvisions          = UserProvisions + ValidatorTotalProvisions + IssuerAmount // ~60,000,000 at genesis
	CurrentProvisionsAsSatoshi = int64(CurrentProvisions * sscq2satoshi)                  // ~60,000,000 at genesis
	TotalLiquid                = MineTotalProvisions + CurrentProvisions                  // 96,000,000
	TotalLiquidAsSatoshi       = int64(TotalLiquid * sscq2satoshi)                        // 96,000,000 * 100,000,000

	sscq2satoshi = 100000000 // 1 sscq = 10 ** 8 satoshi

	AvgBlkReward          = MineTotalProvisions / TotalMineableBlks
	AvgBlkRewardAsSatoshi = sscq2satoshi * AvgBlkReward
	RATIO                 = 0.5

	MAX_AMPLITUDE = AvgBlkReward
	MIN_AMPLITUDE = 0.001
	MAX_CYCLE     = 3000
	MIN_CYCLE     = 100
)

// junying-todo, 2019-12-05
// 60,000,000 + 0.144676 * height
func expectedtotalSupply(blkindex int64) int64 {
	return CurrentProvisionsAsSatoshi +
		int64(float64(blkindex)*AvgBlkRewardAsSatoshi)
}

func randomUint(seed int64) uint64 {
	hash := sha256.Sum256([]byte(strconv.FormatInt(seed, 10)))
	return uint64(hash[:1][0])
}

func randomFloat(seed int64) float64 {
	return float64(randomUint(seed)) / 256.0
}

// junying-todo, 2019-12-06
// rand(0.001,0.144676)
// 0.144676 * rand(0.0,1,0) + 0.001
func randomAmplitude(seed int64) int64 {
	ampf := float64(sscq2satoshi) * ((MAX_AMPLITUDE-MIN_AMPLITUDE)*randomFloat(seed) + MIN_AMPLITUDE)
	return int64(ampf)
}

// rand(100,3000)
// rand(0,2900) + 100
func randomCycle(seed int64) int64 {
	return int64(randomFloat(seed)*float64(MAX_CYCLE-MIN_CYCLE)) + MIN_CYCLE
}

//
func calcRewardFloat(amp int64, cycle int64, step int64) float64 {
	if cycle == 0 {
		return 0.0
	}
	radian := 2.0 * math.Pi * float64(step) / float64(cycle)
	return float64(amp)*math.Sin(radian) + AvgBlkRewardAsSatoshi
}

func calcRewardAsSatoshi(amp int64, cycle int64, step int64) int64 {
	return int64(calcRewardFloat(amp, cycle, step))
}
