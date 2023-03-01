package di

import (
	"github.com/google/wire"
	pubcommands "github.com/planetary-social/scuttlego-pub/service/app/commands"
	pubtransport "github.com/planetary-social/scuttlego-pub/service/domain/messages/transport"
	"github.com/planetary-social/scuttlego/service/adapters/badger"
	"github.com/planetary-social/scuttlego/service/app/commands"
	"github.com/planetary-social/scuttlego/service/domain/feeds"
	"github.com/planetary-social/scuttlego/service/domain/feeds/content"
	"github.com/planetary-social/scuttlego/service/domain/feeds/content/transport"
	"github.com/planetary-social/scuttlego/service/domain/feeds/formats"
)

var formatsSet = wire.NewSet(
	newFormats,

	formats.NewScuttlebutt,

	transport.NewMarshaler,
	wire.Bind(new(content.Marshaler), new(*transport.Marshaler)),
	wire.Bind(new(pubcommands.Marshaler), new(*transport.Marshaler)),

	pubtransport.Mappings,

	formats.NewRawMessageIdentifier,
	wire.Bind(new(commands.RawMessageIdentifier), new(*formats.RawMessageIdentifier)),
	wire.Bind(new(badger.RawMessageIdentifier), new(*formats.RawMessageIdentifier)),
)

func newFormats(
	s *formats.Scuttlebutt,
) []feeds.FeedFormat {
	return []feeds.FeedFormat{
		s,
	}
}
