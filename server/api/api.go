package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nndergunov/auctuionApp/server/pkg/service"
)

// API is the websocket API
type API struct {
	service service.ServiceRepository
	mux     *http.ServeMux
}

// NewAPI creates a new websocket API
func NewAPI(appService service.ServiceRepository) *API {
	api := &API{
		service: appService,
		mux:     http.NewServeMux(),
	}

	api.mux.HandleFunc("/longpoll", api.LongPoll)
	api.mux.HandleFunc("/ws", api.Upgrade)

	api.service.Run()

	return api
}

// ServeHTTP serves the API.
func (api API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.mux.ServeHTTP(w, r)
}

// LongPoll handles long polling requests.
func (api API) LongPoll(w http.ResponseWriter, _ *http.Request) {
	reqTime := time.Now()

	timeout := time.Duration(10) * time.Second

	for {
		if time.Now().After(reqTime.Add(timeout)) {
			break
		}

		lastUpdateTime := api.service.GetLastUpdateTime()

		if lastUpdateTime.After(reqTime) {
			data := api.service.GetAuctionData()

			jsonData, err := json.Marshal(data)
			if err != nil {
				return
			}

			_, err = w.Write(jsonData)
			if err != nil {
				return
			}

			return
		}

		// time.Sleep(100 * time.Microsecond)
	}

	w.WriteHeader(http.StatusNoContent)
}

// Upgrade upgrades the connection to a websocket connection.
func (api API) Upgrade(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		HandshakeTimeout:  0,
		ReadBufferSize:    0,
		WriteBufferSize:   0,
		WriteBufferPool:   nil,
		Subprotocols:      nil,
		Error:             nil,
		CheckOrigin:       func(r *http.Request) bool { return true },
		EnableCompression: false,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	defer func(ws *websocket.Conn) {
		_ = ws.Close()
	}(conn)

	lastUpdateTime := api.service.GetLastUpdateTime()

	for {
		if !lastUpdateTime.Before(api.service.GetLastUpdateTime()) {
			continue
		}

		lastUpdateTime = api.service.GetLastUpdateTime()

		data := api.service.GetAuctionData()

		jsonData, err := json.Marshal(data)
		if err != nil {
			return
		}

		err = conn.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
			return
		}
	}
}
