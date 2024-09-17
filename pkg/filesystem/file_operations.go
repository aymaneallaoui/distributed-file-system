package filesystem

import (
	"distributed-file-system/pkg/storage"
	"distributed-file-system/pkg/types"
	"fmt"
)

const ShardSize = 1024 * 1024

type FileSystem struct {
	nodes             []*storage.Node
	replicationFactor int
}

func NewFileSystem(nodes []*storage.Node, replicationFactor int) *FileSystem {
	return &FileSystem{
		nodes:             nodes,
		replicationFactor: replicationFactor,
	}
}

func (fs *FileSystem) UploadFile(filePath string) error {
	shards, err := ShardFile(filePath, ShardSize)
	if err != nil {
		return err
	}

	for _, shard := range shards {
		err := ReplicateShard(shard, fs.nodes, fs.replicationFactor)
		if err != nil {
			return err
		}
		fmt.Printf("Shard %d replicated across nodes\n", shard.ID)
	}

	return nil
}

func (fs *FileSystem) DownloadFile(shardIDs []int, outputFile string) error {
	var shards []types.Shard
	for _, shardID := range shardIDs {
		shard, err := RetrieveReplicatedShard(shardID, fs.nodes)
		if err != nil {
			return err
		}
		shards = append(shards, shard)
	}

	return CombineShards(shards, outputFile)
}
