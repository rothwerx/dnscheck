# dnscheck

This will send an email if a DNS name doesn't resolve to the IP address you expect.
You'd run it from cron every X minutes. It will send an email on the first failure
and then every hour after that.

It's an admittedly limited use case, but was handy to try to narrow down an internal
vs external DNS issue for a client.
