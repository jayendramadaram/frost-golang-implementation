package rpc

type Parties map[string]string

type RegisterParty struct {
	Address    string `json:"address,strict_check"`
	ReportedIp string `json:"ip,strict_check"`
	Port       string `json:"port,strict_check"`
	Path       string `json:"path,strict_check"`
}

type HealthCheck struct {
	Status string `json:"status"`
}
