package filesystem

import (
	"distributed-file-system/pkg/storage"
	"distributed-file-system/pkg/types"
	"fmt"
	"testing"
)

func TestReplicateShardAndRetrieve(t *testing.T) {

	nodes := createTestNodes(3)

	shard := types.Shard{
		ID:      1,
		Content: []byte("This is a test shard for replication"),
	}

	replicationFactor := 3
	err := ReplicateShard(shard, nodes, replicationFactor)
	if err != nil {
		t.Fatalf("Failed to replicate shard: %v", err)
	}

	retrievedShard, err := RetrieveReplicatedShard(shard.ID, nodes)
	if err != nil {
		t.Fatalf("Failed to retrieve replicated shard: %v", err)
	}

	if string(retrievedShard.Content) != string(shard.Content) {
		t.Fatalf("Shard content mismatch: expected %s, got %s", shard.Content, retrievedShard.Content)
	}
}

func createTestNodes(n int) []*storage.Node {
	nodes := make([]*storage.Node, n)
	for i := 0; i < n; i++ {
		nodes[i] = storage.NewNode(fmt.Sprintf("node-%d", i+1), 8080+i)
	}
	return nodes
}
