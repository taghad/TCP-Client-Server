package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type ghazale struct {
	name string
}
var mapConn map[string][]*net.Conn

func main() {
	mapConn = make(map[string][]*net.Conn)
	listener, err := net.Listen("tcp", "127.0.0.1:8081")
	if err != nil {
		log.Fatal("tcp server listener error:", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("tcp server accept error", err)
		}
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	bufferBytes, err := bufio.NewReader(conn).ReadBytes('\n')

	if err != nil {
		log.Println("client left..")
		conn.Close()
		return
	}

	message := string(bufferBytes)
	clientAddr := conn.RemoteAddr().String()
	response := fmt.Sprintf(message + " from " + clientAddr + "\n")
	messageParam := strings.Split(response," ")
	headerLen:= len(messageParam[0]) + len(messageParam[1]) + 2

	if messageParam[0] == "join" {
		join(messageParam[1],conn)
	}

	if messageParam[0] == "send" {
		send(messageParam[1], response[headerLen:])
	}

	handleConnection(conn)
}

func join(groupName string,conn net.Conn)  {

	groupConn := mapConn[groupName]
	if groupConn == nil {
		_, _ = conn.Write([]byte("group doesn't exist, you created it!"))
		tmp := make([]*net.Conn,0)
		tmp = append(tmp, &conn)
		mapConn[groupName] = tmp
	} else {
		if !Contains(mapConn[groupName],&conn) {
			mapConn[groupName] = append(mapConn[groupName],&conn )
			fmt.Println(len(mapConn[groupName]))

			_, _ = conn.Write([]byte("you added in " + groupName))
		}
	}

}

func send(groupName string, message string)  {

	groupConn := mapConn[groupName]
	fmt.Println(len(mapConn[groupName]))
	for _, p := range groupConn {
		_, _ = (*p).Write([]byte((*p).LocalAddr().String() + ": " + message))
	}
}
func Contains(a []*net.Conn, x *net.Conn) bool {
	if a == nil {
		return false
	}
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}