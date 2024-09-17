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
}

func NewNode(id string, port int) *Node {
	return &Node{
		ID:     id,
		Port:   port,
		shards: make(map[int][]byte),
	}
}

func (n *Node) Start() {
	fmt.Printf("Node %s started on port %d\n", n.ID, n.Port)
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
