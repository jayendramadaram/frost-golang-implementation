package partyclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"frost/internal/party/rpc"
	"frost/pkg/types"
	"net/http"
)

type PartyClient interface {
	ID() string
	Ping() error
	Locate() (string, string)

	NewEpoch(epoch uint) error
}

type partyclient struct {
	id   string
	ip   string
	port string
	path string
}

func New(id, ip, port, path string) PartyClient {
	return &partyclient{id: id, ip: ip, port: port, path: path}
}

func (c *partyclient) Ping() error {
	var PingMessage rpc.PingMessage
	return c.SendRequest("ping", nil, PingMessage)
}

func (c *partyclient) Locate() (string, string) {
	return c.id, fmt.Sprintf("%s:%s:%s", c.ip, c.port, c.path)
}

func (c *partyclient) ID() string {
	return c.id
}

// NewEpoch implements PartyClient.
func (c *partyclient) NewEpoch(epoch uint) error {
	NewEpoch := rpc.NewEpochRequest{
		Epoch: epoch,
	}
	if err := c.SendRequest("new_epoch", NewEpoch, nil); err != nil {
		return err
	}
	return nil
}

func (c *partyclient) SendRequest(method string, params, respType interface{}) error {

	paramsData, err := json.Marshal(params)
	if err != nil {
		return err
	}

	reqObject := types.JSONRequest{
		JSONRPC: types.Version,
		Method:  method,
		Params:  paramsData,
		ID:      1,
	}

	jsonData, err := json.Marshal(reqObject)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s:%s%s", c.ip, c.port, c.path)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	var response types.JSONResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("sigag: %s", response.Error.Message)
	}

	return json.Unmarshal(response.Result, &respType)
}
