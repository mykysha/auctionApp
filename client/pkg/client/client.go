package client

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nndergunov/auctuionApp/client/domain"
)

type Repository interface {
	WSConnect() error
	WSConnectFinite(accepts int) (time.Duration, error)
	LongPollConnect() error
	LongPollConnectFinite(accepts int) (time.Duration, error)
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

		// print in blue color
		// fmt.Printf("\033[34m%v - %v\033[0m\n", time.Now(), len(auctionData))
	}
}

func (c *Client) WSConnectFinite(accepts int) (time.Duration, error) {
	conn, _, err := websocket.DefaultDialer.Dial(c.wsURL, nil)
	if err != nil {
		return 0, err
	}

	start := time.Now()

	// print everything from server with the timestamp
	for i := 0; i < accepts; i++ {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return 0, err
		}

		var auctionData []domain.AuctionData

		err = json.Unmarshal(msg, &auctionData)
		if err != nil {
			return 0, err
		}

		// print in blue color
		// fmt.Printf("\033[34m%v - %v\033[0m\n", time.Now(), len(auctionData))
	}

	return time.Since(start), nil
}

func (c *Client) LongPollConnect() error {
	for {
		resp, err := http.Get(c.lpURL)
		if err != nil {
			return err
		}

		if resp.StatusCode == http.StatusNoContent {
			continue
		}

		var auctionData []domain.AuctionData

		err = json.NewDecoder(resp.Body).Decode(&auctionData)
		if err != nil {
			return err
		}

		// print in yellow color time and length of the auctionData
		// fmt.Printf("\033[33m%v - %v\033[0m\n", time.Now().Format(""), len(auctionData))
	}
}

func (c *Client) LongPollConnectFinite(accepts int) (time.Duration, error) {
	start := time.Now()

	for i := 0; i < accepts; i++ {
		resp, err := http.Get(c.lpURL)
		if err != nil {
			return 0, err
		}

		if resp.StatusCode == http.StatusNoContent {
			i--

			continue
		}

		var auctionData []domain.AuctionData

		err = json.NewDecoder(resp.Body).Decode(&auctionData)
		if err != nil {
			return 0, err
		}

		// print in yellow color time and length of the auctionData
		// fmt.Printf("\033[33m%v - %v\033[0m\n", time.Now().Format(""), len(auctionData))
	}

	return time.Since(start), nil
}
