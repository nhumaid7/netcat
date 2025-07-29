package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	clients = make(map[net.Conn]string)
	channel = make(chan string, 10)
	mu      sync.Mutex
	history  = make([]string, 0, 100) 
)

func main() {
	var port string
	if len(os.Args) < 2 {
		port = ":8989"
	} else {
		port = ":" + os.Args[1]
	}

	ln, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	fmt.Printf("Listening on %s\n", ln.Addr())

	go broadcaster()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)

	fmt.Println([]byte("Welcome to TCP-Chat!
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
 |    `.       | `' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     `-'       `--'
[ENTER YOUR NAME]: "))

	clientName, _ := reader.ReadString('\n')
	cName := strings.TrimSpace(clientName)

	// Add to client list
	mu.Lock()
	clients[conn] = cName
	mu.Unlock()

	joinMsg := fmt.Sprintf("%s has joined our chat...\n", cName)
	channel <- joinMsg

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				leaveMsg := fmt.Sprintf("%s has left our chat...\n", cName)
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
		// Store message in history
		mu.Lock()
		if len(history) >= 100 {
			history = history[1:] // Remove oldest message
			history = append(history, echo)
		} else {
			history = append(history, echo)
		}
		mu.Unlock()

		channel <- history
	}
}

func broadcaster() {
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
