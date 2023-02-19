package commands

import (
	"github.com/boreq/errors"
	"github.com/planetary-social/scuttlego-pub/service/domain"
	known "github.com/planetary-social/scuttlego-pub/service/domain/messages"
	"github.com/planetary-social/scuttlego/service/domain/feeds"
	"github.com/planetary-social/scuttlego/service/domain/feeds/message"
	"github.com/planetary-social/scuttlego/service/domain/identity"
	"github.com/planetary-social/scuttlego/service/domain/refs"
)

type RedeemInvite struct {
	identity     identity.Public
	feedToFollow refs.Feed
}

func NewRedeemInvite(identity identity.Public, feedToFollow refs.Feed) (RedeemInvite, error) {
	if identity.IsZero() {
		return RedeemInvite{}, errors.New("zero value of identity")
	}
	if feedToFollow.IsZero() {
		return RedeemInvite{}, errors.New("zero value of feed to follow")
	}
	return RedeemInvite{identity: identity, feedToFollow: feedToFollow}, nil
}

func (cmd RedeemInvite) Identity() identity.Public {
	return cmd.identity
}

func (cmd RedeemInvite) FeedToFollow() refs.Feed {
	return cmd.feedToFollow
}

func (cmd RedeemInvite) IsZero() bool {
	return cmd.identity.IsZero()
}

type RedeemInviteHandler struct {
	transaction         TransactionProvider
	currentTimeProvider CurrentTimeProvider
	marshaler           Marshaler
	localIdentity       identity.Private
}

func NewRedeemInviteHandler(
	transaction TransactionProvider,
	currentTimeProvider CurrentTimeProvider,
	marshaler Marshaler,
	localIdentity identity.Private,
) *RedeemInviteHandler {
	return &RedeemInviteHandler{
		transaction:         transaction,
		currentTimeProvider: currentTimeProvider,
		marshaler:           marshaler,
		localIdentity:       localIdentity,
	}
}

func (h *RedeemInviteHandler) Handle(cmd RedeemInvite) (refs.Message, error) {
	if cmd.IsZero() {
		return refs.Message{}, errors.New("zero value of cmd")
	}

	localIdentityRef, err := refs.NewIdentityFromPublic(h.localIdentity.Public())
	if err != nil {
		return refs.Message{}, errors.Wrap(err, "could not create the identity ref")
	}

	feedToFollowRef, err := refs.NewIdentityFromPublic(cmd.FeedToFollow().Identity())
	if err != nil {
		return refs.Message{}, errors.Wrap(err, "error creating feed ref")
	}

	msgToPublish, err := h.createMessageToPublish(feedToFollowRef)
	if err != nil {
		return refs.Message{}, errors.Wrap(err, "error creating message to publish")
	}

	var msgId refs.Message

	if err := h.transaction.Transact(func(adapters Adapters) error {
		if err := adapters.Invite.Update(cmd.Identity(), func(invite *domain.Invite) error {
			if err := invite.Redeem(cmd.Identity(), h.currentTimeProvider.Get()); err != nil {
				return errors.Wrap(err, "error redeeming the invite")
			}
			return nil
		}); err != nil {
			return errors.Wrap(err, "error updating the invite")
		}

		graph, err := adapters.SocialGraph.GetSocialGraph()
		if err != nil {
			return errors.Wrap(err, "error getting social graph")
		}

		for _, contact := range graph.Contacts() {
			if contact.Id.Equal(feedToFollowRef) && contact.Hops.Int() == 1 {
				return errors.New("already following this user")
			}
		}

		if err := adapters.Feed.UpdateFeed(localIdentityRef.MainFeed(), func(feed *feeds.Feed) error {
			var err error
			msgId, err = feed.CreateMessage(msgToPublish, h.currentTimeProvider.Get(), h.localIdentity)
			if err != nil {
				return errors.Wrap(err, "failed to create a message")
			}
			return nil
		}); err != nil {
			return errors.Wrap(err, "error updating the local feed")
		}

		return nil
	}); err != nil {
		return refs.Message{}, errors.Wrap(err, "transaction failed")
	}

	return msgId, nil
}

func (h *RedeemInviteHandler) createMessageToPublish(feedRef refs.Identity) (message.RawContent, error) {
	contact, err := known.NewPubFollow(feedRef)
	if err != nil {
		return message.RawContent{}, errors.Wrap(err, "failed to create a message")
	}

	content, err := h.marshaler.Marshal(contact)
	if err != nil {
		return message.RawContent{}, errors.Wrap(err, "error marshaling")
	}

	return content, nil
}
