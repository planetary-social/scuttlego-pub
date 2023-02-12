package adapters_test

import (
	"testing"

	"github.com/planetary-social/scuttlego-pub/service/adapters"
	"github.com/planetary-social/scuttlego/fixtures"
	"github.com/stretchr/testify/require"
)

func TestIdentityStorage(t *testing.T) {
	directory := fixtures.Directory(t)

	storage := adapters.NewIdentityStorage(directory)

	iden := fixtures.SomePrivateIdentity()

	err := storage.Save(iden)
	require.NoError(t, err)

	loadedIden, err := storage.Load()
	require.NoError(t, err)

	require.Equal(t, iden, loadedIden)
}
