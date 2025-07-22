package main

import (
	"fmt"
	"io"
	"net"
	//"os"
	"bufio"
)

func main() {

	// var port string
	// if len(os.Args) == 1 {
	// 	fmt.Print("listinging on port : 9000")
	// } else if len(os.Args) == 2 {
	// 	port = os.Args[1]
	// 	fmt.Sprint("listinging on port :%s", port)
	// }
	ln, err := net.Listen("tcp", "localhost:9000")

	if err != nil {
		panic(err)
	}
	defer ln.Close()

	fmt.Printf("listining on %s\n", ln.Addr())

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("err accepting ", err)
			continue
		}
		conn.Close()
		io.WriteString(conn, fmt.Sprintf("hello from %s\n", conn.RemoteAddr()))

		// fmt.Println(string(bs))

		go handelConnetion(conn)
	}
}

func handelConnetion(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("failed to read from client: connection closed")
				return
			}
			return
		}

		fmt.Printf("received: %s", line)

		echo := fmt.Sprintf("echo: %s", line)

		fmt.Printf("response: %s", echo)

		_, err = conn.Write([]byte(echo))
		if err != nil {
			fmt.Println("failed to write to client:", err)
			return
		}
		fmt.Printf("sent: %s", echo)
		// Read the data from the connection

	}
}
