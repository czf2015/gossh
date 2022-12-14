package api

import (
	"fmt"
	config "gossh/config/v1"
	models "gossh/models/v1"
	"gossh/utils"
	"gossh/vendors/clients"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func VerifyTerminal(c *gin.Context) (models.Terminal, error) {
	var terminal models.Terminal
	if err := c.ShouldBindJSON(&terminal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// if len(name) > 60 || len(name) == 0 {
	// 	return models.Terminal{}, fmt.Errorf("name input error:%s.", name)
	// }

	if len(terminal.Address) > 60 || len(terminal.Address) == 0 {
		return models.Terminal{}, fmt.Errorf("terminal input error")
	}

	if len(terminal.User) > 60 || len(terminal.User) == 0 {
		return models.Terminal{}, fmt.Errorf("user input error")
	}

	if len(terminal.Pwd) > 60 || len(terminal.Pwd) == 0 {
		return models.Terminal{}, fmt.Errorf("pwd input error")
	}
	// p, err := strconv.Atoi(strings.TrimSpace(terminal.Port))
	// if err != nil {
	// 	return models.Terminal{}, fmt.Errorf("port input error")
	// }
	if terminal.Port > 65535 || terminal.Port < 1 {
		return models.Terminal{}, fmt.Errorf("port range input error")
	}

	// fontsize, err := strconv.Atoi(strings.TrimSpace(fontSize))
	// if err != nil {
	terminal.FontSize = 16
	// }
	// if fontsize > 32 || fontsize < 8 {
	// 	fontsize = 16
	// }
	// if len(strings.TrimSpace(background)) == 0 {
	terminal.Background = "#000000"
	// }

	// if len(strings.TrimSpace(foreground)) == 0 {
	terminal.Foreground = "#FFFFFF"
	// }

	// if len(strings.TrimSpace(cursorColor)) == 0 {
	terminal.CursorColor = "#FFFFFF"
	// }

	// if len(strings.TrimSpace(fontFamily)) == 0 {
	terminal.FontFamily = "Courier"
	// }

	// if len(strings.TrimSpace(cursorStyle)) == 0 {
	terminal.CursorStyle = "block"
	// }

	// if len(strings.TrimSpace(shell)) == 0 {
	terminal.Shell = "bash"
	// }

	return terminal, nil
}

func GetAllTerminals(c *gin.Context) {
	var terminal *models.Terminal
	allTerminals, err := terminal.Select()
	if err != nil {
		c.JSON(500, gin.H{"code": config.FAILURE, "msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": config.SUCCEED, "data": allTerminals, "msg": "ok"})
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
