package sender

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"mailgun-relay/utils"
	"mime/multipart"
	"mime/quotedprintable"
	"net/smtp"
	"net/textproto"
	"path/filepath"
	"strings"
	"time"
)

type SmtpEngine struct {
	Host     string
	Port     string
	Username string
	Password string

	Sender    string
	Receivers []string
}

func (s *SmtpEngine) Configure(config SenderConfig) error {
	if config.SmtpConfig == nil {
		return errors.New("smtp config was not provided")
	}

	smtpConfig := config.SmtpConfig

	s.Host = smtpConfig.Host
	s.Port = smtpConfig.Port
	s.Username = smtpConfig.Username
	s.Password = smtpConfig.Password

	s.Sender = config.Sender
	s.Receivers = config.Receivers

	return nil
}

func (s *SmtpEngine) Sendmail(subject string, body string) error {
	var content bytes.Buffer

	writer := multipart.NewWriter(&content)

	header := make(textproto.MIMEHeader)
	header.Set("MIME-Version", "1.0")
	header.Set("Content-Type", fmt.Sprintf("multipart/related; boundary=%s", writer.Boundary()))
  date := time.Now().Format(time.RFC1123Z)

  content.WriteString(fmt.Sprintf("From: %s\r\n", s.Sender))
	content.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(s.Receivers, ",")))
	content.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	content.WriteString(fmt.Sprintf("Date: %s\r\n", date))
	content.WriteString(fmt.Sprintf("MIME-Version: 1.0\r\n"))
	content.WriteString(fmt.Sprintf("Content-Type: multipart/related; boundary=%s\r\n\r\n", writer.Boundary()))

	htmlPart, _ := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type":              {"text/html; charset=\"UTF-8\""},
		"Content-Transfer-Encoding": {"quoted-printable"},
	})

	qpWriter := quotedprintable.NewWriter(htmlPart)
	qpWriter.Write([]byte(body))
	qpWriter.Close()

	imagePath := "./bawana.png"
	imageData, err := utils.Read(imagePath)
	if err != nil {
		return fmt.Errorf("failed to send mail: %v", err)
	}
	imageName := filepath.Base(imagePath)

	imagePart, _ := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type":              {"image/png"},
		"Content-Transfer-Encoding": {"base64"},
		"Content-Disposition":       {fmt.Sprintf(`inline; filename="%s"`, imageName)},
		"Content-ID":                {"<bawana.png>"},
	})

	imagePart.Write([]byte(base64.StdEncoding.EncodeToString(imageData)))

	writer.Close()

	auth := smtp.PlainAuth(s.Username, s.Username, s.Password, s.Host)

	err = smtp.SendMail(s.Host+":"+s.Port, auth, s.Sender, s.Receivers, content.Bytes())
	if err != nil {
		return fmt.Errorf("failed to send mail: %v", err)
	}

	logger.Info(fmt.Sprintf("Mail sent successfully to %v", s.Receivers))
	return nil
}
