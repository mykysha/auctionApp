package main

import (
	"github.com/nndergunov/auctuionApp/server/api"
	"github.com/nndergunov/auctuionApp/server/cmd/server"
	"github.com/nndergunov/auctuionApp/server/pkg/service"
)

func main() {
	appService := service.NewService()

	defaultAPI := api.NewAPI(appService)

	defaultServer := server.NewServer(defaultAPI)

	defaultServer.Start(":7000")
}
