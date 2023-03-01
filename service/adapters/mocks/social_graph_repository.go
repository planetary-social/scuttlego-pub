package mocks

import "github.com/planetary-social/scuttlego/service/domain/graph"

type SocialGraphRepositoryMock struct {
	GetSocialGraphCallsCount  int
	GetSocialGraphReturnValue graph.SocialGraph
}

func NewSocialGraphRepositoryMock() *SocialGraphRepositoryMock {
	return &SocialGraphRepositoryMock{}
}

func (s *SocialGraphRepositoryMock) GetSocialGraph() (graph.SocialGraph, error) {
	s.GetSocialGraphCallsCount++
	return s.GetSocialGraphReturnValue, nil
}
