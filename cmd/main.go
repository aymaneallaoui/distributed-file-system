package main

import (
	"distributed-file-system/pkg/filesystem"
	"distributed-file-system/pkg/storage"
	"fmt"
	"log"
)

func main() {

	nodes := createNodes(5)
	for _, node := range nodes {
		go node.Start()
	}

	fs := filesystem.NewFileSystem(nodes, 3)

	err := fs.UploadFile("example.txt")
	if err != nil {
		log.Fatal("Failed to upload file:", err)
	}
	fmt.Println("File uploaded and sharded successfully")

	err = fs.DownloadFile([]int{0, 1, 2, 3, 4}, "example_reassembled.txt")
	if err != nil {
		log.Fatal("Failed to download file:", err)
	}
	fmt.Println("File reassembled successfully")
}

func createNodes(n int) []*storage.Node {
	nodes := make([]*storage.Node, n)
	for i := 0; i < n; i++ {
		nodes[i] = storage.NewNode(fmt.Sprintf("node-%d", i+1), 8080+i)
	}
	return nodes
}
