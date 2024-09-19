package storage

import (
	"distributed-file-system/pkg/types"
	"errors"
	"fmt"
	"sync"
)

type Node struct {
	ID     string         // Field to store the node's ID
	Port   int            // Field to store the node's port
	shards map[int][]byte // Field to store the node's shards
	mu     sync.Mutex     // Mutex to protect the node's fields
	active bool           // Field to track if the node is active
}

// NewNode function to create a new node
func NewNode(id string, port int) *Node {
	return &Node{
		ID:     id,
		Port:   port,
		shards: make(map[int][]byte),
		active: false, // Initialize as inactive
	}
}

// Start method to start the node
func (n *Node) Start() {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.active = true // Mark the node as active when started
	fmt.Printf("Node %s started on port %d\n", n.ID, n.Port)
}

// a method tp track if the node is active
func (n *Node) IsActive() bool {
	n.mu.Lock()
	defer n.mu.Unlock()

	return n.active
}

// a method to store a shards
func (n *Node) StoreShard(shard types.Shard) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.shards[shard.ID] = shard.Content
	fmt.Printf("Shard %d stored on node %s\n", shard.ID, n.ID)
	return nil
}

// a method to fetch a shard from the id of the shard
func (n *Node) FetchShard(shardID int) (types.Shard, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	content, exists := n.shards[shardID]
	if !exists {
		return types.Shard{}, errors.New("shard not found")
	}

	return types.Shard{
		ID:      shardID,
		Content: content,
	}, nil
}
