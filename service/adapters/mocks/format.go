package mocks

import (
	"github.com/boreq/errors"
	"github.com/planetary-social/scuttlego/service/domain/feeds"
	"github.com/planetary-social/scuttlego/service/domain/feeds/message"
	"github.com/planetary-social/scuttlego/service/domain/identity"
)

type FeedFormatMock struct {
	SignCalls       []FeedFormatMockSignCall
	SignReturnValue message.Message
}

func NewFeedFormatMock() *FeedFormatMock {
	return &FeedFormatMock{}
}

func (f *FeedFormatMock) Load(raw message.VerifiedRawMessage) (message.MessageWithoutId, error) {
	return message.MessageWithoutId{}, errors.New("not implemented")
}

func (f *FeedFormatMock) Verify(raw message.RawMessage) (message.Message, error) {
	return message.Message{}, errors.New("not implemented")
}

func (f *FeedFormatMock) Sign(unsigned message.UnsignedMessage, private identity.Private) (message.Message, error) {
	f.SignCalls = append(f.SignCalls, FeedFormatMockSignCall{
		Unsigned: unsigned,
		Private:  private,
	})
	return f.SignReturnValue, nil
}

func (f *FeedFormatMock) Peek(raw message.RawMessage) (feeds.PeekedMessage, error) {
	return feeds.PeekedMessage{}, errors.New("not implemented")
}

type FeedFormatMockSignCall struct {
	Unsigned message.UnsignedMessage
	Private  identity.Private
}
