package domain_test

import (
	"github.com/boreq/errors"
	"github.com/planetary-social/scuttlego-pub/internal"
	"github.com/planetary-social/scuttlego-pub/service/domain"
	"github.com/planetary-social/scuttlego/service/domain/identity"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestInviteConstructors(t *testing.T) {
	type constructorTestCase struct {
		Name          string
		Seed          domain.SecretKeySeed
		NumberOfUses  *int
		ValidUntil    *time.Time
		ExpectedError error
	}

	commonTestCases := []constructorTestCase{

		{
			Name:          "infinite",
			Seed:          domain.MustNewSecretKeySeed(),
			NumberOfUses:  nil,
			ValidUntil:    nil,
			ExpectedError: nil,
		},
		{
			Name:          "zero_value_of_seed",
			Seed:          domain.SecretKeySeed{},
			NumberOfUses:  nil,
			ValidUntil:    nil,
			ExpectedError: errors.New("zero value of seed"),
		},
		{
			Name:          "limited_number_of_uses",
			Seed:          domain.MustNewSecretKeySeed(),
			NumberOfUses:  internal.Pointer(10),
			ValidUntil:    nil,
			ExpectedError: nil,
		},
		{
			Name:          "limited_time",
			Seed:          domain.MustNewSecretKeySeed(),
			NumberOfUses:  nil,
			ValidUntil:    internal.Pointer(time.Now()),
			ExpectedError: nil,
		},
		{
			Name:          "limited_time_and_number_of_uses",
			Seed:          domain.MustNewSecretKeySeed(),
			NumberOfUses:  internal.Pointer(10),
			ValidUntil:    internal.Pointer(time.Now()),
			ExpectedError: nil,
		},
		{
			Name:          "zero_value_of_time",
			Seed:          domain.MustNewSecretKeySeed(),
			NumberOfUses:  nil,
			ValidUntil:    internal.Pointer(time.Time{}),
			ExpectedError: errors.New("valid until is zero"),
		},
	}

	newInviteTestCases := []constructorTestCase{
		{
			Name:          "zero_uses",
			Seed:          domain.MustNewSecretKeySeed(),
			NumberOfUses:  internal.Pointer(0),
			ValidUntil:    nil,
			ExpectedError: errors.New("number of uses must be positive if set"),
		},
		{
			Name:          "negative_uses",
			Seed:          domain.MustNewSecretKeySeed(),
			NumberOfUses:  internal.Pointer(-1),
			ValidUntil:    nil,
			ExpectedError: errors.New("number of uses must be positive if set"),
		},
	}

	newInviteFromHistoryTestCases := []constructorTestCase{
		{
			Name:          "zero_uses",
			Seed:          domain.MustNewSecretKeySeed(),
			NumberOfUses:  internal.Pointer(0),
			ValidUntil:    nil,
			ExpectedError: nil,
		},
		{
			Name:          "negative_uses",
			Seed:          domain.MustNewSecretKeySeed(),
			NumberOfUses:  internal.Pointer(-1),
			ValidUntil:    nil,
			ExpectedError: errors.New("number of uses can't be negative if set"),
		},
	}

	check := func(t *testing.T, testCase constructorTestCase, invite *domain.Invite, err error) {
		if testCase.ExpectedError == nil {
			require.Equal(t, testCase.Seed, invite.Seed())

			remainingUses, ok := invite.RemainingUses()
			if testCase.NumberOfUses != nil {
				require.True(t, ok)
				require.Equal(t, *testCase.NumberOfUses, remainingUses)
			} else {
				require.False(t, ok)
			}

			validUntil, ok := invite.ValidUntil()
			if testCase.ValidUntil != nil {
				require.True(t, ok)
				require.Equal(t, *testCase.ValidUntil, validUntil)
			} else {
				require.False(t, ok)
			}

			expectedPrivateIdentity, err := identity.NewPrivateFromSeed(testCase.Seed.Bytes())
			require.NoError(t, err)
			require.Equal(t, expectedPrivateIdentity, invite.Identity())
		} else {
			require.EqualError(t, err, testCase.ExpectedError.Error())
		}
	}

	t.Run("NewInvite", func(t *testing.T) {
		for _, testCase := range append(commonTestCases, newInviteTestCases...) {
			t.Run(testCase.Name, func(t *testing.T) {
				invite, err := domain.NewInvite(testCase.Seed, testCase.NumberOfUses, testCase.ValidUntil)
				check(t, testCase, invite, err)
			})
		}
	})

	t.Run("NewInviteFromHistory", func(t *testing.T) {
		for _, testCase := range append(commonTestCases, newInviteFromHistoryTestCases...) {
			t.Run(testCase.Name, func(t *testing.T) {
				invite, err := domain.NewInviteFromHistory(testCase.Seed, testCase.NumberOfUses, testCase.ValidUntil)
				check(t, testCase, invite, err)
			})
		}
	})
}

func TestInvite_CanBeRedeemedUsingMatchingPublicIdentity(t *testing.T) {
	secretKeySeed := domain.MustNewSecretKeySeed()
	privateIdentity := identity.MustNewPrivateFromSeed(secretKeySeed.Bytes())

	invite := domain.MustNewInvite(secretKeySeed, nil, nil)

	err := invite.Redeem(privateIdentity.Public(), time.Now())
	require.NoError(t, err)
}

func TestInvite_CanNotBeRedeemedUsingDifferentPublicIdentity(t *testing.T) {
	secretKeySeed := domain.MustNewSecretKeySeed()

	privateIdentity, err := identity.NewPrivate()
	require.NoError(t, err)

	invite := domain.MustNewInvite(secretKeySeed, nil, nil)

	err = invite.Redeem(privateIdentity.Public(), time.Now())
	require.EqualError(t, err, "given identity doesn't match this invite")
}

func TestInvite_RedeemRespectsValidUntil(t *testing.T) {
	secretKeySeed := domain.MustNewSecretKeySeed()

	validUntil := time.Now()
	beforeValidUntil := validUntil.Add(-1 * time.Minute)
	afterValidUntil := validUntil.Add(1 * time.Minute)

	t.Run("before", func(t *testing.T) {
		invite := domain.MustNewInvite(secretKeySeed, nil, &validUntil)

		err := invite.Redeem(invite.Identity().Public(), beforeValidUntil)
		require.NoError(t, err)
	})

	t.Run("after", func(t *testing.T) {
		invite := domain.MustNewInvite(secretKeySeed, nil, &validUntil)

		err := invite.Redeem(invite.Identity().Public(), afterValidUntil)
		require.EqualError(t, err, "current time is after valid until")
	})
}

func TestInvite_RedeemRespectsNumberOfUses(t *testing.T) {
	secretKeySeed := domain.MustNewSecretKeySeed()

	numberOfUses := 2

	invite := domain.MustNewInvite(secretKeySeed, &numberOfUses, nil)

	// 1
	err := invite.Redeem(invite.Identity().Public(), time.Now())
	require.NoError(t, err)

	remainingUses, ok := invite.RemainingUses()
	require.True(t, ok)
	require.Equal(t, 1, remainingUses)

	// 2
	err = invite.Redeem(invite.Identity().Public(), time.Now())
	require.NoError(t, err)

	remainingUses, ok = invite.RemainingUses()
	require.True(t, ok)
	require.Equal(t, 0, remainingUses)

	// 3
	err = invite.Redeem(invite.Identity().Public(), time.Now())
	require.EqualError(t, err, "invite has no remaining uses")

	remainingUses, ok = invite.RemainingUses()
	require.True(t, ok)
	require.Equal(t, 0, remainingUses)
}
