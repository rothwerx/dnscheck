package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net"
	"net/smtp"
)

func main() {
	hostPtr := flag.String("hostname", "google.com", "Hostname to resolve")
	ipPtr := flag.String("ip", "216.58.216.142", "Expected IP")
	emailPtr := flag.String("email", "", "Email to send failures to")
	flag.Parse()

	netaddr, _ := net.ResolveIPAddr("ip4", *hostPtr)
	if netaddr.String() != *ipPtr {
		fmt.Printf("%s\n", netaddr)

		c, err := smtp.Dial("localhost:25")
		if err != nil {
			log.Fatal(err)
		}
		defer c.Close()

		c.Mail("wupadmin@wildutahproject.org")
		c.Rcpt(*emailPtr)
		wc, err := c.Data()
		if err != nil {
			log.Fatal(err)
		}
		defer wc.Close()
		buf := bytes.NewBufferString(netaddr.String())
		if _, err = buf.WriteTo(wc); err != nil {
			log.Fatal(err)
		}
	}
}
