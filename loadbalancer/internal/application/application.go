package application

import (
	"loadbalancer/internal/proxy"
	"log"
	"net/http"
	"os"
	"strings"
)

// Application is the interface that wraps the Run method
type ApplicationDefault struct {
}

func NewApplicationDefault() *ApplicationDefault {
	return &ApplicationDefault{}
}

func (a *ApplicationDefault) Run() {
	port := os.Getenv("PORT")
	services, servicePaths := LoadServicesFromEnv()
	proxy := proxy.NewProxy(services, servicePaths)
	http.Handle("/", proxy)
	log.Println("Listening on", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func LoadServicesFromEnv() (map[string][]string, map[string]string) {
	services := make(map[string][]string)
	servicePaths := make(map[string]string)

	serviceNames := strings.Split(os.Getenv("SERVICES"), ",")
	for _, serviceName := range serviceNames {
		urls := strings.Split(os.Getenv(serviceName+"_URLS"), ",")
		services[serviceName] = urls

		path := os.Getenv(serviceName + "_PATH")
		servicePaths[path] = serviceName
	}

	return services, servicePaths
}
