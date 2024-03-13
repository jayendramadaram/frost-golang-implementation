package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	client "frost/internal/sigag/sigagclient"
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

	SigAgClient client.SigAgClient
	store       Store
}

type Store interface {
	Lock()
	UnLock()
	IsLocked() bool
	NewEpoch(epoch uint) error
}

func NewServer(store Store, logger *logrus.Logger, SigAgClient client.SigAgClient) *server {
	return &server{store: store, router: gin.New(), logger: logger, SigAgClient: SigAgClient}
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

func (s *server) Ping(_ context.Context, _ *json.RawMessage) (json.RawMessage, error) {
	return json.Marshal(PingMessage{
		Message: "pong",
	})
}

func (s *server) NewEpoch(_ context.Context, params *json.RawMessage) (json.RawMessage, error) {
	if s.store.IsLocked() {
		return nil, fmt.Errorf("Epoch Already in progress")
	}

	if len(*params) == 0 {
		return nil, fmt.Errorf("params is nil")
	}

	var newEpoch NewEpochRequest
	if err := json.Unmarshal(*params, &newEpoch); err != nil {
		return nil, err
	}

	if err := rpc.Validate(newEpoch); err != nil {
		return nil, err
	}

	if err := s.store.NewEpoch(newEpoch.Epoch); err != nil {
		return nil, err
	}

	return json.Marshal(true)
}
