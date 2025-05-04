package main

import (
	"fmt"
	"net/http"

	"github.com/Ale-Cas/marshal"
)

type DomainsResponse struct {
	Total int64           `json:"total"`
	Time      string    `json:"time"`
	NextPage  string    `json:"next_page"`
	Domains   []DomainInfo  `json:"domains"`
}

type DomainInfo struct {
	Domain     string     `json:"domain"`
	CreateDate string     `json:"create_date"`
	UpdateDate string     `json:"update_date"`
	Country    string     `json:"country"`
	IsDead     string     `json:"isDead"`
	A          []string   `json:"A"`
	NS         []string   `json:"NS"`
	CNAME      []string   `json:"CNAME"`
	MX         []MXRecord `json:"MX"`
	TXT        []string   `json:"TXT"`
}

type MXRecord struct {
	Exchange string `json:"exchange"`
	Priority int    `json:"priority"`
}

func main() {
	// Example usage of the marshal package by querying the Domains-Index API
	const endpoint = "https://api.domainsdb.info/v1/domains/search?domain=facebook&zone=com"
	resp, err := marshal.Get[DomainsResponse](http.DefaultClient, endpoint, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", *resp)
}
