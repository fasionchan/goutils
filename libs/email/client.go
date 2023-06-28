/*
 * Author: fasion
 * Created time: 2023-05-14 11:34:25
 * Last Modified by: fasion
 * Last Modified time: 2023-06-28 15:11:09
 */

package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/smtp"
	"os"

	"gopkg.in/gomail.v2"

	"github.com/fasionchan/goutils/baseutils"
)

const (
	DefaultSmtpPort        = 25
	DEfaultSmtpPortWithTls = 465

	EnvNameSmtpServerLoc       = "SMTP_SERVER_LOC"
	EnvNameEmailSenderAccount  = "EMAIL_SENDER_ACCOUNT"
	EnvNameEmailSenderPassword = "EMAIL_SENDER_PASSWORD"
)

type EmailClient struct {
	addr      string
	port      int
	accout    string
	password  string
	tlsConfig *tls.Config
}

func NewEmailClient(loc, account, password string) (*EmailClient, error) {
	addr, port, err := baseutils.ParseNetloc(loc, "localhost", DEfaultSmtpPortWithTls)
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

func NewEmailClientFromEnvPro(locEnvName, senderAccountEnvName, senderPasswordEnvName string, getenv func(string) string) (*EmailClient, error) {
	loc := getenv(locEnvName)
	if loc == "" {
		return nil, baseutils.NewEnvironmentVariableNotFoundError(locEnvName)
	}

	senderAccount := getenv(senderAccountEnvName)
	if senderAccount == "" {
		return nil, baseutils.NewEnvironmentVariableNotFoundError(senderAccountEnvName)
	}

	senderPassword := getenv(senderPasswordEnvName)
	if senderPassword == "" {
		return nil, baseutils.NewEnvironmentVariableNotFoundError(senderPasswordEnvName)
	}

	return NewEmailClient(loc, senderAccount, senderPassword)
}

func NewEmailClientFromEnv(locEnvName, senderAccountEnvName, senderPasswordEnvName string) (*EmailClient, error) {
	return NewEmailClientFromEnvPro(locEnvName, senderAccountEnvName, senderPasswordEnvName, os.Getenv)
}

func NewEmailClientFromDefaultEnvPro(getenv func(string) string) (*EmailClient, error) {
	return NewEmailClientFromEnvPro(EnvNameSmtpServerLoc, EnvNameEmailSenderAccount, EnvNameEmailSenderPassword, getenv)
}

func NewEmailClientFromDefaultEnv() (*EmailClient, error) {
	return NewEmailClientFromEnvPro(EnvNameSmtpServerLoc, EnvNameEmailSenderAccount, EnvNameEmailSenderPassword, os.Getenv)
}

func (client *EmailClient) WithTlsConfig(config *tls.Config) *EmailClient {
	client.tlsConfig = config
	return client
}

func (client *EmailClient) SendMail(to []string, msg []byte) error {
	auth := smtp.PlainAuth("", client.accout, client.password, "")
	return smtp.SendMail(fmt.Sprintf("%s:%d", client.addr, client.port), auth, client.accout, to, msg)
}

func (client *EmailClient) SendMailSmart(to []string, subject, body string) error {
	var b bytes.Buffer
	w := io.Writer(&b)
	fmt.Fprintf(w, "To: %s\r\n", to)
	fmt.Fprintf(w, "Subject: %s\r\n", subject)
	fmt.Fprintf(w, "\r\n")
	fmt.Fprintf(w, "%s\r\n", body)

	return client.SendMail(to, b.Bytes())
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
	msg := gomail.NewMessage()
	msg.SetHeader("From", client.accout)
	msg.SetHeader("To", to...)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	return client.SendMessage(msg)
}
