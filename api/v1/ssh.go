package api

import (
	"fmt"
	config "gossh/config/v1"
	"gossh/libs/sessions"
	"gossh/libs/sse"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"gossh/libs/clients"

	"gossh/libs/websocket"
)

// GET (需要登陆认证)获取已经连接的主机信息
func GetSshStatus(c *gin.Context) {
	session := sessions.Default(c)

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Connection", "keep-alive")
	c.Header("Cache-Control", "no-cache")
	c.Header("Content-Type", "text/event-stream")

	if session.Get("auth") != "Y" {
		c.Render(200, sse.Event{
			Id:    "401",
			Event: "message",
			Retry: 10000,
			Data: map[string]interface{}{
				"code": config.FAILURE,
				"msg":  "Unauthorized",
			},
		})
		return
	}

	for {
		c.Writer.(http.Flusher).Flush()
		time.Sleep(time.Second * 3)
		var data []clients.Ssh
		for _, item := range clients.GetData() {
			data = append(data, *item)
		}
		c.Render(200, sse.Event{
			Id:    "200",
			Event: "message",
			Retry: 10000,
			Data: map[string]interface{}{
				"code": config.SUCCEED,
				"data": data,
				"msg":  "ok",
			},
		})
	}

}

// POST 更新已经连接的主机信息
func UpdateSshStatus(c *gin.Context) {
	ids := c.PostFormArray("ids")
	for _, key := range ids {
		val, ok := clients.GetClientBySessionID(key)
		if ok {
			val.Timeout = time.Now()
		}
	}
	c.JSON(200, gin.H{
		"code": config.SUCCEED,
		"data": ids,
		"msg":  "ok",
	})
}

// DELETE (需要登陆认证)删除已经建立的连接
func DeleteSshConnect(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("auth") != "Y" {
		c.JSON(401, gin.H{
			"code": config.FAILURE,
			"msg":  "Unauthorized",
		})
		return
	}

	defer func() {
		if err := recover(); err != nil {
			c.JSON(500, gin.H{
				"code": config.FAILURE,
				"msg":  "delete connect error",
			})
		}
	}()

	sessionId := c.Query("session_id")
	if sessionId == "" {
		c.JSON(404, gin.H{
			"code": config.FAILURE,
			"msg":  "session not exists",
		})
		return
	}

	sshConn, ok := clients.GetClientBySessionID(sessionId)
	if ok {
		sshConn.SshClient.Close()
		sshConn.SftpClient.Close()
		sshConn.SshSession.Close()
		sshConn.Ws.Close()
		clients.DeleteClientBySessionID(sessionId)
	}

	c.JSON(200, gin.H{
		"code": config.SUCCEED,
		"msg":  "delete connect success",
	})
}

func SshHandler(c *gin.Context) {
	var request = c.Request
	var response = c.Writer

	// 调整窗口大小
	if request.Method == http.MethodPatch {
		Resize(c)
		return
	}

	// WebSock 连接 SSH
	websocket.Handler(func(ws *websocket.Conn) {
		sessionId := ws.Request().URL.Query().Get("session_id")
		w, err := strconv.Atoi(ws.Request().URL.Query().Get("w"))
		if err != nil || (w < 40 || w > 8192) {
			websocket.Message.Send(ws, "connect error window width !!!")
			clients.DeleteClientBySessionID(sessionId)
			ws.Close()
			return
		}
		h, err := strconv.Atoi(ws.Request().URL.Query().Get("h"))
		if err != nil || (h < 2 || h > 4096) {
			websocket.Message.Send(ws, "connect error window height !!!")
			clients.DeleteClientBySessionID(sessionId)
			ws.Close()
			return
		}

		cli, _ := clients.GetClientBySessionID(sessionId)

		err = cli.RunTerminal(cli.Shell, ws, ws, ws, w, h, ws)
		if err != nil {
			websocket.Message.Send(ws, "connect error!!!")
			clients.DeleteClientBySessionID(sessionId)
			ws.Close()
			return
		}
	}).ServeHTTP(response, request)
}

// Resize 调整终端大小
func Resize(c *gin.Context) {
	w, err := strconv.Atoi(c.Query("w"))
	if err != nil || (w < 40 || w > 8192) {
		c.JSON(400, gin.H{
			"code": config.FAILURE,
			"msg":  fmt.Sprintf("connect error window width !!!")})
		return
	}
	h, err := strconv.Atoi(c.Query("h"))
	if err != nil || (h < 2 || h > 4096) {
		c.JSON(400, gin.H{
			"code": config.FAILURE,
			"msg":  fmt.Sprintf("connect error window width !!!")})
		return
	}

	sessionId := c.Query("session_id")

	cli, ok := clients.GetClientBySessionID((sessionId))

	if !ok || cli == nil {
		c.JSON(299, gin.H{"code": config.FAILURE, "msg": "the client is disconnected"})
		return
	}

	if cli.SshSession != nil {
		_ = cli.SshSession.WindowChange(h, w)
		str := fmt.Sprintf("W:%d;H:%d\n", w, h)
		c.JSON(200, gin.H{"code": config.SUCCEED, "data": str, "msg": "ok"})
		return
	}
}
