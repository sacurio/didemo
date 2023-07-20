package main

import (
	"errors"
	"fmt"
)

type (
	emailDriver int

	Mailer interface {
		Send(to, body string)
	}

	Newsletter struct {
		subscribers []Subscriber
	}

	Subscriber struct {
		name   string
		email  string
		mailer Mailer
	}

	MailMonkey struct {
		apiKey    string
		secretKey string
	}

	MailInternal struct {
		host string
		port int
	}

	MailDefault struct {
		host     string
		password string
		port     int
	}
)

const (
	emailProviderMonkey emailDriver = iota
	emailProviderInternal
	emailProviderDefault
	emailUnknownProvider
)

var (
	providers map[emailDriver]Mailer
)

func (mm MailMonkey) Send(to, body string) {
	fmt.Printf(`Sending message:"%s" to %s, using key:%s and secret:"%s`, body, to, mm.apiKey, mm.secretKey)
}

func NewMailMonkey(ak, as string) MailMonkey {
	return MailMonkey{ak, as}
}

func (mi MailInternal) Send(to, body string) {
	fmt.Printf(`Sending message:"%s" to %s, through host:%s and port:"%d %s`, body, to, mi.host, mi.port, "\n")
}

func NewMailInternal(host string, port int) MailInternal {
	return MailInternal{host, port}
}

func (md MailDefault) Send(to, body string) {
	fmt.Printf(`Sending message:"%s" to %s, through host:%s password: %s, and port:"%d %s`, body, to, md.host, md.password, md.port, "\n")
}

func NewMailDefault(host, password string, port int) MailDefault {
	return MailDefault{host, password, port}
}

func (nl Newsletter) Announce(message string) error {
	if len(nl.subscribers) == 0 {
		return errors.New("no subscribers registered")
	}

	for _, s := range nl.subscribers {
		if s.mailer == nil {
			fmt.Printf("No email provider assigned to %s (%s) \n", s.email, s.name)
			continue
		}
		s.mailer.Send(s.email, message)
	}

	return nil
}

func init() {
	providers = map[emailDriver]Mailer{
		emailProviderMonkey:   NewMailMonkey("ak-74894", "as-12345"),
		emailProviderInternal: NewMailInternal("host-4528", 443),
		emailProviderDefault:  NewMailDefault("host-3654", "default-Password", 512),
	}
}

func main() {
	subscribers := []Subscriber{
		{name: "Sandy", email: "sandy@mail.com", mailer: EmailService(emailProviderInternal)},
		{name: "Valentina", email: "valentina@mail.com", mailer: EmailService(emailProviderMonkey)},
		{name: "Martina", email: "martina@mail.com", mailer: EmailService(emailProviderInternal)},
		{name: "Foo", email: "foo@gmail.com"},
		{name: "Boo", email: "boo@gmail.com", mailer: EmailService(emailUnknownProvider)},
	}

	nl := Newsletter{
		subscribers: subscribers,
	}

	if err := nl.Announce("lorem ipsum ad dolorem"); err != nil {
		fmt.Println(err)
	}

	fmt.Printf("\n\nAnnouncements were send successfully!\n")
}

func EmailService(driver emailDriver) Mailer {
	if m, ok := providers[driver]; ok {
		return m
	}
	return providers[emailProviderDefault]
}
