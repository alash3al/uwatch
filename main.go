package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"gopkg.in/gomail.v2"
)

var (
	flagSMTPHost     = flag.String("smtp-host", "smtp.gmail.com", "the smtp host")
	flagSMTPPort     = flag.Int("smtp-port", 587, "the smtp server port")
	flagSMTPUser     = flag.String("smtp-user", "yourmail@gmail.com", "the smtp server user")
	flagSMTPPassword = flag.String("smtp-password", "secret", "the smtp server password")
	flagInterval     = flag.Int64("interval", 60, "interval in seconds")
	flagURL          = flag.String("url", "https://facebook.com/alash3al", "the webpage url to watch")
	flagChrome       = flag.String("chrome", "chrome", "the path to chrome binary")
	flagFromEmail    = flag.String("email-from", "no-reply@your.domain", "the from email address")
	flagToEmail      = flag.String("email-to", "email@example.com,email2@example.com", "the target email address(es)")
)

func init() {
	flag.Parse()
}

func main() {
	ticker := time.NewTicker(time.Duration(*flagInterval) * time.Second)

	for range ticker.C {
		log.Println("Checking ...")

		res, err := http.Get(*flagURL)
		if err != nil {
			log.Println("Net::Error::" + err.Error())
			continue
		}

		if res.StatusCode >= 200 && res.StatusCode < 300 {
			log.Println("ONLINE !")
			img, err := screenshot(*flagURL, "./screenshot.png", "1280,1796")
			if err != nil {
				log.Println("The Page Is Online, But I cannot take a screenshot")
			} else {
				log.Println("Sending Mail ...")
				err := mail("BackOnline | "+*flagURL, "The URL ("+*flagURL+") is now online, the screent in the attachments", img)
				if err != nil {
					log.Println("SMTP::Error::" + err.Error())
				} else {
					log.Println("Mail Sent Successfully!")
				}
			}
		} else {
			log.Println("OFFLINE !")
		}

		defer res.Body.Close()
	}
}

func screenshot(url string, filename, dimensions string) (string, error) {
	err := exec.Command(
		*flagChrome,
		"--no-sandbox",
		"--window-size="+dimensions,
		"--headless",
		"--disable-gpu",
		"--screenshot="+filename,
		url,
	).Run()

	if err != nil {
		return "", err
	}

	return filename, nil
}

func mail(subject string, body string, filename string) error {
	d := gomail.NewDialer(*flagSMTPHost, *flagSMTPPort, *flagSMTPUser, *flagSMTPPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	m := gomail.NewMessage()

	m.SetHeader("From", *flagFromEmail)
	m.SetHeader("To", strings.Split(*flagToEmail, ",")...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	m.Attach(filename)

	return d.DialAndSend(m)
}
