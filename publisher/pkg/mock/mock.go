// Package mock provides mock data that should otherwise be obtained by the outer service.
package mock

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/nndergunov/auctuionApp/publisher/domain"
)

// AuctionRepository is a mock implementation of the auction client.
type AuctionRepository interface {
	GetAuctionDataContinuous() (chan []domain.AuctionData, error)
}

// AuctionMock is a mock implementation of the auction client.
type AuctionMock struct {
	auctionData map[int]domain.AuctionData
	lastID      int
}

// NewAuctionMock creates a new mock auction client.
func NewAuctionMock() *AuctionMock {
	auctionData := make(map[int]domain.AuctionData)

	return &AuctionMock{
		auctionData: auctionData,
		lastID:      9,
	}
}

func (a *AuctionMock) GetAuctionDataContinuous() (chan []domain.AuctionData, error) {
	dataChan := make(chan []domain.AuctionData)

	go func() {
		for {
			time.Sleep(time.Millisecond)

			for id, auction := range a.auctionData {
				if !auction.Ongoing {
					delete(a.auctionData, id)

					continue
				}

				if randInt(0, 100) < 10 {
					auction.Ongoing = false

					delete(a.auctionData, id)

					continue
				}

				auction.HighestBid += float64(randInt(100, 1000))
			}

			if randInt(0, 10) > 7 {
				a.lastID++

				startingPrice := randFloat(1000, 10000)

				a.auctionData[a.lastID] = domain.AuctionData{
					Ongoing:       true,
					ID:            a.lastID,
					StartingPrice: startingPrice,
					HighestBid:    startingPrice + randFloat(100, 1000),
					Product:       "Product " + strconv.Itoa(a.lastID),
					Owner:         "Owner " + strconv.Itoa(a.lastID),
				}
			}

			data := make([]domain.AuctionData, 0, len(a.auctionData))

			for _, auction := range a.auctionData {
				data = append(data, auction)
			}

			dataChan <- data
		}
	}()

	return dataChan, nil
}

// randFloat returns a random float64 between min and max.
func randFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// randInt returns a random int between min and max.
func randInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func chanceMachine(chance int) bool {
	percentNumber := 100

	return randInt(1, percentNumber) > chance
}
