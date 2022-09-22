package logger

import "gossh/vendors/log"

var logLevel = log.Error

var Logger = log.NewLogger("GoSSH", logLevel, true, true)
