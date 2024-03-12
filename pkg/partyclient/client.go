package partyclient

type PartyClient interface {
	ping() error
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

func (c *partyclient) ping() error {
	return nil
}
