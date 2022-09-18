package api

import (
	"fmt"
	config "gossh/config/v1"
	"gossh/libs/clients"
	"gossh/libs/gin"
	models "gossh/models/v1"
	"gossh/utils"
	"strconv"
	"strings"
)

func VerifyHost(c *gin.Context) (models.Host, error) {
	name := c.PostForm("name")
	address := c.PostForm("address")
	user := c.PostForm("user")
	pwd := c.PostForm("pwd")
	port := c.PostForm("port")
	fontSize := c.PostForm("font_size")
	background := c.PostForm("background")
	foreground := c.PostForm("foreground")
	cursorColor := c.PostForm("cursor_color")
	fontFamily := c.PostForm("font_family")
	cursorStyle := c.PostForm("cursor_style")
	shell := c.PostForm("shell")

	if len(name) > 60 || len(name) == 0 {
		return models.Host{}, fmt.Errorf("name input error:%s.", name)
	}

	if len(address) > 60 || len(address) == 0 {
		return models.Host{}, fmt.Errorf("host input error")
	}

	if len(user) > 60 || len(user) == 0 {
		return models.Host{}, fmt.Errorf("user input error")
	}

	if len(pwd) > 60 || len(pwd) == 0 {
		return models.Host{}, fmt.Errorf("pwd input error")
	}
	p, err := strconv.Atoi(strings.TrimSpace(port))
	if err != nil {
		return models.Host{}, fmt.Errorf("port input error")
	}
	if p > 65535 || p < 1 {
		return models.Host{}, fmt.Errorf("port range input error")
	}

	fontsize, err := strconv.Atoi(strings.TrimSpace(fontSize))
	if err != nil {
		fontsize = 16
	}
	if fontsize > 32 || fontsize < 8 {
		fontsize = 16
	}
	if len(strings.TrimSpace(background)) == 0 {
		background = "#000000"
	}

	if len(strings.TrimSpace(foreground)) == 0 {
		foreground = "#FFFFFF"
	}

	if len(strings.TrimSpace(cursorColor)) == 0 {
		cursorColor = "#FFFFFF"
	}

	if len(strings.TrimSpace(fontFamily)) == 0 {
		fontFamily = "Courier"
	}

	if len(strings.TrimSpace(cursorStyle)) == 0 {
		cursorStyle = "block"
	}

	if len(strings.TrimSpace(shell)) == 0 {
		shell = "bash"
	}

	hostInfo := models.Host{
		Name:        name,
		Address:     address,
		User:        user,
		Pwd:         pwd,
		Port:        p,
		FontSize:    fontsize,
		Background:  background,
		Foreground:  foreground,
		CursorColor: cursorColor,
		FontFamily:  fontFamily,
		CursorStyle: cursorStyle,
		Shell:       shell,
	}
	return hostInfo, nil
}

func GetAllHost(c *gin.Context) {
	var host *models.Host
	allHost, err := host.Select()
	if err != nil {
		c.JSON(500, gin.H{"code": config.FAILURE, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": config.SUCCEED, "data": allHost, "msg": "ok"})
}

func AddHost(c *gin.Context) {
	var host *models.Host
	h, err := VerifyHost(c)
	if err != nil {
		c.JSON(400, gin.H{"code": config.FAILURE, "msg": err.Error()})
		return
	}
	_, err = host.Insert(h.Name, h.Address, h.User, h.Pwd, h.Port, h.FontSize, h.Background, h.Foreground, h.CursorColor, h.FontFamily, h.CursorStyle, h.Shell)
	if err != nil {
		c.JSON(500, gin.H{"code": config.FAILURE, "msg": err.Error()})
		return
	}
	GetAllHost(c)
}

func UpdateHost(c *gin.Context) {
	var host *models.Host
	h, err := VerifyHost(c)
	if err != nil {
		c.JSON(400, gin.H{"code": config.FAILURE, "msg": err.Error()})
		return
	}
	id, err := strconv.Atoi(strings.TrimSpace(c.PostForm("id")))
	if err != nil {
		c.JSON(400, gin.H{"code": config.FAILURE, "msg": err.Error()})
		return
	}

	_, err = host.Update(id, h.Name, h.Address, h.User, h.Pwd, h.Port, h.FontSize, h.Background, h.Foreground, h.CursorColor, h.FontFamily, h.CursorStyle, h.Shell)
	if err != nil {
		c.JSON(500, gin.H{"code": config.FAILURE, "msg": err.Error()})
		return
	}
	GetAllHost(c)
}

func DeleteHost(c *gin.Context) {
	var host *models.Host
	id, err := strconv.Atoi(strings.TrimSpace(c.PostForm("id")))
	if err != nil {
		c.JSON(400, gin.H{"code": config.FAILURE, "msg": err.Error()})
		return
	}
	_, err = host.Delete(id)
	if err != nil {
		c.JSON(500, gin.H{"code": config.FAILURE, "msg": err.Error()})
		return
	}
	GetAllHost(c)
}

func GetSessionId(c *gin.Context) {
	h, err := VerifyHost(c)
	if err != nil {
		c.JSON(400, gin.H{"code": config.FAILURE, "msg": err.Error()})
		return
	}
	sessionId := utils.RandString(15)
	client := NewClient(h.Address, h.User, h.Pwd, h.Port, h.Shell, sessionId)

	clients.Lock()
	clients.SetData(sessionId, client)
	clients.Unlock()
	c.JSON(200, gin.H{"code": config.SUCCEED, "data": sessionId, "msg": "ok"})
}
