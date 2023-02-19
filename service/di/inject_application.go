package di

import (
	"github.com/google/wire"
	"github.com/planetary-social/scuttlego-pub/service/app"
	"github.com/planetary-social/scuttlego-pub/service/app/commands"
	ebtadapters "github.com/planetary-social/scuttlego/service/adapters/ebt"
	scuttlegoapp "github.com/planetary-social/scuttlego/service/app"
	scuttlegocommands "github.com/planetary-social/scuttlego/service/app/commands"
	"github.com/planetary-social/scuttlego/service/app/queries"
	"github.com/planetary-social/scuttlego/service/domain/replication"
	"github.com/planetary-social/scuttlego/service/ports/network"
	"github.com/planetary-social/scuttlego/service/ports/pubsub"
	portsrpc "github.com/planetary-social/scuttlego/service/ports/rpc"
)

var applicationSet = wire.NewSet(
	wire.Struct(new(app.Application), "*"),

	commandsSet,
	queriesSet,
)

var commandsSet = wire.NewSet(
	wire.Struct(new(app.Commands), "*"),

	commands.NewRedeemInviteHandler,
	commands.NewCreateInviteHandler,
)

var queriesSet = wire.NewSet(
	wire.Struct(new(app.Queries), "*"),
)

var scuttlegoApplicationSet = wire.NewSet(
	scuttlegoCommandsSet,
	scuttlegoQueriesSet,
)

var scuttlegoCommandsSet = wire.NewSet(
	wire.Struct(new(scuttlegoapp.Commands), "*"),

	scuttlegocommands.NewRedeemInviteHandler,
	scuttlegocommands.NewFollowHandler,
	scuttlegocommands.NewConnectHandler,
	scuttlegocommands.NewDisconnectAllHandler,
	scuttlegocommands.NewPublishRawHandler,
	scuttlegocommands.NewPublishRawAsIdentityHandler,
	scuttlegocommands.NewDownloadBlobHandler,
	scuttlegocommands.NewCreateBlobHandler,
	scuttlegocommands.NewDownloadFeedHandler,
	scuttlegocommands.NewRoomsAliasRegisterHandler,
	scuttlegocommands.NewRoomsAliasRevokeHandler,
	scuttlegocommands.NewAddToBanListHandler,
	scuttlegocommands.NewRemoveFromBanListHandler,
	scuttlegocommands.NewRunMigrationsHandler,

	scuttlegocommands.NewProcessNewLocalDiscoveryHandler,
	wire.Bind(new(network.ProcessNewLocalDiscoveryCommandHandler), new(*scuttlegocommands.ProcessNewLocalDiscoveryHandler)),

	scuttlegocommands.NewAcceptNewPeerHandler,
	wire.Bind(new(network.AcceptNewPeerCommandHandler), new(*scuttlegocommands.AcceptNewPeerHandler)),

	scuttlegocommands.NewEstablishNewConnectionsHandler,
	wire.Bind(new(network.EstablishNewConnectionsCommandHandler), new(*scuttlegocommands.EstablishNewConnectionsHandler)),

	scuttlegocommands.NewRawMessageHandler,
	wire.Bind(new(replication.RawMessageHandler), new(*scuttlegocommands.RawMessageHandler)),

	scuttlegocommands.NewCreateWantsHandler,
	wire.Bind(new(portsrpc.CreateWantsCommandHandler), new(*scuttlegocommands.CreateWantsHandler)),

	scuttlegocommands.NewHandleIncomingEbtReplicateHandler,
	wire.Bind(new(portsrpc.EbtReplicateCommandHandler), new(*scuttlegocommands.HandleIncomingEbtReplicateHandler)),

	scuttlegocommands.NewProcessRoomAttendantEventHandler,
	wire.Bind(new(pubsub.ProcessRoomAttendantEventHandler), new(*scuttlegocommands.ProcessRoomAttendantEventHandler)),

	scuttlegocommands.NewAcceptTunnelConnectHandler,
	wire.Bind(new(portsrpc.AcceptTunnelConnectHandler), new(*scuttlegocommands.AcceptTunnelConnectHandler)),
)

var scuttlegoQueriesSet = wire.NewSet(
	wire.Struct(new(scuttlegoapp.Queries), "*"),

	queries.NewReceiveLogHandler,
	queries.NewPublishedLogHandler,
	queries.NewStatusHandler,
	queries.NewBlobDownloadedEventsHandler,
	queries.NewRoomsListAliasesHandler,
	queries.NewGetMessageBySequenceHandler,

	queries.NewCreateHistoryStreamHandler,
	wire.Bind(new(portsrpc.CreateHistoryStreamQueryHandler), new(*queries.CreateHistoryStreamHandler)),
	wire.Bind(new(ebtadapters.CreateHistoryStreamHandler), new(*queries.CreateHistoryStreamHandler)),

	queries.NewGetBlobHandler,
	wire.Bind(new(portsrpc.GetBlobQueryHandler), new(*queries.GetBlobHandler)),
)
