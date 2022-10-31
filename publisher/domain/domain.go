package domain

type AuctionData struct {
	Ongoing       bool
	ID            int     `json:"id"`
	StartingPrice float64 `json:"startingPrice"`
	HighestBid    float64 `json:"highestBid"`
	Product       string  `json:"product"`
	Owner         string  `json:"owner"`
}
