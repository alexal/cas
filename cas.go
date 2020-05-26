package cas

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

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

// Connected node count.
func (c *Client) NodeCount() (int, error) {
	req, err := http.NewRequest("GET", c.createURL("cas/nodeCount"), nil)

	if err != nil {
		return -1, err
	}

	bytes, err := c.doRequest(req)
	if err != nil {
		return -1, err
	}

	return strconv.Atoi(string(bytes[0]))
}

// List connected nodes.
func (c *Client) NodeNames() ([]string, error) {
	req, err := http.NewRequest("GET", c.createURL("cas/nodeNames"), nil)

	if err != nil {
		//return -1, err
	}

	bytes, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	var nodeNames []string
	if err := json.Unmarshal(bytes, &nodeNames); err != nil {
		return nil, fmt.Errorf("could not decode response JSON, %s: %v", string(bytes), err)
	}

	return nodeNames, nil
}

// List of CAS nodes.
func (c *Client) Nodes() ([]Node, error) {
	req, err := http.NewRequest("GET", c.createURL("cas/nodes"), nil)

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
