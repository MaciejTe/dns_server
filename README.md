# Simple DNS filter

This repository contains simple PoC of Go-based DNS recursive server with filtering capabilities.

### Required apps

1. Docker
2. Docker-compose

### How to run

Execute `make up` command in terminal to launch DNS server and Redis cache containers.


### How to test

Use `dig` to send DNS queries to local DNS server:

1. Filtered DNS response (local DNS server responds with 0.0.0.0 and NOERROR, similar to nextDNS.io)
```
maciek@maciek-dell:~/Git/dns_server$ dig @localhost -p 53000 cdn.cookielaw.org

; <<>> DiG 9.18.18-0ubuntu0.22.04.1-Ubuntu <<>> @localhost -p 53000 cdn.cookielaw.org
; (1 server found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 48418
;; flags: qr aa rd ra; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0

;; QUESTION SECTION:
;cdn.cookielaw.org.		IN	A

;; ANSWER SECTION:
cdn.cookielaw.org.	300	IN	A	0.0.0.0     <<<<---- local DNS server answer

;; Query time: 0 msec
;; SERVER: 127.0.0.1#53000(localhost) (UDP)
;; WHEN: Sun Feb 11 14:12:16 CET 2024
;; MSG SIZE  rcvd: 68

```

2. Unfiltered DNS query
```
maciek@maciek-dell:~/Git/dns_server$ dig @localhost -p 53000 wp.pl

; <<>> DiG 9.18.18-0ubuntu0.22.04.1-Ubuntu <<>> @localhost -p 53000 wp.pl
; (1 server found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 61983
;; flags: qr aa rd ra; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0

;; QUESTION SECTION:
;wp.pl.				IN	A

;; ANSWER SECTION:
wp.pl.			62	IN	A	212.77.98.9

;; Query time: 48 msec
;; SERVER: 127.0.0.1#53000(localhost) (UDP)
;; WHEN: Sun Feb 11 14:12:09 CET 2024
;; MSG SIZE  rcvd: 44
```
