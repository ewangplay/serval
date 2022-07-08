package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/ewangplay/serval/io"
)

type HttpClient struct {
	c *http.Client
}

func NewHttpClient() (*HttpClient, error) {
	c := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Second*3)
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(time.Second * 3))
				return conn, nil
			},
			ResponseHeaderTimeout: time.Second * 3,
		},
	}
	return &HttpClient{c}, nil
}

func (c *HttpClient) Post(url string, data []byte) ([]byte, error) {
	body := bytes.NewBuffer(data)
	resp, err := c.c.Post(url, "application/json;charset=utf-8", body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return c.parseResponse(resp)
}

func (c *HttpClient) Get(url string) ([]byte, error) {
	resp, err := c.c.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return c.parseResponse(resp)
}

func (c *HttpClient) parseResponse(resp *http.Response) ([]byte, error) {
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var r io.Response
	err = json.Unmarshal(respBody, &r)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%v: %v - %v", resp.Status, r.Code, r.Msg)
	}

	if r.Code != 0 {
		return nil, fmt.Errorf("%v - %v", r.Code, r.Msg)
	}

	rData, err := json.Marshal(r.Data)
	if err != nil {
		return nil, err
	}

	return rData, nil
}
