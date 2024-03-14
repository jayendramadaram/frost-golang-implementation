package dkg

import (
	sigagrpc "frost/internal/sigag/rpc"
	sss "frost/pkg/SSS"
	"math/big"
	"math/rand"
)

func DkGRound1(Parties sigagrpc.Parties, Threshold uint,
) error {
	a0_secret := new(big.Int).Rand(rand.New(rand.NewSource(9)), big.NewInt(256))
	_ = sss.MakePolynomial(a0_secret, Threshold-1)
	return nil
}
