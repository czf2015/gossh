package config

import (
	"gossh/libs/configfile"
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
	WorkDir, _ = os.Getwd()
	configFilePath := path.Join(WorkDir, "/config/v1/config.ini")
	Config = configfile.Parse(configFilePath)
}
