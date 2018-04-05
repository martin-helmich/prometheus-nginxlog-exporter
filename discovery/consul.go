package discovery

import "github.com/martin-helmich/prometheus-nginxlog-exporter/config"
import "github.com/hashicorp/consul/api"

// ConsulRegistrator is a helper struct that handles Consul service registration
type ConsulRegistrator struct {
	config    *config.Config
	client    *api.Client
	serviceID string
}

func getDefault(a string, b string) string {
	if a == "" {
		return b
	}
	return a
}

// NewConsulRegistrator is a constructor function for building a new ConsulRegistrator
func NewConsulRegistrator(cfg *config.Config) (*ConsulRegistrator, error) {
	config := api.Config{
		Address:    getDefault(cfg.Consul.Address, "localhost:8500"),
		Datacenter: getDefault(cfg.Consul.Datacenter, "dc1"),
		Scheme:     getDefault(cfg.Consul.Scheme, "http"),
		Token:      cfg.Consul.Token,
	}

	client, err := api.NewClient(&config)
	if err != nil {
		return nil, err
	}

	name := getDefault(cfg.Consul.Service.Name, "nginx-exporter")
	serviceID := getDefault(cfg.Consul.Service.ID, name)

	return &ConsulRegistrator{
		config:    cfg,
		client:    client,
		serviceID: serviceID,
	}, nil
}

// RegisterConsul registers the exporter instance at Consul
func (r *ConsulRegistrator) RegisterConsul() error {
	registration := api.AgentServiceRegistration{
		ID:   r.serviceID,
		Port: r.config.Listen.Port,
		Name: getDefault(r.config.Consul.Service.Name, "nginx-exporter"),
		Tags: r.config.Consul.Service.Tags,
	}

	err := r.client.Agent().ServiceRegister(&registration)
	if err != nil {
		return err
	}

	return nil
}

// UnregisterConsul deregisters the exporter from Consul again
func (r *ConsulRegistrator) UnregisterConsul() error {
	return r.client.Agent().ServiceDeregister(r.serviceID)
}
