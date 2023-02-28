package commands_test

import (
	"fmt"
	"testing"

	"github.com/planetary-social/scuttlego-pub/internal/fixtures"
	"github.com/planetary-social/scuttlego-pub/service/adapters/mocks"
	"github.com/planetary-social/scuttlego-pub/service/app/commands"
	"github.com/planetary-social/scuttlego-pub/service/di"
	"github.com/planetary-social/scuttlego-pub/service/domain"
	known "github.com/planetary-social/scuttlego-pub/service/domain/messages"
	"github.com/planetary-social/scuttlego/service/domain/feeds"
	"github.com/planetary-social/scuttlego/service/domain/feeds/message"
	"github.com/planetary-social/scuttlego/service/domain/graph"
	"github.com/planetary-social/scuttlego/service/domain/identity"
	"github.com/planetary-social/scuttlego/service/domain/refs"
	"github.com/stretchr/testify/require"
)

func TestRedeemInviteHandler_UpdatesAnInviteAndPublishesPubFollow(t *testing.T) {
	ts, err := di.BuildTestApplication(t)
	require.NoError(t, err)

	secretKeySeed := fixtures.SomeSecretKeySeed()
	numberOfUses := fixtures.SomePositiveInt()
	invite := domain.MustNewInvite(secretKeySeed, &numberOfUses, nil)

	privateIdentity, err := identity.NewPrivateFromSeed(secretKeySeed.Bytes())
	require.NoError(t, err)

	localFeed := refs.MustNewIdentityFromPublic(ts.LocalIdentity.Public()).MainFeed()
	feedToFollow := fixtures.SomeRefFeed()

	rawContent := fixtures.SomeRawContent()
	ts.Marshaler.MarshalReturnValue = rawContent

	msg := fixtures.SomeMessageWithFeedSequence(localFeed, message.NewFirstSequence())
	ts.FeedFormat.SignReturnValue = msg

	ts.InviteRepository.MockInvite(invite)

	cmd, err := commands.NewRedeemInvite(privateIdentity.Public(), feedToFollow)
	require.NoError(t, err)

	msgRef, err := ts.Commands.RedeemInvite.Handle(cmd)
	require.NoError(t, err)
	require.False(t, msgRef.IsZero())

	// command redeems an invite
	remainingUses, ok := invite.RemainingUses()
	require.True(t, ok)
	require.Equal(t, numberOfUses-1, remainingUses)

	// command checks the social graph
	require.Equal(t, ts.SocialGraphRepository.GetSocialGraphCallsCount, 1)

	// command publishes a follow message
	require.Equal(t,
		[]mocks.MarshalerMockMarshalCall{
			{
				Content: known.MustNewPubFollow(refs.MustNewIdentityFromPublic(feedToFollow.Identity())),
			},
		},
		ts.Marshaler.MarshalCalls,
	)

	require.Len(t, ts.FeedRepository.UpdateFeedResults, 1)
	require.Equal(t,
		localFeed,
		ts.FeedRepository.UpdateFeedResults[0].Id,
	)
	require.Equal(t,
		[]feeds.MessageToPersist{
			feeds.MustNewMessageToPersist(
				msg,
				nil,
				nil,
				nil,
			),
		},
		ts.FeedRepository.UpdateFeedResults[0].Result.PopForPersisting(),
	)
}

func TestRedeemInviteHandler_ReturnsAnErrorIfTheUserIsAlreadyBeingFollowed(t *testing.T) {
	ts, err := di.BuildTestApplication(t)
	require.NoError(t, err)

	secretKeySeed := fixtures.SomeSecretKeySeed()
	numberOfUses := fixtures.SomePositiveInt()
	invite := domain.MustNewInvite(secretKeySeed, &numberOfUses, nil)

	privateIdentity, err := identity.NewPrivateFromSeed(secretKeySeed.Bytes())
	require.NoError(t, err)

	feed := fixtures.SomeRefFeed()
	feedIdentity := refs.MustNewIdentityFromPublic(feed.Identity())

	ts.SocialGraphRepository.GetSocialGraphReturnValue = graph.NewSocialGraph(map[string]graph.Hops{
		feedIdentity.String(): graph.MustNewHops(1),
	})

	rawContent := fixtures.SomeRawContent()
	ts.Marshaler.MarshalReturnValue = rawContent

	msg := fixtures.SomeMessageWithFeedSequence(feed, message.NewFirstSequence())
	ts.FeedFormat.SignReturnValue = msg

	ts.InviteRepository.MockInvite(invite)

	cmd, err := commands.NewRedeemInvite(privateIdentity.Public(), feed)
	require.NoError(t, err)

	_, err = ts.Commands.RedeemInvite.Handle(cmd)
	require.EqualError(t, err, "transaction failed: already following this user")
}

func TestRedeemInviteHandler_UsesProvidedIdentityToGetAnInvite(t *testing.T) {
	ts, err := di.BuildTestApplication(t)
	require.NoError(t, err)

	invalidPublicIdentity := fixtures.SomePublicIdentity()

	ts.Marshaler.MarshalReturnValue = fixtures.SomeRawContent()

	cmd, err := commands.NewRedeemInvite(invalidPublicIdentity, fixtures.SomeRefFeed())
	require.NoError(t, err)

	_, err = ts.Commands.RedeemInvite.Handle(cmd)
	require.EqualError(t,
		err,
		fmt.Sprintf("transaction failed: error updating the invite: invite not mocked for this public identity: %s", invalidPublicIdentity.String()),
	)

}
