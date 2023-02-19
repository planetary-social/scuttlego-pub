package transport

import (
	"encoding/json"

	"github.com/boreq/errors"
	known "github.com/planetary-social/scuttlego-pub/service/domain/messages"
	scuttlegoknown "github.com/planetary-social/scuttlego/service/domain/feeds/content/known"
	"github.com/planetary-social/scuttlego/service/domain/feeds/content/transport"
)

var ContactMapping = transport.MessageContentMapping{
	Marshal: func(con scuttlegoknown.KnownMessageContent) ([]byte, error) {
		pubFollow, ok := con.(known.PubFollow)
		if ok {
			t := transportContact{
				MessageContentType: transport.NewMessageContentType(pubFollow),
				Contact:            pubFollow.Contact().String(),
				Following:          true,
				Pub:                true,
			}
			return json.Marshal(t)
		}

		_, ok = con.(scuttlegoknown.Contact)
		if ok {
			return transport.ContactMapping.Marshal(con)
		}

		return nil, errors.New("unknown type")

	},
	Unmarshal: func(b []byte) (scuttlegoknown.KnownMessageContent, error) {
		var t transportContact

		if err := json.Unmarshal(b, &t); err != nil {
			return nil, errors.Wrap(err, "json unmarshal failed")
		}

		knownMessageContent, err := transport.ContactMapping.Unmarshal(b)
		if err != nil {
			return nil, errors.Wrap(err, "error unmarshaling contact message")
		}

		if t.Pub {
			contact, ok := knownMessageContent.(scuttlegoknown.Contact)
			if !ok {
				return nil, errors.New("mapping didn't return a contact message")
			}

			actions := contact.Actions().List()
			if len(actions) == 1 && actions[0] == scuttlegoknown.ContactActionFollow {
				return known.NewPubFollow(contact.Contact())
			}
		}

		return knownMessageContent, nil
	},
}

type transportContact struct {
	transport.MessageContentType
	Contact   string `json:"contact"`
	Following bool   `json:"following"`
	Pub       bool   `json:"pub"`
}
