package cas

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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

type Node struct {
	Name         string
	Port         int
	Type         string
	Connected    bool
	PID          int
	HTTPPort     int
	HTTPProtocol string
	UUID         string
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
	responseBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if len(responseBody) == 0 {
		return nil, fmt.Errorf("HTTP %d: %s (body empty)", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", responseBody)
	}

	return responseBody, nil
}

func (c *Client) NodeCount() (int, error) {
	url := fmt.Sprintf(c.baseURL + "cas/nodeCount")
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return -1, err
	}

	bytes, err := c.doRequest(req)
	if err != nil {
		return -1, err
	}

	return strconv.Atoi(string(bytes[0]))
}

func (c *Client) Nodes() ([]Node, error) {
	url := fmt.Sprintf(c.baseURL + "cas/nodes")
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	bytes, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var nodes []Node
	if err := json.Unmarshal(bytes, &nodes); err != nil {
		return nil, fmt.Errorf("could not decode response JSON, %s: %v", string(bytes), err)
	}

	return nodes, nil
}
