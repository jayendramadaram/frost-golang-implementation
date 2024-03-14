package rpc

import (
	sigagrpc "frost/internal/sigag/rpc"
)

type PingMessage struct {
	Message string `json:"message"`
}

type NewEpochRequest struct {
	Epoch uint `json:"epoch,strict_check"`
}

type DKGInitRequest struct {
	Parties   sigagrpc.Parties `json:"epoch,strict_check"`
	Threshold uint             `json:"threshold,strict_check"`
}
