package filesystem

import (
	"distributed-file-system/pkg/storage"
	"distributed-file-system/pkg/types"
	"fmt"
)

// ShardSize defines the size of each shard
const ShardSize = 1024 * 1024

// FileMetadata tracks metadata about each uploaded file
type FileMetadata struct {
	FileName string // the name of the file
	ShardIDs []int  // the shard IDs that make up the file
}

// FileSystem represents the distributed file system
type FileSystem struct {
	nodes             []*storage.Node // the nodes in the distributed file system
	replicationFactor int
	files             map[string]FileMetadata // New field to track files
}

// NewFileSystem creates a new distributed file system
func NewFileSystem(nodes []*storage.Node, replicationFactor int) *FileSystem {
	return &FileSystem{
		nodes:             nodes,
		replicationFactor: replicationFactor,
		files:             make(map[string]FileMetadata), // Initialize file tracking map
	}
}

// UploadFile shards a file and distributes its shards across the nodes
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

// DownloadFile reassembles the shards into the original file
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

// ListFiles returns a list of files currently stored in the distributed file system
func (fs *FileSystem) ListFiles() []FileMetadata {
	files := []FileMetadata{}
	for _, metadata := range fs.files {
		files = append(files, metadata)
	}
	return files
}

// GetFileMetadata returns the metadata for a specific file by filename
func (fs *FileSystem) GetFileMetadata(fileName string) (FileMetadata, bool) {
	fileMeta, exists := fs.files[fileName]
	return fileMeta, exists
}
