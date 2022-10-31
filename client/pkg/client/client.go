package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nndergunov/auctuionApp/client/domain"
)

type Repository interface {
	WSConnect() error
	LongPollConnect() error
}

type Client struct {
	wsURL string
	lpURL string
}

func NewClient(wsURL, lpURL string) *Client {
	return &Client{
		wsURL: wsURL,
		lpURL: lpURL,
	}
}

func (c *Client) WSConnect() error {
	conn, _, err := websocket.DefaultDialer.Dial(c.wsURL, nil)
	if err != nil {
		return err
	}

	// print everything from server with the timestamp
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return err
		}

		var auctionData []domain.AuctionData

		err = json.Unmarshal(msg, &auctionData)
		if err != nil {
			return err
		}

		fmt.Println(auctionData, time.Now())
	}
}

func (c *Client) LongPollConnect() error {
	for {
		time.Sleep(100 * time.Millisecond)

		resp, err := http.Get(c.lpURL)
		if err != nil {
			return err
		}

		var auctionData []domain.AuctionData

		err = json.NewDecoder(resp.Body).Decode(&auctionData)
		if err != nil {
			return err
		}

		fmt.Println(auctionData, time.Now())
	}
}
