# Distributed File System

A lightweight, sharded, and replicated distributed file system built in Go. This project demonstrates how files can be split into shards, stored across multiple nodes, and replicated for redundancy. The system allows for file uploading, listing, and downloading from a distributed set of storage nodes.

## Features

- **Sharding**: Files are split into smaller chunks (shards) for distributed storage.
- **Replication**: Each shard is replicated across multiple nodes to ensure high availability and fault tolerance.
- **Node Management**: Storage nodes can be added to the system, each responsible for storing and serving shards.
- **File Uploading**: Files can be uploaded into the distributed system where they are split and stored.
- **File Downloading**: Reconstruct files from shards stored across the distributed nodes.
- **User Interaction**: Allows the user to confirm whether they want to download the file after upload.

## Project Structure

**this diagram has been made using my tool [dirscanner](https://github.com/aymaneallaoui/dirscanner)**

```
├── .env
├── .gitignore
├── README.md
├── cmd
│   └── main.go
├── config
│   └── config.go
├── go.mod
├── go.sum
├── pkg
│   ├── consensus
│   │   ├── raft.go
│   │   └── service_discovery.go
│   ├── filesystem
│   │   ├── file_operations.go
│   │   ├── replication.go
│   │   └── sharding.go
│   ├── storage
│   │   └── node.go
│   ├── transport
│   │   └── rpc_server.go
│   └── types
│       └── shard.go
└── test_file.bin
```

## How It Works

### Sharding

Files are split into chunks called shards, each with a maximum size defined in filesystem/file_operations.go. These shards are then distributed across multiple storage nodes.

### Replication

Each shard is replicated across multiple nodes to ensure redundancy. If a node fails, the file can still be reconstructed from shards stored on other nodes.

### Node Storage

Each node has its own storage where it holds multiple shards. Shards are identified by a unique ShardID, and nodes store and retrieve shards upon reques
