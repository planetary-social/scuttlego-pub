package commands_test

import (
	"testing"

	"github.com/planetary-social/scuttlego-pub/internal/fixtures"
	"github.com/planetary-social/scuttlego-pub/service/adapters/mocks"
	"github.com/planetary-social/scuttlego-pub/service/app/commands"
	"github.com/planetary-social/scuttlego-pub/service/di"
	"github.com/planetary-social/scuttlego-pub/service/domain"
	"github.com/stretchr/testify/require"
)

func TestCreateInviteHandler(t *testing.T) {
	ts, err := di.BuildTestApplication(t)
	require.NoError(t, err)

	numberOfUses := fixtures.SomePositiveInt()
	validUntil := fixtures.SomeTime()

	cmd, err := commands.NewCreateInvite(&numberOfUses, &validUntil)
	require.NoError(t, err)

	secretKeySeed, err := ts.Commands.CreateInvite.Handle(cmd)
	require.NoError(t, err)
	require.False(t, secretKeySeed.IsZero())

	require.Equal(t,
		[]mocks.InviteRepositoryPutCall{
			{
				Invite: domain.MustNewInvite(secretKeySeed, &numberOfUses, &validUntil),
			},
		},
		ts.InviteRepository.PutCalls,
	)
}
