package clients

import (
	"encoding/json"
	"fmt"
	"gossh/libs/logger"
	"io"
	"net"
	"sync"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"

	"gossh/libs/websocket"
)

type Ssh struct {
	IP         string          `json:"ip"`         //IP地址
	Username   string          `json:"username"`   //用户名
	Password   string          `json:"-"`          //密码
	Port       int             `json:"port"`       //端口号
	SessionId  string          `json:"session_id"` //会话ID
	Shell      string          `json:"shell"`
	Timeout    time.Time       `json:"timeout"`
	StartTime  time.Time       `json:"start_time"` // 建立连接的时间
	SshClient  *ssh.Client     //ssh客户端
	SftpClient *sftp.Client    //sftp客户端
	SshSession *ssh.Session    //ssh会话
	Ws         *websocket.Conn // websocket 连接
}

// MarshalJSON 重写序列化方法
func (s Ssh) MarshalJSON() ([]byte, error) {
	type Alias Ssh
	return json.Marshal(&struct {
		Alias
		Timeout   string `json:"timeout"`
		StartTime string `json:"start_time"`
	}{
		Alias:     (Alias)(s),
		Timeout:   s.Timeout.Format("2006-01-02 15:04:05"),
		StartTime: s.StartTime.Format("2006-01-02 15:04:05"),
	})
}

// 连接主机
func (s *Ssh) connect() error {
	defer func() {
		if err := recover(); err != nil {
			logger.Logger.Error(err)
		}
	}()

	config := ssh.ClientConfig{
		User: s.Username,
		Auth: []ssh.AuthMethod{ssh.Password(s.Password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 30 * time.Second,
	}
	addr := fmt.Sprintf("%s:%d", s.IP, s.Port)
	sshClient, err := ssh.Dial("tcp", addr, &config)
	if err != nil {
		return err
	}

	s.SshClient = sshClient
	//使用sshClient构建sftpClient
	var sftpClient *sftp.Client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		logger.Logger.Error("create sftp sshClient error:", err)
	}
	s.SftpClient = sftpClient
	return nil
}

// Clients 存储的客户端信息
type Clients struct {
	lock sync.RWMutex
	data map[string]*Ssh
}

var clients = Clients{
	lock: sync.RWMutex{},
	data: make(map[string]*Ssh),
}

// RunTerminal 运行一个终端
func (s *Ssh) RunTerminal(shell string, stdout, stderr io.Writer, stdin io.Reader, w, h int, ws *websocket.Conn) error {
	if s.SshClient == nil {
		if err := s.connect(); err != nil {
			logger.Logger.Error(err)
			return err
		}
	}

	s.Ws = ws

	sshSession, err := s.SshClient.NewSession()
	if err != nil {
		s.SshClient.Close()
		logger.Logger.Error(err.Error())
		return err
	}
	// defer func() {
	// 	DeleteClientBySessionID((s.SessionId))
	// }()

	sshSession.Stdout = stdout
	sshSession.Stderr = stderr
	sshSession.Stdin = stdin
	s.SshSession = sshSession

	modes := ssh.TerminalModes{}
	if err := sshSession.RequestPty("xterm-256color", h, w, modes); err != nil {
		return err
	}

	err = sshSession.Run(shell)
	if err != nil {
		logger.Logger.Error(err.Error())
		return err
	}

	return nil
}

func GetClientBySessionID(sessionID string) (*Ssh, bool) {
	clients.lock.RLock()
	defer clients.lock.RUnlock()
	client, err := clients.data[sessionID]
	return client, err
}

func GetData() map[string]*Ssh {
	clients.lock.RLock()
	defer clients.lock.RUnlock()
	return clients.data
}

func SetData(sessionId string, client *Ssh) {
	clients.data[sessionId] = client
}

func AddData(ip string, username string, password string, port int, shell, sessionId string) {
	now := time.Now()
	clientsSsh := &Ssh{
		IP:        ip,
		Username:  username,
		Password:  password,
		Port:      port,
		Shell:     shell,
		SessionId: sessionId,
		Timeout:   now,
		StartTime: now,
	}
	clients.lock.Lock()
	clients.data[sessionId] = clientsSsh
	clients.lock.Unlock()
}

func DeleteClientBySessionID(sessionId string) error {
	clients.lock.Lock()
	defer clients.lock.Unlock()
	delete(clients.data, sessionId)
	return nil
}

func Lock() {
	clients.lock.Lock()
}

func Unlock() {
	clients.lock.Unlock()
}

func RLock() {
	clients.lock.RLock()
}

func RUnlock() {
	clients.lock.RUnlock()
}

// ConnectGC 清理已经断开的连接
func ConnectGC() {
	defer func() {
		if err := recover(); err != nil {
			logger.Logger.Error(err)
		}
	}()

	for {
		time.Sleep(time.Second)
		duration, _ := time.ParseDuration("-1m")
		longAgo := time.Now().Add(duration)
		for key, item := range clients.data {
			if item.Timeout.Before(longAgo) {
				item.SshSession.Close()
				item.SshClient.Close()
				item.SftpClient.Close()
				item.Ws.Close()
				DeleteClientBySessionID(key)
			}
		}
	}
}
