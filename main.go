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
const clientNum=10
var (
	
	clients = make(map[net.Conn]string, clientNum)
	channel = make(chan string)
	mu      sync.Mutex
	history = make([]string, 0, 1000)
)



func main() {
	var port string
	if len(os.Args) < 2 {
		port = ":8989"
	} else if len(os.Args) == 2 {
		port = ":" + os.Args[1]
	} else{
		fmt.Println("Usage: go run main.go [port]")
		return
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
	
	if len(clients)<clientNum{
		mu.Lock()
		clients[conn] = cName
		
		mu.Unlock()
	}else{
		fmt.Print("cannot add u to the server, try again later")
		conn.Close()
	}
	

	joinMsg := fmt.Sprintf("%s has joined our chat...\n", cName)
		channel <- joinMsg
	r:=strings.Join(history, " ")
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
		history=append(history, echo)
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