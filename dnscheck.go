package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net"
	"net/smtp"
)

type Opts struct {
	email        string
	sender       string
	hostname     string
	ipAddr       string
	expectedAddr string
}

func sendemail(s Opts) int {
	c, err := smtp.Dial("localhost:25")
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
	flag.Parse()

	netaddr, _ := net.ResolveIPAddr("ip4", *hostPtr)

	ns := Opts{
		email:        *emailPtr,
		sender:       *senderPtr,
		hostname:     *hostPtr,
		ipAddr:       *ipPtr,
		expectedAddr: netaddr.String(),
	}

	if ns.expectedAddr != ns.ipAddr {
		fmt.Printf("%s\n", ns.expectedAddr)
		sendemail(ns)
	}

}
