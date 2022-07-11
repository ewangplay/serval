package client

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ewangplay/serval/io"
)

type Client struct {
	addr string
	c    *HttpClient
}

func NewClient(addr string) (*Client, error) {
	if len(addr) == 0 {
		return nil, fmt.Errorf("addr must be set")
	}
	c, err := NewHttpClient()
	if err != nil {
		return nil, err
	}

	return &Client{addr, c}, nil
}

func (c *Client) Ping() error {
	url := fmt.Sprintf("http://%s/api/v1/ping", c.addr)
	respBody, err := c.c.Get(url)
	if err != nil {
		return err
	}

	fmt.Println("Ping response: ", string(respBody))

	var m struct {
		Message string
	}
	err = json.Unmarshal(respBody, &m)
	if strings.Compare(strings.ToUpper(m.Message), "PONG") != 0 {
		return fmt.Errorf("ping-pong testing failed")
	}

	return nil
}

func (c *Client) CreateDid(req *io.CreateDidReq) error {
	url := fmt.Sprintf("http://%s/api/v1/did/create", c.addr)

	reqBody, err := json.Marshal(req)
	if err != nil {
		return err
	}

	respBody, err := c.c.Post(url, reqBody)
	if err != nil {
		return err
	}

	fmt.Println("CreateDid response: ", string(respBody))

	return nil
}

func (c *Client) ResolveDid(did string) (*io.DDO, error) {
	if did == "" {
		return nil, fmt.Errorf("did cannot be empty")
	}

	url := fmt.Sprintf("http://%s/api/v1/did/resolve/%s", c.addr, did)

	respBody, err := c.c.Get(url)
	if err != nil {
		return nil, err
	}

	fmt.Println("ResolveDid response: ", string(respBody))

	var resp io.ResolveDidResp
	err = json.Unmarshal(respBody, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Document, nil
}

func (c *Client) RevokeDid(req *io.RevokeDidReq) error {
	url := fmt.Sprintf("http://%s/api/v1/did/revoke", c.addr)

	reqBody, err := json.Marshal(req)
	if err != nil {
		return err
	}

	respBody, err := c.c.Post(url, reqBody)
	if err != nil {
		return err
	}

	fmt.Println("RevokeDid response: ", string(respBody))

	return nil
}
