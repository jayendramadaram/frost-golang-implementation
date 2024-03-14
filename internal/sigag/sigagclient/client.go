package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"frost/internal/sigag/rpc"
	"frost/pkg/types"
	"net/http"
)

type SigAgClient interface {
	Register(id, url string, noTLS bool) error
	GetParticipants() (rpc.Parties, error)
	CheckUptime() (bool, error)
}

type client struct {
	url string
	jwt string
}

func New(url string) SigAgClient {
	return &client{
		url: url,
	}
}

func (c *client) Register(id, url string, noTLS bool) error {
	var params = rpc.RegisterParty{
		Address: id,
		Url:     url,
		NoTLS:   noTLS,
	}
	err := c.SendRequest("register", params, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *client) CheckUptime() (bool, error) {
	var reponse rpc.HealthCheck
	err := c.SendRequest("health", nil, &reponse)
	if err != nil {
		return false, err
	}
	fmt.Printf("sigag: %s\n", reponse.Status)

	return reponse.Status == "ok", nil
}

func (c *client) GetParticipants() (rpc.Parties, error) {
	var reponse rpc.Parties
	err := c.SendRequest("get_parties", nil, &reponse)
	if err != nil {
		return nil, err
	}

	return reponse, nil
}

func (c *client) SendRequest(method string, params, respType interface{}) error {

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

	// url := fmt.Sprintf("http://%s:%s%s", c.Ip, c.Port, c.Path)
	req, err := http.NewRequest(http.MethodPost, c.url, bytes.NewReader(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.jwt)

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
