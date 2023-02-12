package service

import (
	"github.com/planetary-social/scuttlego/service/domain/feeds/formats"
	"github.com/planetary-social/scuttlego/service/domain/transport/boxstream"
)

type Config struct {
	// DataDirectory specifies where the primary database and other data
	// will be stored.
	DataDirectory string

	// ListenAddress for the TCP listener in the format accepted by the
	// standard library.
	// Optional, defaults to ":8008".
	ListenAddress string

	// Setting NetworkKey is mainly useful for test networks.
	// Optional, defaults to boxstream.NewDefaultNetworkKey().
	NetworkKey boxstream.NetworkKey

	// Setting MessageHMAC is mainly useful for test networks.
	// Optional, defaults to formats.NewDefaultMessageHMAC().
	MessageHMAC formats.MessageHMAC
}

func NewDefaultConfig() Config {
	return Config{
		DataDirectory: "/some/data/directory",
		ListenAddress: ":8008",
		NetworkKey:    boxstream.NewDefaultNetworkKey(),
		MessageHMAC:   formats.NewDefaultMessageHMAC(),
	}
}
