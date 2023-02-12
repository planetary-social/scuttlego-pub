package adapters

import (
	"os"
	"path/filepath"

	"github.com/boreq/errors"
	"github.com/pelletier/go-toml/v2"
	"github.com/planetary-social/scuttlego-pub/service"
	"github.com/planetary-social/scuttlego/service/domain/feeds/formats"
	"github.com/planetary-social/scuttlego/service/domain/transport/boxstream"
)

type ConfigStorage struct {
	directory string
}

func NewConfigStorage(directory string) *ConfigStorage {
	return &ConfigStorage{directory: directory}
}

func (s *ConfigStorage) Save(config service.Config) error {
	storedConfig := storedConfig{
		DataDirectory: config.DataDirectory,
		ListenAddress: config.ListenAddress,
		NetworkKey:    config.NetworkKey.Bytes(),
		MessageHMAC:   config.MessageHMAC.Bytes(),
	}

	f, err := os.OpenFile(s.configFilePath(), os.O_WRONLY|os.O_CREATE, 0o700)
	if err != nil {
		return errors.Wrap(err, "open file error")
	}

	if err := toml.NewEncoder(f).Encode(storedConfig); err != nil {
		return errors.Wrap(err, "error encoding toml")
	}

	return nil
}

func (s *ConfigStorage) Load() (service.Config, error) {
	var storedConfig storedConfig

	f, err := os.Open(s.configFilePath())
	if err != nil {
		return service.Config{}, errors.Wrap(err, "open file error")
	}

	if err = toml.NewDecoder(f).DisallowUnknownFields().Decode(&storedConfig); err != nil {
		return service.Config{}, errors.Wrap(err, "error decoding toml")
	}

	networkKey, err := boxstream.NewNetworkKey(storedConfig.NetworkKey)
	if err != nil {
		return service.Config{}, errors.Wrap(err, "error creating a network key")
	}

	messageHMAC, err := formats.NewMessageHMAC(storedConfig.MessageHMAC)
	if err != nil {
		return service.Config{}, errors.Wrap(err, "error creating message HMAC")
	}

	config := service.Config{
		DataDirectory: storedConfig.DataDirectory,
		ListenAddress: storedConfig.ListenAddress,
		NetworkKey:    networkKey,
		MessageHMAC:   messageHMAC,
	}

	return config, nil
}

func (s *ConfigStorage) configFilePath() string {
	return filepath.Join(s.directory, "config.toml")
}

type storedConfig struct {
	DataDirectory string `toml:"data_directory" comment:"Directory for data storage. Can be the same as config directory."`
	ListenAddress string `toml:"listen_address" comment:"Listen address for the Secure Scuttlebutt RPC TCP listener in the format accepted by the Go programming language standard library."`
	NetworkKey    []byte `toml:"network_key" comment:"Secure Scuttlebutt network key. Used to create networks separate from the Secure Scuttlebutt mainnet."`
	MessageHMAC   []byte `toml:"message_hmac" comment:"Secure Scuttlebutt message HMAC. Used mostly for testing to make messages incompatibile with the Secure Scuttlebutt mainnet."`
}
