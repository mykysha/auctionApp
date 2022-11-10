package service

import (
	"fmt"

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

	dataChan, err := s.auctionMock.GetAuctionDataContinuous()
	if err != nil {
		return err
	}

	for data := range dataChan {
		fmt.Println("posting data to message broker")

		err = broker.Publish("auction", data)
		if err != nil {
			return err
		}
	}

	return nil
}
