package main

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"bbcli/pkg/server"
	"bbcli/pkg/types"
)

func main() {
	mcpServer := server.NewMCPServer()

	// Read from stdin and write to stdout (STDIO transport)
	decoder := json.NewDecoder(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)

	for {
		var request types.MCPRequest
		if err := decoder.Decode(&request); err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("Error decoding request: %v", err)
			continue
		}

		response := mcpServer.HandleRequest(&request)
		// Only send response if it's not nil (notifications don't get responses)
		if response != nil {
			if err := encoder.Encode(response); err != nil {
				log.Printf("Error encoding response: %v", err)
			}
		}
	}
}
