//go:build wireinject
// +build wireinject

package di

import (
	"context"
	"path/filepath"

	"github.com/boreq/errors"
	"github.com/dgraph-io/badger/v3"
	"github.com/google/wire"
	"github.com/planetary-social/scuttlego-pub/service"
	"github.com/planetary-social/scuttlego/fixtures"
	"github.com/planetary-social/scuttlego/logging"
	badgeradapters "github.com/planetary-social/scuttlego/service/adapters/badger"
	"github.com/planetary-social/scuttlego/service/adapters/badger/notx"
	"github.com/planetary-social/scuttlego/service/app/commands"
	"github.com/planetary-social/scuttlego/service/domain"
	blobReplication "github.com/planetary-social/scuttlego/service/domain/blobs/replication"
	"github.com/planetary-social/scuttlego/service/domain/feeds/content/transport"
	"github.com/planetary-social/scuttlego/service/domain/feeds/formats"
	"github.com/planetary-social/scuttlego/service/domain/graph"
	"github.com/planetary-social/scuttlego/service/domain/identity"
	"github.com/planetary-social/scuttlego/service/domain/network/local"
	"github.com/planetary-social/scuttlego/service/domain/replication"
	"github.com/planetary-social/scuttlego/service/domain/replication/ebt"
	"github.com/planetary-social/scuttlego/service/domain/replication/gossip"
	"github.com/planetary-social/scuttlego/service/domain/rooms"
	"github.com/planetary-social/scuttlego/service/domain/rooms/tunnel"
	"github.com/sirupsen/logrus"
)

var hops = graph.MustNewHops(0)

func buildTestBadgerNoTxTxAdapters(*badger.Txn, badgeradapters.TestAdaptersDependencies) (notx.TxAdapters, error) {
	wire.Build(
		wire.Struct(new(notx.TxAdapters), "*"),

		badgerRepositoriesSet,
		badgerUnpackTestDependenciesSet,

		formats.NewDefaultMessageHMAC,
		formats.NewScuttlebutt,
		transport.DefaultMappings,

		transport.NewMarshaler,
		wire.Bind(new(formats.Marshaler), new(*transport.Marshaler)),

		fixtures.SomeLogger,

		wire.Value(hops),
	)

	return notx.TxAdapters{}, nil
}

func buildBadgerNoTxTxAdapters(*badger.Txn, identity.Public, service.Config, logging.Logger) (notx.TxAdapters, error) {
	wire.Build(
		wire.Struct(new(notx.TxAdapters), "*"),

		badgerRepositoriesSet,
		formatsSet,
		extractFromConfigSet,
		adaptersSet,

		wire.Value(hops),
	)

	return notx.TxAdapters{}, nil
}

func buildBadgerTestAdapters(*badger.Txn, badgeradapters.TestAdaptersDependencies) (badgeradapters.TestAdapters, error) {
	wire.Build(
		wire.Struct(new(badgeradapters.TestAdapters), "*"),

		badgerRepositoriesSet,
		badgerUnpackTestDependenciesSet,

		formats.NewDefaultMessageHMAC,
		formats.NewScuttlebutt,
		transport.DefaultMappings,

		transport.NewMarshaler,
		wire.Bind(new(formats.Marshaler), new(*transport.Marshaler)),

		fixtures.SomeLogger,

		wire.Value(hops),
	)

	return badgeradapters.TestAdapters{}, nil
}

func buildBadgerTransactableAdapters(*badger.Txn, identity.Public, service.Config, logging.Logger) (commands.Adapters, error) {
	wire.Build(
		wire.Struct(new(commands.Adapters), "*"),

		badgerRepositoriesSet,
		formatsSet,
		extractFromConfigSet,
		adaptersSet,

		wire.Value(hops),
	)

	return commands.Adapters{}, nil
}

// BuildService creates a new service which uses the provided context as a long-term context used as a base context for
// e.g. established connections.
func BuildService(context.Context, identity.Private, service.Config) (service.Service, func(), error) {
	wire.Build(
		service.NewService,

		domain.NewPeerManager,
		wire.Bind(new(commands.NewPeerHandler), new(*domain.PeerManager)),
		wire.Bind(new(commands.PeerManager), new(*domain.PeerManager)),

		newBadger,

		newAdvertiser,
		privateIdentityToPublicIdentity,

		commands.NewMessageBuffer,

		rooms.NewScanner,
		wire.Bind(new(domain.RoomScanner), new(*rooms.Scanner)),

		rooms.NewPeerRPCAdapter,
		wire.Bind(new(rooms.MetadataGetter), new(*rooms.PeerRPCAdapter)),
		wire.Bind(new(rooms.AttendantsGetter), new(*rooms.PeerRPCAdapter)),

		tunnel.NewDialer,
		wire.Bind(new(domain.RoomDialer), new(*tunnel.Dialer)),

		newContextLogger,
		newLoggingSystem,
		wire.Bind(new(logging.LoggingSystem), new(logging.LogrusLoggingSystem)),

		newPeerManagerConfig,

		portsSet,
		applicationSet,
		scuttlegoApplicationSet,
		replicatorSet,
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
	)
	return service.Service{}, nil, nil
}

var replicatorSet = wire.NewSet(
	gossip.NewManager,
	wire.Bind(new(gossip.ReplicationManager), new(*gossip.Manager)),

	gossip.NewGossipReplicator,
	wire.Bind(new(replication.CreateHistoryStreamReplicator), new(*gossip.GossipReplicator)),
	wire.Bind(new(ebt.SelfCreateHistoryStreamReplicator), new(*gossip.GossipReplicator)),

	ebt.NewReplicator,
	wire.Bind(new(replication.EpidemicBroadcastTreesReplicator), new(ebt.Replicator)),

	replication.NewWantedFeedsCache,
	wire.Bind(new(gossip.ContactsStorage), new(*replication.WantedFeedsCache)),
	wire.Bind(new(ebt.ContactsStorage), new(*replication.WantedFeedsCache)),

	ebt.NewSessionTracker,
	wire.Bind(new(ebt.Tracker), new(*ebt.SessionTracker)),

	ebt.NewSessionRunner,
	wire.Bind(new(ebt.Runner), new(*ebt.SessionRunner)),

	replication.NewNegotiator,
	wire.Bind(new(domain.MessageReplicator), new(*replication.Negotiator)),
)

var blobReplicatorSet = wire.NewSet(
	blobReplication.NewManager,
	wire.Bind(new(blobReplication.ReplicationManager), new(*blobReplication.Manager)),
	wire.Bind(new(commands.BlobReplicationManager), new(*blobReplication.Manager)),

	blobReplication.NewReplicator,
	wire.Bind(new(domain.BlobReplicator), new(*blobReplication.Replicator)),

	blobReplication.NewBlobsGetDownloader,
	wire.Bind(new(blobReplication.Downloader), new(*blobReplication.BlobsGetDownloader)),

	blobReplication.NewHasHandler,
	wire.Bind(new(blobReplication.HasBlobHandler), new(*blobReplication.HasHandler)),
)

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
			logger.WithError(err).Error("error closing the database")
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
