package router

import (
	"gowebssh/api"

	"github.com/gin-gonic/gin"
)

// router
var engine = gin.Default()

func Run() {
	engine.POST("/api/getSessionId", api.GetSessionId)
	engine.GET("/ws/ssh", api.SshHandler)

	engine.Run("0.0.0.0:8080")
}
