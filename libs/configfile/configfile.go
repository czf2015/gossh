//###############################
// 读取配置文件功能
//###############################
package configfile

import (
	"bufio"
	"errors"
	"gossh/libs/logger"
	"io"
	"os"
	"strconv"
	"strings"
)

type configFile struct {
	fileName string
	comment  []string
}

// Section 配置片段类型
type Section map[string]string

// GetInt 获取片段中的值,(转换成int)
func (s Section) GetInt(key string) (int, error) {
	data := 0
	if val, ok := s[key]; ok {
		data, err := strconv.Atoi(val)
		if err == nil {
			return data, nil
		}
	}
	return data, errors.New("GetInt Error")
}

// GetFloat 获取片段中的值,(转换成float)
func (s Section) GetFloat(key string) (float64, error) {
	data := 0.0
	if val, ok := s[key]; ok {
		data, err := strconv.ParseFloat(val, 64)
		if err == nil {
			return data, nil
		}
	}
	return data, errors.New("GetFloat Error")
}

// GetString 获取片段中的值,(默认值字符串)
func (s Section) GetString(key string) (string, error) {
	if val, ok := s[key]; ok {
		return val, nil
	}
	return "", errors.New("GetString Error")
}

// GetBool 获取片段中的值,(转换成bool)
func (s Section) GetBool(key string) (bool, error) {
	if val, ok := s[key]; ok {
		data, err := strconv.ParseBool(val)
		if err == nil {
			return data, nil
		}
	}
	return false, errors.New("GetBool Error")
}

// ReadLines 读取配置文件的每一行
func (c configFile) ReadLines() (lines []string, err error) {
	fd, err := os.Open(c.fileName)
	if err != nil {
		return
	}
	defer func() {
		_ = fd.Close()
	}()
	lines = make([]string, 0)
	reader := bufio.NewReader(fd)
	prefix := ""
	var isLongLine bool
	for {
		byteLine, isPrefix, er := reader.ReadLine()
		if er != nil && er != io.EOF {
			return nil, er
		}
		if er == io.EOF {
			break
		}
		line := string(byteLine)
		if isPrefix {
			prefix += line
			continue
		} else {
			isLongLine = true
		}

		line = prefix + line
		if isLongLine {
			prefix = ""
		}
		line = strings.TrimSpace(line)
		// 跳过空白行
		if len(line) == 0 {
			continue
		}
		// 跳过注释行
		var breakLine = false
		for _, v := range c.comment {
			if strings.HasPrefix(line, v) {
				breakLine = true
				break
			}
		}
		if breakLine {
			continue
		}
		lines = append(lines, line)
	}
	return lines, nil
}

// GetConfig 获取所有配置
func (c configFile) GetConfig() map[string]map[string]string {
	config := make(map[string]map[string]string)
	lines, err := c.ReadLines()
	if err != nil {
		logger.Logger.Error(err)
	}
	var section = make(map[string]string, 1)

	for _, line := range lines {
		if line[0] == '[' && line[len(line)-1] == ']' {
			sectionName := line[1 : len(line)-1]
			section = make(map[string]string, 1)
			config[sectionName] = section
		} else {
			configKeyVal := strings.Split(line, "=")
			key := strings.TrimSpace(configKeyVal[0])
			val := strings.TrimSpace(strings.Join(configKeyVal[1:], "="))
			section[key] = val
		}
	}
	return config
}

// GetSection 获取某一段配置
func (c configFile) GetSection(section string) (Section, error) {
	if data, ok := c.GetConfig()[section]; ok {
		return data, nil
	}
	return map[string]string{}, nil
}

// LoadConfigFile 加载配置文件
func LoadConfigFile(filename string, comment []string) (configFile, error) {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Logger.Error("file not exist:", err)
			return configFile{}, err
		}
	}
	return configFile{
		fileName: filename,
		comment:  comment,
	}, nil
}

func Parse(filename string) (config map[string]map[string]string) {
	configFile, err := LoadConfigFile(filename, []string{"#", ";"})
	if err == nil {
		logger.Logger.Debug("读取配置文件成功,使用系统配置文件")
		config = configFile.GetConfig()
	}
	return
}
