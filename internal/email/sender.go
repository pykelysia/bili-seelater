package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"strings"

	"bili-seelater/internal/bilibili"
	"bili-seelater/internal/config"

	"gopkg.in/gomail.v2"
)

type Sender struct {
	dialer *gomail.Dialer
	from   string
	to     string
}

func NewSender(cfg *config.EmailConfig) *Sender {
	ssl := false
	useTLS := false

	switch cfg.SMTPPort {
	case 465:
		ssl = true
	case 587:
		useTLS = true
	case 25:
		// no ssl/tls
	}

	dialer := gomail.NewDialer(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.Username,
		cfg.Password,
	)
	dialer.SSL = ssl
	if useTLS {
		dialer.TLSConfig = &tls.Config{ServerName: cfg.SMTPHost}
	}

	return &Sender{
		dialer: dialer,
		from:   cfg.From,
		to:     cfg.To,
	}
}

func (s *Sender) SendVideoList(videos []bilibili.Video) error {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", s.from, "B站稍后再看")
	m.SetAddressHeader("To", s.to, "Gmail")
	m.SetHeader("Subject", fmt.Sprintf("B站稍后再看 - 共%d个视频", len(videos)))
	m.SetBody("text/html", s.buildHTML(videos))

	if err := s.dialer.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func (s *Sender) buildHTML(videos []bilibili.Video) string {
	var buf strings.Builder
	buf.WriteString("<html><body><h2>B站稍后再看列表</h2>")
	buf.WriteString(fmt.Sprintf("<p>共 %d 个视频</p>", len(videos)))
	buf.WriteString("<ul>")

	for _, v := range videos {
		link := fmt.Sprintf("https://www.bilibili.com/video/%s", v.Bvid)
		duration := fmt.Sprintf("%d:%02d", v.Duration/60, v.Duration%60)
		buf.WriteString(fmt.Sprintf(
			`<li><a href="%s">%s</a> - %s - %s</li>`,
			link, v.Title, v.Author, duration,
		))
	}

	buf.WriteString("</ul></body></html>")
	return buf.String()
}

func (s *Sender) SendPlainText(videos []bilibili.Video) error {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("B站稍后再看列表 - 共%d个视频\n\n", len(videos)))

	for i, v := range videos {
		link := fmt.Sprintf("https://www.bilibili.com/video/%s", v.Bvid)
		buf.WriteString(fmt.Sprintf("%d. %s - %s\n%s\n\n", i+1, v.Title, v.Author, link))
	}

	m := gomail.NewMessage()
	m.SetAddressHeader("From", s.from, "B站稍后再看")
	m.SetAddressHeader("To", s.to, "Gmail")
	m.SetHeader("Subject", fmt.Sprintf("B站稍后再看 - 共%d个视频", len(videos)))
	m.SetBody("text/plain", buf.String())

	if err := s.dialer.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
