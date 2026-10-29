## PeerCache

**A Peer-to-Peer Content Distribution System for Software Packages**

PeerCache is a learning and development project built with Go to explore peer-to-peer (P2P) networking and reliable software package distribution. It aims to allow peers to exchange software packages directly, verify downloaded content, cache packages, and redistribute them to other peers.

The initial version runs multiple peer processes on a single machine. The long-term goal is to develop a system that can distribute packages efficiently and reliably without requiring one central server to serve every download.

## Project Objectives

- Build a TCP-based communication system between peers.
- Transfer software packages between nodes.
- Detect incomplete file transfers.
- Verify package integrity using SHA-256 hashes.
- Allow peers to cache and redistribute received packages.
- Explore fault tolerance, concurrency, and performance in distributed systems.

## Current Features

- **TCP communication:** Nodes communicate using Go's `net` package.
- **Package transfer:** Node A sends a software archive to Node B.
- **Size-based framing:** An 8-byte header communicates the package size before the file contents.
- **Incomplete-transfer handling:** The receiver detects failed transfers and removes partial files.
- **SHA-256 hashing:** The receiver calculates a hash to help verify the integrity of the downloaded package.
- **Concurrent connections:** Node A handles connections in separate goroutines.

> **Integrity note:** Calculating a SHA-256 hash is not enough by itself to authenticate a package. The receiver must compare it with a trusted expected hash to verify that the downloaded bytes match the expected content.

## Technology Stack

- **Language:** Go
- **Networking:** TCP
- **Concurrency:** Goroutines
- **File handling:** Go standard library (`os`, `io`, `path/filepath`)
- **Binary encoding:** `encoding/binary`
- **Integrity checking:** `crypto/sha256`

## Project Structure

```text
tinyp2p/
├── cmd/
│   ├── node/
│   │   └── main.go
│   ├── node-a/
│   │   └── main.go
│   ├── node-b/
│   │   └── main.go
│   └── file-reader/
│       └── main.go
├── packages/
│   └── demo-app/
│       ├── VERSION
│       └── README.txt
├── dist/
│   └── demo-app-v1.0.0.tar.gz
├── received/
└── README.md
```

The `received/` directory is created when Node B downloads the package. Its contents may vary depending on your tests.

## Prerequisites

- Go installed on your system.
- A terminal.
- Git (optional, for version control).

Check your Go installation:

```bash
go version
```

## Getting Started

### 1. Clone the repository

If the project is hosted on GitHub:

```bash
git clone <your-repository-url>
cd tinyp2p
```

Alternatively, navigate to your existing project directory:

```bash
cd ~/tinyp2p
```

### 2. Prepare the package

Ensure the demo application directory exists and contains the package files.

Create the archive:

```bash
mkdir -p dist
tar -czf dist/demo-app-v1.0.0.tar.gz -C packages demo-app
```

### 3. Start Node A

Open the first terminal in the project root:

```bash
go run ./cmd/node-a
```

Node A listens for incoming TCP connections on port `8001` and serves the package.

### 4. Start Node B

Open a second terminal in the same project root:

```bash
go run ./cmd/node-b
```

Node B connects to Node A, receives the package, saves it under `received/`, and calculates its SHA-256 hash.

### 5. Compare the package hashes

Run:

```bash
sha256sum dist/demo-app-v1.0.0.tar.gz received/demo-app-v1.0.0.tar.gz
```

Matching hashes indicate that both files contain the same bytes.

## How the Current Transfer Works

1. Node B establishes a TCP connection to Node A.
2. Node A opens the package archive.
3. Node A sends the package size as an 8-byte, big-endian unsigned integer.
4. Node B reads the size header.
5. Node B checks the size against its configured maximum.
6. Node A sends the package contents.
7. Node B writes the received bytes to disk.
8. Node B removes the partial file if the transfer fails.
9. Node B calculates the SHA-256 hash of the received package.

## Development Roadmap

### Phase 1 — Basic Networking
- [x] Start a TCP listener.
- [x] Accept peer connections.
- [x] Handle connections using goroutines.

### Phase 2 — File Transfer
- [x] Send package size before file contents.
- [x] Transfer a software archive.
- [x] Save the received archive.
- [x] Detect incomplete transfers and remove partial files.
- [x] Calculate the received package's SHA-256 hash.
- [ ] Automate verification against a trusted expected hash.

### Phase 3 — Request Protocol
- [ ] Allow peers to request a specific package.
- [ ] Return explicit success or error responses.
- [ ] Validate package identifiers and requests.

### Phase 4 — Peer Caching and Distribution
- [ ] Cache successfully verified packages.
- [ ] Allow Node B to serve packages to Node C.
- [ ] Support transfers across multiple peers.

### Phase 5 — Reliability and Performance
- [ ] Add transfer timeouts and connection limits.
- [ ] Test interrupted transfers and peer failures.
- [ ] Measure transfer speed and resource usage.
- [ ] Test with three or four local nodes.

### Phase 6 — Peer Discovery
- [ ] Introduce a bootstrap or discovery mechanism.
- [ ] Allow peers to find available package sources.
- [ ] Keep package data transfers peer-to-peer.

## Design Principles

- **Decentralized transfer:** Peers should exchange package data directly.
- **Integrity checking:** Verify downloaded bytes against a trusted expected hash.
- **Failure handling:** Incomplete downloads must not be treated as successful packages.
- **Modularity:** Keep networking, transfer protocols, storage, and discovery logically separated.
- **Incremental development:** Test each capability before adding the next one.

## Current Limitations

- The current Node A and Node B configuration uses `localhost`, so it is designed for testing on the same machine.
- Node A currently serves a fixed package.
- Package requests, caching, and peer discovery are still under development.
- Hash calculation is implemented, but trusted hash verification is not yet automatic.
- The project has not yet been benchmarked under substantial network load.

## Contributing

PeerCache is an educational project focused on learning Go, networking, and distributed systems. Contributions, suggestions, and improvements are welcome.

## License

No license has been selected yet. Add a license before distributing the project for reuse by others.

