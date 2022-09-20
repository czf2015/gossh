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

func VerifyTerminal(c *gin.Context) (models.Terminal, error) {
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
		return models.Terminal{}, fmt.Errorf("name input error:%s.", name)
	}

	if len(address) > 60 || len(address) == 0 {
		return models.Terminal{}, fmt.Errorf("terminal input error")
	}

	if len(user) > 60 || len(user) == 0 {
		return models.Terminal{}, fmt.Errorf("user input error")
	}

	if len(pwd) > 60 || len(pwd) == 0 {
		return models.Terminal{}, fmt.Errorf("pwd input error")
	}
	p, err := strconv.Atoi(strings.TrimSpace(port))
	if err != nil {
		return models.Terminal{}, fmt.Errorf("port input error")
	}
	if p > 65535 || p < 1 {
		return models.Terminal{}, fmt.Errorf("port range input error")
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

	terminal := models.Terminal{
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
	return terminal, nil
}

func GetAllTerminals(c *gin.Context) {
	var terminal *models.Terminal
	allTerminal, err := terminal.Select()
	if err != nil {
		c.JSON(500, gin.H{"code": config.FAILURE, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": config.SUCCEED, "data": allTerminal, "msg": "ok"})
}

func AddTerminal(c *gin.Context) {
	var terminal *models.Terminal
	h, err := VerifyTerminal(c)
	if err != nil {
		c.JSON(400, gin.H{"code": config.FAILURE, "msg": err.Error()})
		return
	}
	_, err = terminal.Insert(h.Name, h.Address, h.User, h.Pwd, h.Port, h.FontSize, h.Background, h.Foreground, h.CursorColor, h.FontFamily, h.CursorStyle, h.Shell)
	if err != nil {
		c.JSON(500, gin.H{"code": config.FAILURE, "msg": err.Error()})
		return
	}
	GetAllTerminals(c)
}

func UpdateTerminal(c *gin.Context) {
	var terminal *models.Terminal
	h, err := VerifyTerminal(c)
	if err != nil {
		c.JSON(400, gin.H{"code": config.FAILURE, "msg": err.Error()})
		return
	}
	id, err := strconv.Atoi(strings.TrimSpace(c.PostForm("id")))
	if err != nil {
		c.JSON(400, gin.H{"code": config.FAILURE, "msg": err.Error()})
		return
	}

	_, err = terminal.Update(id, h.Name, h.Address, h.User, h.Pwd, h.Port, h.FontSize, h.Background, h.Foreground, h.CursorColor, h.FontFamily, h.CursorStyle, h.Shell)
	if err != nil {
		c.JSON(500, gin.H{"code": config.FAILURE, "msg": err.Error()})
		return
	}
	GetAllTerminals(c)
}

func DeleteTerminal(c *gin.Context) {
	var terminal *models.Terminal
	id, err := strconv.Atoi(strings.TrimSpace(c.PostForm("id")))
	if err != nil {
		c.JSON(400, gin.H{"code": config.FAILURE, "msg": err.Error()})
		return
	}
	_, err = terminal.Delete(id)
	if err != nil {
		c.JSON(500, gin.H{"code": config.FAILURE, "msg": err.Error()})
		return
	}
	GetAllTerminals(c)
}

func GetSessionId(c *gin.Context) {
	h, err := VerifyTerminal(c)
	if err != nil {
		c.JSON(400, gin.H{"code": config.FAILURE, "msg": err.Error()})
		return
	}
	sessionId := utils.RandString(15)
	clients.AddData(h.Address, h.User, h.Pwd, h.Port, h.Shell, sessionId)
	c.JSON(200, gin.H{"code": config.SUCCEED, "data": sessionId, "msg": "ok"})
}
