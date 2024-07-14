package models

type Response struct {
	TimeStamp int64  `json:"timestamp"`
	Asks      []Asks `json:"asks"`
	Bids      []Bids `json:"bids"`
}

type Asks struct {
	Price string `json:"price"`
}

type Bids struct {
	Price string `json:"price"`
}

type Health struct {
	AppStatus string
	DBStatus  string
	APIStatus string
}
