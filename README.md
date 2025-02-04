# Stealthy - Go Library for persistent identity Web Requests

A lightweight Go library for making stealthy HTTP requests with persistent identities, proxy support, and serialization capabilities.

## Features

- 🛡️ **Persistent Session Identity**  
  Consistent User-Agent, Headers, and Session-ID across requests
- 🔄 **Proxy Rotation Support**  
  Built-in support for SmartProxy
- 📦 **State Serialization**  
  Save/restore client state to Base64 strings
- ⏱️ **Request Spoofing**  
  Auto-header generation with device fingerprinting

## Usage

```go
package main

import (
  "github.com/m0ttii/stealthy"
  "net/http"
)

func main() {
  // Create client with SmartProxy
  client, _ := stealthy.NewStealthClient(
    stealthclient.WithProxy(stealthclient.ProxyConfig{
      Host:         "gate.smartproxy.com",
      User:         "USER",
      ZonePassword: "PASSWORD",
      Port:         10001,
    }),
  )

  // Make request
  req, _ := http.NewRequest("GET", "https://api.example.com/data", nil)
  resp, _ := client.Do(req)
}
```

### State Management

```go
// Serialize client
data, _ := client.Serialize() 
// -> "eyJVc2VyQWdlbnQiOiJNb3ppbGxhLzUuMC4u..."

// Restore client
restoredClient, _ := stealthclient.DeserializeClient(data)

// Rotate SmartProxy session ID
restoredClient.RotateSmartProxySession("new_session_id")
```

 
