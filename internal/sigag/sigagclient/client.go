package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"frost/pkg/types"
	"net/http"
)

type SigAgClient interface {
	Register() error
	CheckUptime() (bool, error)
}

type client struct {
	Ip   string
	Port string
	Path string
	jwt  string
}

func New(ip, port, path string) SigAgClient {
	return &client{
		Ip:   ip,
		Port: port,
		Path: path,
	}
}

func (c *client) Register() error {
	return nil
}

func (c *client) CheckUptime() (bool, error) {
	resp, err := c.SendRequest("health", nil)
	if err != nil {
		return false, err
	}

	health := resp.(map[string]interface{})
	return health["status"] == "ok", nil
}

func (c *client) SendRequest(method string, params []byte) (any, error) {

	reqObject := types.JSONRequest{
		JSONRPC: types.Version,
		Method:  method,
		Params:  params,
		ID:      1,
	}

	jsonData, err := json.Marshal(reqObject)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("http://%s:%s%s", c.Ip, c.Port, c.Path)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.jwt)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var response types.JSONResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("sigag: %s", response.Error.Message)
	}

	return response.Result, nil
}
