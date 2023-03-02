package badger_test

import (
	"testing"
	"time"

	"github.com/planetary-social/scuttlego-pub/internal"
	"github.com/planetary-social/scuttlego-pub/service/di"
	"github.com/planetary-social/scuttlego-pub/service/domain"
	"github.com/planetary-social/scuttlego/service/domain/identity"
	"github.com/stretchr/testify/require"
)

func TestInviteRepository_PutTwiceIsAnErrorIGuess(t *testing.T) {
	ts, err := di.BuildBadgerTestAdapters(t)
	require.NoError(t, err)

	secretKeySeed := domain.MustNewSecretKeySeed()
	invite := domain.MustNewInvite(secretKeySeed, nil, nil)

	err = ts.TransactionProvider.Update(func(adapters di.TestAdapters) error {
		return adapters.InviteRepository.Put(invite)
	})
	require.NoError(t, err)

	err = ts.TransactionProvider.Update(func(adapters di.TestAdapters) error {
		return adapters.InviteRepository.Put(invite)
	})
	require.EqualError(t, err, "this invite was already saved")
}

func TestInviteRepository_PutInsertsTheInviteAndUpdateLoadsItCorrectly(t *testing.T) {
	testCases := []struct {
		Name         string
		NumberOfUses *int
		ValidUntil   *time.Time
	}{
		{
			Name:         "nil",
			NumberOfUses: nil,
			ValidUntil:   nil,
		},
		{
			Name:         "not_nil",
			NumberOfUses: internal.Pointer(123),
			ValidUntil:   internal.Pointer(time.Now()),
		},
	}

	ts, err := di.BuildBadgerTestAdapters(t)
	require.NoError(t, err)

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			secretKeySeed := domain.MustNewSecretKeySeed()

			privateIdentity, err := identity.NewPrivateFromSeed(secretKeySeed.Bytes())
			require.NoError(t, err)

			publicIdentity := privateIdentity.Public()

			err = ts.TransactionProvider.Update(func(adapters di.TestAdapters) error {
				invite := domain.MustNewInvite(secretKeySeed, testCase.NumberOfUses, testCase.ValidUntil)
				return adapters.InviteRepository.Put(invite)
			})
			require.NoError(t, err)

			err = ts.TransactionProvider.Update(func(adapters di.TestAdapters) error {
				return adapters.InviteRepository.Update(publicIdentity, func(invite *domain.Invite) error {
					require.Equal(t, secretKeySeed, invite.Seed())

					validUntil, ok := invite.ValidUntil()
					if testCase.ValidUntil != nil {
						require.True(t, ok)
						require.Equal(t, testCase.ValidUntil.UTC().Round(time.Second), validUntil.UTC().Round(time.Second))
					} else {
						require.False(t, ok)
					}

					remainingUses, ok := invite.RemainingUses()
					if testCase.ValidUntil != nil {
						require.True(t, ok)
						require.Equal(t, *testCase.NumberOfUses, remainingUses)
					} else {
						require.False(t, ok)
					}

					return nil
				})
			})
			require.NoError(t, err)
		})
	}
}

func TestInviteRepository_UpdateSavesTheInvite(t *testing.T) {
	ts, err := di.BuildBadgerTestAdapters(t)
	require.NoError(t, err)

	secretKeySeed := domain.MustNewSecretKeySeed()
	numberOfUses := 123
	redeemTime := time.Now()
	validUntil := redeemTime.Add(10 * time.Second)

	privateIdentity, err := identity.NewPrivateFromSeed(secretKeySeed.Bytes())
	require.NoError(t, err)

	publicIdentity := privateIdentity.Public()

	err = ts.TransactionProvider.Update(func(adapters di.TestAdapters) error {
		invite := domain.MustNewInvite(secretKeySeed, &numberOfUses, &validUntil)
		return adapters.InviteRepository.Put(invite)
	})
	require.NoError(t, err)

	err = ts.TransactionProvider.Update(func(adapters di.TestAdapters) error {
		return adapters.InviteRepository.Update(publicIdentity, func(invite *domain.Invite) error {
			return invite.Redeem(publicIdentity, redeemTime)
		})
	})
	require.NoError(t, err)

	err = ts.TransactionProvider.Update(func(adapters di.TestAdapters) error {
		return adapters.InviteRepository.Update(publicIdentity, func(invite *domain.Invite) error {
			require.Equal(t, secretKeySeed, invite.Seed())

			loadedValidUntil, ok := invite.ValidUntil()
			require.True(t, ok)
			require.Equal(t, validUntil.UTC().Round(time.Second), loadedValidUntil.UTC().Round(time.Second))

			remainingUses, ok := invite.RemainingUses()
			require.True(t, ok)
			require.Equal(t, numberOfUses-1, remainingUses)

			return nil
		})
	})
	require.NoError(t, err)
}
