/*
 * Author: fasion
 * Created time: 2023-05-14 11:34:25
 * Last Modified by: fasion
 * Last Modified time: 2024-12-20 10:44:24
 */

package email

import (
	"crypto/tls"
	"fmt"
	"os"

	"gopkg.in/gomail.v2"

	"github.com/fasionchan/goutils/baseutils"
	"github.com/fasionchan/goutils/stl"
)

const (
	DefaultSmtpPort        = 25
	DefaultSmtpPortWithTls = 465

	EnvNameSmtpServerLoc = "SMTP_SERVER_LOC"
	EnvNameEmailAccount  = "EMAIL_ACCOUNT"
	EnvNameEmailPassword = "EMAIL_PASSWORD"
)

type EmailClient struct {
	addr      string
	port      int
	accout    string
	password  string
	tlsConfig *tls.Config
}

func NewEmailClient(loc, account, password string) (*EmailClient, error) {
	addr, port, err := baseutils.ParseNetloc(loc, "localhost", DefaultSmtpPortWithTls)
	if err != nil {
		return nil, err
	}

	return &EmailClient{
		addr:     addr,
		port:     port,
		accout:   account,
		password: password,
	}, nil
}

func NewEmailClientFromEnvPro(locEnvName, accountEnvName, passwordEnvName string, getenv func(string) string) (*EmailClient, error) {
	loc := getenv(locEnvName)
	if loc == "" {
		return nil, baseutils.NewEnvironmentVariableNotFoundError(locEnvName)
	}

	account := getenv(accountEnvName)
	if account == "" {
		return nil, baseutils.NewEnvironmentVariableNotFoundError(accountEnvName)
	}

	password := getenv(passwordEnvName)
	if password == "" {
		return nil, baseutils.NewEnvironmentVariableNotFoundError(passwordEnvName)
	}

	return NewEmailClient(loc, account, password)
}

func NewEmailClientFromEnv(locEnvName, accountEnvName, passwordEnvName string) (*EmailClient, error) {
	return NewEmailClientFromEnvPro(locEnvName, accountEnvName, passwordEnvName, os.Getenv)
}

func NewEmailClientFromDefaultEnvPro(getenv func(string) string) (*EmailClient, error) {
	return NewEmailClientFromEnvPro(EnvNameSmtpServerLoc, EnvNameEmailAccount, EnvNameEmailPassword, getenv)
}

func NewEmailClientFromDefaultEnv() (*EmailClient, error) {
	return NewEmailClientFromEnvPro(EnvNameSmtpServerLoc, EnvNameEmailAccount, EnvNameEmailPassword, os.Getenv)
}

func (client *EmailClient) Addr() string {
	if client == nil {
		return ""
	}
	return client.addr
}

func (client *EmailClient) Port() int {
	if client == nil {
		return 0
	}
	return client.port
}

func (client *EmailClient) Account() string {
	if client == nil {
		return ""
	}
	return client.accout
}

func (client *EmailClient) Dup() *EmailClient {
	return stl.Dup(client)
}

func (client *EmailClient) WithAccount(account, password string) *EmailClient {
	if client == nil {
		return nil
	}
	client.accout = account
	client.password = password
	return client
}

func (client *EmailClient) WithTlsConfig(config *tls.Config) *EmailClient {
	if client == nil {
		return nil
	}
	client.tlsConfig = config
	return client
}

func (client *EmailClient) ForkWithAccount(account, password string) *EmailClient {
	if client == nil {
		return nil
	}
	return client.Dup().WithAccount(account, password)
}

func (client *EmailClient) NetLoc() string {
	if client == nil {
		return ""
	}
	return fmt.Sprintf("%s:%d", client.addr, client.port)
}

func (client *EmailClient) SendMessage(m *gomail.Message) error {
	d := gomail.NewDialer(client.addr, client.port, client.accout, client.password)
	if client.tlsConfig != nil {
		d.TLSConfig = client.tlsConfig
	}
	m.SetHeader("From", client.accout)
	return d.DialAndSend(m)
}

func (client *EmailClient) SendMessageSmart(to []string, subject, body string) error {
	// todo: make encoding configurable
	msg := gomail.NewMessage(gomail.SetEncoding(gomail.Base64))
	msg.SetHeader("From", client.accout)
	msg.SetHeader("To", to...)
	msg.SetHeader("Subject", subject)

	if len(body) > 0 {
		msg.SetBody("text/html", body)
	}

	return client.SendMessage(msg)
}
