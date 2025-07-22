package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

type Client struct {
	conn net.Conn
	name string
}

var (
	clients   = make(map[net.Conn]*Client)
	mutex     sync.Mutex
	broadcast = make(chan string)
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	nameLine, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	name := strings.TrimSpace(nameLine)

	client := &Client{conn: conn, name: name}

	mutex.Lock()
	clients[conn] = client
	mutex.Unlock()

	broadcast <- fmt.Sprintf("%s has joined the chat", name)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		message = strings.TrimSpace(message)
		if message == "" {
			continue
		}
		broadcast <- fmt.Sprintf("[%s]: %s", name, message)
	}

	mutex.Lock()
	delete(clients, conn)
	mutex.Unlock()

	broadcast <- fmt.Sprintf("%s has left the chat", name)
}

func handleBroadcast() {
	for msg := range broadcast {
		mutex.Lock()
		for conn := range clients {
			fmt.Fprintln(conn, msg)
		}
		mutex.Unlock()
		fmt.Println(msg)
	}
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	fmt.Println("Server listening on :8080")

	go handleBroadcast()

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn)
	}
}
