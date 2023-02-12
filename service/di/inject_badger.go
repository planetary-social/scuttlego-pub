package di

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/google/wire"
	"github.com/planetary-social/scuttlego-pub/service"
	"github.com/planetary-social/scuttlego/logging"
	badgeradapters "github.com/planetary-social/scuttlego/service/adapters/badger"
	"github.com/planetary-social/scuttlego/service/adapters/badger/notx"
	"github.com/planetary-social/scuttlego/service/adapters/mocks"
	"github.com/planetary-social/scuttlego/service/app/commands"
	"github.com/planetary-social/scuttlego/service/app/queries"
	blobReplication "github.com/planetary-social/scuttlego/service/domain/blobs/replication"
	"github.com/planetary-social/scuttlego/service/domain/identity"
	"github.com/planetary-social/scuttlego/service/domain/replication"
)

//nolint:unused
var badgerUnpackTestDependenciesSet = wire.NewSet(
	wire.FieldsOf(new(badgeradapters.TestAdaptersDependencies),
		"BanListHasher",
		"CurrentTimeProvider",
		"RawMessageIdentifier",
		"LocalIdentity",
	),
	wire.Bind(new(badgeradapters.BanListHasher), new(*mocks.BanListHasherMock)),
	wire.Bind(new(commands.CurrentTimeProvider), new(*mocks.CurrentTimeProviderMock)),
	wire.Bind(new(badgeradapters.RawMessageIdentifier), new(*mocks.RawMessageIdentifierMock)),
)

//nolint:unused
var badgerAdaptersSet = wire.NewSet(
	badgeradapters.NewGarbageCollector,
)

//nolint:unused
var badgerNoTxRepositoriesSet = wire.NewSet(
	notx.NewNoTxBlobWantListRepository,
	wire.Bind(new(blobReplication.WantListStorage), new(*notx.NoTxBlobWantListRepository)),
	wire.Bind(new(blobReplication.WantListRepository), new(*notx.NoTxBlobWantListRepository)),

	notx.NewNoTxFeedRepository,
	wire.Bind(new(queries.FeedRepository), new(*notx.NoTxFeedRepository)),

	notx.NewNoTxMessageRepository,
	wire.Bind(new(queries.MessageRepository), new(*notx.NoTxMessageRepository)),

	notx.NewNoTxReceiveLogRepository,
	wire.Bind(new(queries.ReceiveLogRepository), new(*notx.NoTxReceiveLogRepository)),

	notx.NewNoTxWantedFeedsRepository,
	wire.Bind(new(replication.WantedFeedsRepository), new(*notx.NoTxWantedFeedsRepository)),
)

//nolint:unused
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

	badgeradapters.NewSocialGraphRepository,
	wire.Bind(new(commands.SocialGraphRepository), new(*badgeradapters.SocialGraphRepository)),

	badgeradapters.NewFeedRepository,
	wire.Bind(new(commands.FeedRepository), new(*badgeradapters.FeedRepository)),

	badgeradapters.NewWantedFeedsRepository,
	badgeradapters.NewPubRepository,
	badgeradapters.NewMessageRepository,
	badgeradapters.NewBlobRepository,
)

//nolint:unused
var badgerTransactionProviderSet = wire.NewSet(
	badgeradapters.NewTransactionProvider,
	wire.Bind(new(commands.TransactionProvider), new(*badgeradapters.TransactionProvider)),

	badgerTransactableAdaptersFactory,
)

//nolint:unused
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

func badgerTransactableAdaptersFactory(config service.Config, local identity.Public, logger logging.Logger) badgeradapters.AdaptersFactory {
	return func(tx *badger.Txn) (commands.Adapters, error) {
		return buildBadgerTransactableAdapters(tx, local, config, logger)
	}
}
