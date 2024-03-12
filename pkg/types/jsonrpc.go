package types

import (
	"encoding/json"
	"net/http"
)

const (
	// Version is JSON-RPC 2.0.
	Version = "2.0"

	BatchRequestKey  = '['
	ContentTypeKey   = "Content-Type"
	ContentTypeValue = "application/json"
)

type ErrCode int

const (
	ErrDefault ErrCode = iota

	RpcParseError     ErrCode = -32700
	RpcInvalidRequest ErrCode = -32600
	RpcMethodNotFound ErrCode = -32601
	RpcInvalidParams  ErrCode = -32602
)

var ErrToStatusCode = map[ErrCode]int{
	ErrDefault:        http.StatusInternalServerError,
	RpcParseError:     http.StatusInternalServerError,
	RpcInvalidRequest: http.StatusBadRequest,
	RpcMethodNotFound: http.StatusNotFound,
	RpcInvalidParams:  http.StatusBadRequest,
}

type JSONError struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type JSONRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
	ID      interface{}     `json:"id,omitempty"`
}

type JSONResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Error   *JSONError  `json:"error,omitempty"`
	Result  any         `json:"result,omitempty"`
	ID      interface{} `json:"id"`
}

func DecodeParams[T any](p json.RawMessage) (T, error) {
	var t T
	err := json.Unmarshal(p, &t)
	return t, err
}

func NewErrorResp(ID interface{}, err error, errCode ErrCode) *JSONResponse {
	return &JSONResponse{
		JSONRPC: "2.0",
		ID:      ID,
		Error: &JSONError{
			Code:    int(errCode),
			Message: err.Error(),
		},
	}
}

func NewResponse(r *JSONRequest) *JSONResponse {
	return &JSONResponse{
		JSONRPC: Version,
		ID:      r.ID,
	}
}
