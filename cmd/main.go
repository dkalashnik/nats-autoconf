package main

import (
	"flag"
	"log"

	"github.com/dkalashnik/nats-autoconf/pkg/nats"
	"github.com/dkalashnik/nats-autoconf/pkg/operator"
	"github.com/dkalashnik/nats-autoconf/pkg/service"
)

func main() {
	mode := flag.String("mode", "", "Operation mode: operator or service")
	natsURL := flag.String("nats-url", "nats://localhost:4222", "NATS server URL")
	serviceName := flag.String("service-name", "", "Service name (required for service mode)")
	flag.Parse()

	if *mode == "" {
		log.Fatal("Mode must be specified: -mode=operator or -mode=service")
	}

	natsClient, err := nats.NewClient(*natsURL, "service-configs")
	if err != nil {
		log.Fatalf("Failed to initialize NATS client: %v", err)
	}
	defer natsClient.Close()

	switch *mode {
	case "operator":
		if err := operator.Run(natsClient); err != nil {
			log.Fatalf("Operator failed: %v", err)
		}
	case "service":
		if *serviceName == "" {
			log.Fatal("Service name must be specified in service mode")
		}
		if err := service.Run(*serviceName, natsClient); err != nil {
			log.Fatalf("Service failed: %v", err)
		}
	default:
		log.Fatalf("Unknown mode: %s", *mode)
	}
}
