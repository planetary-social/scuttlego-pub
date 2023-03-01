package known

import (
	"github.com/boreq/errors"
	"github.com/planetary-social/scuttlego/service/domain/feeds/content/known"
	"github.com/planetary-social/scuttlego/service/domain/refs"
)

type PubFollow struct {
	contact refs.Identity
}

func NewPubFollow(contact refs.Identity) (PubFollow, error) {
	if contact.IsZero() {
		return PubFollow{}, errors.New("zero value of contact")
	}

	return PubFollow{
		contact: contact,
	}, nil
}

func MustNewPubFollow(contact refs.Identity) PubFollow {
	c, err := NewPubFollow(contact)
	if err != nil {
		panic(err)
	}
	return c
}

func (c PubFollow) Type() known.MessageContentType {
	return "contact"
}

func (c PubFollow) Contact() refs.Identity {
	return c.contact
}

func (c PubFollow) Actions() known.ContactActions {
	return known.MustNewContactActions([]known.ContactAction{known.ContactActionFollow})
}
