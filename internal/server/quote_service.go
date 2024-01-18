package server

import (
	"math/rand"
)

// QuoteService provides quotes for the Word of Wisdom server.
type QuoteService interface {
	GetRandomQuote() string
}

// ConcreteQuoteService is the concrete implementation of QuoteService.
type ConcreteQuoteService struct {
	Quotes []string
}

// NewQuoteService creates a new instance of QuoteService.
func NewQuoteService(quotes []string) QuoteService {
	return &ConcreteQuoteService{
		Quotes: quotes,
	}
}

// GetRandomQuote returns a random quote from the collection.
func (qs *ConcreteQuoteService) GetRandomQuote() string {
	index := rand.Intn(len(qs.Quotes))
	return qs.Quotes[index]
}
