package secure

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type Key struct {
	Encryption *big.Int
	Decryption *big.Int
	Prime      *big.Int
}

func KeyFromPrime(prime *big.Int) (Key, error) {
	if len(prime.Bytes()) == 0 {
		return Key{}, fmt.Errorf("empty prime")
	}
	for {
		n, err := randomBigInt(len(prime.Bytes()) * 8)
		if err != nil {
			return Key{}, err
		}
		primeMinusOne := &big.Int{}
		primeMinusOne.Sub(prime, big.NewInt(1))
		gcd := &big.Int{}
		gcd.GCD(nil, nil, n, primeMinusOne)
		if gcd.Cmp(big.NewInt(1)) == 0 {
			mi := &big.Int{}
			mi.ModInverse(n, primeMinusOne)
			return Key{
				Encryption: n,
				Decryption: mi,
				Prime:      prime,
			}, nil
		}
	}
}

func random2048() (*big.Int, error) {
	// (2^2048 - 1) - 2 ^ 2047
	powTwo := &big.Int{}
	powTwo.Exp(big.NewInt(2), big.NewInt(2047), nil)

	size := &big.Int{}
	size.Exp(big.NewInt(2), big.NewInt(2048), nil)
	size.Sub(size, big.NewInt(1))
	size.Sub(size, powTwo)
	random, err := rand.Int(rand.Reader, size)
	if err != nil {
		return nil, err
	}
	random.Add(random, powTwo)
	return random, nil
}

func randomBigInt(n int) (*big.Int, error) {
	powTwo := &big.Int{}
	powTwo.Exp(big.NewInt(2), big.NewInt(int64(n-1)), nil)

	size := &big.Int{}
	size.Exp(big.NewInt(2), big.NewInt(int64(n)), nil)
	size.Sub(size, big.NewInt(1))
	size.Sub(size, powTwo)
	random, err := rand.Int(rand.Reader, size)
	if err != nil {
		return nil, err
	}
	random.Add(random, powTwo)
	return random, nil
}

func RandomPrime(n int) (*big.Int, error) {
	prime, err := rand.Prime(rand.Reader, n)
	if err != nil {
		return nil, err
	}
	return prime, nil
}

func ShamirEncrypt(message *big.Int, key Key) *big.Int {
	result := &big.Int{}
	result.Exp(message, key.Encryption, key.Prime)
	return result
}

func ShamirDecrypt(ciphertext *big.Int, key Key) *big.Int {
	result := &big.Int{}
	result.Exp(ciphertext, key.Decryption, key.Prime)
	return result
}
