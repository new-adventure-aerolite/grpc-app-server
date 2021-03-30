package istio

import (
	"fmt"
	"net/http"
	"time"

	"k8s.io/klog"
)

// Proxy contains methods for interacting with the istio-proxy container.
type Proxy struct {
	client     *http.Client
	maxRetries int
	retryDelay time.Duration
	quitURL    string
	healthzURL string
}

// New returns a Proxy.
func New(maxRetries int, timeout, retryDelay time.Duration) *Proxy {
	return &Proxy{
		client:     &http.Client{Timeout: timeout},
		maxRetries: maxRetries,
		retryDelay: retryDelay,
		quitURL:    "http://127.0.0.1:15000/quitquitquit",
		healthzURL: "http://127.0.0.1:15020/healthz/ready",
	}
}

// Wait waits until the istio-proxy container is ready before returning.
func (p *Proxy) Wait() error {
	retries := 0

	for {
		retries++
		if retries > p.maxRetries {
			return fmt.Errorf("proxy: max retries reached for Wait(), maxRetries = %d", p.maxRetries)
		}

		response, err := p.client.Get(p.healthzURL)
		if err != nil {
			klog.Infof("proxy: wait client get failed, retries = %d, maxRetries = %d, error = %s", retries, p.maxRetries, err)
			klog.Infof("retrying in %.2f seconds...", p.retryDelay.Seconds())
			time.Sleep(p.retryDelay)
			continue
		}

		if response.StatusCode == http.StatusOK {
			break
		} else {
			klog.Infof("proxy: wait unexpected response code, retries = %d, maxRetries = %d, responseCode = %d", retries, p.maxRetries, response.StatusCode)
			klog.Infof("retrying in %.2f seconds...", p.retryDelay.Seconds())
			time.Sleep(p.retryDelay)
			continue
		}
	}

	return nil
}

// Close closes the istio-proxy container.
func (p *Proxy) Close() error {
	retries := 0

	for {
		retries++
		if retries > p.maxRetries {
			return fmt.Errorf("proxy: max retries reached for Close(), maxRetries = %d", p.maxRetries)
		}

		response, err := p.client.Post(p.quitURL, "application/json", nil)
		if err != nil {
			klog.Infof("proxy: close client post failed, retries = %d, maxRetries = %d, error = %s", retries, p.maxRetries, err)
			klog.Infof("retrying in %.2f seconds...", p.retryDelay.Seconds())
			time.Sleep(p.retryDelay)
			continue
		}

		if response.StatusCode == http.StatusOK {
			break
		} else {
			klog.Infof("proxy: close unexpected response code, retries = %d, maxRetries = %d, responseCode = %d", retries, p.maxRetries, response.StatusCode)
			klog.Infof("retrying in %.2f seconds...", p.retryDelay.Seconds())
			time.Sleep(p.retryDelay)
			continue
		}
	}

	return nil
}
