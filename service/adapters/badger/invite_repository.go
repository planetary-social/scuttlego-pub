package badger

import (
	"encoding/base64"
	"encoding/json"
	"github.com/boreq/errors"
	"github.com/dgraph-io/badger/v3"
	"github.com/planetary-social/scuttlego-pub/service/domain"
	"github.com/planetary-social/scuttlego/service/adapters/badger/utils"
	"github.com/planetary-social/scuttlego/service/domain/identity"
	"time"
)

type InviteRepository struct {
	tx *badger.Txn
}

func NewInviteRepository(tx *badger.Txn) *InviteRepository {
	return &InviteRepository{tx: tx}
}

func (i *InviteRepository) Put(invite *domain.Invite) error {
	return i.save(invite, false)
}

func (i *InviteRepository) Update(publicIdentity identity.Public, fn func(invite *domain.Invite) error) error {
	v, err := i.load(publicIdentity)
	if err != nil {
		return errors.Wrap(err, "error loading the invite")
	}

	if err := fn(v); err != nil {
		return errors.Wrap(err, "provided function returned an error")
	}

	if err := i.save(v, true); err != nil {
		return errors.Wrap(err, "error saving the invite")
	}

	return nil
}

func (i *InviteRepository) load(publicIdentity identity.Public) (*domain.Invite, error) {
	key := i.newKey(publicIdentity)
	b := i.getInvitesBucket()

	item, err := b.Get(key)
	if err != nil {
		return nil, errors.Wrap(err, "get error")
	}

	value, err := item.ValueCopy(nil)
	if err != nil {
		return nil, errors.Wrap(err, "error getting value")
	}

	var v persistedInvite
	if err := json.Unmarshal(value, &v); err != nil {
		return nil, errors.Wrap(err, "error unmarshaling the invite")

	}

	invite, err := domain.NewInviteFromHistory(v.Seed, v.RemainingUses, v.ValidUntil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating the invite")
	}

	return invite, nil
}

func (i *InviteRepository) save(invite *domain.Invite, canExist bool) error {
	b := i.getInvitesBucket()
	key := i.newKey(invite.Identity().Public())

	if !canExist {
		if _, err := b.Get(key); err != nil {
			if !errors.Is(err, badger.ErrKeyNotFound) {
				return errors.Wrap(err, "error checking if invite exists")
			}
		} else {
			return errors.New("this invite was already saved")
		}
	}

	value, err := json.Marshal(newPersistedInvite(invite))
	if err != nil {
		return errors.Wrap(err, "error persisting the invite")

	}

	if err := b.Set(key, value); err != nil {
		return errors.Wrap(err, "set error")
	}

	return nil
}

func (i *InviteRepository) newKey(publicIdentity identity.Public) []byte {
	return []byte(base64.StdEncoding.EncodeToString(publicIdentity.PublicKey()))
}

func (i *InviteRepository) getInvitesBucket() utils.Bucket {
	return utils.MustNewBucket(i.tx, utils.MustNewKey(
		utils.MustNewKeyComponent([]byte("invites")),
	))
}

type persistedInvite struct {
	RemainingUses *int                 `json:"remaining_uses,omitempty"`
	ValidUntil    *time.Time           `json:"valid_until,omitempty"`
	Seed          domain.SecretKeySeed `json:"seed"`
}

func newPersistedInvite(invite *domain.Invite) *persistedInvite {
	v := &persistedInvite{
		Seed: invite.Seed(),
	}

	remainingUses, ok := invite.RemainingUses()
	if ok {
		v.RemainingUses = &remainingUses
	}

	validUntil, ok := invite.ValidUntil()
	if ok {
		v.ValidUntil = &validUntil
	}

	return v
}
