package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync/atomic"
)

// Proxy es un servidor proxy que redirige las solicitudes a un servicio
type Proxy struct {
	serviceTargets map[string][]*url.URL
	currentIndexes map[string]uint32
}

// NewProxy crea una instancia de Proxy
func NewProxy(targets map[string][]string) *Proxy {
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

	return &Proxy{serviceTargets: serviceTargets, currentIndexes: currentIndexes}
}

// ServeHTTP redirige la solicitud a un servicio
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var target *url.URL
	var service string
	path := r.URL.Path

	switch {
	case strings.HasPrefix(path, "/auth"):
		service = "auth_service"
	case strings.HasPrefix(path, "/users"):
		service = "user_service"
		// r.URL.Path = strings.TrimPrefix(path, "/users")
	case strings.HasPrefix(path, "/products"):
		service = "product_service"
		// r.URL.Path = strings.TrimPrefix(path, "/products")
	case strings.HasPrefix(path, "/orders"):
		service = "order_service"
		// r.URL.Path = strings.TrimPrefix(path, "/orders")
	case strings.HasPrefix(path, "/notifications"):
		service = "notification_service"
		// r.URL.Path = strings.TrimPrefix(path, "/notifications")
	default:
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	// Selecciona una instancia del servicio
	targets := p.serviceTargets[service]
	currIndex := p.currentIndexes[service]
	index := atomic.AddUint32(&currIndex, 1) % uint32(len(targets))
	target = targets[index]
	// Redirige la solicitud al servicio
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ServeHTTP(w, r)
}

// Service representa un servicio
type Service struct {
	Name string
	Host string
	Port string
}

// URL devuelve la URL del servicio
func (s Service) URL() string {
	return "http://" + s.Host + ":" + s.Port
}

func main() {
	auth := Service{Name: "auth_service", Host: "auth_service", Port: "5000"}
	user := Service{Name: "user_service", Host: "user_service", Port: "5001"}
	product := Service{Name: "product_service", Host: "product_service", Port: "5002"}
	order := Service{Name: "order_service", Host: "order_service", Port: "5003"}
	notification := Service{Name: "notification_service", Host: "notification_service", Port: "5004"}

	targets := map[string][]string{
		auth.Name:         {auth.URL()},
		user.Name:         {user.URL()},
		product.Name:      {product.URL()},
		order.Name:        {order.URL()},
		notification.Name: {notification.URL()},
	}

	proxy := NewProxy(targets)
	http.Handle("/", proxy)
	log.Println("Listening on :80")
	log.Fatal(http.ListenAndServe(":80", nil))
}
