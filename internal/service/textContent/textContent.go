package textContent

import (
	"errors"
	"log"
)

var (
	ErrFetchingTextContent = errors.New("could not fetch text content by slug")
)

type TextContent struct {
	Content string `json:"content"`
}

type TextContentStore interface {
	GetTextContent(string) (TextContent, error)
}

type Service struct {
	Store TextContentStore
}

func NewService(store TextContentStore) *Service {
	return &Service{
		Store: store,
	}
}

func (s *Service) GetTextContent(slug string) (TextContent, error) {
	// calls store passing in the context
	textContent, err := s.Store.GetTextContent(slug)
	if err != nil {
		log.Fatalf("an error occured fetching the comment: %s", err.Error())
		return TextContent{}, ErrFetchingTextContent
	}
	return textContent, nil
}
