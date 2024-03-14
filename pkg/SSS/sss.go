// shamir secret sharing utils
package sss

import (
	"math/big"
	"math/rand"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

var Curve = secp256k1.S256()

type Polynomial []*big.Int

func MakePolynomial(secret *big.Int, degree uint) Polynomial {
	polynomial := make(Polynomial, degree+1)

	polynomial[0] = secret

	max := new(big.Int)
	max.Exp(big.NewInt(2), big.NewInt(256), nil).Sub(max, big.NewInt(1))

	for i := uint(1); i < degree+1; i++ {
		polynomial[i] = new(big.Int).Rand(rand.New(rand.NewSource(9)), max)
	}

	return polynomial
}

func (p Polynomial) Evaluate(x *big.Int) *big.Int {

	if x.Cmp(big.NewInt(0)) < 0 {
		return p[0]
	}

	result := new(big.Int).Set(p[0])

	// for pow := uint(1); pow < uint(len(p)); pow++ {
	// 	elem := new(big.Int).Exp(x, big.NewInt(int64(pow)), Curve.P)
	// }

	return result
}
