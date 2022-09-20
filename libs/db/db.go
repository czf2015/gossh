package db

import (
	"database/sql"
	"fmt"
	"gossh/config/v1"
	"gossh/libs/logger"
	"os"
	"path"
)

var db *sql.DB

func init() {
	var err error

	db, err = sql.Open("sqlite3", path.Join(config.WorkDir, config.ProjectName+".db"))
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("创建数据库文件:%s失败\n", path.Join(config.WorkDir, config.ProjectName+".db")))
		os.Exit(1)
		return
	}

	// fileInfo, err := os.Stat(config.WorkDir)

	// if os.IsNotExist(err) {
	// 	err = os.Mkdir(config.WorkDir, fs.ModePerm)
	// 	if err != nil {
	// 		logger.Logger.Error(fmt.Sprintf("创建目录:%s 失败,%s\n", config.WorkDir, err))
	// 		os.Exit(1)
	// 		return
	// 	}

	// } else {
	// 	if !fileInfo.IsDir() {
	// 		logger.Logger.Error(fmt.Sprintf("请删除:%s文件\n", config.WorkDir))
	// 		os.Exit(1)
	// 		return
	// 	}
	// }

	// configFilePath := path.Join(config.WorkDir, config.ProjectName+".cnf")
	// _, err = os.Stat(configFilePath)
	// if os.IsNotExist(err) {
	// 	file, err := os.Create(configFilePath)
	// 	if err != nil {
	// 		logger.Logger.Error(fmt.Sprintf("创建默认配置文件:%s 失败,%s\n", configFilePath, err))
	// 		os.Exit(1)
	// 		return
	// 	}
	// 	defer func() {
	// 		_ = file.Close()
	// 	}()

	// 		configContent := `
	// [app]
	// AppName=GoWebSSH
	// [server]
	// Address=0.0.0.0
	// Port=8899
	// CertFile=` + path.Join(config.WorkDir, "cert.pem") + `
	// KeyFile=` + path.Join(config.WorkDir, "key.key") + `
	// [session]
	// Secret=` + utils.RandString(64) + `
	// Name=session_id
	// Path=/
	// Domain=
	// MaxAge=86400
	// Secure=false
	// HttpOnly=true
	// SameSite=2
	// `
	// 		_, err = file.WriteString(configContent)
	// 		if err != nil {
	// 			logger.Logger.Error(fmt.Sprintf("写入配置文件:%s 失败,%s\n", configFilePath, err))
	// 			os.Exit(1)
	// 			return
	// 		}
	// 		_ = file.Sync()
	// 	}

	// 	createTerminalTable := `
	// CREATE TABLE IF NOT EXISTS 'terminal'
	// (
	//     'Id'          INTEGER PRIMARY KEY AUTOINCREMENT,
	//     'Name'        VARCHAR(32) NOT NULL UNIQUE,
	//     'Address'     VARCHAR(64) NULL,
	//     'User'        VARCHAR(64) NULL,
	//     'Pwd'         VARCHAR(64) NULL,
	//     'Port'        INT         NOT NULL DEFAULT 22,
	//     'FontSize'    INT         NOT NULL DEFAULT 14,
	//     'Background'  VARCHAR(32) NOT NULL DEFAULT '#000000',
	//     'Foreground'  VARCHAR(32) NOT NULL DEFAULT '#FFFFFF',
	//     'CursorColor' VARCHAR(32) NOT NULL DEFAULT '#FFFFFF',
	//     'FontFamily'  VARCHAR(32) NOT NULL DEFAULT 'Courier',
	//     'CursorStyle' VARCHAR(32) NOT NULL DEFAULT 'block',
	//     'Shell'       VARCHAR(32) NOT NULL DEFAULT 'bash'
	// );
	// `
	// 	_, err = db.Exec(createTerminalTable)
	// 	if err != nil {
	// 		logger.Logger.Error(err)
	// 	}

	// 	createConfigTable := `
	// CREATE TABLE IF NOT EXISTS 'config'
	// (
	//     'Id'          INTEGER PRIMARY KEY AUTOINCREMENT,
	//     'Pwd'         VARCHAR(64) NOT NULL DEFAULT 'admin'
	// );
	// `
	// 	_, err = db.Exec(createConfigTable)
	// 	if err != nil {
	// 		logger.Logger.Error(err)
	// 	}

	// 	insertSql := `INSERT INTO config(Id,Pwd)  values(?,?)`
	// 	stmt, _ := db.Prepare(insertSql)
	// 	_, err = stmt.Exec(1, "admin")

}

func GetDB() *sql.DB {
	return db
}
