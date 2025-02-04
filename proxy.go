package stealthy

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type ProxyConfig struct {
	Host            string
	User            string
	ZonePassword    string
	SessionDuration int
	Port            int
}

// WithProxy adds a proxy to the client
func WithProxy(cfg ProxyConfig) func(*StealthClient) {
	return func(c *StealthClient) {
		user := fmt.Sprintf("user-%s-session-%s-sessionduration-%d",
			cfg.User,
			c.sessionID,
			cfg.SessionDuration)

		proxyUrl := fmt.Sprintf("http://%s:%s@%s:%d",
			user,
			cfg.ZonePassword,
			cfg.Host,
			cfg.Port)

		u, _ := url.Parse(proxyUrl)
		c.proxyURL = u

		c.serializable.Proxy = fmt.Sprintf(
			"http://%s:%s:%d:%d",
			cfg.User,
			cfg.ZonePassword,
			cfg.SessionDuration,
			cfg.Port)

	}
}

func (c *StealthClient) RotateProxySession() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	newSessionID := generateSessionID(5)

	if strings.HasPrefix(c.serializable.Proxy, "http://") {

		proxyURL, err := url.Parse(c.serializable.Proxy)
		if err != nil {
			return err
		}

		userInfo := strings.Split(proxyURL.User.String(), ":")
		if len(userInfo) != 2 {
			return fmt.Errorf("invalid user info format")
		}

		usernameParts := strings.Split(userInfo[0], "-")
		if len(usernameParts) < 5 {
			return fmt.Errorf("invalid smartproxy username format")
		}

		sessionDuration, err := strconv.Atoi(usernameParts[5])

		port, err := strconv.Atoi(proxyURL.Port())
		if err != nil {
			return fmt.Errorf("invalid proxy port: %v", err)
		}

		cfg := ProxyConfig{
			Host:            proxyURL.Hostname(),
			User:            strings.Join(usernameParts[1:2], "-"), // "sp463pynue"
			ZonePassword:    userInfo[1],
			SessionDuration: sessionDuration,
			Port:            port,
		}

		// Neue Session ID setzen
		c.sessionID = newSessionID
		c.serializable.SessionID = newSessionID

		// Proxy neu konfigurieren
		WithProxy(cfg)(c)
	}

	return nil
}
