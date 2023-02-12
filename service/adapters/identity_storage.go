package adapters

import (
	"os"
	"path/filepath"

	"github.com/boreq/errors"
	"github.com/pelletier/go-toml/v2"
	"github.com/planetary-social/scuttlego/service/domain/identity"
)

type IdentityStorage struct {
	directory string
}

func NewIdentityStorage(directory string) *IdentityStorage {
	return &IdentityStorage{directory: directory}
}

func (s *IdentityStorage) Save(iden identity.Private) error {
	storedIden := storedIdentity{
		PrivateKey: iden.PrivateKey(),
	}

	f, err := os.OpenFile(s.identityFilePath(), os.O_WRONLY|os.O_CREATE, 0o700)
	if err != nil {
		return errors.Wrap(err, "open file error")
	}

	if err := toml.NewEncoder(f).Encode(storedIden); err != nil {
		return errors.Wrap(err, "error encoding toml")
	}

	return nil
}

func (s *IdentityStorage) Load() (identity.Private, error) {
	var storedIden storedIdentity

	f, err := os.Open(s.identityFilePath())
	if err != nil {
		return identity.Private{}, errors.Wrap(err, "open file error")
	}

	if err = toml.NewDecoder(f).Decode(&storedIden); err != nil {
		return identity.Private{}, errors.Wrap(err, "error decoding toml")
	}

	iden, err := identity.NewPrivateFromBytes(storedIden.PrivateKey)
	if err != nil {
		return identity.Private{}, errors.Wrap(err, "error creating private identity from bytes")
	}

	return iden, nil
}

func (s *IdentityStorage) identityFilePath() string {
	return filepath.Join(s.directory, "identity.toml")
}

type storedIdentity struct {
	PrivateKey []byte `toml:"private_key" comment:"Never share your private key with anyone."`
}
