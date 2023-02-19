package di

import (
	"github.com/google/wire"
	"github.com/planetary-social/scuttlego-pub/service"
	"github.com/planetary-social/scuttlego/service/domain/feeds/formats"
	"github.com/planetary-social/scuttlego/service/domain/graph"
	"github.com/planetary-social/scuttlego/service/domain/transport/boxstream"
)

var extractFromConfigSet = wire.NewSet(
	extractNetworkKeyFromConfig,
	extractMessageHMACFromConfig,
	extractHopsFromConfig,
)

func extractNetworkKeyFromConfig(config service.Config) boxstream.NetworkKey {
	return config.NetworkKey
}

func extractMessageHMACFromConfig(config service.Config) formats.MessageHMAC {
	return config.MessageHMAC
}

func extractHopsFromConfig(config service.Config) graph.Hops {
	return config.Hops
}
