package rpc

import (
	"fmt"
	"frost/pkg/collections"
	"frost/pkg/partyclient"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type server struct {
	router     *gin.Engine
	peerIpList *collections.OrderedList[partyclient.PartyClient]
}

type Server interface {
}

func NewServer(peerIpList *collections.OrderedList[partyclient.PartyClient]) Server {
	return &server{peerIpList: peerIpList}
}

func (s *server) Run(port string) error {
	s.router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	return s.router.Run(fmt.Sprintf(":%s", port))
}
