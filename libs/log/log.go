package log

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type LogLevel uint8

const (
	Emergency LogLevel = iota
	Alert
	Critical
	Error
	Warning
	Notice
	Info
	Debug
)

type Log struct {
	*log.Logger
	logFile    *os.File
	Name       string
	Level      LogLevel
	OutFile    bool
	OutConsole bool
}

// brush is a color join function
type brush func(string) string

// newBrush return a fix color Brush
func newBrush(color string) brush {
	return func(text string) string {
		return "\033[" + color + "m" + text + "\033[0m"
	}
}

var colors = []brush{
	newBrush("1;41"), // Emergency          white
	newBrush("1;36"), // Alert              cyan
	newBrush("1;35"), // Critical           magenta
	newBrush("1;31"), // Error              red
	newBrush("1;33"), // Warning            yellow
	newBrush("1;32"), // Notice             green
	newBrush("1;34"), // Informational      blue
	newBrush("1;38"), // Debug              white
}

// NewLogger
// name 日志文件名称
// level 日志级别
// outFile 是否把日志输出到文件
// outConsole 是否把日志输出到控制台
func NewLogger(name string, level LogLevel, outFile, outConsole bool) *Log {
	if strings.TrimSpace(name) == "" {
		log.Println("Panic:Log name cannot be empty")
	}

	logFile, err := os.OpenFile(name+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Panic:Open log file error")
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)

	logObject := &Log{
		logger,
		logFile,
		name,
		level,
		outFile,
		outConsole,
	}
	return logObject
}

func (l *Log) Emergency(v ...interface{}) {
	if l.Level < 0 {
		return
	}
	l.SetPrefix("[M]: ")
	l.write(fmt.Sprint(v...), 0)
}

func (l *Log) Alert(v ...interface{}) {
	if l.Level < 1 {
		return
	}
	l.SetPrefix("[A]: ")
	l.write(fmt.Sprint(v...), 1)
}

func (l *Log) Critical(v ...interface{}) {
	if l.Level < 2 {
		return
	}
	l.SetPrefix("[C]: ")
	l.write(fmt.Sprint(v...), 2)
}

func (l *Log) Error(v ...interface{}) {
	if l.Level < 3 {
		return
	}
	l.SetPrefix("[E]: ")
	l.write(fmt.Sprint(v...), 3)
}

func (l *Log) Warning(v ...interface{}) {
	if l.Level < 4 {
		return
	}
	l.SetPrefix("[W]: ")
	l.write(fmt.Sprint(v...), 4)
}

func (l *Log) Notice(v ...interface{}) {
	if l.Level < 5 {
		return
	}
	l.SetPrefix("[N]: ")
	l.write(fmt.Sprint(v...), 5)
}

func (l *Log) Info(v ...interface{}) {
	if l.Level < 6 {
		return
	}
	l.SetPrefix("[I]: ")
	l.write(fmt.Sprint(v...), 6)
}

func (l *Log) Debug(v ...interface{}) {
	if l.Level < 7 {
		return
	}
	l.SetPrefix("[D]: ")
	l.write(fmt.Sprint(v...), 7)
}

func (l *Log) write(msg string, level int) {
	if l.OutConsole {
		l.SetOutput(os.Stdout)
		_ = l.Output(3, colors[level](msg))
	}

	if l.OutFile {
		l.SetOutput(l.logFile)
		_ = l.Output(3, msg)
	}
}

func (l *Log) SetLogLevel(level LogLevel) {
	l.Level = level
}
