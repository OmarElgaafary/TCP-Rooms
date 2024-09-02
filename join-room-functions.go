type Room struct {
	ID    string
	Users []User
	Connections []net.Conn
}

var globalRooms []Room

func (r *Room) addRoomUsers(conn net.Conn) {
	r.Connections = append(r.Connections, conn)
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
