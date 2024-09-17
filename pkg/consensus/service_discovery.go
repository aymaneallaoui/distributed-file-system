package consensus

import (
	"context"
	"log"
	"time"

	"go.etcd.io/etcd/clientv3"
)

func DiscoverNodes(client *clientv3.Client) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.Get(ctx, "nodes/", clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	var nodes []string
	for _, ev := range resp.Kvs {
		nodes = append(nodes, string(ev.Value))
	}
	return nodes, nil
}

func RegisterNode(client *clientv3.Client, nodeID string, address string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Put(ctx, "nodes/"+nodeID, address)
	if err != nil {
		log.Println("Failed to register node:", err)
		return err
	}

	log.Printf("Node %s registered with address %s\n", nodeID, address)
	return nil
}
