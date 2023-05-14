/*
 * Author: fasion
 * Created time: 2023-05-14 11:34:25
 * Last Modified by: fasion
 * Last Modified time: 2023-05-14 12:39:24
 */

package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/smtp"

	"gopkg.in/gomail.v2"

	"github.com/fasionchan/goutils/baseutils"
)

const (
	DefaultSmtpPort        = 25
	DEfaultSmtpPortWithTls = 465
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

type SmtpPlainAuth struct {
	username, password string
}

func NewSmtpPlainAuth(username, password string) smtp.Auth {
	return &SmtpPlainAuth{username, password}
}

func (a *SmtpPlainAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *SmtpPlainAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		}
	}
	return nil, nil
}
