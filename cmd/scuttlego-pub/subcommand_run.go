package main

import (
	"context"

	"github.com/boreq/errors"
	"github.com/boreq/guinea"
	"github.com/planetary-social/scuttlego-pub/service/adapters"
	"github.com/planetary-social/scuttlego-pub/service/di"
)

var runCommand = guinea.Command{
	Run:         runFn,
	Subcommands: nil,
	Options:     nil,
	Arguments: []guinea.Argument{
		{
			Name:        "config_directory",
			Multiple:    false,
			Optional:    false,
			Description: "Path to the directory containing the configuration.",
		},
	},
	ShortDescription: "runs the pub",
	Description:      "Runs the pub with the provided configuration.",
}

func runFn(cliContext guinea.Context) error {
	configDirectory := cliContext.Arguments[0]

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	identityStorage := adapters.NewIdentityStorage(configDirectory)
	configStorage := adapters.NewConfigStorage(configDirectory)

	iden, err := identityStorage.Load()
	if err != nil {
		return errors.Wrap(err, "error loading identity")
	}

	config, err := configStorage.Load()
	if err != nil {
		return errors.Wrap(err, "error loading config")
	}

	service, cleanup, err := di.BuildService(iden, config)
	if err != nil {
		return errors.Wrap(err, "error building the service")
	}
	defer cleanup()

	if err := service.RunMigrations(ctx); err != nil {
		return errors.Wrap(err, "error running migrations")
	}

	if err := service.Run(ctx); err != nil {
		return errors.Wrap(err, "error running the service")
	}

	return nil
}
