package logger

import "gossh/libs/log"

var logLevel = log.Error

var Logger = log.NewLogger("GoSSH", logLevel, true, true)
