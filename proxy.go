package main

import (
	"fmt"
	"net/url"
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
