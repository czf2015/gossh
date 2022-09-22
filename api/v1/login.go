package api

import (
	config "gossh/config/v1"
	"gossh/libs/db"
	"gossh/libs/logger"
	"gossh/libs/sessions"
	models "gossh/models/v1"

	"github.com/gin-gonic/gin"
)

// POST 登陆管理页面
func Login(c *gin.Context) {
	pwd := c.DefaultPostForm("pwd", "")
	row := db.GetDB().QueryRow("select Id,Pwd from login where Id = 1")
	login := new(models.Login)
	err := row.Scan(&login.Id, &login.Pwd)
	if err != nil {
		logger.Logger.Error(err)
		c.JSON(500, gin.H{
			"code": config.FAILURE,
			"msg":  "login error",
		})
		return
	}

	if login.Pwd != pwd {
		c.JSON(401, gin.H{
			"code": config.FAILURE,
			"msg":  "login password error",
		})
		return
	}

	session := sessions.Default(c)

	session.Set("auth", "Y")
	//记着调用save方法，写入session
	_ = session.Save()
	c.JSON(200, gin.H{
		"code": config.SUCCEED,
		"msg":  "login success",
	})
}

// PATCH (需要登陆认证)修改登陆密码
func RevisePassword(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("auth") != "Y" {
		c.JSON(401, gin.H{
			"code": config.FAILURE,
			"msg":  "Unauthorized",
		})
		return
	}

	oldPwd := c.PostForm("old_pwd")
	newPwd := c.PostForm("new_pwd")
	row := db.GetDB().QueryRow(`select Id,Pwd from login where Id = 1`)
	login := new(models.Login)
	err := row.Scan(&login.Id, &login.Pwd)
	if err != nil {
		logger.Logger.Error(err)
		c.JSON(500, gin.H{
			"code": config.FAILURE,
			"msg":  "change password error",
		})
		return
	}

	if login.Pwd != oldPwd {
		c.JSON(401, gin.H{
			"code": config.FAILURE,
			"msg":  "change password error",
		})
		return
	}

	stmt, _ := db.GetDB().Prepare(`update config set Pwd=? where id=1`)
	_, err = stmt.Exec(newPwd)
	if err != nil {
		c.JSON(401, gin.H{
			"code": config.FAILURE,
			"msg":  "modify password config.FAILURE",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": config.SUCCEED,
		"msg":  "modify password success",
	})
}
