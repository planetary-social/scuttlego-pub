package mocks

import (
	"github.com/boreq/errors"
	"github.com/planetary-social/scuttlego/service/app/commands"
	"github.com/planetary-social/scuttlego/service/domain/feeds"
	"github.com/planetary-social/scuttlego/service/domain/refs"
)

type FeedRepositoryMock struct {
	UpdateFeedResults []FeedRepositoryMockUpdateFeedCall

	feedFormat *FeedFormatMock
}

func NewFeedRepositoryMock(feedFormat *FeedFormatMock) *FeedRepositoryMock {
	return &FeedRepositoryMock{
		feedFormat: feedFormat,
	}
}

func (m *FeedRepositoryMock) UpdateFeed(ref refs.Feed, fn commands.UpdateFeedFn) error {
	feed := feeds.NewFeed(m.feedFormat)
	if err := fn(feed); err != nil {
		return errors.Wrap(err, "provided function returned an error")
	}
	m.UpdateFeedResults = append(m.UpdateFeedResults, FeedRepositoryMockUpdateFeedCall{
		Id:     ref,
		Result: feed,
	})
	return nil
}

type FeedRepositoryMockUpdateFeedCall struct {
	Id     refs.Feed
	Result *feeds.Feed
}
