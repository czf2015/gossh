package main

import (
	_ "github.com/mattn/go-sqlite3"

	"gossh/libs/logger"
	router "gossh/router/v1"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			logger.Logger.Emergency(err)
		}
	}()
	// go clients.ConnectGC()
	router.Run()
}
