package mocks

import (
	"github.com/planetary-social/scuttlego/service/domain/feeds/content/known"
	"github.com/planetary-social/scuttlego/service/domain/feeds/message"
)

type MarshalerMock struct {
	MarshalCalls       []MarshalerMockMarshalCall
	MarshalReturnValue message.RawContent
}

func NewMarshalerMock() *MarshalerMock {
	return &MarshalerMock{}
}

func (m *MarshalerMock) Marshal(content known.KnownMessageContent) (message.RawContent, error) {
	m.MarshalCalls = append(m.MarshalCalls, MarshalerMockMarshalCall{Content: content})
	if m.MarshalReturnValue.IsZero() {
		panic("zero value")
	}
	return m.MarshalReturnValue, nil
}

type MarshalerMockMarshalCall struct {
	Content known.KnownMessageContent
}
