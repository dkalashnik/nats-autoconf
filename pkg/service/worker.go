package service

import (
	"context"
	"sync"

	"github.com/dkalashnik/nats-autoconf/pkg/apis/config/v1alpha1"
	"github.com/dkalashnik/nats-autoconf/pkg/nats"
)

type Worker struct {
	config     *v1alpha1.VPPCaptureConfig
	cancelFunc context.CancelFunc
}

type Service struct {
	serviceName string
	natsClient  *nats.Client
	workers     map[string]*Worker
	mu          sync.RWMutex
}

func NewService(serviceName string, natsClient *nats.Client) *Service {
	return &Service{
		serviceName: serviceName,
		natsClient:  natsClient,
		workers:     make(map[string]*Worker),
	}
}

func (s *Service) startWorker(key string, config *v1alpha1.VPPCaptureConfig) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if worker, exists := s.workers[key]; exists {
		worker.cancelFunc()
	}

	ctx, cancel := context.WithCancel(context.Background())

	s.workers[key] = &Worker{
		config:     config,
		cancelFunc: cancel,
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}()
}

func Run(serviceName string, natsClient *nats.Client) error {
	service := NewService(serviceName, natsClient)

	err := natsClient.WatchConfigs(serviceName, service.startWorker)
	if err != nil {
		return err
	}

	select {}
}
