package storage

import (
	"distributed-file-system/pkg/types"
	"errors"
	"fmt"
	"sync"
)

type Node struct {
	ID     string
	Port   int
	shards map[int][]byte
	mu     sync.Mutex
	active bool // Field to track if the node is active
}

func NewNode(id string, port int) *Node {
	return &Node{
		ID:     id,
		Port:   port,
		shards: make(map[int][]byte),
		active: false, // Initialize as inactive
	}
}

func (n *Node) Start() {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.active = true // Mark the node as active when started
	fmt.Printf("Node %s started on port %d\n", n.ID, n.Port)
}

func (n *Node) IsActive() bool {
	n.mu.Lock()
	defer n.mu.Unlock()

	return n.active // Return the node's active status
}

func (n *Node) StoreShard(shard types.Shard) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.shards[shard.ID] = shard.Content
	fmt.Printf("Shard %d stored on node %s\n", shard.ID, n.ID)
	return nil
}

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
