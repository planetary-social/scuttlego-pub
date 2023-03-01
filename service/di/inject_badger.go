package di

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/google/wire"
	"github.com/planetary-social/scuttlego-pub/service"
	pubbadgeradapters "github.com/planetary-social/scuttlego-pub/service/adapters/badger"
	pubcommands "github.com/planetary-social/scuttlego-pub/service/app/commands"
	"github.com/planetary-social/scuttlego/logging"
	scuttlegobadgeradapters "github.com/planetary-social/scuttlego/service/adapters/badger"
	"github.com/planetary-social/scuttlego/service/adapters/badger/notx"
	scuttlegocommands "github.com/planetary-social/scuttlego/service/app/commands"
	scuttlegoqueries "github.com/planetary-social/scuttlego/service/app/queries"
	blobReplication "github.com/planetary-social/scuttlego/service/domain/blobs/replication"
	"github.com/planetary-social/scuttlego/service/domain/identity"
)

var badgerAdaptersSet = wire.NewSet(
	scuttlegobadgeradapters.NewGarbageCollector,
)

var badgerNoTxRepositoriesSet = wire.NewSet(
	notx.NewNoTxBlobWantListRepository,
	wire.Bind(new(blobReplication.WantedBlobsProvider), new(*notx.NoTxBlobWantListRepository)),
	wire.Bind(new(blobReplication.WantListRepository), new(*notx.NoTxBlobWantListRepository)),

	notx.NewNoTxBlobsRepository,
	wire.Bind(new(blobReplication.BlobsRepository), new(*notx.NoTxBlobsRepository)),

	notx.NewNoTxFeedWantListRepository,
)

var badgerRepositoriesSet = wire.NewSet(
	scuttlegobadgeradapters.NewBanListRepository,
	wire.Bind(new(scuttlegocommands.BanListRepository), new(*scuttlegobadgeradapters.BanListRepository)),
	wire.Bind(new(scuttlegoqueries.BanListRepository), new(*scuttlegobadgeradapters.BanListRepository)),

	scuttlegobadgeradapters.NewBlobWantListRepository,
	wire.Bind(new(scuttlegocommands.BlobWantListRepository), new(*scuttlegobadgeradapters.BlobWantListRepository)),
	wire.Bind(new(blobReplication.WantListRepository), new(*scuttlegobadgeradapters.BlobWantListRepository)),

	scuttlegobadgeradapters.NewFeedWantListRepository,
	wire.Bind(new(scuttlegocommands.FeedWantListRepository), new(*scuttlegobadgeradapters.FeedWantListRepository)),
	wire.Bind(new(scuttlegoqueries.FeedWantListRepository), new(*scuttlegobadgeradapters.FeedWantListRepository)),

	scuttlegobadgeradapters.NewReceiveLogRepository,
	wire.Bind(new(scuttlegocommands.ReceiveLogRepository), new(*scuttlegobadgeradapters.ReceiveLogRepository)),
	wire.Bind(new(scuttlegoqueries.ReceiveLogRepository), new(*scuttlegobadgeradapters.ReceiveLogRepository)),

	scuttlegobadgeradapters.NewSocialGraphRepository,
	wire.Bind(new(scuttlegocommands.SocialGraphRepository), new(*scuttlegobadgeradapters.SocialGraphRepository)),
	wire.Bind(new(scuttlegoqueries.SocialGraphRepository), new(*scuttlegobadgeradapters.SocialGraphRepository)),
	wire.Bind(new(pubcommands.SocialGraphRepository), new(*scuttlegobadgeradapters.SocialGraphRepository)),

	scuttlegobadgeradapters.NewFeedRepository,
	wire.Bind(new(scuttlegocommands.FeedRepository), new(*scuttlegobadgeradapters.FeedRepository)),
	wire.Bind(new(scuttlegoqueries.FeedRepository), new(*scuttlegobadgeradapters.FeedRepository)),
	wire.Bind(new(pubcommands.FeedRepository), new(*scuttlegobadgeradapters.FeedRepository)),

	scuttlegobadgeradapters.NewMessageRepository,
	wire.Bind(new(scuttlegoqueries.MessageRepository), new(*scuttlegobadgeradapters.MessageRepository)),

	scuttlegobadgeradapters.NewPubRepository,
	scuttlegobadgeradapters.NewBlobRepository,

	pubbadgeradapters.NewInviteRepository,
	wire.Bind(new(pubcommands.InviteRepository), new(*pubbadgeradapters.InviteRepository)),
)

var badgerTransactionProviderSet = wire.NewSet(
	scuttlegobadgeradapters.NewCommandsTransactionProvider,
	wire.Bind(new(scuttlegocommands.TransactionProvider), new(*scuttlegobadgeradapters.CommandsTransactionProvider)),

	badgerScuttlegoCommandsAdaptersFactory,

	scuttlegobadgeradapters.NewQueriesTransactionProvider,
	wire.Bind(new(scuttlegoqueries.TransactionProvider), new(*scuttlegobadgeradapters.QueriesTransactionProvider)),

	badgerScuttlegoQueriesAdaptersFactory,

	pubbadgeradapters.NewCommandsTransactionProvider,
	wire.Bind(new(pubcommands.TransactionProvider), new(*pubbadgeradapters.CommandsTransactionProvider)),

	badgerPubCommandsAdaptersFactory,
)

var badgerNoTxTransactionProviderSet = wire.NewSet(
	notx.NewTxAdaptersFactoryTransactionProvider,
	wire.Bind(new(notx.TransactionProvider), new(*notx.TxAdaptersFactoryTransactionProvider)),

	noTxTxAdaptersFactory,
)

func noTxTxAdaptersFactory(local identity.Public, conf service.Config, logger logging.Logger) notx.TxAdaptersFactory {
	return func(tx *badger.Txn) (notx.TxAdapters, error) {
		return buildBadgerNoTxTxAdapters(tx, local, conf, logger)
	}
}

func badgerScuttlegoCommandsAdaptersFactory(config service.Config, local identity.Public, logger logging.Logger) scuttlegobadgeradapters.CommandsAdaptersFactory {
	return func(tx *badger.Txn) (scuttlegocommands.Adapters, error) {
		return buildBadgerScuttlegoCommandsAdapters(tx, local, config, logger)
	}
}

func badgerScuttlegoQueriesAdaptersFactory(config service.Config, local identity.Public, logger logging.Logger) scuttlegobadgeradapters.QueriesAdaptersFactory {
	return func(tx *badger.Txn) (scuttlegoqueries.Adapters, error) {
		return buildBadgerScuttlegoQueriesAdapters(tx, local, config, logger)
	}
}

func badgerPubCommandsAdaptersFactory(config service.Config, local identity.Public, logger logging.Logger) pubbadgeradapters.CommandsAdaptersFactory {
	return func(tx *badger.Txn) (pubcommands.Adapters, error) {
		return buildBadgerPubCommandsAdapters(tx, local, config, logger)
	}
}
