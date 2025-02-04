package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

// StealthClient, represents a http client with stealth capabilities
type StealthClient struct {
	userAgent  string
	sessionID  string
	headers    map[string]string
	proxyURL   *url.URL
	httpClient *http.Client
	mu         sync.Mutex

	// Embedded configuration that needs to be serialized
	serializable struct {
		UserAgent string            `json:"user_agent"`
		SessionID string            `json:"session_id"`
		Headers   map[string]string `json:"headers"`
		Proxy     string            `json:"proxy,omitempty"`
	}
}

// NewStealthClient creates a new StealthClient
func NewStealthClient(options ...func(*StealthClient)) (*StealthClient, error) {
	client := &StealthClient{
		headers: make(map[string]string),
	}

	// Generate random values
	client.serializable.UserAgent = randomUserAgent()
	client.serializable.SessionID = generateSessionID(5)
	client.sessionID = client.serializable.SessionID

	// Apply options
	for _, option := range options {
		option(client)
	}

	// Initialize HTTP client
	transport := &http.Transport{}
	if client.proxyURL != nil {
		transport.Proxy = http.ProxyURL(client.proxyURL)
	}

	client.httpClient = &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	// Sync serializable data
	client.userAgent = client.serializable.UserAgent
	client.headers = client.serializable.Headers

	return client, nil
}

// Serialize serializes the client to a string
func (c *StealthClient) Serialize() (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.proxyURL != nil {
		c.serializable.Proxy = c.proxyURL.String()
	}

	data, err := json.Marshal(c.serializable)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

// DeserializeClient deserializes a string into a StealthClient
func DeserializeClient(data string) (*StealthClient, error) {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	var serializable struct {
		UserAgent string            `json:"user_agent"`
		SessionID string            `json:"session_id"`
		Headers   map[string]string `json:"headers"`
		Proxy     string            `json:"proxy,omitempty"`
	}

	if err := json.Unmarshal(decoded, &serializable); err != nil {
		return nil, err
	}

	client := &StealthClient{
		serializable: serializable,
		userAgent:    serializable.UserAgent,
		sessionID:    serializable.SessionID,
		headers:      serializable.Headers,
	}

	fmt.Println("Proxy: ", serializable.Proxy)

	// Parse proxy URL
	if strings.HasPrefix(serializable.Proxy, "http://") {

		proxyURL, err := url.Parse(serializable.Proxy)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy URL: %v", err)
		}

		userInfo := strings.Split(proxyURL.User.String(), ":")
		if len(userInfo) != 2 {
			return nil, fmt.Errorf("invalid user info format")
		}

		usernameParts := strings.Split(userInfo[0], "-")
		if len(usernameParts) < 5 {
			return nil, fmt.Errorf("invalid smartproxy username format")
		}

		sessionDuration, err := strconv.Atoi(usernameParts[5])

		port, err := strconv.Atoi(proxyURL.Port())
		if err != nil {
			return nil, fmt.Errorf("invalid proxy port: %v", err)
		}

		cfg := ProxyConfig{
			Host:            proxyURL.Hostname(),
			User:            strings.Join(usernameParts[1:2], "-"), // "sp463pynue"
			ZonePassword:    userInfo[1],
			SessionDuration: sessionDuration,
			Port:            port,
		}

		WithProxy(cfg)(client)
	}

	transport := &http.Transport{}
	if client.proxyURL != nil {
		transport.Proxy = http.ProxyURL(client.proxyURL)
	}

	client.httpClient = &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	return client, nil
}

// WithCustomHeaders adds custom headers to the client
func WithCustomHeaders(headers map[string]string) func(*StealthClient) {
	return func(c *StealthClient) {
		for k, v := range headers {
			c.serializable.Headers[k] = v
		}
	}
}

// Do makes a HTTP request
func (c *StealthClient) Do(req *http.Request) (*http.Response, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	req.Header.Set("User-Agent", c.userAgent)
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	return c.httpClient.Do(req)
}
