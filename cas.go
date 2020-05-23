package cas

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	HTTPClientTimeout = 15 * time.Second
)

type Client struct {
	http     *http.Client
	baseURL  string
	userName string
	password string
	debugLog *log.Logger
}

func NewBasicAuthClient(url, username, password string) *Client {
	//TODO: For development purposes only.
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	return &Client{
		baseURL:  url,
		userName: username,
		password: password,
		http: &http.Client{
			Timeout:   HTTPClientTimeout,
			Transport: customTransport,
		},
	}
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.SetBasicAuth(c.userName, c.password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if 200 != resp.StatusCode {
		return nil, fmt.Errorf("%s", body)
	}

	return body, nil
}

func (c *Client) NodeCount() {
	url := fmt.Sprintf(c.baseURL + "cas/nodeCount")
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatal(err)
	}

	bytes, err := c.doRequest(req)
	if err != nil {
		log.Fatal(err)
	}
}
