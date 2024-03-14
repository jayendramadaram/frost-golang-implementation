package rpc

type Parties map[string]string

type RegisterParty struct {
	Address string `json:"address,strict_check"`
	Url     string `json:"url,strict_check"`

	NoTLS bool `json:"no_tls"`
}

type HealthCheck struct {
	Status string `json:"status"`
}
