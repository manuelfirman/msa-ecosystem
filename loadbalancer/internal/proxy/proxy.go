package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync/atomic"

	"github.com/google/uuid"
)

// Proxy is a reverse proxy that load balances the requests
type Proxy struct {
	// serviceTargets contains the URLs of the services with the service name as the key
	serviceTargets map[string][]*url.URL
	// currentIndexes contains the index of the last selected service for each service
	currentIndexes map[string]uint32
	// servicePaths contains the path of the service with the service name as the key
	servicePaths map[string]string
}

// NewProxy creates a new Proxy
func NewProxy(targets map[string][]string, paths map[string]string) *Proxy {
	serviceTargets := make(map[string][]*url.URL)
	currentIndexes := make(map[string]uint32)

	for service, urls := range targets {
		serviceURLs := make([]*url.URL, len(urls))
		for i, urlStr := range urls {
			url, err := url.Parse(urlStr)
			if err != nil {
				log.Fatal(err)
			}
			serviceURLs[i] = url
		}
		serviceTargets[service] = serviceURLs
		currentIndexes[service] = 0
	}

	return &Proxy{serviceTargets: serviceTargets, currentIndexes: currentIndexes, servicePaths: paths}
}

// ServeHTTP redirige la solicitud a un servicio
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var target *url.URL
	var service string

	// find the service that the request is for
	for path, svc := range p.servicePaths {
		if strings.HasPrefix(r.URL.Path, path) {
			service = svc
			break
		}
	}

	if service == "" {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	// select an instance of the service to forward the request to
	targets := p.serviceTargets[service]
	currIndex := p.currentIndexes[service]
	index := atomic.AddUint32(&currIndex, 1) % uint32(len(targets))
	target = targets[index]

	// Set the X-Request-ID header for tracking the request
	r.Header.Set("X-Request-ID", uuid.New().String())

	// create a reverse proxy that forwards the request to the selected service
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ServeHTTP(w, r)
}
