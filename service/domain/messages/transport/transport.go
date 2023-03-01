package transport

import (
	"github.com/planetary-social/scuttlego/service/domain/feeds/content/known"
	"github.com/planetary-social/scuttlego/service/domain/feeds/content/transport"
)

func Mappings() transport.MessageContentMappings {
	return transport.MessageContentMappings{
		known.Contact{}.Type(): ContactMapping,
		known.Pub{}.Type():     transport.PubMapping,
	}
}
