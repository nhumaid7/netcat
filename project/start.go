package project

import 
(
	"net"
	"fmt"
	"sync"
)

// start initializes the server and listens for incoming connections
type Server struct {
	port string	
	clients map[net.Conn]string
	channel chan string
	mu      sync.Mutex
	history []string
}
	
func Start(port string){

	serve:= &Server{
		port : port,
		clients : make(map[net.Conn]string, 10),
		channel : make(chan string),
		mu : sync.Mutex{},
		history : make([]string, 0, 1000),
	}
	

	ln, err := net.Listen("tcp", port)
		if err != nil {
			panic(err)
		}
	defer ln.Close()

	fmt.Printf("Listening on %s\n", ln.Addr())

	go broadcaster(serve.clients, serve.channel, serve.mu)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn, serve.clients, serve.channel, serve.mu, serve.history)
	}
}