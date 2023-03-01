package app

import "github.com/planetary-social/scuttlego-pub/service/app/commands"

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateInvite *commands.CreateInviteHandler
	RedeemInvite *commands.RedeemInviteHandler
}

type Queries struct {
}
