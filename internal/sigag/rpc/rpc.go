package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"frost/pkg/collections"
	"frost/pkg/partyclient"
	"frost/pkg/rpc"
	"reflect"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type server struct {
	logger     *zap.Logger
	router     *gin.Engine
	peerIpList *collections.OrderedList[partyclient.PartyClient]
}

type Server interface {
	Run(port string) error
	RegisterParty(_ context.Context, params *json.RawMessage) (any, error)
}

func NewServer(peerIpList *collections.OrderedList[partyclient.PartyClient], logger *zap.Logger) Server {
	return &server{peerIpList: peerIpList, router: gin.Default(), logger: logger}
}

func (s *server) Run(port string) error {
	s.router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	mr := rpc.NewMethodRecord()

	s.logger.Info("registering rpc methods", zap.Any("methods", mr))

	serverType := reflect.TypeOf(s)
	for i := 0; i < serverType.NumMethod(); i++ {
		method := serverType.Method(i)
		if method.Type.NumIn() != 3 || method.Type.NumOut() != 2 {
			continue
		}

		methodName := ConvertToSnakeCase(method.Name)

		handler := func(c context.Context, params *json.RawMessage) (any, error) {
			result := method.Func.Call([]reflect.Value{reflect.ValueOf(s), reflect.ValueOf(c), reflect.ValueOf(params)})
			return result[0].Interface(), result[1].Interface().(error)
		}

		if err := mr.RegisterMethod(methodName, handler); err != nil {
			return err
		}

		s.logger.Info("registered rpc method", zap.String("method", methodName))

	}

	s.logger.Info("listening on port", zap.String("port", port))
	return s.router.Run(fmt.Sprintf("127.0.0.1:%s", port))
}

type RegisterParty struct {
	Address    string `json:"address"`
	ReportedIp string `json:"ip"`
	Port       string `json:"port"`
	Path       string `json:"path"`
}

func (s *server) RegisterParty(_ context.Context, params *json.RawMessage) (any, error) {
	if params == nil {
		return false, fmt.Errorf("params is nil")
	}

	var registerParty RegisterParty
	if err := json.Unmarshal(*params, &registerParty); err != nil {
		return false, err
	}

	if registerParty.Address == "" || registerParty.ReportedIp == "" || registerParty.Port == "" || registerParty.Path == "" {
		return false, fmt.Errorf("invalid params")
	}

	participant := partyclient.New(registerParty.Address, registerParty.ReportedIp, registerParty.Port, registerParty.Path)

	if s.peerIpList.Contains(participant) {
		return false, fmt.Errorf("address already registered")
	}

	s.peerIpList.Add(participant)

	return true, nil
}

func ConvertToSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}
