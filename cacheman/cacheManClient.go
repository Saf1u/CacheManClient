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

func BuildClient(address string) *Client {
	return &Client{Address: address}
}

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

func (c *Client) Put(key string, value []byte) error {
	url := fmt.Sprint("http://", c.Address, "/", key)
	resp, err := http.Post(url, "-", bytes.NewReader(value))
	if err != nil {
		return errors.New("could not add record")
	}
	if resp.StatusCode != 200 {
		return errors.New("could not add record")
	}
	return nil
}

func (c *Client) Get(key string) ([]byte, error) {
	url := fmt.Sprint("http://", c.Address, "/", key)
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.New("could not add record")
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("could not get record")
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return data, nil
}
