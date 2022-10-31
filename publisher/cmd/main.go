package main

import (
	"github.com/nndergunov/auctuionApp/publisher/pkg/mock"
	"github.com/nndergunov/auctuionApp/publisher/pkg/service"
)

func main() {
	mockAuction := mock.NewAuctionMock()

	publisherService := service.NewService(mockAuction)

	publisherService.PostAuctionData()
}
