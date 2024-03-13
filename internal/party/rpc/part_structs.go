package rpc

type PingMessage struct {
	Message string `json:"message"`
}

type NewEpochRequest struct {
	Epoch uint `json:"epoch,strict_check"`
}
