package main

import (
	"log"

	"github.com/nndergunov/auctuionApp/publisher/pkg/mock"
	"github.com/nndergunov/auctuionApp/publisher/pkg/service"
)

func main() {
	mockAuction := mock.NewAuctionMock()

	publisherService := service.NewService(mockAuction)

	err := publisherService.PostAuctionData()
	if err != nil {
		log.Println(err)

		return
	}
}
