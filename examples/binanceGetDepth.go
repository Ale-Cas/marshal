package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Ale-Cas/marshal"
)

type Resp struct {
	LastUpdateId int64 `json:"lastUpdateId"`
	Bids 	 [][]json.Number `json:"bids"`
	Asks 	 [][]json.Number `json:"asks"`
}

func main() {
	// Example usage of the marshal package
	baseUrl := "https://api.binance.com"

	// Perform a GET request
	resp, err := marshal.Get[Resp](http.DefaultClient, baseUrl + "/api/v3/depth?symbol=BTCUSDT&limit=5")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", *resp)
	
	// resp, err = marshal.Post[Body, Resp](http.DefaultClient, baseUrl, body) 
}