package api

import (
	config "gossh/config/v1"
	"gossh/libs/clients"
	"gossh/libs/gin"
	"io"
	"path"
	"strconv"
	"strings"
)

// GET sftp 获取指定目录下文件信息
func SftpDir(c *gin.Context) {
	dirPath := c.Query("path")
	sessionId := c.Query("session_id")
	cli, ok := clients.GetClientBySessionID(sessionId)
	if !ok {
		c.JSON(400, gin.H{"code": config.FAILURE, "msg": "sftpClient error"})
		return
	}
	files, err := cli.SftpClient.ReadDir(dirPath)
	if err != nil {
		c.JSON(400, gin.H{"code": config.FAILURE, "msg": "list Folder error"})
		return
	}

	fileCount := 0
	dirCount := 0
	var fileList []interface{}
	for _, file := range files {
		fileInfo := map[string]interface{}{}
		fileInfo["path"] = path.Join(dirPath, file.Name())
		fileInfo["name"] = file.Name()
		fileInfo["mode"] = file.Mode().String()
		fileInfo["size"] = file.Size()
		fileInfo["mod_time"] = file.ModTime().Format("2006-01-02 15:04:05")
		if file.IsDir() {
			fileInfo["type"] = "d"
			dirCount += 1
		} else {
			fileInfo["type"] = "f"
			fileCount += 1
		}
		fileList = append(fileList, fileInfo)
	}

	// 内部方法,处理路径信息
	pathHandler := func(dirPath string) (paths []map[string]string) {
		tmp := strings.Split(dirPath, "/")

		var dirs []string
		if strings.HasPrefix(dirPath, "/") {
			dirs = append(dirs, "/")
		}

		for _, item := range tmp {
			name := strings.TrimSpace(item)
			if len(name) > 0 {
				dirs = append(dirs, name)
			}
		}

		for i, item := range dirs {
			fullPath := path.Join(dirs[:i+1]...)
			pathInfo := map[string]string{}
			pathInfo["name"] = item
			pathInfo["dir"] = fullPath
			paths = append(paths, pathInfo)
		}
		return paths
	}

	data := map[string]interface{}{
		"files":       fileList,
		"file_count":  fileCount,
		"dir_count":   dirCount,
		"paths":       pathHandler(dirPath),
		"current_dir": dirPath,
	}

	c.JSON(200, gin.H{"code": config.SUCCEED, "data": data, "msg": "ok"})
}

// POST sftp 下载文件
func SftpDownload(c *gin.Context) {
	sessionId := c.PostForm("session_id")
	fullPath := c.PostForm("path")
	cli, ok := clients.GetClientBySessionID(sessionId)
	if ok {
		file, _ := cli.SftpClient.Open(fullPath)
		defer func() {
			_ = file.Close()
		}()
		_, _ = io.Copy(c.Writer, file)
	}
}

// PUT sftp 上传文件
func SftpUpload(c *gin.Context) {
	sessionId := c.PostForm("session_id")
	dstPath := c.PostForm("path")
	//获取上传的文件组
	files := c.Request.MultipartForm.File["file"]

	for _, file := range files {
		cli, ok := clients.GetClientBySessionID(sessionId)
		if ok {
			srcFile, _ := file.Open()
			dstFile, _ := cli.SftpClient.Create(path.Join(dstPath, file.Filename))
			_, _ = io.Copy(dstFile, srcFile)
			_ = srcFile.Close()
			_ = dstFile.Close()
		}
	}
	msg := strconv.Itoa(len(files)) + "文件上传成功"
	c.JSON(200, gin.H{"code": config.SUCCEED, "msg": msg})
}
