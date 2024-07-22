// load_balancer/main.go
package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Proxy struct {
	targets []*url.URL
}

func NewProxy(targets []string) *Proxy {
	urls := make([]*url.URL, len(targets))
	for i, target := range targets {
		url, err := url.Parse(target)
		if err != nil {
			log.Fatal(err)
		}
		urls[i] = url
	}
	return &Proxy{targets: urls}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	target := p.targets[0] // Aquí podrías implementar tu lógica de balanceo de carga
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ServeHTTP(w, r)
}

func main() {
	targets := []string{
		"http://auth_service:5000",
		"http://user_service:5001",
		"http://product_service:5002",
		"http://order_service:5003",
		"http://notification_service:5004",
	}
	proxy := NewProxy(targets)
	http.Handle("/", proxy)
	log.Println("Listening on :80")
	log.Fatal(http.ListenAndServe(":80", nil))
}
