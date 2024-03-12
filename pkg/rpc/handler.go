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

type Handler func(c context.Context, params *json.RawMessage) (any, error)

type methodRecord struct {
	m sync.RWMutex
	r map[string]Handler
}

func NewMethodRecord() *methodRecord {
	return &methodRecord{
		m: sync.RWMutex{},
		r: map[string]Handler{},
	}
}

func (mr *methodRecord) RegisterMethod(method string, h Handler) error {
	if method == "" || h == nil {
		return errors.New("jsonrpc: method name and function should not be empty")
	}
	mr.m.Lock()
	mr.r[method] = h
	mr.m.Unlock()

	return nil
}

func (mr *methodRecord) InvokeMethod(c context.Context, r *types.JSONRequest) (*types.JSONResponse, int) {
	if r.ID == nil {
		return types.NewErrorResp(1, errors.New("jsonrpc: id cannot be nil"), types.RpcInvalidRequest), types.ErrToStatusCode[types.RpcInvalidRequest]
	}

	if r.Method == "" || r.JSONRPC != types.Version {
		return types.NewErrorResp(r.ID, errors.New("jsonrpc: invalid request"), types.RpcInvalidRequest), types.ErrToStatusCode[types.RpcInvalidRequest]
	}

	var md Handler
	res := types.NewResponse(r)

	mr.m.RLock()
	md, ok := mr.r[r.Method]
	mr.m.RUnlock()

	if !ok {
		return types.NewErrorResp(r.ID, errors.New("jsonrpc: method not found"), types.RpcMethodNotFound), types.ErrToStatusCode[types.RpcMethodNotFound]
	}

	resp, err := md(c, &r.Params)
	if err != nil {
		return types.NewErrorResp(r.ID, err, types.RpcParseError), types.ErrToStatusCode[types.RpcParseError]
	}

	res.Result = resp
	return res, http.StatusOK
}

func (mr *methodRecord) ServeHTTP(c *gin.Context) {
	r, batch, err := ParseRequest(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, types.NewErrorResp(nil, err, types.RpcParseError))
		return
	}

	partialSuccess := false
	var statusCode int

	resp := make([]*types.JSONResponse, len(r))
	for i := range r {
		resp[i], statusCode = mr.InvokeMethod(c.Request.Context(), r[i])
		if statusCode == http.StatusOK {
			partialSuccess = true
		}
	}

	if batch || len(resp) > 1 {
		if partialSuccess {
			c.JSON(http.StatusPartialContent, resp)
		} else {
			c.JSON(http.StatusOK, resp)
		}
	}
	c.JSON(statusCode, resp[0])
}

func ParseRequest(r *http.Request) ([]*types.JSONRequest, bool, error) {
	var rerr error
	if !strings.HasPrefix(r.Header.Get(types.ContentTypeKey), types.ContentTypeValue) {
		return nil, false, fmt.Errorf("jsonrpc: invalid content type: %s", r.Header.Get(types.ContentTypeKey))
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
	if f != types.BatchRequestKey {
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
