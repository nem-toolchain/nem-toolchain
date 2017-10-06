package util

import (
	"encoding/hex"
	"fmt"
	"math"
	"runtime"
	"time"

	"github.com/nem-toolchain/nem-toolchain/pkg/core"
	"github.com/nem-toolchain/nem-toolchain/pkg/keypair"
)

// Calculate generation rate of CPU
func CPUKeyPairsInSeconds() (rate float64) {
	res := make(chan int, runtime.NumCPU())
	for i := 0; i < cap(res); i++ {
		go func(res chan<- int) {
			res <- countKeyPairs(3200)
		}(res)
	}
	for i := 0; i < cap(res); i++ {
		rate += float64(<-res) / 3.2
	}
	return
}

// Calculate amount of keypairs to be generated
// to find account with probability `pbty`
func NumberOfKeyPairs(pbty, prec float64) float64 {
	return math.Log(1-prec) / math.Log(1-pbty)
}

// Format estimated time
func TimeInSeconds(val float64) string {
	val = 1e9 * math.Trunc(val)
	if val >= math.MaxInt64 || math.IsInf(val, 0) {
		return "Inf"
	}
	return time.Duration(val).String()
}

// Pretty print account details
func PrintAccountDetails(chain core.Chain, pair keypair.KeyPair) {
	fmt.Println("Address:", pair.Address(chain).PrettyString())
	fmt.Println("Public key:", hex.EncodeToString(pair.Public))
	fmt.Println("Private key:", hex.EncodeToString(pair.Private))
}

func countKeyPairs(milliseconds time.Duration) int {
	timeout := time.After(time.Millisecond * milliseconds)
	for count := 0; ; count++ {
		keypair.Gen().Address(core.Mainnet)
		select {
		case <-timeout:
			return count
		default:
			continue
		}
	}
}
