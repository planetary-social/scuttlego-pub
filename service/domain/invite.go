package domain

import (
	"time"

	"github.com/boreq/errors"
	"github.com/planetary-social/scuttlego-pub/internal"
	"github.com/planetary-social/scuttlego/service/domain/identity"
)

type Invite struct {
	remainingUses *int
	validUntil    *time.Time
	seed          SecretKeySeed
}

func NewInvite(
	seed SecretKeySeed,
	numberOfUses *int,
	validUntil *time.Time,
) (*Invite, error) {
	if numberOfUses != nil && *numberOfUses <= 0 {
		return nil, errors.New("number of uses must be positive if set")
	}

	return newInvite(seed, numberOfUses, validUntil)
}

func MustNewInvite(
	seed SecretKeySeed,
	numberOfUses *int,
	validUntil *time.Time,
) *Invite {
	v, err := newInvite(seed, numberOfUses, validUntil)
	if err != nil {
		panic(err)
	}
	return v
}

func NewInviteFromHistory(
	seed SecretKeySeed,
	numberOfUses *int,
	validUntil *time.Time,
) (*Invite, error) {
	if numberOfUses != nil && *numberOfUses < 0 {
		return nil, errors.New("number of uses can't be negative if set")
	}

	return newInvite(seed, numberOfUses, validUntil)
}

func newInvite(
	seed SecretKeySeed,
	numberOfUses *int,
	validUntil *time.Time,
) (*Invite, error) {
	if seed.IsZero() {
		return nil, errors.New("zero value of seed")
	}

	if validUntil != nil && validUntil.IsZero() {
		return nil, errors.New("valid until is zero")
	}

	invite := &Invite{seed: seed}

	if numberOfUses != nil {
		invite.remainingUses = internal.Pointer(*numberOfUses)
	}

	if validUntil != nil {
		invite.validUntil = internal.Pointer(*validUntil)
	}

	return invite, nil
}

func (i *Invite) Redeem(publicIdentity identity.Public, currentTime time.Time) error {
	if i.remainingUses != nil && *i.remainingUses <= 0 {
		return errors.New("invite has no remaining uses")
	}

	if i.validUntil != nil && i.validUntil.Before(currentTime) {
		return errors.New("current time is after valid until")
	}

	if !publicIdentity.Equal(i.Identity().Public()) {
		return errors.New("given identity doesn't match this invite")
	}

	if i.remainingUses != nil {
		*i.remainingUses -= 1
	}

	return nil
}

func (i *Invite) Identity() identity.Private {
	privateIdentity, err := identity.NewPrivateFromSeed(i.seed.Bytes())
	if err != nil {
		panic(err)
	}
	return privateIdentity
}

func (i *Invite) RemainingUses() (int, bool) {
	if i.remainingUses == nil {
		return 0, false
	}
	return *i.remainingUses, true
}

func (i *Invite) ValidUntil() (time.Time, bool) {
	if i.validUntil == nil {
		return time.Time{}, false
	}
	return *i.validUntil, true
}

func (i *Invite) Seed() SecretKeySeed {
	return i.seed
}
