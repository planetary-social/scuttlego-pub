package domain_test

import (
	"testing"

	"github.com/planetary-social/scuttlego-pub/internal/fixtures"
	"github.com/planetary-social/scuttlego-pub/service/domain"
	"github.com/stretchr/testify/require"
)

func TestNewSecretKeySeed(t *testing.T) {
	seed, err := domain.NewSecretKeySeed()
	require.NoError(t, err)
	require.NotEmpty(t, seed.Bytes())
	require.False(t, seed.IsZero())
}

func TestNewSecretKeySeedFromBytes(t *testing.T) {
	b := fixtures.SomeBytesOfLength(32)
	seed, err := domain.NewSecretKeySeedFromBytes(b)
	require.NoError(t, err)
	require.Equal(t, b, seed.Bytes())
}

func TestNewSecretKeySeedFromBytes_ReturnsAnErrorForSlicesOfWrongLength(t *testing.T) {
	b := fixtures.SomeBytesOfLength(10)
	_, err := domain.NewSecretKeySeedFromBytes(b)
	require.EqualError(t, err, "invalid seed size")
}
