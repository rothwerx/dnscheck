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

type Opts struct {
	email        string
	sender       string
	hostname     string
	ipAddr       string
	expectedAddr string
	tSmtp        string
}

func sendemail(s *Opts) int {
	c, err := smtp.Dial(s.tSmtp)
	if err != nil {
		log.Fatal(err)
		return 1
	}
	defer c.Close()

	c.Mail(s.sender)
	c.Rcpt(s.email)
	wc, err := c.Data()
	if err != nil {
		log.Fatal(err)
		return 1
	}
	defer wc.Close()
	ns := fmt.Sprintf("Address %s is resolving to %s\n",
		s.hostname, s.expectedAddr)
	buf := bytes.NewBufferString(ns)
	if _, err = buf.WriteTo(wc); err != nil {
		log.Fatal(err)
		return 1
	}
	return 0
}

func main() {
	hostPtr := flag.String("hostname", "google.com", "Hostname to resolve")
	ipPtr := flag.String("ip", "216.58.216.142", "Expected IP")
	emailPtr := flag.String("email", "", "Email address to send failures")
	senderPtr := flag.String("sender", "example@example.com", "Email address to send as")
	smtpPtr := flag.String("smtp", "localhost:25", "SMTP address:port to use")
	flag.Parse()

	netaddr, err := net.ResolveIPAddr("ip4", *hostPtr)
	if err != nil {
		log.Fatal(err)
	}
	ns := &Opts{
		email:        *emailPtr,
		sender:       *senderPtr,
		hostname:     *hostPtr,
		ipAddr:       *ipPtr,
		tSmtp:        *smtpPtr,
		expectedAddr: netaddr.String(),
	}

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
			}
		}
	} else {
		if _, err := os.Stat(timefile); err == nil {
			os.Remove(timefile)
		}
	}
}
