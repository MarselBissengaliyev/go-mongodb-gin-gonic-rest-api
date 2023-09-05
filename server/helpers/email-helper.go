package helpers

import (
	"bytes"
	"crypto/tls"
	"html/template"
	"log"

	"os"
	"path/filepath"

	"github.com/MarselBisengaliev/go-react-blog/config"
	"github.com/MarselBisengaliev/go-react-blog/models"
	"github.com/k3a/html2text"
	"gopkg.in/gomail.v2"
)

type EmailHelper struct{}

type EmailData struct {
	URL       string
	FirstName string
	Subject   string
}

// Email template parser
func (h *EmailHelper) ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}

func (h *EmailHelper) SendEmail(user *models.User, data *EmailData) {
	conf, err := config.LoadConfig("./")

	if err != nil {
		log.Fatal("could not load config", err)
	}

	from := conf.EmailFrom
	smtpPass := conf.SMTPPass
	smtpUser := conf.SMTPUser
	to := user.Email
	smtpHost := conf.SMTPHost
	smtpPort := conf.SMTPPort

	var body bytes.Buffer

	template, err := h.ParseTemplateDir("templates")
	if err != nil {
		log.Fatal("Could not parse template", err)
	}

	template.ExecuteTemplate(&body, "verificationCode.html", &data)

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", *to)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send email
	if err := d.DialAndSend(m); err != nil {
		log.Fatal("Could not send email: ", err)
	}
}
