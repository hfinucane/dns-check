package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

func Lookup(ctx context.Context, hostname, dns_server string) ([]string, error) {
	r := &net.Resolver{
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Duration(time.Millisecond * 500),
			}
			return d.DialContext(ctx, "udp", dns_server)
		},
	}

	return r.LookupHost(ctx, hostname)
}

func main() {
	hostname := flag.String("host", "", "host name to look up")
	flag.Parse()

	if *hostname == "" {
		fmt.Println("-host is required")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*1))
	defer cancel()

	for _, dns_server := range flag.Args() {
		strings, err := Lookup(ctx, *hostname, fmt.Sprintf("%s:53", dns_server))
		fmt.Println(strings)
		fmt.Println(err)
		if err != nil {
			os.Exit(2)
		}
	}
}
