package main

import (
	"context"
	"log"
	"net"
	"net/url"
	"os"
	"strings"

	"github.com/things-go/go-socks5"
	"golang.org/x/net/proxy"
)

// LoadBalancer is a simple round-robin load balancer for proxies
type LoadBalancer struct {
	index int
	proxies []proxy.Dialer
}

// Dial returns the next proxy in the list
func (lb *LoadBalancer) Dial(ctx context.Context, network string, address string) (net.Conn, error) {
	dialer := lb.proxies[lb.index]

	socket, err := dialer.Dial(network, address)
	if (err != nil) {
		return nil, err
	}

	lb.index = (lb.index + 1) % len(lb.proxies)

	return socket, nil
}

func getenv(key, fallback string) string {
    value := os.Getenv(key)
    if len(value) == 0 {
        return fallback
    }
    return value
}

func main() {
	host := getenv("HOST", "")
	port := getenv("PORT", "1080")
	username := getenv("USERNAME", "")
	password := getenv("PASSWORD", "")	

	addr := host + ":" + port

	socks_proxies := getenv("PROXIES", "")
	proxies := make([]proxy.Dialer, 0)
	
	for _, p := range strings.Split(socks_proxies, ",") {
		url, err := url.Parse(p);
		if err != nil {
			log.Fatal(err)
		}

		dialer, err := proxy.FromURL(url, proxy.Direct)
		if err != nil {
			log.Fatal(err)
		}

		proxies = append(proxies, dialer)
	}

	lb := &LoadBalancer{
		index: 0,
		proxies: proxies,
	}

	opts := []socks5.Option{
		socks5.WithDial(lb.Dial),
	}

	if username != "" && password != "" {
		opts = append(opts, socks5.WithCredential(socks5.StaticCredentials{
			username: username,
			password: password,
		}))
	}

	server := socks5.NewServer(opts...)

	if err := server.ListenAndServe("tcp", addr); err != nil {
		panic(err)
	}
}
