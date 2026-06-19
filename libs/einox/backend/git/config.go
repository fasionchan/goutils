package git

import (
	"fmt"
	"strings"
)

// Transport 表示 Git 仓库访问协议。
type Transport string

const (
	TransportHTTP  Transport = "HTTP"
	TransportSSH   Transport = "SSH"
	TransportLocal Transport = "LOCAL"
)

// AuthConfig 按 Transport 选用对应认证字段；敏感信息不得写入日志或 error 字符串。
type AuthConfig struct {
	// HTTP 认证
	Username string
	Password string
	Token    string // 优先于 Username/Password

	// SSH 认证
	PrivateKey []byte
	Passphrase string
	KnownHosts []byte // 为空时仅允许在测试场景配合 InsecureIgnoreHostKey
}

// RefConfig 指定仓库快照，Commit > Tag > Branch > HEAD。
type RefConfig struct {
	Branch string
	Tag    string
	Commit string
}

// Config 是 GitBackend 的完整配置。
type Config struct {
	URL       string
	Transport Transport
	Auth      AuthConfig
	Ref       RefConfig
	BaseDir   string
}

// Validate 校验配置合法性。
func (c Config) Validate() error {
	if strings.TrimSpace(c.URL) == "" {
		return fmt.Errorf("git backend: URL is required")
	}

	transport := c.Transport
	if transport == "" {
		transport = inferTransport(c.URL)
	}

	switch transport {
	case TransportHTTP, TransportSSH, TransportLocal:
	default:
		return fmt.Errorf("git backend: unsupported transport %q", transport)
	}

	if transport == TransportSSH && len(c.Auth.PrivateKey) == 0 {
		return fmt.Errorf("git backend: SSH transport requires PrivateKey")
	}

	if refCount(c.Ref) > 1 {
		return fmt.Errorf("git backend: only one of Branch, Tag, Commit may be set")
	}

	return nil
}

func inferTransport(url string) Transport {
	switch {
	case strings.HasPrefix(url, "file://"), strings.HasPrefix(url, "/"):
		return TransportLocal
	case strings.HasPrefix(url, "git@"), strings.HasPrefix(url, "ssh://"):
		return TransportSSH
	default:
		return TransportHTTP
	}
}

func refCount(ref RefConfig) int {
	n := 0
	if ref.Branch != "" {
		n++
	}
	if ref.Tag != "" {
		n++
	}
	if ref.Commit != "" {
		n++
	}
	return n
}
