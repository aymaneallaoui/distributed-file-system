package filesystem

import (
	"bytes"
	"os"
	"testing"
)

func TestShardFileAndCombineShards(t *testing.T) {

	testData := []byte("This is a test for sharding and combining.")

	filePath := "test_file.txt"
	err := createTempFile(filePath, testData)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer deleteTempFile(filePath)

	shardSize := 10
	shards, err := ShardFile(filePath, shardSize)
	if err != nil {
		t.Fatalf("Failed to shard file: %v", err)
	}

	expectedShards := (len(testData) + shardSize - 1) / shardSize
	if len(shards) != expectedShards {
		t.Fatalf("Expected %d shards, got %d", expectedShards, len(shards))
	}

	outputFilePath := "reassembled_file.txt"
	err = CombineShards(shards, outputFilePath)
	if err != nil {
		t.Fatalf("Failed to combine shards: %v", err)
	}
	defer deleteTempFile(outputFilePath)

	reassembledData, err := readFile(outputFilePath)
	if err != nil {
		t.Fatalf("Failed to read reassembled file: %v", err)
	}

	if !bytes.Equal(testData, reassembledData) {
		t.Fatalf("Reassembled file content does not match original")
	}
}

func createTempFile(filePath string, content []byte) error {
	return os.WriteFile(filePath, content, 0644)
}

func deleteTempFile(filePath string) {
	os.Remove(filePath)
}

func readFile(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}
