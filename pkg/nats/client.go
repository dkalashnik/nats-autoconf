package nats

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dkalashnik/nats-autoconf/pkg/apis/config/v1alpha1"
	"github.com/nats-io/nats.go"
)

type Client struct {
	nc     *nats.Conn
	kv     nats.KeyValue
	bucket string
}

func NewClient(url, bucket string) (*Client, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %v", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, fmt.Errorf("failed to create JetStream context: %v", err)
	}

	kv, err := js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket: bucket,
	})
	if err != nil {
		kv, err = js.KeyValue(bucket)
		if err != nil {
			return nil, fmt.Errorf("failed to access KV bucket: %v", err)
		}
	}

	return &Client{
		nc:     nc,
		kv:     kv,
		bucket: bucket,
	}, nil
}

func (c *Client) Close() {
	if c.nc != nil {
		c.nc.Close()
	}
}

func GetConfigPath(org, product, version, service, config string) string {
	return fmt.Sprintf("%s.%s.%s.%s.%s", org, product, version, service, config)
}

func (c *Client) PutConfig(key string, config *v1alpha1.VPPCaptureConfig) error {
	data, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	_, err = c.kv.Put(key, data)
	if err != nil {
		return fmt.Errorf("failed to store config: %v", err)
	}

	return nil
}

func (c *Client) GetConfig(key string) (*v1alpha1.VPPCaptureConfig, error) {
	entry, err := c.kv.Get(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %v", err)
	}

	var config v1alpha1.VPPCaptureConfig
	if err := json.Unmarshal(entry.Value(), &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	return &config, nil
}

func (c *Client) WatchConfigs(serviceName string, callback func(key string, config *v1alpha1.VPPCaptureConfig)) error {
	watcher, err := c.kv.WatchAll()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %v", err)
	}

	go func() {
		for entry := range watcher.Updates() {
			if entry == nil {
				continue
			}

			parts := strings.Split(entry.Key(), ".")
			if len(parts) != 5 || parts[3] != serviceName {
				continue
			}

			var config v1alpha1.VPPCaptureConfig
			if err := json.Unmarshal(entry.Value(), &config); err != nil {
				continue
			}

			callback(entry.Key(), &config)
		}
	}()

	return nil
}

func (c *Client) DeleteKey(key string) error {
	err := c.kv.Delete(key)
	if err != nil {
		return fmt.Errorf("failed to delete key %s: %v", key, err)
	}
	return nil
}

func (c *Client) ListKeys() ([]string, error) {
	keys, err := c.kv.Keys()
	if err != nil {
		return nil, fmt.Errorf("failed to list keys: %v", err)
	}
	return keys, nil
}
