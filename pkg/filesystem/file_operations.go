package filesystem

import (
	"distributed-file-system/pkg/storage"
	"distributed-file-system/pkg/types"
	"fmt"
)

const ShardSize = 1024 * 1024

type FileMetadata struct {
	FileName string
	ShardIDs []int
}

type FileSystem struct {
	nodes             []*storage.Node
	replicationFactor int
	files             map[string]FileMetadata
}

func NewFileSystem(nodes []*storage.Node, replicationFactor int) *FileSystem {
	return &FileSystem{
		nodes:             nodes,
		replicationFactor: replicationFactor,
		files:             make(map[string]FileMetadata),
	}
}

func (fs *FileSystem) UploadFile(filePath string) error {
	shards, err := ShardFile(filePath, ShardSize)
	if err != nil {
		return err
	}

	var shardIDs []int
	for _, shard := range shards {
		err := ReplicateShard(shard, fs.nodes, fs.replicationFactor)
		if err != nil {
			return err
		}
		fmt.Printf("Shard %d replicated across nodes\n", shard.ID)
		shardIDs = append(shardIDs, shard.ID)
	}

	fs.files[filePath] = FileMetadata{
		FileName: filePath,
		ShardIDs: shardIDs,
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

func (fs *FileSystem) ListFiles() []FileMetadata {
	files := []FileMetadata{}
	for _, metadata := range fs.files {
		files = append(files, metadata)
	}
	return files
}

func (fs *FileSystem) GetFileMetadata(fileName string) (FileMetadata, bool) {
	fileMeta, exists := fs.files[fileName]
	return fileMeta, exists
}
