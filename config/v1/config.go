package config

import (
	"fmt"
	"gossh/libs/configfile"
	"gossh/libs/logger"
	"os"
	"path"
)

var WorkDir string

var Config map[string]map[string]string

const (
	SUCCEED = 0
	FAILURE = 1
)

func init() {
	WorkDir, err := os.Getwd()
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("获取运行目录:%s 失败", err))
	}
	configFilePath := path.Join(WorkDir, "/config/v1/config.ini")
	Config = configfile.Parse(configFilePath)
}
