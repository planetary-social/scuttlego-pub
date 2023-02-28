package domain

import (
	"crypto/ed25519"
	cryptorand "crypto/rand"
	"io"

	"github.com/boreq/errors"
	"github.com/planetary-social/scuttlego-pub/internal"
)

type SecretKeySeed struct {
	seed []byte
}

func NewSecretKeySeed() (SecretKeySeed, error) {
	seed := make([]byte, ed25519.SeedSize)
	if _, err := io.ReadFull(cryptorand.Reader, seed); err != nil {
		return SecretKeySeed{}, errors.Wrap(err, "error reading random bytes")
	}
	return NewSecretKeySeedFromBytes(seed)
}

func MustNewSecretKeySeed() SecretKeySeed {
	v, err := NewSecretKeySeed()
	if err != nil {
		panic(err)
	}
	return v
}

func NewSecretKeySeedFromBytes(seed []byte) (SecretKeySeed, error) {
	if len(seed) != ed25519.SeedSize {
		return SecretKeySeed{}, errors.New("invalid seed size")
	}
	return SecretKeySeed{seed: seed}, nil
}

func MustNewSecretKeySeedFromBytes(seed []byte) SecretKeySeed {
	v, err := NewSecretKeySeedFromBytes(seed)
	if err != nil {
		panic(err)
	}
	return v
}

func (s SecretKeySeed) Bytes() []byte {
	return internal.CopySlice(s.seed)
}

func (s SecretKeySeed) IsZero() bool {
	return len(s.seed) == 0
}
