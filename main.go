package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

// Heritage Nagios Return Codes
const (
	OK = iota
	WARNING
	CRITICAL
	UNKNOWN
)

type LookupResult struct {
	Results []string
	Server  string
	Time    time.Duration
	Error   error
}

func (lr *LookupResult) Print() {
	if lr.Error == nil {
		fmt.Println(lr.Server, lr.Results, "in", lr.Time)
	} else {
		fmt.Println("Failed:", lr.Error)
	}
}

func Lookup(ctx context.Context, hostname, dns_server string) *LookupResult {
	start := time.Now()
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			timeout := time.Second * 1
			deadline, ok := ctx.Deadline()
			if ok {
				timeout = time.Until(deadline)
			}

			return (&net.Dialer{
				Timeout: timeout,
			}).DialContext(ctx, "udp", dns_server)
		},
	}

	ips, err := r.LookupHost(ctx, hostname)
	return &LookupResult{
		Server:  dns_server,
		Results: ips,
		Time:    time.Since(start),
		Error:   err,
	}
}

func main() {
	hostname := flag.String("host", "", "host name to look up")
	deadline := flag.Int("deadline", 5, "deadline in seconds")
	warn_threshold := flag.Int("w", 1, "warn threshold - percent of servers failing")
	critical_threshold := flag.Int("c", 1, "critical threshold - percent of servers failing")
	full_deadline := time.Second * time.Duration(*deadline)
	flag.Parse()

	if *hostname == "" {
		fmt.Println("-host is required")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(full_deadline))
	defer cancel()

	results := make(chan *LookupResult, len(flag.Args()))

	for _, dns_server := range flag.Args() {
		go func(dns_server string) {
			results <- Lookup(ctx, *hostname, fmt.Sprintf("%s:53", dns_server))
		}(dns_server)
	}

	errors := 0
	start := time.Now()
	for i := 0; i < len(flag.Args()); i++ {
		select {
		case lookup := <-results:
			lookup.Print()
			if lookup.Error != nil {
				errors++
			}
		case <-time.After(full_deadline):
			break
		}
	}
	fmt.Println("Done in", time.Since(start))

	if errors > 0 {
		error_percent := errors * 100 / len(flag.Args())
		if error_percent >= *critical_threshold {
			os.Exit(CRITICAL)
		} else if error_percent >= *warn_threshold {
			os.Exit(WARNING)
		}
	}
	os.Exit(OK)
}
