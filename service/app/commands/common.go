package commands

import (
	"time"

	"github.com/planetary-social/scuttlego-pub/service/domain"
	scuttlegocommands "github.com/planetary-social/scuttlego/service/app/commands"
	"github.com/planetary-social/scuttlego/service/domain/feeds/content/known"
	"github.com/planetary-social/scuttlego/service/domain/feeds/message"
	"github.com/planetary-social/scuttlego/service/domain/graph"
	"github.com/planetary-social/scuttlego/service/domain/identity"
	"github.com/planetary-social/scuttlego/service/domain/refs"
)

type TransactionProvider interface {
	Update(func(adapters Adapters) error) error
}

type Adapters struct {
	SocialGraph SocialGraphRepository
	Invite      InviteRepository
	Feed        FeedRepository
}

type InviteRepository interface {
	Put(invite *domain.Invite) error
	Update(publicIdentity identity.Public, fn func(invite *domain.Invite) error) error
}

type SocialGraphRepository interface {
	GetSocialGraph() (graph.SocialGraph, error)
}

type CurrentTimeProvider interface {
	Get() time.Time
}

type Marshaler interface {
	Marshal(content known.KnownMessageContent) (message.RawContent, error)
}

type FeedRepository interface {
	// UpdateFeed updates the specified feed by calling the provided function on
	// it. Feed is never nil.
	UpdateFeed(ref refs.Feed, fn scuttlegocommands.UpdateFeedFn) error
}
