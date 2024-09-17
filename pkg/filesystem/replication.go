package filesystem

import (
	"distributed-file-system/pkg/storage"
	"distributed-file-system/pkg/types"
	"errors"
	"time"
)

const maxRetries = 3
const retryInterval = 2 * time.Second

func ReplicateShard(shard types.Shard, nodes []*storage.Node, replicationFactor int) error {
	if len(nodes) < replicationFactor {
		return errors.New("not enough nodes for replication")
	}

	for i := 0; i < replicationFactor; i++ {
		retries := 0
		for retries < maxRetries {
			err := nodes[i].StoreShard(shard)
			if err != nil {
				retries++
				time.Sleep(retryInterval)
				continue
			}
			break
		}
		if retries == maxRetries {
			return errors.New("failed to replicate shard after retries")
		}
	}

	return nil
}

func RetrieveReplicatedShard(shardID int, nodes []*storage.Node) (types.Shard, error) {
	for _, node := range nodes {
		shard, err := node.FetchShard(shardID)
		if err == nil {
			return shard, nil
		}
	}
	return types.Shard{}, errors.New("shard not found on any node")
}
