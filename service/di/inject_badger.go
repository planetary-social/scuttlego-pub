package di

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/google/wire"
	"github.com/planetary-social/scuttlego-pub/service"
	"github.com/planetary-social/scuttlego/logging"
	badgeradapters "github.com/planetary-social/scuttlego/service/adapters/badger"
	"github.com/planetary-social/scuttlego/service/adapters/badger/notx"
	"github.com/planetary-social/scuttlego/service/app/commands"
	"github.com/planetary-social/scuttlego/service/app/queries"
	blobReplication "github.com/planetary-social/scuttlego/service/domain/blobs/replication"
	"github.com/planetary-social/scuttlego/service/domain/identity"
	"github.com/planetary-social/scuttlego/service/domain/replication"
)

var badgerAdaptersSet = wire.NewSet(
	badgeradapters.NewGarbageCollector,
)

var badgerNoTxRepositoriesSet = wire.NewSet(
	notx.NewNoTxBlobWantListRepository,
	wire.Bind(new(blobReplication.WantedBlobsProvider), new(*notx.NoTxBlobWantListRepository)),
	wire.Bind(new(blobReplication.WantListRepository), new(*notx.NoTxBlobWantListRepository)),

	notx.NewNoTxWantedFeedsRepository,
	wire.Bind(new(replication.WantedFeedsRepository), new(*notx.NoTxWantedFeedsRepository)),

	notx.NewNoTxBlobsRepository,
	wire.Bind(new(blobReplication.BlobsRepository), new(*notx.NoTxBlobsRepository)),
)

var badgerRepositoriesSet = wire.NewSet(
	badgeradapters.NewBanListRepository,
	wire.Bind(new(commands.BanListRepository), new(*badgeradapters.BanListRepository)),

	badgeradapters.NewBlobWantListRepository,
	wire.Bind(new(commands.BlobWantListRepository), new(*badgeradapters.BlobWantListRepository)),
	wire.Bind(new(blobReplication.WantListRepository), new(*badgeradapters.BlobWantListRepository)),

	badgeradapters.NewFeedWantListRepository,
	wire.Bind(new(commands.FeedWantListRepository), new(*badgeradapters.FeedWantListRepository)),

	badgeradapters.NewReceiveLogRepository,
	wire.Bind(new(commands.ReceiveLogRepository), new(*badgeradapters.ReceiveLogRepository)),
	wire.Bind(new(queries.ReceiveLogRepository), new(*badgeradapters.ReceiveLogRepository)),

	badgeradapters.NewSocialGraphRepository,
	wire.Bind(new(commands.SocialGraphRepository), new(*badgeradapters.SocialGraphRepository)),

	badgeradapters.NewFeedRepository,
	wire.Bind(new(commands.FeedRepository), new(*badgeradapters.FeedRepository)),
	wire.Bind(new(queries.FeedRepository), new(*badgeradapters.FeedRepository)),

	badgeradapters.NewMessageRepository,
	wire.Bind(new(queries.MessageRepository), new(*badgeradapters.MessageRepository)),

	badgeradapters.NewWantedFeedsRepository,
	badgeradapters.NewPubRepository,
	badgeradapters.NewBlobRepository,
)

var badgerTransactionProviderSet = wire.NewSet(
	badgeradapters.NewCommandsTransactionProvider,
	wire.Bind(new(commands.TransactionProvider), new(*badgeradapters.CommandsTransactionProvider)),

	badgerCommandsAdaptersFactory,

	badgeradapters.NewQueriesTransactionProvider,
	wire.Bind(new(queries.TransactionProvider), new(*badgeradapters.QueriesTransactionProvider)),

	badgerQueriesAdaptersFactory,
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

func badgerCommandsAdaptersFactory(config service.Config, local identity.Public, logger logging.Logger) badgeradapters.CommandsAdaptersFactory {
	return func(tx *badger.Txn) (commands.Adapters, error) {
		return buildBadgerCommandsAdapters(tx, local, config, logger)
	}
}

func badgerQueriesAdaptersFactory(config service.Config, local identity.Public, logger logging.Logger) badgeradapters.QueriesAdaptersFactory {
	return func(tx *badger.Txn) (queries.Adapters, error) {
		return buildBadgerQueriesAdapters(tx, local, config, logger)
	}
}
