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

	var cName string
	for {
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
		cName = strings.TrimSpace(clientName)

		invalid := false
		for i := 0; i < len(cName); i++ {
			if cName[i] > 127 {
				invalid = true
				break
			}
		}
		if invalid {
			conn.Write([]byte("Only ASCII characters are allowed in names.\n"))
			continue
		}

		server.mu.Lock()

		if cName == "" {
			for i := 1; i <= 10; i++ {
				Anonymous := fmt.Sprintf("Anonymous%d", i)
				exists := false
				for _, name := range server.clients {
					if name == Anonymous {
						exists = true
						break
					}
				}
				if !exists {
					cName = Anonymous
					break
				}
			}
		}

		taken := false
		for _, name := range server.clients {
			if name == cName {
				taken = true
				break
			}
		}
		server.mu.Unlock()

		if taken {
			conn.Write([]byte("This name is already taken. Please choose a different name.\n"))
			continue
		}

		break
	}

	server.mu.Lock()
	if len(server.clients) >= 2 {
		server.mu.Unlock()
		conn.Write([]byte("Server is full\n"))
		conn.Close()
		return
	}
	server.clients[conn] = cName
	server.mu.Unlock()

	joinMsg := fmt.Sprintf("%s has joined our chat...\n", cName)
	server.channel <- joinMsg

	r := strings.Join(server.history, " ")
	conn.Write([]byte(r))

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				leaveMsg := fmt.Sprintf("%s has left our chat...\n", cName)
				server.mu.Lock()
				delete(server.clients, conn)
				server.mu.Unlock()
				server.channel <- leaveMsg
				conn.Close()
				return
			}
			fmt.Println("Error reading from client:", err)
			return
		}

		emptyMsg := strings.TrimSpace(msg)
		if emptyMsg == "" {
			conn.Write([]byte("Empty messages are not allowed. Please write something meaningful.\n"))
			continue
		}
		invalid := false
		for i := 0; i < len(emptyMsg); i++ {
			if emptyMsg[i] > 127 {
				invalid = true
				break
			}
		}
		if invalid {
			conn.Write([]byte("Only ASCII characters are allowed in messages.\n"))
			continue
		}

		echo := fmt.Sprintf("[%s][%s]: %s", time.Now().Format("2006-01-02 15:04:05"), cName, msg)
		server.channel <- echo

		server.mu.Lock()
		server.history = append(server.history, echo)
		server.mu.Unlock()
	}
}
