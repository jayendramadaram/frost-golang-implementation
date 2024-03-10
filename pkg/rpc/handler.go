package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"frost/pkg/types"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	ServeJSONRPC(c context.Context, params *json.RawMessage) (result any, err error)
}

type methodRecord struct {
	m sync.RWMutex
	r map[string]Metadata
}

type Metadata struct {
	Handler Handler
	Params  any
	Result  any
}

func NewMethodRecord() *methodRecord {
	return &methodRecord{
		m: sync.RWMutex{},
		r: map[string]Metadata{},
	}
}

func (mr *methodRecord) RegisterMethod(method string, h Handler, params, result any) error {
	if method == "" || h == nil {
		return errors.New("jsonrpc: method name and function should not be empty")
	}
	mr.m.Lock()
	mr.r[method] = Metadata{
		Handler: h,
		Params:  params,
		Result:  result,
	}
	mr.m.Unlock()

	return nil
}

func (mr *methodRecord) InvokeMethod(c context.Context, r *types.JSONRequest) *types.JSONResponse {
	if r.ID == nil {
		return types.NewErrorResp(1, errors.New("jsonrpc: id cannot be nil"), types.RpcInvalidRequest)
	}

	if r.Method == "" || r.JSONRPC != Version {
		return types.NewErrorResp(r.ID, errors.New("jsonrpc: invalid request"), types.RpcInvalidRequest)
	}

	var md Metadata
	res := NewResponse(r)

	mr.m.RLock()
	md, ok := mr.r[r.Method]
	mr.m.RUnlock()

	if !ok {
		return types.NewErrorResp(r.ID, errors.New("jsonrpc: method not found"), types.RpcMethodNotFound)
	}

	resp, err := md.Handler.ServeJSONRPC(c, &r.Params)
	if err != nil {
		return types.NewErrorResp(r.ID, err, types.RpcParseError)
	}

	res.Result = resp
	return res
}

func (mr *methodRecord) ServeHTTP(c *gin.Context) {
	r, batch, err := ParseRequest(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.NewErrorResp(nil, err, types.RpcParseError))
		return
	}

	resp := make([]*types.JSONResponse, len(r))
	for i := range r {
		resp[i] = mr.InvokeMethod(c.Request.Context(), r[i])
	}

	if batch || len(resp) > 1 {
		c.JSON(http.StatusOK, resp)
		return
	}
	c.JSON(http.StatusOK, resp[0])
}

func ParseRequest(r *http.Request) ([]*types.JSONRequest, bool, error) {
	var rerr error
	if !strings.HasPrefix(r.Header.Get(contentTypeKey), contentTypeValue) {
		return nil, false, fmt.Errorf("jsonrpc: invalid content type: %s", r.Header.Get(contentTypeKey))
	}

	buf := bytes.NewBuffer(make([]byte, 0, r.ContentLength))
	if _, err := buf.ReadFrom(r.Body); err != nil {
		return nil, false, fmt.Errorf("jsonrpc: failed to read body: %w", err)
	}
	defer func(r *http.Request) {
		err := r.Body.Close()
		if err != nil {
			rerr = fmt.Errorf("jsonrpc: failed to close body: %w", err)
		}
	}(r)

	if buf.Len() == 0 {
		return nil, false, fmt.Errorf("jsonrpc: empty request")
	}

	f, _, err := buf.ReadRune()
	if err != nil {
		return nil, false, fmt.Errorf("jsonrpc: failed to read request: %w", err)
	}
	if err := buf.UnreadRune(); err != nil {
		return nil, false, fmt.Errorf("jsonrpc: failed to read request: %w", err)
	}

	var rs []*types.JSONRequest
	if f != batchRequestKey {
		var req *types.JSONRequest
		if err := json.NewDecoder(buf).Decode(&req); err != nil {
			return nil, false, fmt.Errorf("jsonrpc: failed to decode request: %w", err)
		}

		return append(rs, req), false, nil
	}

	if err := json.NewDecoder(buf).Decode(&rs); err != nil {
		return nil, false, fmt.Errorf("jsonrpc: failed to decode request: %w", err)
	}

	return rs, true, rerr
}
