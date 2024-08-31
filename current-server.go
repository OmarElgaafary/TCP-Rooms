package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
)

type User struct {
	Nickname string
}

type Room struct {
	ID              string
	Users           []User
	UserConnections []net.Conn
}

var globalRooms []Room

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

var connections map[net.Conn]*User

func handleConnection(conn net.Conn) {
	remoteAddr := conn.RemoteAddr().String()
	fmt.Println("Client connected from " + remoteAddr)

	scanner := bufio.NewScanner(conn)

	for {
		ok := scanner.Scan()

		if !ok {
			break
		}

		data := scanner.Text()

		var userMsg CommunicationMessage

		json.Unmarshal([]byte(data), &userMsg)

		switch userMsg.Op {
		case "msg":
			SendMessage(conn, userMsg.Data)
		case "identify":
			IdentifyConnection(conn, userMsg.Data)
		case CreateRoom:
			DoCreateRoom(conn, userMsg.RoomId)
		case ListRoom:
			DoListRoom(conn)
		case JoinRoom:
			DoJoinRoom(conn, userMsg.RoomId)
		case LeaveRoom:
			DoLeaveRoom(conn, userMsg.RoomId)
		}

		// tokens := strings.Split(data, " ")
		// if len(tokens) == 0 {
		// 	continue
		// }

		// commandDataTokens := tokens[1:]

		// switch tokens[0] {
		// case "\\msg":
		// 	SendMessage(conn, strings.Join(commandDataTokens, " "))
		// case "\\identify":
		// 	IdentifyConnection(conn, strings.Join(commandDataTokens, " "))
		// }

	}

	fmt.Println("Client at " + remoteAddr + " disconnected.")
}

func DoCreateRoom(conn net.Conn, roomId string) {
	r := Room{
		ID: roomId,
	}
	globalRooms = append(globalRooms, r)
}

func DoListRoom(conn net.Conn) {
	jsonString, _ := json.Marshal(globalRooms)
	conn.Write(jsonString)
}

func DoJoinRoom(conn net.Conn, roomID string) {

	targetRoom, err := findRoom(conn, roomID)
	targetRoom.addRoomUsers(conn)

	if err != nil {
		nillRoomError := fmt.Sprintf("Room %v not found: %v", roomID, err)
		conn.Write([]byte(nillRoomError))
		return
	}

	targetRoom.addRoomUsers(conn)

	user := GetUserByConnection(conn)
	if user == nil {
		fmt.Printf("User not found")
	}
	joinMsg := fmt.Sprintf("%s has joined the room\r\n", user)
	conn.Write([]byte(joinMsg))

}

func DoLeaveRoom(conn net.Conn, roomId string) {}

func SendMessage(conn net.Conn, msg string) {
	user := GetUserByConnection(conn)
	if user == nil {
		return
	}

	targetUserName := ""
	var targetConnection net.Conn = nil

	if msg[0] == '@' {
		data := strings.Split(msg, " ")
		targetUserName = data[0]
		targetConnection = GetUserConnection(targetUserName[1:])
		msg = strings.Join(data[1:], " ")

	}

	for k, _ := range connections {
		if targetConnection == nil || targetConnection == k {
			message := fmt.Sprintf("%s: %s\r\n", user.Nickname, msg)
			k.Write([]byte(message))
		}

	}
}

func IdentifyConnection(conn net.Conn, msg string) {
	connections[conn] = &User{Nickname: msg}
}

func GetUserByConnection(conn net.Conn) *User {
	for k, v := range connections {
		if k == conn {
			return v
		}
	}

	return nil
}

func GetUserConnection(username string) net.Conn {
	for k, v := range connections {
		if v.Nickname == username {
			return k
		}
	}

	return nil
}

func findRoom(conn net.Conn, roomID string) (*Room, error) {
	for _, room := range globalRooms {
		if room.ID == roomID {
			return &room, nil
		} else {
			fmt.Println("Room Not Found")
		}
	}

	return nil, errors.New("room not found")
}

func (r *Room) addRoomUsers(conn net.Conn) {
	r.UserConnections = append(r.UserConnections, conn)
}

func main() {

	connections = make(map[net.Conn]*User, 0)

	address := "192.168.1.7:90"

	listener, _ := net.Listen("tcp", address)
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		connections[conn] = &User{}

		if err != nil {
			fmt.Printf("Some connection error: %s\n", err)
		}
		go handleConnection(conn)
	}
}
