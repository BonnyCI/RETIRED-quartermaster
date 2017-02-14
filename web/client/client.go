package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	jww "github.com/spf13/jwalterweatherman"
)

type ClientInterface interface {
	Auth(string, interface{}) (string, error)
	Get(string, []byte, interface{}) error
	Put(string, string, string, []byte, interface{}) error
	Delete(string, string, string, []byte) error
	Token(string, *http.Request)
}

type Client struct {
	Client *http.Client
}

func (c *Client) Auth(url string, in interface{}) (string, error) {
	body, _ := json.Marshal(in)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(body))
	resp, err := c.Client.Do(req)
	if err != nil {
		jww.ERROR.Println(err)
		return "", err
	}

	tokenBytes, _ := ioutil.ReadAll(resp.Body)
	return string(tokenBytes), nil
}

func (c *Client) Token(token string, r *http.Request) {
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
}

func (c *Client) Get(url string, in interface{}, response interface{}) error {
	body, _ := json.Marshal(in)
	req, _ := http.NewRequest("GET", url, bytes.NewReader(body))
	resp, err := c.Client.Do(req)
	if err != nil {
		jww.ERROR.Println(err)
		return err
	}

	if resp.StatusCode != 200 {
		jww.ERROR.Printf("HTTP Error: %s", resp.Status)
		return fmt.Errorf("%s", resp.Status)
	}

	json.NewDecoder(resp.Body).Decode(&response)
	return nil
}

func (c *Client) Put(token string, url string, in interface{}, response interface{}) error {
	body, _ := json.Marshal(in)
	req, _ := http.NewRequest("PUT", url, bytes.NewReader(body))
	c.Token(token, req)
	resp, err := c.Client.Do(req)
	if err != nil {
		jww.ERROR.Println(err)
		return err
	}

	if resp.StatusCode != 200 {
		jww.ERROR.Printf("HTTP Error: %s", resp.Status)
		return fmt.Errorf("%s", resp.Status)
	}

	json.NewDecoder(resp.Body).Decode(&response)
	return nil
}

func (c *Client) Delete(token string, url string, in interface{}) error {
	body, _ := json.Marshal(in)
	req, _ := http.NewRequest("DELETE", url, bytes.NewReader(body))
	c.Token(token, req)
	resp, err := c.Client.Do(req)
	if err != nil {
		jww.ERROR.Println(err)
		return err
	}

	if resp.StatusCode != 200 {
		jww.ERROR.Printf("HTTP Error: %s", resp.Status)
		return fmt.Errorf("%s", resp.Status)
	}

	return nil
}

func NewClient() *Client {
	return &Client{
		&http.Client{
			Timeout: time.Second * 10,
		},
	}
}
