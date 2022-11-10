package main

import (
	"log"
	"sync"
	"time"

	"github.com/nndergunov/auctuionApp/client/pkg/client"
	"github.com/spf13/viper"
)

type result struct {
	wsIsFaster bool
	timeDiff   float64
}

func main() {
	viper.SetConfigFile("config.yaml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	host := viper.GetString("host")

	numberOfMeasurements := 10

	numberOfRuns := 1000
	numberOfTests := 2

	clientService := client.NewClient("ws://"+host+":7000/ws", "http://"+host+":7000/longpoll")

	results := make([]result, numberOfMeasurements)

	for i := 0; i < numberOfMeasurements; i++ {
		startFlag := make(chan struct{})

		startWG := new(sync.WaitGroup)

		startWG.Add(numberOfTests)

		endWG := new(sync.WaitGroup)

		endWG.Add(numberOfTests)

		wsTime := time.Duration(0)
		lpTime := time.Duration(0)

		go func() {
			startWG.Done()

			<-startFlag

			took, err := clientService.WSConnectFinite(numberOfRuns)
			if err != nil {
				log.Println(err)
			}

			wsTime = took

			endWG.Done()
		}()

		go func() {
			startWG.Done()

			<-startFlag

			took, err := clientService.LongPollConnectFinite(numberOfRuns)
			if err != nil {
				log.Println(err)
			}

			lpTime = took

			endWG.Done()
		}()

		startWG.Wait()

		close(startFlag)

		endWG.Wait()

		results[i] = result{
			wsIsFaster: wsTime < lpTime,
			timeDiff:   float64(lpTime.Milliseconds()-wsTime.Milliseconds()) / float64(wsTime.Milliseconds()) * 100,
		}
	}

	wsFaster := 0

	averageDiff := 0.0

	for _, res := range results {
		if res.wsIsFaster {
			wsFaster++
		}

		averageDiff += res.timeDiff
	}

	averageDiff /= float64(numberOfMeasurements)

	log.Printf("Websocket is faster in %d/%d cases, average difference is %f%%", wsFaster, numberOfMeasurements, averageDiff)
}
