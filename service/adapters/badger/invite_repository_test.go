package badger_test

import (
	"testing"

	"github.com/planetary-social/scuttlego-pub/service/di"
	"github.com/stretchr/testify/require"
)

func TestInviteRepository(t *testing.T) {
	_, err := di.BuildBadgerTestAdapters(t)
	require.NoError(t, err)
}
