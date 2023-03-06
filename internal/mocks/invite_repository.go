package mocks

import (
	"encoding/hex"
	"fmt"

	"github.com/planetary-social/scuttlego-pub/service/domain"
	"github.com/planetary-social/scuttlego/service/domain/identity"
)

type InviteRespositoryMock struct {
	PutCalls []InviteRepositoryPutCall

	updateInvites map[string]*domain.Invite
}

func NewInviteRespositoryMock() *InviteRespositoryMock {
	return &InviteRespositoryMock{
		updateInvites: make(map[string]*domain.Invite),
	}
}

func (i *InviteRespositoryMock) Put(invite *domain.Invite) error {
	i.PutCalls = append(i.PutCalls, InviteRepositoryPutCall{
		Invite: invite,
	})
	return nil
}

func (i *InviteRespositoryMock) Update(publicIdentity identity.Public, fn func(invite *domain.Invite) error) error {
	v, ok := i.updateInvites[hex.EncodeToString(publicIdentity.PublicKey())]
	if !ok {
		return fmt.Errorf("invite not mocked for this public identity: %s", publicIdentity.String())
	}
	return fn(v)
}

func (i *InviteRespositoryMock) MockInvite(invite *domain.Invite) {
	privateIdentity, err := identity.NewPrivateFromSeed(invite.Seed().Bytes())
	if err != nil {
		panic(err)
	}
	i.updateInvites[hex.EncodeToString(privateIdentity.Public().PublicKey())] = invite
}

type InviteRepositoryPutCall struct {
	Invite *domain.Invite
}
