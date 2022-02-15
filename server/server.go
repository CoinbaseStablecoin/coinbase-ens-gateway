package server

import (
	"fmt"

	"github.cbhq.net/pete/coinbase-ens-gateway/gateway"
	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine
	port   uint16
	gw     *gateway.Gateway
}

func New(engine *gin.Engine, port uint16, privKeyHex string) (*Server, error) {
	gw, err := gateway.New(privKeyHex)
	if err != nil {
		return nil, err
	}

	srv := &Server{engine, port, gw}

	engine.GET("/r/:sender/:data", srv.resolve)

	return srv, nil
}

func (s *Server) Start() {
	fmt.Println("Signer address:", s.gw.SignerAddress())
	s.engine.Run(fmt.Sprintf("0.0.0.0:%d", s.port))
}
