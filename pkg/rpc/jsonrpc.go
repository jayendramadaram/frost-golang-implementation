package rpc

import (
	ty "frost/pkg/types"
)

const (
	// Version is JSON-RPC 2.0.
	Version = "2.0"

	batchRequestKey  = '['
	contentTypeKey   = "Content-Type"
	contentTypeValue = "application/json"
)

func NewResponse(r *ty.JSONRequest) *ty.JSONResponse {
	return &ty.JSONResponse{
		JSONRPC: Version,
		ID:      r.ID,
	}
}
