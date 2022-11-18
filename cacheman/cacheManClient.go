package cacheman

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	Address string
}

// BuildClient returns a cacheMan client to access a cache at the passed address
func BuildClient(address string) *Client {
	return &Client{Address: address}
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
	resp, err := http.Post(url, "-", bytes.NewReader(value))
	if err != nil {
		return errors.New("could not add record")
	}
	if resp.StatusCode != 200 {
		return errors.New("could not add record")
	}
	resp.Body.Close()
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
