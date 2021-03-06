package internal

import (
	"math/big"
)

// This package does not support parameter hot change by now.
const (
	DefaultBestContentIndexN = 5
	// parameters of earlier and more you donate, higher reputation you got.
	// unit: coin
	KeyPriceC = 1000 // C = 0.01 lino, C must be larger than 2.
	// The K parameter is set to be always 1, where it means
	// K = 0.00001 lino. The reason is that making K = 1 can optimize the computation a lot.
	// we can still change the key distribution rate by changing C.
	RoundDuration        = 25                                   // how many hours does game last.
	SampleWindowSize     = 10                                   // how many rounds is used to sample out user's customer score.
	DecayFactor          = 100 - (100 / (3 * SampleWindowSize)) // reduce to ~97% at most each time.
	MaxNumKeysEachTime   = 10000000000
	OneLinoCoin          = 100000
	InitialCustomerScore = OneLinoCoin // initial and minimum score is 1 lino.
)

var BigIntZero bigInt = big.NewInt(0)
