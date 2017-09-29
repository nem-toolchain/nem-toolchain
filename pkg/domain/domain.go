package domain

import (
	"runtime"
	"time"

	"github.com/r8d8/nem-toolchain/pkg/core"
	"github.com/r8d8/nem-toolchain/pkg/keypair"
)

// EstimateRate estimates search rate in `accounts/sec`
func EstimateRate(dt float64) uint {
	start := time.Now()
	count := 0

	for time.Since(start).Seconds() <= dt {
		keypair.Gen().Address(core.Testnet)
		count++
	}

	rate := count / int(dt)
	return uint(rate * runtime.NumCPU())
}
