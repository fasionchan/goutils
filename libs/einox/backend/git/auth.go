package git

import (
	"fmt"
	"os"

	git "github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/client"
	githttp "github.com/go-git/go-git/v6/plumbing/transport/http"
	gitssh "github.com/go-git/go-git/v6/plumbing/transport/ssh"
	"github.com/go-git/go-git/v6/plumbing/transport/ssh/knownhosts"
	gossh "golang.org/x/crypto/ssh"
)
func buildClientOptions(cfg Config) ([]client.Option, error) {
	transportType := cfg.Transport
	if transportType == "" {
		transportType = inferTransport(cfg.URL)
	}

	switch transportType {
	case TransportHTTP:
		if auth := buildHTTPAuth(cfg.Auth); auth != nil {
			return []client.Option{client.WithHTTPAuth(auth)}, nil
		}
		return nil, nil
	case TransportSSH:
		auth, err := buildSSHAuth(cfg.Auth)
		if err != nil {
			return nil, err
		}
		return []client.Option{client.WithSSHAuth(auth)}, nil
	case TransportLocal:
		return nil, nil
	default:
		return nil, fmt.Errorf("git backend: unsupported transport %q", transportType)
	}
}

func buildHTTPAuth(auth AuthConfig) client.HTTPAuth {
	if auth.Token != "" {
		username := auth.Username
		if username == "" {
			username = "oauth2"
		}
		return &githttp.BasicAuth{Username: username, Password: auth.Token}
	}
	if auth.Username != "" || auth.Password != "" {
		return &githttp.BasicAuth{Username: auth.Username, Password: auth.Password}
	}
	return nil
}

func buildSSHAuth(auth AuthConfig) (*gitssh.PublicKeys, error) {
	publicKeys, err := gitssh.NewPublicKeys("git", auth.PrivateKey, auth.Passphrase)
	if err != nil {
		return nil, fmt.Errorf("git backend: invalid SSH private key")
	}

	if len(auth.KnownHosts) > 0 {
		file, err := os.CreateTemp("", "known_hosts-*")
		if err != nil {
			return nil, fmt.Errorf("git backend: prepare known_hosts")
		}
		path := file.Name()
		if _, err := file.Write(auth.KnownHosts); err != nil {
			_ = file.Close()
			_ = os.Remove(path)
			return nil, fmt.Errorf("git backend: prepare known_hosts")
		}
		if err := file.Close(); err != nil {
			_ = os.Remove(path)
			return nil, fmt.Errorf("git backend: prepare known_hosts")
		}

		db, err := knownhosts.NewDB(path)
		if err != nil {
			_ = os.Remove(path)
			return nil, fmt.Errorf("git backend: invalid known_hosts")
		}
		publicKeys.HostKeyCallback = db.HostKeyCallback()
	} else {
		publicKeys.HostKeyCallback = gossh.InsecureIgnoreHostKey()
	}

	return publicKeys, nil
}

// cloneOptions 构造 clone 选项。
func cloneOptions(cfg Config) (*git.CloneOptions, error) {
	clientOpts, err := buildClientOptions(cfg)
	if err != nil {
		return nil, err
	}

	opts := &git.CloneOptions{
		URL:           cfg.URL,
		Bare:          true,
		NoCheckout:    true,
		ClientOptions: clientOpts,
	}
	return opts, nil
}
