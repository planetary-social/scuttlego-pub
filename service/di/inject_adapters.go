package di

import (
	"path"

	"github.com/google/wire"
	"github.com/planetary-social/scuttlego-pub/service"
	pubcommands "github.com/planetary-social/scuttlego-pub/service/app/commands"
	"github.com/planetary-social/scuttlego/logging"
	"github.com/planetary-social/scuttlego/service/adapters"
	"github.com/planetary-social/scuttlego/service/adapters/badger"
	"github.com/planetary-social/scuttlego/service/adapters/blobs"
	ebtadapters "github.com/planetary-social/scuttlego/service/adapters/ebt"
	invitesadapters "github.com/planetary-social/scuttlego/service/adapters/invites"
	"github.com/planetary-social/scuttlego/service/app/commands"
	"github.com/planetary-social/scuttlego/service/app/queries"
	blobreplication "github.com/planetary-social/scuttlego/service/domain/blobs/replication"
	"github.com/planetary-social/scuttlego/service/domain/invites"
	"github.com/planetary-social/scuttlego/service/domain/replication/ebt"
	"github.com/planetary-social/scuttlego/service/domain/transport/boxstream"
)

var blobsAdaptersSet = wire.NewSet(
	newFilesystemStorage,
	wire.Bind(new(blobreplication.BlobStorage), new(*blobs.FilesystemStorage)),
	wire.Bind(new(blobreplication.BlobStorer), new(*blobs.FilesystemStorage)),
	wire.Bind(new(queries.BlobStorage), new(*blobs.FilesystemStorage)),
	wire.Bind(new(blobreplication.BlobSizeRepository), new(*blobs.FilesystemStorage)),
	wire.Bind(new(commands.BlobCreator), new(*blobs.FilesystemStorage)),
)

func newFilesystemStorage(logger logging.Logger, config service.Config) (*blobs.FilesystemStorage, error) {
	return blobs.NewFilesystemStorage(path.Join(config.DataDirectory, "blobs"), logger)
}

var adaptersSet = wire.NewSet(
	adapters.NewCurrentTimeProvider,
	wire.Bind(new(commands.CurrentTimeProvider), new(*adapters.CurrentTimeProvider)),
	wire.Bind(new(boxstream.CurrentTimeProvider), new(*adapters.CurrentTimeProvider)),
	wire.Bind(new(invitesadapters.CurrentTimeProvider), new(*adapters.CurrentTimeProvider)),
	wire.Bind(new(blobreplication.CurrentTimeProvider), new(*adapters.CurrentTimeProvider)),
	wire.Bind(new(pubcommands.CurrentTimeProvider), new(*adapters.CurrentTimeProvider)),

	adapters.NewBanListHasher,
	wire.Bind(new(badger.BanListHasher), new(*adapters.BanListHasher)),

	ebtadapters.NewCreateHistoryStreamHandlerAdapter,
	wire.Bind(new(ebt.MessageStreamer), new(*ebtadapters.CreateHistoryStreamHandlerAdapter)),

	invitesadapters.NewInviteDialer,
	wire.Bind(new(invites.InviteDialer), new(*invitesadapters.InviteDialer)),
)
