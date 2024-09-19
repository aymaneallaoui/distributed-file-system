// cmd/main.go
package main

import (
	"bufio"
	"distributed-file-system/pkg/filesystem"
	"distributed-file-system/pkg/storage"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	// Step 1: Initialize storage nodes
	nodes := []*storage.Node{
		storage.NewNode("node-1", 8081),
		storage.NewNode("node-2", 8082),
		storage.NewNode("node-3", 8083),
	}

	// Start nodes
	for _, node := range nodes {
		node.Start()
	}

	// Step 2: Initialize the distributed file system with replication factor 2
	replicationFactor := 2
	fs := filesystem.NewFileSystem(nodes, replicationFactor)

	// Step 3: Upload a file
	filePath := "test_file.bin" // Path to the file to be uploaded
	err := fs.UploadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to upload file: %v", err)
	}
	fmt.Printf("File '%s' uploaded successfully\n", filePath)

	// Step 4: List files in the system
	files := fs.ListFiles()
	fmt.Println("Files stored in the distributed system:")
	for _, file := range files {
		fmt.Printf("File: %s, Shards: %v\n", file.FileName, file.ShardIDs)
	}

	// Step 5: Ask user if they want to download the file
	if len(files) > 0 {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Do you want to download the file? (yes/y to confirm): ")
		userInput, _ := reader.ReadString('\n')
		userInput = strings.TrimSpace(strings.ToLower(userInput))

		if userInput == "yes" || userInput == "y" {
			// Assuming we download the first file in the list
			fileMeta := files[0]
			outputFile := "downloaded_" + fileMeta.FileName

			err = fs.DownloadFile(fileMeta.ShardIDs, outputFile)
			if err != nil {
				log.Fatalf("Failed to download file: %v", err)
			}
			fmt.Printf("File '%s' downloaded successfully as '%s'\n", fileMeta.FileName, outputFile)
		} else {
			fmt.Println("Download skipped.")
		}
	}
}
