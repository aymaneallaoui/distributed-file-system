package storage

import (
	"distributed-file-system/pkg/types"
	"testing"
)

func TestNodeStoreAndFetchShard(t *testing.T) {

	node := NewNode("test-node", 8080)

	shard := types.Shard{
		ID:      1,
		Content: []byte("This is a test shard"),
	}

	err := node.StoreShard(shard)
	if err != nil {
		t.Fatalf("Failed to store shard: %v", err)
	}

	retrievedShard, err := node.FetchShard(shard.ID)
	if err != nil {
		t.Fatalf("Failed to fetch shard: %v", err)
	}

	if string(retrievedShard.Content) != string(shard.Content) {
		t.Fatalf("Shard content mismatch: expected %s, got %s", shard.Content, retrievedShard.Content)
	}
}
