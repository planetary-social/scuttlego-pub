package commands

import (
	"time"

	"github.com/boreq/errors"
	"github.com/planetary-social/scuttlego-pub/internal"
	"github.com/planetary-social/scuttlego-pub/service/domain"
)

type CreateInvite struct {
	numberOfUses *int
	validUntil   *time.Time
}

func NewCreateInvite(numberOfUses *int, validUntil *time.Time) (CreateInvite, error) {
	if numberOfUses != nil && *numberOfUses == 0 {
		return CreateInvite{}, errors.New("number of uses is zero")
	}

	if validUntil != nil && validUntil.IsZero() {
		return CreateInvite{}, errors.New("valid until is zero")
	}

	return CreateInvite{numberOfUses: numberOfUses, validUntil: validUntil}, nil
}

func (c CreateInvite) NumberOfUses() *int {
	if c.numberOfUses == nil {
		return nil
	}
	return internal.Pointer(*c.numberOfUses)
}

func (c CreateInvite) ValidUntil() *time.Time {
	if c.validUntil == nil {
		return nil
	}
	return internal.Pointer(*c.validUntil)
}

type CreateInviteHandler struct {
	transaction TransactionProvider
}

func NewCreateInviteHandler(transaction TransactionProvider) *CreateInviteHandler {
	return &CreateInviteHandler{transaction: transaction}
}

func (h *CreateInviteHandler) Handle(cmd CreateInvite) (domain.SecretKeySeed, error) {
	secretKeySeed, err := domain.NewSecretKeySeed()
	if err != nil {
		return domain.SecretKeySeed{}, errors.Wrap(err, "error creating a secret key seed")
	}

	invite, err := domain.NewInvite(secretKeySeed, cmd.NumberOfUses(), cmd.ValidUntil())
	if err != nil {
		return domain.SecretKeySeed{}, errors.Wrap(err, "error creating an invite")
	}

	if err := h.transaction.Transact(func(adapters Adapters) error {
		if err := adapters.Invite.Put(invite); err != nil {
			return errors.Wrap(err, "error saving the invite")
		}

		return nil
	}); err != nil {
		return domain.SecretKeySeed{}, errors.Wrap(err, "transaction failed")
	}

	return secretKeySeed, nil
}
