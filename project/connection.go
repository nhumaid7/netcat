package project

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

func handleConnection(conn net.Conn, server *Server) {
	reader := bufio.NewReader(conn)

	welcome, err := os.ReadFile("welcome.txt")
	if err != nil {
		fmt.Println("Reading welcome file error:", err)
		return
	}

	 _, err = conn.Write(welcome)
    if err != nil {
        fmt.Println("Error writing welcome message:", err)
        return
    }

	clientName, err := reader.ReadString('\n')
	 	if err != nil {
    fmt.Println("Error reading client name:", err)
    return
	}
	cName := strings.TrimSpace(clientName)
	if cName == "" {
		cName = "Anonymous"
	}

	// Check if client name already exists
	// server.mu.Lock()
	// if _, exists := server.clients[conn]; exists {
	// 	server.mu.Unlock()
	// 	conn.Write([]byte("Client name already exists, please choose a different name.\n"))
	// 	conn.Close()
	// 	return
	// }
	// server.mu.Unlock()

	
	server.mu.Lock()

	if len(server.clients) >= 10 {
		server.mu.Unlock()
		// conn.Write("cannot add u to the server, try again later")
		conn.Close()
		return
	}else{
		server.clients[conn] = cName
		
		server.mu.Unlock()
	}
	

	joinMsg := fmt.Sprintf("%s has joined our chat...\n", cName)
	server.channel <- joinMsg
	
	r := strings.Join(server.history, " ")
	conn.Write([]byte(r))
    

	for {

		fixed:= fmt.Sprintf("[%s][%s]: ", time.Now().Format("2006-01-02 15:04:05"), cName)
		conn.Write([]byte(fixed))

		line, err := reader.ReadString('\n')
		
		if err != nil {
			if err == io.EOF {
				leaveMsg := fmt.Sprintf("%s has left our chat...\n", cName)
				// history=append(history, leaveMsg)
				// fmt.Print(leaveMsg)
				server.mu.Lock()
				delete(server.clients, conn)
				server.mu.Unlock()
				server.channel <- leaveMsg
				conn.Close()
				return
			} else {
				
				fmt.Println("Error reading from client:", err)
				return
			}
		}
		
		echo := fmt.Sprintf("[%s][%s]: %s", time.Now().Format("2006-01-02 15:04:05"), cName, line )
	
		server.channel <- echo
		server.mu.Lock()
		server.history = append(server.history, echo)
		server.mu.Unlock()

	}
}
