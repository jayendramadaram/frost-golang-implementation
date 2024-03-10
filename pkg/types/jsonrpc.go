package types

import (
	"encoding/json"
)

type ErrCode int

const (
	ErrCodeParseError ErrCode = iota

	RpcParseError     ErrCode = -32700
	RpcInvalidRequest ErrCode = -32600
	RpcMethodNotFound ErrCode = -32601
	RpcInvalidParams  ErrCode = -32602
)

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
