package main

import (
	"context"
	"dns_server/loader"
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/miekg/dns"
	"github.com/redis/go-redis/v9"
)

const (
	recursiveDNS = "1.1.1.1:53"
	redisAddr    = "localhost:6379"
	port         = ":53000"
)

var (
	blacklistSources = []string{
		"./blacklists/domains.txt", // source: "https://hole.cert.pl/domains/v2/domains.txt",
	}
)

func resolve(domain string, qtype uint16) []dns.RR {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), qtype)
	m.RecursionDesired = true

	c := new(dns.Client)
	in, _, err := c.Exchange(m, recursiveDNS)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return in.Answer
}

type dnsHandler struct {
	redisClient *redis.Client
}

func (h *dnsHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := new(dns.Msg)
	msg.SetReply(r)
	msg.Authoritative = true
	msg.RecursionAvailable = true

	for _, question := range r.Question {
		fmt.Println("Received query: ", spew.Sdump(question))
		if err := h.filterFunc(context.Background(), question); err != nil {
			fmt.Print(err.Error())
			// nextdns answers with: cdn.cookielaw.org.	300	IN	A	0.0.0.0
			fakeRR, err := dns.NewRR(fmt.Sprintf("%s	300	IN	A	0.0.0.0", question.Name))
			if err != nil {
				fmt.Printf("Fake resource record creation error: %s\n", err.Error())
			}
			msg.Answer = append(msg.Answer, fakeRR)
		} else {
			answers := resolve(question.Name, question.Qtype)
			msg.Answer = append(msg.Answer, answers...)
		}
	}

	w.WriteMsg(msg)
}

func (h *dnsHandler) filterFunc(ctx context.Context, question dns.Question) error {
	fqdn := strings.Split(question.Name, ".")
	fqdn = fqdn[:len(fqdn)-1]
	val, err := h.redisClient.Get(ctx, strings.Join(fqdn[:len(fqdn)-1], ".")).Result()
	if val != "" {
		return fmt.Errorf("blocking domain %s", question.Name)
	}
	if err != redis.Nil {
		return err
	}
	return nil
}

func main() {
	// feed redis with blacklist data
	rdb := redis.NewClient(&redis.Options{
		Addr:         redisAddr,
		Password:     "",
		DB:           0,
		MaxIdleConns: 5,
	})
	kv := loader.NewKeyValueLoader(rdb)
	ctx := context.Background()
	for _, source := range blacklistSources {
		if err := kv.Load(ctx, source); err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("Blacklists uploaded to Redis")

	handler := new(dnsHandler)
	handler.redisClient = rdb
	server := &dns.Server{
		Addr:      port,
		Net:       "udp",
		Handler:   handler,
		UDPSize:   65535,
		ReusePort: true,
	}

	fmt.Println("Starting DNS server on port ", port)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Failed to start server: %s\n", err.Error())
	}
}
