package logger

import "gowebssh/vendors/log"

var logLevel = log.Error

var Logger = log.NewLogger("GoSSH", logLevel, true, true)
