package project

import (
	"fmt"

)

func broadcaster(server *Server) {
	for msg := range server.channel {
		server.mu.Lock()
		for conn := range server.clients {
			_, err := conn.Write([]byte(msg))
			if err != nil {
				fmt.Println("Broadcast error:", err)
				conn.Close()
				delete(server.clients, conn)
			}
		}
		server.mu.Unlock()
	}
}
