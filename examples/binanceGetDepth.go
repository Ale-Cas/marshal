package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Ale-Cas/marshal"
)

type BinanceOrderBook struct {
	LastUpdateId int64           `json:"lastUpdateId"`
	Bids         [][]json.Number `json:"bids"`
	Asks         [][]json.Number `json:"asks"`
}

func main() {
	// Example usage of the marshal package by querying the Binance API
	const endpoint = "https://api.binance.com/api/v3/depth?symbol=BTCUSDT&limit=5"
	resp, err := marshal.Get[BinanceOrderBook](http.DefaultClient, endpoint, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", *resp)
}
