package filesystem

import (
	"distributed-file-system/pkg/types"
	"io"
	"log"
	"os"
)

func ShardFile(filePath string, shardSize int) ([]types.Shard, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var shards []types.Shard
	buffer := make([]byte, shardSize)
	shardID := 0

	for {
		bytesRead, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if bytesRead == 0 {
			break
		}

		shardContent := make([]byte, bytesRead)
		copy(shardContent, buffer[:bytesRead])

		shards = append(shards, types.Shard{
			ID:      shardID,
			Content: shardContent,
		})
		log.Printf("Shard %d: %v", shardID, shardContent)
		shardID++
	}

	return shards, nil
}

func CombineShards(shards []types.Shard, outputFile string) error {
	output, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer output.Close()

	for _, shard := range shards {
		_, err := output.Write(shard.Content)
		if err != nil {
			return err
		}
		log.Printf("Writing shard %d: %v", shard.ID, shard.Content)
	}
	return nil
}
