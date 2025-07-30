package project

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
	"sync"
)

func handleConnection(conn net.Conn, clients map[net.Conn]string,channel chan string, mu sync.Mutex, history []string) {
	reader := bufio.NewReader(conn)

	welcome, err := os.ReadFile("welcome.txt")
	if err != nil {
		fmt.Println("Reading welcome file error:", err)
		return
	}
	// conn.Write([]byte(fmt.Sprintf("Connected users: %d\n", len(clients))))
	conn.Write(welcome)

	clientName, _ := reader.ReadString('\n')
	cName := strings.TrimSpace(clientName)

	// Add to client list
	mu.Lock()
	defer mu.Unlock()
	 
	if len(clients) >= 10 {
		
		fmt.Print("cannot add u to the server, try again later")
		conn.Close()
		return
	}else{
		clients[conn] = cName
		mu.Unlock()
	}

	joinMsg := fmt.Sprintf("%s has joined our chat...\n", cName)
	channel <- joinMsg
	r := strings.Join(history, " ")
	conn.Write([]byte(r))

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				leaveMsg := fmt.Sprintf("%s has left our chat...\n", cName)
				// history=append(history, leaveMsg)
				fmt.Print(leaveMsg)
				mu.Lock()
				delete(clients, conn)
				mu.Unlock()
				channel <- leaveMsg
				conn.Close()
				return
			} else {
				fmt.Println("Error reading from client:", err)
				return
			}
		}

		echo := fmt.Sprintf("[%s][%s]: %s", time.Now().Format("2006-01-02 15:04:05"), cName, line)
		channel <- echo
		history = append(history, echo)
	}
}
