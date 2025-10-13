# MailStream

MailStream is a Go library that provides an efficient interface to interact with IMAP servers. It enables the streaming of mail updates in real-time. This library is a wrapper around the [github.com/emersion/go-imap/v2](https://github.com/emersion/go-imap) library, which provides the IMAP client implementation.

## Features

- Connect to IMAP servers using secure protocols.
- Subscribe to real-time mail updates.
- Efficient handling of concurrent mail streams.
- Fetching Mails using IMAP IDLE. - No steady polling.
- Mails parsed into structured objects seperating the html and text parts.
- SOCKS5 Proxy support.

## Installation

To install MailStream, use the `go get` command:

```bash
go get github.com/nukilabs/mailstream
```

This will retrieve the library from GitHub and install it in your Go workspace.

## Usage

Below are two basic examples of how to use the MailStream library: one for default usage and another for concurrent processing. Please look at the examples directory for more detailed explanations.

### Default Usage

The default example demonstrates how to set up a simple mail listening service that logs all incoming mails.

**File: `examples/default/main.go`**

```go
package main

import (
    "context"
    "fmt"
    "log"
    "github.com/nukilabs/mailstream"
)

func main() {
    config := mailstream.Config{
        Host:     "imap.example.com",
        Port:     993,
        Email:    "mymail@example.com",
        Password: "password1234",
    }
    client, err := mailstream.New(config)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

	done := client.WaitForUpdates(context.Background())
    listener := client.Subscribe()
    for {
        select {
        case mail := <-listener:
            fmt.Println(mail.Subject)
        case err := <-done:
            log.Fatal(err)
        }
    }
}
```

### Concurrent Usage

The concurrent example demonstrates how to handle multiple subscribers that listen for new mails concurrently.

**File: `examples/concurrent/main.go`**

```go
package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/nukilabs/mailstream"
)

func main() {
	config := mailstream.Config{
		Host:     "imap.example.com",
		Port:     993,
		Email:    "mymail@example.com",
		Password: "password1234",
	}
	client, err := mailstream.New(config)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		listener := client.Subscribe()
		defer client.Unsubscribe(listener)
		mail := <-listener
		fmt.Printf("Task 1 - Received mail: %s\n", mail.Subject)
	}()

	go func() {
		defer wg.Done()
		listener := client.Subscribe()
		defer client.Unsubscribe(listener)
		mail := <-listener
		fmt.Printf("Task 2 - Received mail: %s\n", mail.Subject)
	}()

	ctx, cancel := context.WithCancel(context.Background())
	client.WaitForUpdates(ctx)

	wg.Wait()
	cancel()
}

```

### Proxy Support

MailStream supports connecting to IMAP servers through proxy servers using SOCKS5.

**File: `examples/proxy/main.go`**
```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/nukilabs/mailstream"
)

func main() {
	config := mailstream.Config{
    	Host:     "imap.example.com",
    	Port:     993,
    	Email:    "mymail@example.com",
    	Password: "password1234",
		ProxyURL: &url.URL{
			Scheme: "socks5h",
			Host: "localhost:1080",
			User: url.UserPassword("username", "password"),
		},
	}
	client, err := mailstream.New(config)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	done := client.WaitForUpdates(context.Background())
    listener := client.Subscribe()
    for {
        select {
        case mail := <-listener:
            fmt.Println(mail.Subject)
        case err := <-done:
            log.Fatal(err)
        }
    }
}

```

### Supported Proxy Schemes

- `socks5://` - SOCKS5 proxy with local DNS resolution
- `socks5h://` - SOCKS5 proxy with remote DNS resolution

### Proxy URL Format

Proxy URLs can include authentication credentials:

```
scheme://[username:password@]host:port
```

Examples:
- `socks5://localhost:1080`
- `socks5://user:pass@proxy.example.com:1080`
```
