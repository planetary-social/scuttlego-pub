package transport_test

import (
	"testing"

	known "github.com/planetary-social/scuttlego-pub/service/domain/messages"
	"github.com/planetary-social/scuttlego-pub/service/domain/messages/transport"
	"github.com/planetary-social/scuttlego/fixtures"
	msgcontents "github.com/planetary-social/scuttlego/service/domain/feeds/content"
	scuttlegoknown "github.com/planetary-social/scuttlego/service/domain/feeds/content/known"
	scuttlegotransport "github.com/planetary-social/scuttlego/service/domain/feeds/content/transport"
	"github.com/planetary-social/scuttlego/service/domain/feeds/message"
	"github.com/planetary-social/scuttlego/service/domain/refs"
	"github.com/stretchr/testify/require"
)

func TestMappingContactUnmarshal(t *testing.T) {
	makeContactWithActions := func(actions []scuttlegoknown.ContactAction) scuttlegoknown.Contact {
		return scuttlegoknown.MustNewContact(
			refs.MustNewIdentity("@sxlUkN7dW/qZ23Wid6J1IAnqWEJ3V13dT6TaFtn5LTc=.ed25519"),
			scuttlegoknown.MustNewContactActions(actions),
		)
	}

	pubFollow := known.MustNewPubFollow(
		refs.MustNewIdentity("@sxlUkN7dW/qZ23Wid6J1IAnqWEJ3V13dT6TaFtn5LTc=.ed25519"),
	)

	testCases := []struct {
		Name            string
		Content         string
		ExpectedMessage scuttlegoknown.KnownMessageContent
	}{
		{
			Name: "missing_action",
			Content: `
{
	"type": "contact",
	"contact": "@sxlUkN7dW/qZ23Wid6J1IAnqWEJ3V13dT6TaFtn5LTc=.ed25519"
}`,
			ExpectedMessage: nil,
		},
		{
			Name: "following",
			Content: `
{
	"type": "contact",
	"contact": "@sxlUkN7dW/qZ23Wid6J1IAnqWEJ3V13dT6TaFtn5LTc=.ed25519",
	"following": true
}`,
			ExpectedMessage: makeContactWithActions([]scuttlegoknown.ContactAction{
				scuttlegoknown.ContactActionFollow,
			}),
		},
		{
			Name: "unfollowing",
			Content: `
{
	"type": "contact",
	"contact": "@sxlUkN7dW/qZ23Wid6J1IAnqWEJ3V13dT6TaFtn5LTc=.ed25519",
	"following": false
}`,
			ExpectedMessage: makeContactWithActions([]scuttlegoknown.ContactAction{
				scuttlegoknown.ContactActionUnfollow,
			}),
		},
		{
			Name: "blocking",
			Content: `
{
	"type": "contact",
	"contact": "@sxlUkN7dW/qZ23Wid6J1IAnqWEJ3V13dT6TaFtn5LTc=.ed25519",
	"blocking": true
}`,
			ExpectedMessage: makeContactWithActions([]scuttlegoknown.ContactAction{
				scuttlegoknown.ContactActionBlock,
			}),
		},
		{
			Name: "unblocking",
			Content: `
{
	"type": "contact",
	"contact": "@sxlUkN7dW/qZ23Wid6J1IAnqWEJ3V13dT6TaFtn5LTc=.ed25519",
	"blocking": false
}`,
			ExpectedMessage: makeContactWithActions([]scuttlegoknown.ContactAction{
				scuttlegoknown.ContactActionUnblock,
			}),
		},
		{
			Name: "following_and_unblocking",
			Content: `
{
	"type": "contact",
	"contact": "@sxlUkN7dW/qZ23Wid6J1IAnqWEJ3V13dT6TaFtn5LTc=.ed25519",
	"following": true,
	"blocking": false
}`,
			ExpectedMessage: makeContactWithActions([]scuttlegoknown.ContactAction{
				scuttlegoknown.ContactActionFollow,
				scuttlegoknown.ContactActionUnblock,
			}),
		},
		{
			Name: "unfollowing_and_blocking",
			Content: `
{
	"type": "contact",
	"contact": "@sxlUkN7dW/qZ23Wid6J1IAnqWEJ3V13dT6TaFtn5LTc=.ed25519",
	"following": false,
	"blocking": true
}`,
			ExpectedMessage: makeContactWithActions([]scuttlegoknown.ContactAction{
				scuttlegoknown.ContactActionUnfollow,
				scuttlegoknown.ContactActionBlock,
			}),
		},
		{
			Name: "following_and_blocking",
			Content: `
{
	"type": "contact",
	"contact": "@sxlUkN7dW/qZ23Wid6J1IAnqWEJ3V13dT6TaFtn5LTc=.ed25519",
	"following": true,
	"blocking": true
}`,
			ExpectedMessage: nil,
		},
		{
			Name: "unfollowing_and_unblocking",
			Content: `
{
	"type": "contact",
	"contact": "@sxlUkN7dW/qZ23Wid6J1IAnqWEJ3V13dT6TaFtn5LTc=.ed25519",
	"following": false,
	"blocking": false
}`,
			ExpectedMessage: makeContactWithActions([]scuttlegoknown.ContactAction{
				scuttlegoknown.ContactActionUnfollow,
				scuttlegoknown.ContactActionUnblock,
			}),
		},
		{
			Name: "pub_follow",
			Content: `
{
	"type": "contact",
	"contact": "@sxlUkN7dW/qZ23Wid6J1IAnqWEJ3V13dT6TaFtn5LTc=.ed25519",
	"following": true,
	"pub": true
}`,
			ExpectedMessage: pubFollow,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			marshaler := newMarshaler(t)

			msg, err := marshaler.Unmarshal(message.MustNewRawContent([]byte(testCase.Content)))
			if testCase.ExpectedMessage != nil {
				require.NoError(t, err)
				require.Equal(t, testCase.ExpectedMessage, msg)
			} else {
				require.ErrorIs(t, err, msgcontents.ErrUnknownContent)
			}
		})
	}
}

func TestMappingContactMarshal(t *testing.T) {
	iden := refs.MustNewIdentity("@sxlUkN7dW/qZ23Wid6J1IAnqWEJ3V13dT6TaFtn5LTc=.ed25519")

	testCases := []struct {
		Name            string
		Actions         []scuttlegoknown.ContactAction
		ExpectedContent string
	}{
		{
			Name: "following",
			Actions: []scuttlegoknown.ContactAction{
				scuttlegoknown.ContactActionFollow,
			},
			ExpectedContent: `{"type":"contact","contact":"@sxlUkN7dW/qZ23Wid6J1IAnqWEJ3V13dT6TaFtn5LTc=.ed25519","following":true}`,
		},
		{
			Name: "unfollowing",
			Actions: []scuttlegoknown.ContactAction{
				scuttlegoknown.ContactActionUnfollow,
			},
			ExpectedContent: `{"type":"contact","contact":"@sxlUkN7dW/qZ23Wid6J1IAnqWEJ3V13dT6TaFtn5LTc=.ed25519","following":false}`,
		},
		{
			Name: "blocking",
			Actions: []scuttlegoknown.ContactAction{
				scuttlegoknown.ContactActionBlock,
			},
			ExpectedContent: `{"type":"contact","contact":"@sxlUkN7dW/qZ23Wid6J1IAnqWEJ3V13dT6TaFtn5LTc=.ed25519","blocking":true}`,
		},
		{
			Name: "unblocking",
			Actions: []scuttlegoknown.ContactAction{
				scuttlegoknown.ContactActionUnblock,
			},
			ExpectedContent: `{"type":"contact","contact":"@sxlUkN7dW/qZ23Wid6J1IAnqWEJ3V13dT6TaFtn5LTc=.ed25519","blocking":false}`,
		},
		{
			Name: "unfollowing_and_blocking",
			Actions: []scuttlegoknown.ContactAction{
				scuttlegoknown.ContactActionUnfollow,
				scuttlegoknown.ContactActionBlock,
			},
			ExpectedContent: `{"type":"contact","contact":"@sxlUkN7dW/qZ23Wid6J1IAnqWEJ3V13dT6TaFtn5LTc=.ed25519","following":false,"blocking":true}`,
		},
		{
			Name: "following_and_unblocking",
			Actions: []scuttlegoknown.ContactAction{
				scuttlegoknown.ContactActionFollow,
				scuttlegoknown.ContactActionUnblock,
			},
			ExpectedContent: `{"type":"contact","contact":"@sxlUkN7dW/qZ23Wid6J1IAnqWEJ3V13dT6TaFtn5LTc=.ed25519","following":true,"blocking":false}`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			msg := scuttlegoknown.MustNewContact(iden, scuttlegoknown.MustNewContactActions(testCase.Actions))

			marshaler := newMarshaler(t)

			raw, err := marshaler.Marshal(msg)
			require.NoError(t, err)

			require.Equal(
				t,
				testCase.ExpectedContent,
				string(raw.Bytes()),
			)
		})
	}
}

func TestMappingPubFollowMarshal(t *testing.T) {
	pubFollow := known.MustNewPubFollow(refs.MustNewIdentity("@sxlUkN7dW/qZ23Wid6J1IAnqWEJ3V13dT6TaFtn5LTc=.ed25519"))

	marshaler := newMarshaler(t)

	raw, err := marshaler.Marshal(pubFollow)
	require.NoError(t, err)

	require.Equal(
		t,
		`{"type":"contact","contact":"@sxlUkN7dW/qZ23Wid6J1IAnqWEJ3V13dT6TaFtn5LTc=.ed25519","following":true,"pub":true}`,
		string(raw.Bytes()),
	)
}

func newMarshaler(t *testing.T) *scuttlegotransport.Marshaler {
	marshaler, err := scuttlegotransport.NewMarshaler(transport.Mappings(), fixtures.SomeLogger())
	require.NoError(t, err)

	return marshaler
}
