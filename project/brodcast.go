package project

import (
	"fmt"
	"net"
	"sync"

)

func broadcaster(clients map[net.Conn]string, channel chan string, mu sync.Mutex) {
	for msg := range channel {
		mu.Lock()
		for conn := range clients {
			_, err := conn.Write([]byte(msg))
			if err != nil {
				fmt.Println("Broadcast error:", err)
				conn.Close()
				delete(clients, conn)
			}
		}
		mu.Unlock()
	}
}
