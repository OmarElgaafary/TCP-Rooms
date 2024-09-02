func (r *Room) removeRoomUsers(conn net.Conn) error {
	index := -1

	for i, j := range r.UserConnections {
		if j == conn {
			index = i
			break
		}
	}

	if index == -1 {
		return errors.New("connection not found")
	}

	r.UserConnections = append(r.UserConnections[:index], r.UserConnections[:index+1]...)
	return nil
}


func removeNilRoom(roomID string) {

	index := -1

	for i, targetRoom := range globalRooms {
		if targetRoom.ID == roomID {
			index = i
			break
		}
	}

	if index == -1 {
		return
	}

	globalRooms = append(globalRooms[:index], globalRooms[:index+1]...)
}


func DoLeaveRoom(conn net.Conn, roomId string) {

	targetRoom, err := findRoom(roomId)

	if err != nil {
		nillRoomError := fmt.Sprintf("Room %v not found", roomId)
		conn.Write([]byte(nillRoomError))
	}

	err = targetRoom.removeRoomUsers(conn)

	if err != nil {
		removeUserError := fmt.Sprintf("Error while leaving room %v", roomId)
		conn.Write([]byte(removeUserError))
	}

	if len(targetRoom.UserConnections) == 0 {
		removeNilRoom(roomId)
	}

	user := GetUserByConnection(conn)
	if user == nil {
		userNotFoundError := fmt.Sprintf("User not found")
		conn.Write([]byte(userNotFoundError))
	}

	userLeft := fmt.Sprintf("%v has left the room", user)
	conn.Write([]byte(userLeft))

}
