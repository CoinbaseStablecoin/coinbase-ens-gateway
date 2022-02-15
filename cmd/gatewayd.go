package main

import (
	"log"
	"os"
	"strconv"

	"github.cbhq.net/pete/coinbase-ens-gateway/server"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	port, _ := strconv.ParseUint(os.Getenv("PORT"), 10, 16)
	if port == 0 {
		port = 3000
	}

	privKeyHex := os.Getenv("PRIVATE_KEY")

	srv, err := server.New(r, uint16(port), privKeyHex)
	if err != nil {
		log.Fatalln("Failed to initialize server", err)
	}

	srv.Start()
}
