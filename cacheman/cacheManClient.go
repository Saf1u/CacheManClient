package cacheman

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Client struct {
	Address string
	client  *http.Client
}

// BuildClient returns a cacheMan client to access a cache at the passed address
func BuildClient(address string) *Client {
	return &Client{Address: address, client: &http.Client{
		Timeout: time.Second * 20,
	}}
}

// Ping to check if server is live returns error if unavailable
func (c *Client) Ping() error {
	url := fmt.Sprint("http://", c.Address, "/ping")
	resp, err := http.Get(url)
	if err != nil {
		return errors.New("could not reach cacheman server")
	}
	if resp.StatusCode != 200 {
		return errors.New("could not reach cacheman server")
	}
	return nil
}

// Put adds the passed value into the cache and evicts records to make space
func (c *Client) Put(key string, value []byte) error {
	url := fmt.Sprint("http://", c.Address, "/", key)
	resp, err := c.client.Post(url, "-", bytes.NewReader(value))
	defer resp.Body.Close()
	if err != nil {
		return errors.New("could not add record")
	}
	if resp.StatusCode != 200 {
		return errors.New("could not add record")
	}

	return nil
}

// Get fetches stored record from cache
func (c *Client) Get(key string) ([]byte, error) {
	url := fmt.Sprint("http://", c.Address, "/", key)
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.New("could not get record")
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("could not get record")
	}
	data, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}
	return data, nil
}
