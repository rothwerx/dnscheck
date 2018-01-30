package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"os"
	"time"
)

type opts struct {
	email        string
	sender       string
	hostname     string
	ipAddr       string
	expectedAddr string
	smtp         string
}

func sendemail(s *opts) {
	c, err := smtp.Dial(s.smtp)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	c.Mail(s.sender)
	c.Rcpt(s.email)
	wc, err := c.Data()
	if err != nil {
		log.Fatal(err)
	}
	defer wc.Close()
	ns := fmt.Sprintf("Address %s is resolving to %s\n",
		s.hostname, s.expectedAddr)
	buf := bytes.NewBufferString(ns)
	if _, err = buf.WriteTo(wc); err != nil {
		log.Fatal(err)
	}
}

func main() {
	hostname := flag.String("hostname", "google.com", "Hostname to resolve")
	ip := flag.String("ip", "216.58.216.142", "Expected IP")
	email := flag.String("email", "", "Email address to send failures")
	sender := flag.String("sender", "example@example.com", "Email address to send as")
	smtp := flag.String("smtp", "localhost:25", "SMTP address:port to use")
	flag.Parse()

	netaddr, err := net.ResolveIPAddr("ip4", *hostname)
	if err != nil {
		log.Fatal(err)
	}
	ns := &opts{
		email:        *email,
		sender:       *sender,
		hostname:     *hostname,
		ipAddr:       *ip,
		smtp:         *smtp,
		expectedAddr: netaddr.String(),
	}

	// To keep track of time, create a hidden file
	timefile := ".dnsdrop"
	if ns.expectedAddr != ns.ipAddr {
		// Send an email once when it begins to fail
		if fl, err := os.Stat(timefile); err != nil {
			os.Create(timefile)
			sendemail(ns)
		} else {
			// if it's still failing, send another email after an hour
			duration := time.Since(fl.ModTime())
			since := duration.Minutes()
			if since > 60 {
				sendemail(ns)
				// Reset mod time to start the 60 minute timer
				now := time.Now().Local()
				err := os.Chtimes(timefile, now, now)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	} else {
		// Remove the file if it's no longer failing
		if _, err := os.Stat(timefile); err == nil {
			os.Remove(timefile)
		}
	}
}
