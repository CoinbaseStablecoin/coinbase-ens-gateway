package server

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type Server struct {
	engine  *gin.Engine
	port    uint16
	privKey *ecdsa.PrivateKey
}

func New(engine *gin.Engine, port uint16, privKeyHex string) (*Server, error) {
	privKeyBytes, err := hexutil.Decode(privKeyHex)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode private key")
	}
	privKey, err := crypto.ToECDSA(privKeyBytes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize private key")
	}

	srv := &Server{engine, port, privKey}

	engine.GET("/r/:sender/:data", srv.GetResolve)

	return srv, nil
}

func (s *Server) Start() {
	fmt.Println("Signer address:", s.SignerAddress())
	s.engine.Run(fmt.Sprintf("0.0.0.0:%d", s.port))
}

func (s *Server) SignerAddress() string {
	pubKey := s.privKey.Public()
	ecdsaPubKey, _ := pubKey.(*ecdsa.PublicKey)

	return crypto.PubkeyToAddress(*ecdsaPubKey).Hex()
}
