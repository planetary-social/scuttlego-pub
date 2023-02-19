package commands

import (
	"time"

	"github.com/planetary-social/scuttlego-pub/service/domain"
	"github.com/planetary-social/scuttlego/service/app/commands"
	"github.com/planetary-social/scuttlego/service/domain/feeds/content/known"
	"github.com/planetary-social/scuttlego/service/domain/feeds/message"
	"github.com/planetary-social/scuttlego/service/domain/graph"
	"github.com/planetary-social/scuttlego/service/domain/identity"
)

type TransactionProvider interface {
	Transact(func(adapters Adapters) error) error
}

type Adapters struct {
	SocialGraph SocialGraphRepository
	Invite      InviteRepository
	Feed        commands.FeedRepository
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
