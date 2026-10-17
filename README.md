## TinyP2P

TinyP2P is a simple P2P networking library for Go.it is a simple P2P networking library for Go.

## Install

```bash
go get github.com/tinyp2p/tinyp2p
```

## Usage

```go
package main

import (
	"fmt"
	"github.com/tinyp2p/tinyp2p"
)

func main() {
	// Create a new P2P network
	p2p := tinyp2p.New()

	// Listen for incoming connections
	p2p.Listen("127.0.0.1:1234")

	// Handle incoming connections
	p2p.OnConnect(func(conn *tinyp2p.Conn) {
		fmt.Println("Connected to", conn.RemoteAddr())
	})

	// Start the network
	p2p.Start()
}
```

## License

MIT

