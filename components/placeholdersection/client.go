package placeholdersection

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

type PlaceholderClient struct {
	client  *http.Client
	baseURL string
}

func NewPlaceholderClient(baseURL string) *PlaceholderClient {
	var client *http.Client

	clientOnce := sync.Once{}
	clientOnce.Do(func() {
		client = &http.Client{
			Timeout: time.Second * 10,
			Transport: &http.Transport{
				MaxIdleConns:       10,
				IdleConnTimeout:    30 * time.Second,
				DisableCompression: true,
				Dial: (&net.Dialer{
					Timeout: 5 * time.Second,
				}).Dial,
				MaxIdleConnsPerHost: 10,
			},
		}
	})

	return &PlaceholderClient{
		client:  client,
		baseURL: baseURL,
	}
}

func (p *PlaceholderClient) Get(query string) []PlaceholderModel {

	req, _ := http.NewRequest("GET", p.baseURL+"/"+query, nil)
	resp, err := p.client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	respString, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("An error happened while reading response body : %v\n", err)
	}
	var result []PlaceholderModel
	err = json.Unmarshal(respString, &result)

	if err != nil {
		log.Printf("An error happened while parsing response body : %v\n", err)
	}

	return result
}

func (p *PlaceholderClient) GetBaseURL() string {
	return p.baseURL
}
