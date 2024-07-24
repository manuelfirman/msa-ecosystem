package main

import "loadbalancer/internal/application"

func main() {
	app := application.NewApplicationDefault()
	app.Run()
}
