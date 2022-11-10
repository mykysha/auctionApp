package main

import (
	"log"

	"github.com/nndergunov/auctuionApp/client/pkg/client"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("config.yaml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	host := viper.GetString("host")

	clientService := client.NewClient("ws://"+host+":7000/ws", "http://"+host+":7000/longpoll")

	err = clientService.WSConnect()
	if err != nil {
		log.Println(err)
	}
}
