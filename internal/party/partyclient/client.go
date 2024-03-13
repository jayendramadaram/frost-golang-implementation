package partyclient

import "fmt"

type PartyClient interface {
	Ping() error
	Locate() (string, string)
	ID() string
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
	return nil
}

func (c *partyclient) Locate() (string, string) {
	return c.id, fmt.Sprintf("%s:%s:%s", c.ip, c.port, c.path)
}

func (c *partyclient) ID() string {
	return c.id
}
