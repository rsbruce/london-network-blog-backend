package models

import (
	"errors"
	"log"
)

var (
	ErrFetchingSlugsForUser = errors.New("Could not retrieve slugs for user")
)

type Slugs struct {
	Slugs []string `json:"slugs"`
}

type SlugsStore interface {
	GetSlugsForUser(string) (Slugs, error)
}

type SlugsService struct {
	Store SlugsStore
}

func NewSlugsService(store SlugsStore) *SlugsService {
	return &SlugsService{
		Store: store,
	}
}

func (s *SlugsService) GetSlugsForUser(slug string) (Slugs, error) {
	slugs, err := s.Store.GetSlugsForUser(slug)
	if err != nil {
		log.Fatalf("An error occured fetching the latest post feed: %s", err.Error())
		return Slugs{}, ErrFetchingSlugsForUser
	}
	return slugs, nil
}
