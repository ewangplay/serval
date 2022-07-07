package client

import (
	"encoding/json"
	"fmt"

	"github.com/ewangplay/serval/io"
)

type Client struct {
	addr string
}

func NewClient(addr string) (*Client, error) {
	if len(addr) == 0 {
		return nil, fmt.Errorf("addr must be set")
	}
	return &Client{addr}, nil
}

func (c *Client) Ping() (err error) {
	url := fmt.Sprintf("http://%s/api/v1/ping", c.addr)
	data, err := httpGet(url)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func (c *Client) CreateDid(req *io.CreateDidReq) (err error) {
	url := fmt.Sprintf("http://%s/api/v1/did/create", c.addr)

	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	_, err = httpPost(url, data)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RevoleDid(did string) (ddo *io.DDO, err error) {
	url := fmt.Sprintf("http://%s/api/v1/did/resove/%s", c.addr, did)

	data, err := httpGet(url)
	if err != nil {
		return nil, err
	}

	var resp io.ResolveDidResp
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Document, nil
}

func (c *Client) RevokeDid(req *io.RevokeDidReq) error {
	url := fmt.Sprintf("http://%s/api/v1/did/revoke", c.addr)

	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	_, err = httpPost(url, data)
	if err != nil {
		return err
	}

	return nil
}
