package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var conn *net.TCPConn

type Operation string

const Identify Operation = "identify"
const Msg Operation = "msg"
const CreateRoom Operation = "create_room"
const ListRoom Operation = "list_rooms"
const JoinRoom Operation = "join_room"
const LeaveRoom Operation = "leave_room"

type CommunicationMessage struct {
	Op     Operation
	Data   string
	RoomId string
}

func handleRead(conn *net.TCPConn) {

	response := make([]byte, 20000)
	for {
		_, err := conn.Read(response)
		if err != nil {
			println(err.Error())
		}
		println(string(response))
	}

}

func identify(nickname string) {
	data := CommunicationMessage{
		Op:   Identify,
		Data: nickname}
	jsonString, _ := json.Marshal(data)
	jsonString = append(jsonString, '\n')
	conn.Write(jsonString)
}

func sendMessage(msg string) {
	data := CommunicationMessage{
		Op:   Msg,
		Data: msg}
	jsonString, _ := json.Marshal(data)
	jsonString = append(jsonString, '\n')
	conn.Write([]byte(jsonString))
}

func DoJoinRoom(roomId string) {
	data := CommunicationMessage{
		Op:     JoinRoom,
		RoomId: roomId,
		Data:   ""}

	jsonString, _ := json.Marshal(data)
	jsonString = append(jsonString, '\n')
	conn.Write(jsonString)
}

func DoListRooms() {

	data := CommunicationMessage{
		Op:   ListRoom,
		Data: ""}

	jsonString, _ := json.Marshal(data)
	jsonString = append(jsonString, '\n')
	conn.Write(jsonString)

}

func DoCreateRoom(roomId string) {
	data := CommunicationMessage{
		Op:     CreateRoom,
		RoomId: roomId,
		Data:   "",
	}

	jsonString, _ := json.Marshal(data)
	jsonString = append(jsonString, '\n')
	conn.Write(jsonString)
}

func processCommand(data string) {
	if len(data) == 0 {
		return
	}
	if data[0] == '\\' {
		// this is a command
		tokens := strings.Split(data, " ")
		command := tokens[0][1:]
		switch command {
		case "list":
			{
				DoListRooms()
			}
		case "join":
			{
				DoJoinRoom(tokens[1])
			}
		case "leave":
			{
			}
		case "create":
			DoCreateRoom(tokens[1])

		}
		return
	}

	//send chat message
}

func main() {
	address := ":90"
	tcpAddress, err := net.ResolveTCPAddr("tcp", address)

	if err != nil {
		println(err.Error())
		return
	}

	conn, err = net.DialTCP("tcp", nil, tcpAddress)
	if err != nil {
		println(err.Error())
		return
	}
	go handleRead(conn)

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("Enter your nickname: ")
	scanner.Scan()
	nickname := scanner.Text()
	identify(nickname)

	for scanner.Scan() {
		data := scanner.Text()
		processCommand(data)
	}

	// close app then CTRL + C is pressed

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
	fmt.Println("Adios!")
}
