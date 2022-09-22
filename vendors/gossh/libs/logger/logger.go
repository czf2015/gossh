package logger

import "gossh/vendors/libs/log"

var logLevel = log.Error

var Logger = log.NewLogger("GoSSH", logLevel, true, true)
