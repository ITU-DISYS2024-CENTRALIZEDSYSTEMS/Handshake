package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand/v2"
	"strconv"
	"net"
	"strings"
	"time"
)

var sequenceNumberServer int = 0
var sequenceNumberClient int = 0

func main() {
	go server();	
	time.Sleep(1 * time.Second)
	go client();

	select {} // Prevent Main from instantly exiting
}

// Client side program

func client() {
	conn := clientHandshake("127.0.0.1:8080")

	defer conn.Close();
}

func clientHandshake(ip string) net.Conn {
	conn, err := net.Dial("udp", ip)
    if err != nil {
        fmt.Println("Error happened on dialup: ", err)
        return nil
    }

	clientSendMessage(conn, "SYN")
	response := clientRecieveMessage(conn)

	if (response != "SYN ACK") { return nil }
	clientSendMessage(conn, "ACK")

	return conn
}

func clientSendMessage(conn net.Conn, message string) {
	if sequenceNumberClient == 0 {sequenceNumberClient = rand.Int()} else {sequenceNumberClient++}
	fmt.Println("Client sending:", message, "; seq: ", sequenceNumberClient)
	fmt.Fprintf(conn, "%v;%d", message, sequenceNumberClient)
}

func clientRecieveMessage(conn net.Conn) string {
	p :=  make([]byte, 2048)

	_, err := bufio.NewReader(conn).Read(p)
    if err == nil {
		var response = string(bytes.Trim(p, "\x00")[:])
		var responseArr = strings.Split(response, ";")  // Splitting message and sequence number
		
		seq, _ := strconv.Atoi(responseArr[1])
		if (sequenceNumberClient == 0){
			temp,_ := strconv.Atoi(responseArr[1])
			sequenceNumberClient = temp
		}else if seq == sequenceNumberClient + 1{
			sequenceNumberClient++
		}else {
			fmt.Printf("%v%d\n", "Client: ", sequenceNumberClient)
			fmt.Printf("%v%v\n", "Server: ", responseArr[1])
			panic("Client recieved an incorrect sequence number")
		}
		fmt.Println("Client received:", responseArr[0], "; seq: ", responseArr[1])
		return responseArr[0]
	} else {
		return "Error :("
    }
}

// Server side program

func server() {
    addr := net.UDPAddr{
        Port: 8080,
        IP: net.ParseIP("127.0.0.1"),
    }

	serverHandshake(&addr)
}

func serverHandshake(addr *net.UDPAddr) (*net.UDPConn, *net.UDPAddr) {
	conn, err := net.ListenUDP("udp", addr)
    if err != nil {
        fmt.Printf("Some error %v\n", err)
    }

	response, remoteaddr := serverRecieveMessage(conn)
	if (response != "SYN") { panic(response) }
	serverSendMessage(conn, remoteaddr, "SYN ACK")

	response, remoteaddr = serverRecieveMessage(conn)
	if (response != "ACK") { panic(response) }

	return conn, remoteaddr
}

func serverSendMessage(conn *net.UDPConn, addr *net.UDPAddr, message string) {
	if sequenceNumberServer == 0 {sequenceNumberServer = rand.Int()} else {sequenceNumberServer++}
	var seqString = strconv.Itoa(sequenceNumberServer)
	fmt.Println("Server sending:", message, "; seq: ", seqString)
    _,err := conn.WriteToUDP([]byte(message + ";" + seqString), addr)
    if err != nil {
        fmt.Printf("Couldn't send response %v", err)
    }
}

func serverRecieveMessage(conn *net.UDPConn) (string, *net.UDPAddr) {
	p := make([]byte, 2048)

	_,remoteaddr,err := conn.ReadFromUDP(p)
	if err != nil {
        fmt.Printf("Server receiving error: %v\n", err)
    }

	var response = string(bytes.Trim(p, "\x00")[:])

	var responseArr = strings.Split(response, ";")  // Splitting message and sequence number
	seq, _ := strconv.Atoi(responseArr[1])
	if (sequenceNumberServer == 0){
		temp,_ := strconv.Atoi(responseArr[1])
		sequenceNumberServer = temp
	}else if seq == sequenceNumberServer + 1{
		sequenceNumberServer++
	}else {
		fmt.Printf("%v%d\n", "Server: ", sequenceNumberServer)
		fmt.Printf("%v%v\n", "Client: ", responseArr[1])
		panic("Server recieved an incorrect sequence number")
	}

	fmt.Println("Server received:", responseArr[0], "; seq: ", responseArr[1])
	return responseArr[0], remoteaddr
}