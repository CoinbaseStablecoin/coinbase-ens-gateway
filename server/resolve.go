package server

import "github.com/gin-gonic/gin"

type resolveParams struct {
	Sender string `uri:"sender" binding:"required"`
	Data   string `uri:"data" binding:"required"`
}

func (s *Server) resolve(c *gin.Context) {
	var params resolveParams

	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(400, gin.H{"error": "invalid params"})
		return
	}

	c.JSON(400, gin.H{"message": "not implemented"})
}
