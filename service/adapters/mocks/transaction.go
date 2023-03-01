package mocks

import (
	"github.com/planetary-social/scuttlego-pub/service/app/commands"
)

type MockCommandsTransactionProvider struct {
	adapters commands.Adapters
}

func NewMockCommandsTransactionProvider(adapters commands.Adapters) *MockCommandsTransactionProvider {
	return &MockCommandsTransactionProvider{adapters: adapters}
}

func (p *MockCommandsTransactionProvider) Transact(f func(adapters commands.Adapters) error) error {
	return f(p.adapters)
}
