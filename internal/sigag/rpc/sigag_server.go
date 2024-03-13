package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"frost/internal/party/partyclient"
	"frost/pkg/collections"
	"frost/pkg/rpc"
	"net"
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
	Register(_ context.Context, params *json.RawMessage) (bool, error)
	Health(_ context.Context, params *json.RawMessage) (any, error)
	GetParties(_ context.Context, params *json.RawMessage) (Parties, error)
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
			if result[1].IsNil() {
				return result[0].Interface(), nil
			}
			return result[0], result[1].Interface().(error)
		}

		if err := mr.RegisterMethod(methodName, handler); err != nil {
			return err
		}

		s.logger.Info("registered rpc method", zap.String("method", methodName))

	}

	s.router.POST("/", func(c *gin.Context) {
		mr.ServeHTTP(c)
	})

	s.logger.Info("listening on port", zap.String("port", port))
	return s.router.Run(fmt.Sprintf("127.0.0.1:%s", port))
}

// concurrent safe
func (s *server) Register(_ context.Context, params *json.RawMessage) (bool, error) {
	if len(*params) == 0 {
		return false, fmt.Errorf("params is nil")
	}

	var registerParty RegisterParty
	if err := json.Unmarshal(*params, &registerParty); err != nil {
		return false, err
	}

	if err := rpc.Validate(registerParty); err != nil {
		return false, err
	}

	ip := net.ParseIP(registerParty.ReportedIp)
	if ip == nil {
		return false, fmt.Errorf("invalid ip address: %s", registerParty.ReportedIp)
	}

	participant := partyclient.New(registerParty.Address, registerParty.ReportedIp, registerParty.Port, registerParty.Path)

	containsID := func(item, element partyclient.PartyClient) bool {
		return item.ID() == element.ID()
	}

	if s.peerIpList.Contains(participant, containsID) {
		return false, fmt.Errorf("address already registered")
	}

	if err := participant.Ping(); err != nil {
		return false, err
	}

	s.peerIpList.Add(participant)

	return true, nil
}

func (s *server) Health(_ context.Context, params *json.RawMessage) (any, error) {
	health := HealthCheck{
		Status: "ok",
	}

	return health, nil
}

func (s *server) GetParties(_ context.Context, params *json.RawMessage) (Parties, error) {
	Parties := make(Parties)
	for _, v := range s.peerIpList.Items {
		id, url := v.Locate()
		Parties[id] = url
	}

	return Parties, nil
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
