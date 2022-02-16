package server

import (
	"log"

	"github.com/gin-gonic/gin"
)

type resolveParams struct {
	Sender string `uri:"sender" binding:"required"`
	Data   string `uri:"data" binding:"required"`
}

func (s *Server) resolve(c *gin.Context) {
	var params resolveParams

	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(400, gin.H{"error": "invalid_params"})
		return
	}

	result, err := s.gw.Resolve(params.Sender, params.Data)
	if err != nil {
		log.Println("failed to resolve:", err)
		c.JSON(400, gin.H{"error": "invalid_params"})
		return
	}
	if len(result) == 0 {
		c.JSON(404, gin.H{"error": "not_found"})
		return
	}

	c.JSON(400, gin.H{"error": "not_implemented"})
}
