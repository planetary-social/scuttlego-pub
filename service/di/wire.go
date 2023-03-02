//go:build wireinject
// +build wireinject

package di

import (
	"path/filepath"
	"testing"

	"github.com/boreq/errors"
	"github.com/dgraph-io/badger/v3"
	"github.com/google/wire"
	"github.com/planetary-social/scuttlego-pub/internal/fixtures"
	"github.com/planetary-social/scuttlego-pub/service"
	pubmocks "github.com/planetary-social/scuttlego-pub/service/adapters/mocks"
	"github.com/planetary-social/scuttlego-pub/service/app"
	"github.com/planetary-social/scuttlego-pub/service/app/commands"
	"github.com/planetary-social/scuttlego/logging"
	badgeradapters "github.com/planetary-social/scuttlego/service/adapters/badger"
	"github.com/planetary-social/scuttlego/service/adapters/badger/notx"
	scuttlegocommands "github.com/planetary-social/scuttlego/service/app/commands"
	scuttlegoqueries "github.com/planetary-social/scuttlego/service/app/queries"
	"github.com/planetary-social/scuttlego/service/domain"
	"github.com/planetary-social/scuttlego/service/domain/identity"
	"github.com/planetary-social/scuttlego/service/domain/network/local"
	"github.com/planetary-social/scuttlego/service/domain/rooms/tunnel"
	"github.com/sirupsen/logrus"
)

func BuildService(identity.Private, service.Config) (service.Service, func(), error) {
	wire.Build(
		service.NewService,

		domain.NewPeerManager,
		wire.Bind(new(scuttlegocommands.PeerManager), new(*domain.PeerManager)),

		newBadger,

		newAdvertiser,
		privateIdentityToPublicIdentity,

		scuttlegocommands.NewMessageBuffer,

		tunnel.NewDialer,
		wire.Bind(new(domain.RoomDialer), new(*tunnel.Dialer)),

		newContextLogger,
		newLoggingSystem,
		wire.Bind(new(logging.LoggingSystem), new(logging.LogrusLoggingSystem)),

		newPeerManagerConfig,

		portsSet,
		applicationSet,
		scuttlegoApplicationSet,
		replicationSet,
		blobReplicatorSet,
		formatsSet,
		pubSubSet,
		badgerNoTxRepositoriesSet,
		badgerTransactionProviderSet,
		badgerNoTxTransactionProviderSet,
		badgerAdaptersSet,
		blobsAdaptersSet,
		adaptersSet,
		extractFromConfigSet,
		networkingSet,
		migrationsSet,
		contentSet,
	)
	return service.Service{}, nil, nil
}

type TestApplication struct {
	Commands app.Commands

	SocialGraphRepository *pubmocks.SocialGraphRepositoryMock
	InviteRepository      *pubmocks.InviteRespositoryMock
	FeedRepository        *pubmocks.FeedRepositoryMock
	Marshaler             *pubmocks.MarshalerMock
	FeedFormat            *pubmocks.FeedFormatMock
	LocalIdentity         identity.Private
	CurrentTimeProvider   *pubmocks.CurrentTimeProviderMock
}

func BuildTestApplication(testing.TB) (TestApplication, error) {
	wire.Build(
		wire.Struct(new(TestApplication), "*"),

		commandsSet,

		pubmocks.NewMockCommandsTransactionProvider,
		wire.Bind(new(commands.TransactionProvider), new(*pubmocks.MockCommandsTransactionProvider)),

		wire.Struct(new(commands.Adapters), "*"),

		pubmocks.NewSocialGraphRepositoryMock,
		wire.Bind(new(commands.SocialGraphRepository), new(*pubmocks.SocialGraphRepositoryMock)),

		pubmocks.NewInviteRespositoryMock,
		wire.Bind(new(commands.InviteRepository), new(*pubmocks.InviteRespositoryMock)),

		pubmocks.NewFeedRepositoryMock,
		wire.Bind(new(commands.FeedRepository), new(*pubmocks.FeedRepositoryMock)),

		pubmocks.NewCurrentTimeProviderMock,
		wire.Bind(new(commands.CurrentTimeProvider), new(*pubmocks.CurrentTimeProviderMock)),

		pubmocks.NewMarshalerMock,
		wire.Bind(new(commands.Marshaler), new(*pubmocks.MarshalerMock)),

		pubmocks.NewFeedFormatMock,

		fixtures.SomePrivateIdentity,
	)

	return TestApplication{}, nil
}

type BadgerTestAdapters struct {
	TransactionProvider *TestTransactionProvider
}

func BuildBadgerTestAdapters(testing.TB) (BadgerTestAdapters, error) {
	wire.Build(
		wire.Struct(new(BadgerTestAdapters), "*"),

		badgerTestTransactionProviderSet,
		fixtures.Badger,
	)

	return BadgerTestAdapters{}, nil
}

func buildBadgerNoTxTxAdapters(*badger.Txn, identity.Public, service.Config, logging.Logger) (notx.TxAdapters, error) {
	wire.Build(
		wire.Struct(new(notx.TxAdapters), "*"),

		badgerRepositoriesSet,
		formatsSet,
		extractFromConfigSet,
		adaptersSet,
		contentSet,
	)

	return notx.TxAdapters{}, nil
}

func buildBadgerScuttlegoCommandsAdapters(*badger.Txn, identity.Public, service.Config, logging.Logger) (scuttlegocommands.Adapters, error) {
	wire.Build(
		wire.Struct(new(scuttlegocommands.Adapters), "*"),

		badgerRepositoriesSet,
		formatsSet,
		extractFromConfigSet,
		adaptersSet,
		contentSet,
	)

	return scuttlegocommands.Adapters{}, nil
}

func buildBadgerScuttlegoQueriesAdapters(*badger.Txn, identity.Public, service.Config, logging.Logger) (scuttlegoqueries.Adapters, error) {
	wire.Build(
		wire.Struct(new(scuttlegoqueries.Adapters), "*"),

		badgerRepositoriesSet,
		formatsSet,
		extractFromConfigSet,
		adaptersSet,
		contentSet,
	)

	return scuttlegoqueries.Adapters{}, nil
}

func buildBadgerPubCommandsAdapters(*badger.Txn, identity.Public, service.Config, logging.Logger) (commands.Adapters, error) {
	wire.Build(
		wire.Struct(new(commands.Adapters), "*"),

		badgerRepositoriesSet,
		formatsSet,
		extractFromConfigSet,
		adaptersSet,
		contentSet,
	)

	return commands.Adapters{}, nil
}

func buildBadgerTestAdapters(*badger.Txn) (TestAdapters, error) {
	wire.Build(
		wire.Struct(new(TestAdapters), "*"),

		badgerRepositoriesSet,
	)

	return TestAdapters{}, nil
}

func newAdvertiser(l identity.Public, config service.Config) (*local.Advertiser, error) {
	return local.NewAdvertiser(l, config.ListenAddress)
}

func newBadger(system logging.LoggingSystem, logger logging.Logger, config service.Config) (*badger.DB, func(), error) {
	badgerDirectory := filepath.Join(config.DataDirectory, "badger")

	options := badger.DefaultOptions(badgerDirectory)
	options.Logger = badgeradapters.NewLogger(system, badgeradapters.LoggerLevelWarning)
	options.SyncWrites = true

	db, err := badger.Open(options)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to open the database")
	}

	return db, func() {
		if err := db.Close(); err != nil {
			logger.Error().WithError(err).Message("error closing the database")
		}
	}, nil

}

func privateIdentityToPublicIdentity(p identity.Private) identity.Public {
	return p.Public()
}

func newContextLogger(loggingSystem logging.LoggingSystem) logging.Logger {
	return logging.NewContextLogger(loggingSystem, "scuttlego")
}

func newLoggingSystem() logging.LogrusLoggingSystem {
	logger := logrus.New()

	return logging.NewLogrusLoggingSystem(logger)
}

func newPeerManagerConfig() domain.PeerManagerConfig {
	return domain.PeerManagerConfig{
		PreferredPubs: nil,
	}
}
