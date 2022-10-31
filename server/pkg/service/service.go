// Package service is a main application logic.
package service

import (
	"encoding/json"

	"github.com/nndergunov/auctuionApp/server/domain"
	"github.com/nndergunov/auctuionApp/server/pkg/messagebroker"
)

type ServiceRepository interface {
	GetAuctionData() ([]domain.AuctionData, error)
}

type Service struct {
	auctionChan chan []byte
}

func NewService() *Service {
	auctionChan := make(chan []byte)

	subscriber, err := messagebroker.NewEventSubscriber("amqp://guest:guest@localhost:5672/", auctionChan)
	if err != nil {
		return nil
	}

	err = subscriber.SubscribeToTopic("auction")
	if err != nil {
		return nil
	}

	return &Service{
		auctionChan: auctionChan,
	}
}

func (s *Service) GetAuctionData() ([]domain.AuctionData, error) {
	auctionData := <-s.auctionChan

	var data []domain.AuctionData

	err := json.Unmarshal(auctionData, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
