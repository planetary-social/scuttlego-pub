package adapters_test

import (
	"testing"

	"github.com/planetary-social/scuttlego-pub/service"
	"github.com/planetary-social/scuttlego-pub/service/adapters"
	"github.com/planetary-social/scuttlego/fixtures"
	"github.com/stretchr/testify/require"
)

func TestConfigStorage(t *testing.T) {
	directory := fixtures.Directory(t)

	storage := adapters.NewConfigStorage(directory)

	config := service.NewDefaultConfig()

	err := storage.Save(config)
	require.NoError(t, err)

	loadedConfig, err := storage.Load()
	require.NoError(t, err)

	require.Equal(t, config, loadedConfig)
}
