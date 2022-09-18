package router

import (
	"fmt"
	api "gossh/api/v1"
	config "gossh/config/v1"
	"gossh/libs/gin"
	"gossh/libs/logger"
	"gossh/middlewares"
	"net/http"
	"os"
)

var engine = gin.Default()

func Run() {
	var err error

	//加入session中间件
	engine.Use(middlewares.Session())
	engine.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/gogossh/")
	})

	engine.POST("/api/login", api.Login)
	engine.PATCH("/api/revise-password", api.RevisePassword)

	engine.GET("/api/ssh/status", api.GetSshStatus)
	engine.POST("/api/ssh/status", api.UpdateSshStatus)
	engine.DELETE("/api/ssh/status", api.DeleteSshConnect)

	engine.GET("/api/host", api.GetAllHost)
	engine.POST("/api/host", api.AddHost)
	engine.PUT("/api/host", api.UpdateHost)
	engine.DELETE("/api/host", api.DeleteHost)

	engine.POST("/api/getSessionId", api.GetSessionId)

	engine.GET("/api/sftp-dir", api.SftpDir)
	engine.POST("/api/sftp-download", api.SftpDownload)
	engine.POST("/api/sftp-upload", api.SftpUpload)

	engine.Any("/ws/ssh", api.SshHandler)

	// 证书加密
	address := fmt.Sprintf("%s:%s", config.Config["server"]["Address"], config.Config["server"]["Port"])

	certFile := config.Config["server"]["CertFile"]
	keyFile := config.Config["server"]["KeyFile"]

	_, certErr := os.Open(certFile)
	_, keyErr := os.Open(keyFile)

	// 如果证书和私钥文件存在,就使用https协议,否则使用http协议
	if certErr == nil && keyErr == nil {
		logger.Logger.Debug("https://{IP}:" + config.Config["server"]["Port"])
		err = engine.RunTLS(address, certFile, keyFile)
		if err != nil {
			logger.Logger.Error("RunServeTLSError:", err.Error())
			os.Exit(1)
			return
		}
	} else {
		logger.Logger.Debug("http://{IP}:" + config.Config["server"]["Port"])
		err = engine.Run(address)
		if err != nil {
			logger.Logger.Error("RunServeError:", err.Error())
			os.Exit(1)
			return
		}
	}
}
