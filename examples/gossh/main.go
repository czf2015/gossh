package main

import (
	"fmt"
	"net/http"
	"strconv"

	"math/rand"
	"time"

	clients "github.com/czf2015/gopkg-clients"
	websocket "github.com/czf2015/gopkg-websocket"
	"github.com/gin-gonic/gin"
)

// config
const (
	SUCCEED = 0
	FAILURE = 1
)

// models
type Terminal struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	User        string `json:"user"`
	Pwd         string `json:"pwd"`
	Port        int    `json:"port"`
	Shell       string `json:"shell"`
	FontFamily  string `json:"font_family"`
	FontSize    int    `json:"font_size"`
	Foreground  string `json:"foreground"`
	Background  string `json:"background"`
	CursorColor string `json:"cursor_color"`
	CursorStyle string `json:"cursor_style"`
}

// /* utils. */ RandString 生成指定长度随机字符串
func RandString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	data := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, data[r.Intn(len(data))])
	}
	return string(result)
}

// api
func SshHandler(c *gin.Context) {
	var request = c.Request
	var response = c.Writer

	if request.Method == http.MethodPatch {
		Resize(c)
		return
	}
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
			"code":/* config. */ FAILURE,
			"msg": fmt.Sprintf("connect error window width !!!")})
		return
	}
	h, err := strconv.Atoi(c.Query("h"))
	if err != nil || (h < 2 || h > 4096) {
		c.JSON(400, gin.H{
			"code":/* config. */ FAILURE,
			"msg": fmt.Sprintf("connect error window width !!!")})
		return
	}

	sessionId := c.Query("session_id")

	cli, ok := clients.GetClientBySessionID((sessionId))

	if !ok || cli == nil {
		c.JSON(299, gin.H{"code": /* config. */ FAILURE, "msg": "the client is disconnected"})
		return
	}

	if cli.SshSession != nil {
		_ = cli.SshSession.WindowChange(h, w)
		str := fmt.Sprintf("W:%d;H:%d\n", w, h)
		c.JSON(200, gin.H{"code": /* config. */ SUCCEED, "data": str, "msg": "ok"})
		return
	}
}

func VerifyTerminal(c *gin.Context) ( /* models. */ Terminal, error) {
	var terminal /* models. */ Terminal
	if err := c.ShouldBindJSON(&terminal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if len(terminal.Address) > 60 || len(terminal.Address) == 0 {
		return /* models. */ Terminal{}, fmt.Errorf("terminal input error")
	}

	if len(terminal.User) > 60 || len(terminal.User) == 0 {
		return /* models. */ Terminal{}, fmt.Errorf("user input error")
	}

	if len(terminal.Pwd) > 60 || len(terminal.Pwd) == 0 {
		return /* models. */ Terminal{}, fmt.Errorf("pwd input error")
	}

	if terminal.Port > 65535 || terminal.Port < 1 {
		return /* models. */ Terminal{}, fmt.Errorf("port range input error")
	}

	terminal.FontFamily = "Courier"
	terminal.FontSize = 16
	terminal.Foreground = "#FFFFFF"
	terminal.Background = "#000000"
	terminal.CursorColor = "#FFFFFF"
	terminal.CursorStyle = "block"
	terminal.Shell = "bash"

	return terminal, nil
}

func GetSessionId(c *gin.Context) {
	h, err := VerifyTerminal(c)
	if err != nil {
		c.JSON(400, gin.H{"code": /* config. */ FAILURE, "msg": err.Error()})
		return
	}
	sessionId := /* utils. */ RandString(15)
	clients.AddData(h.Address, h.User, h.Pwd, h.Port, h.Shell, sessionId)
	c.JSON(200, gin.H{"code": /* config. */ SUCCEED, "data": sessionId, "msg": "ok"})
}

// router
func Run() {
	var engine = gin.Default()
	engine.POST("/api/getSessionId" /* api. */, GetSessionId)
	engine.GET("/ws/ssh" /* api. */, SshHandler)

	engine.Run("0.0.0.0:8888")
}

func main() {
	/* router. */ Run()
}
