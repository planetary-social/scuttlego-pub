package main

import (
	"os"

	"github.com/boreq/errors"
	"github.com/boreq/guinea"
	"github.com/planetary-social/scuttlego-pub/service"
	"github.com/planetary-social/scuttlego-pub/service/adapters"
	"github.com/planetary-social/scuttlego/service/domain/identity"
)

var initCommand = guinea.Command{
	Run:         initFn,
	Subcommands: nil,
	Options:     nil,
	Arguments: []guinea.Argument{
		{
			Name:        "config_directory",
			Multiple:    false,
			Optional:    false,
			Description: "directory in which the configuration will be initialized",
		},
	},
	ShortDescription: "initializes the configuration",
	Description:      "Initializes the configuration in the provided directory.",
}

func initFn(cliContext guinea.Context) error {
	directory := cliContext.Arguments[0]

	if err := checkDirectoryExists(directory); err != nil {
		return errors.Wrap(err, "error checking if directory exists")
	}

	privateIdentity, err := identity.NewPrivate()
	if err != nil {
		return errors.Wrap(err, "error creating private identity")
	}

	config := service.NewDefaultConfig()

	identityStorage := adapters.NewIdentityStorage(directory)
	if err := identityStorage.Save(privateIdentity); err != nil {
		return errors.Wrap(err, "error saving identity")
	}

	configStorage := adapters.NewConfigStorage(directory)
	if err := configStorage.Save(config); err != nil {
		return errors.Wrap(err, "error saving config")
	}

	return nil
}

func checkDirectoryExists(directory string) error {
	stat, err := os.Stat(directory)
	if err != nil {
		return errors.Wrap(err, "stat error")
	}

	if !stat.IsDir() {
		return errors.New("not a directory")
	}

	return nil
}
