package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"frost/pkg/rpc"
	"reflect"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

type server struct {
	logger *logrus.Logger
	router *gin.Engine
	store  Store
}

type Store interface {
	AddParticipant(RegisterParty) error
	GetParties() Parties
	IsLocked() bool

	GetEpochParties() Parties
}

func NewServer(store Store, logger *logrus.Logger) *server {
	return &server{store: store, router: gin.New(), logger: logger}
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

		methodName := rpc.ConvertToSnakeCase(method.Name)

		handler := func(c context.Context, params *json.RawMessage) (json.RawMessage, error) {
			result := method.Func.Call([]reflect.Value{reflect.ValueOf(s), reflect.ValueOf(c), reflect.ValueOf(params)})
			if result[1].IsNil() {
				return result[0].Interface().(json.RawMessage), nil
			}
			return result[0].Interface().(json.RawMessage), result[1].Interface().(error)
		}

		if err := mr.RegisterMethod(methodName, handler); err != nil {
			return err
		}

		s.logger.Info("registered rpc method", zap.String("method", methodName))

	}

	s.router.Use(rpc.RequestLoggingMiddleware(s.logger))
	// s.router.Use(gin.LoggerWithWriter(s.logger.Writer()))

	s.router.POST("/", func(c *gin.Context) {
		mr.ServeHTTP(c)
	})

	s.logger.Info("listening on port", zap.String("port", port))
	return s.router.Run(fmt.Sprintf("127.0.0.1:%s", port))
}

// concurrent safe
func (s *server) Register(_ context.Context, params *json.RawMessage) (json.RawMessage, error) {

	if s.store.IsLocked() {
		return nil, fmt.Errorf("DKG in progress, cant accept registration ATM")
	}

	if len(*params) == 0 {
		return nil, fmt.Errorf("params is nil")
	}

	var registerParty RegisterParty
	if err := json.Unmarshal(*params, &registerParty); err != nil {
		return nil, err
	}

	if err := rpc.Validate(registerParty); err != nil {
		return nil, err
	}

	// ip := net.ParseIP(registerParty.ReportedIp)
	// if ip == nil {
	// 	return nil, fmt.Errorf("invalid ip address: %s", registerParty.ReportedIp)
	// }

	if err := s.store.AddParticipant(registerParty); err != nil {
		return nil, err
	}

	return json.Marshal(true)
}

func (s *server) Health(_ context.Context, params *json.RawMessage) (json.RawMessage, error) {
	health := HealthCheck{
		Status: "ok",
	}

	return json.Marshal(health)
}

func (s *server) GetParties(_ context.Context, params *json.RawMessage) (json.RawMessage, error) {
	return json.Marshal(s.store.GetParties())
}

func (s *server) GetEpochParties(_ context.Context, params *json.RawMessage) (json.RawMessage, error) {
	return json.Marshal(s.store.GetEpochParties())
}
