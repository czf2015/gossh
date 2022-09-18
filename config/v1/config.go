package config

import (
	"fmt"
	"gossh/libs/configfile"
	"gossh/utils"
	"os"
	"path"
)

var userHomeDir, _ = os.UserHomeDir()

var ProjectName = "GoWebSSH"

// WorkDir 程序默认工作目录,在用户的home目录下 .GoWebSSH 目录
var WorkDir = path.Join(userHomeDir, fmt.Sprintf("/.%s/", ProjectName))

// Config 默认配置,当配置文件不存在的时候,就使用这个默认配置
var Config = map[string]map[string]string{
	"app": {
		"AppName": "GoWebSSH",
	},
	"server": {
		"Address":  "0.0.0.0",
		"Port":     "8899",
		"CertFile": path.Join(WorkDir, "cert.pem"),
		"KeyFile":  path.Join(WorkDir, "key.key"),
	},
	"session": {
		"Secret":   utils.RandString(64),
		"Name":     "session_id",
		"Path":     "/",
		"Domain":   "",
		"MaxAge":   "86400",
		"Secure":   "false",
		"HttpOnly": "true",
		"SameSite": "2",
	},
}

const (
	SUCCEED = 0
	FAILURE = 1
)

func init() {
	Config = configfile.Parse(path.Join(WorkDir, "GoWebSSH.cnf"))
}
