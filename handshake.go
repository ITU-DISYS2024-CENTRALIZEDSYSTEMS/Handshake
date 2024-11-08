package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand/v2"
	"net"
	"strconv"
	"strings"
	"time"
)

var sequenceNumberServer int = 0
var sequenceNumberClient int = 0

func main() {
	go server()
	time.Sleep(1 * time.Second)
	go client()

	select {} // Prevent Main from instantly exiting
}

// Client side program

func client() {
	conn := clientHandshake("127.0.0.1:8080")

	defer conn.Close()
}

func clientHandshake(ip string) net.Conn {
	conn, err := net.Dial("udp", ip)
	if err != nil {
		fmt.Println("Error happened on dialup: ", err)
		return nil
	}

	clientSendMessage(conn, "SYN")
	response := clientRecieveMessage(conn)

	if response != "SYN ACK" {
		return nil
	}
	clientSendMessage(conn, "ACK")

	return conn
}

func clientSendMessage(conn net.Conn, message string) {
	if message == "SYN" {
		var seq = 0
		if sequenceNumberClient == 0 {
			seq = rand.Int()
		} else {
			seq = sequenceNumberClient
		}

		fmt.Println("Client sending:", "SYN", "| seq: ", seq)
		fmt.Fprintf(conn, "%v;%d", message, seq)
	} else if message == "ACK" {
		sequenceNumberServer++
		fmt.Println("Client sending:", "ACK", "| seq:", sequenceNumberClient, "| ack:", sequenceNumberServer)
		fmt.Fprintf(conn, "%v;%s", message, strconv.Itoa(sequenceNumberClient)+","+strconv.Itoa(sequenceNumberServer))
	}
}

func clientRecieveMessage(conn net.Conn) string {
	p := make([]byte, 2048)

	_, err := bufio.NewReader(conn).Read(p)
	if err != nil {
		return "Error :("
	}

	var response = string(bytes.Trim(p, "\x00")[:])
	var responseArr = strings.Split(response, ";") // Splitting message and sequence number

	var seq, _ = strconv.Atoi(strings.Split(responseArr[1], ",")[0])
	var ack, _ = strconv.Atoi(strings.Split(responseArr[1], ",")[1])

	sequenceNumberServer = seq
	sequenceNumberClient = ack

	fmt.Println("Client received:", responseArr[0], "| seq: ", seq, "| ack: ", ack)
	return responseArr[0]
}

// Server side program

func server() {
	addr := net.UDPAddr{
		Port: 8080,
		IP:   net.ParseIP("127.0.0.1"),
	}

	serverHandshake(&addr)
}

func serverHandshake(addr *net.UDPAddr) (*net.UDPConn, *net.UDPAddr) {
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
	}

	response, remoteaddr := serverRecieveMessage(conn)
	if response != "SYN" {
		panic(response)
	}
	serverSendMessage(conn, remoteaddr, "SYN ACK")

	response, remoteaddr = serverRecieveMessage(conn)
	if response != "ACK" {
		panic(response)
	}

	return conn, remoteaddr
}

func serverSendMessage(conn *net.UDPConn, addr *net.UDPAddr, message string) {
	var seq = 0

	if sequenceNumberServer == 0 {
		seq = rand.Int()
	} else {
		seq = sequenceNumberServer
	}

	fmt.Println("Server sending:", message, "| seq:", seq, "| ack:", sequenceNumberClient+1)
	_, err := conn.WriteToUDP([]byte(message+";"+strconv.Itoa(seq)+","+strconv.Itoa(sequenceNumberClient+1)), addr)
	if err != nil {
		fmt.Printf("Couldn't send response %v", err)
	}
}

func serverRecieveMessage(conn *net.UDPConn) (string, *net.UDPAddr) {
	p := make([]byte, 2048)

	_, remoteaddr, err := conn.ReadFromUDP(p)
	if err != nil {
		fmt.Printf("Server receiving error: %v\n", err)
	}

	var response = string(bytes.Trim(p, "\x00")[:])

	var responseArr = strings.Split(response, ";") // Splitting message and sequence number

	if responseArr[0] == "SYN" {
		sequenceNumberClient, _ = strconv.Atoi(responseArr[1])
		fmt.Println("Server received:", "SYN", "| seq:", responseArr[1])
	} else if responseArr[0] == "ACK" {
		var seq, _ = strconv.Atoi(strings.Split(responseArr[1], ",")[0])
		var ack, _ = strconv.Atoi(strings.Split(responseArr[1], ",")[1])

		sequenceNumberClient = seq
		sequenceNumberServer = ack

		fmt.Println("Server received:", responseArr[0], "| seq:", seq, "| ack:", ack)
	}

	return responseArr[0], remoteaddr
}
