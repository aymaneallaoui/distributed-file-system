package consensus

import (
	"context"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Raft struct {
	nodes      []string
	leader     string
	etcdClient *clientv3.Client
}

func NewRaft(etcdEndpoints []string) (*Raft, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdEndpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	return &Raft{
		etcdClient: client,
	}, nil
}

func (r *Raft) ElectLeader() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.etcdClient.Put(ctx, "leader", "node-1")
	if err != nil {
		log.Println("Leader election failed:", err)
	}
}
