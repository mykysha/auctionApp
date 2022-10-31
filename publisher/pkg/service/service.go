package service

import (
	"encoding/json"

	"github.com/nndergunov/auctuionApp/publisher/pkg/messagebroker"
	"github.com/nndergunov/auctuionApp/publisher/pkg/mock"
)

type Repository interface {
	PostAuctionData() error
}

type Service struct {
	auctionMock mock.AuctionRepository
}

func NewService(auctionMock mock.AuctionRepository) *Service {
	return &Service{
		auctionMock: auctionMock,
	}
}

func (s *Service) PostAuctionData() error {
	broker, err := messagebroker.NewEventPublisher("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return err
	}

	for {
		data, err := s.auctionMock.GetAuctionData()
		if err != nil {
			return err
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			return err
		}

		err = broker.Publish("auction", jsonData)
		if err != nil {
			return err
		}
	}
}
