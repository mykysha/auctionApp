// Package service is a main application logic.
package service

import (
	"encoding/json"
	"log"
	"time"

	"github.com/nndergunov/auctuionApp/server/domain"
	"github.com/nndergunov/auctuionApp/server/pkg/messagebroker"
)

type ServiceRepository interface {
	Run()
	GetAuctionData() []domain.AuctionData
	GetLastUpdateTime() time.Time
}

type Service struct {
	auctionChan chan []byte
	auctionData []domain.AuctionData
	lastUpdate  time.Time
}

func (s *Service) AuctionData() []domain.AuctionData {
	return s.auctionData
}

func NewService() *Service {
	auctionChan := make(chan []byte)
	auctionData := make([]domain.AuctionData, 0)

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
		auctionData: auctionData,
		lastUpdate:  time.Unix(0, 0),
	}
}

func (s *Service) Run() {
	go func(s *Service) {
		for data := range s.auctionChan {
			var auctionData []domain.AuctionData

			// log.Println(string(data))

			err := json.Unmarshal(data, &auctionData)
			if err != nil {
				log.Println(err)
			}

			s.auctionData = auctionData

			s.lastUpdate = time.Now()
			// log.Println(s.auctionData)
		}
	}(s)
}

func (s *Service) GetAuctionData() []domain.AuctionData {
	return s.auctionData
}

func (s *Service) GetLastUpdateTime() time.Time {
	return s.lastUpdate
}
