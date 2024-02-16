package ws

import (
	"fmt"
	"github.com/EraldCaka/chat-room/internal/types"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

type Handler struct {
	hub *Hub
}

func NewHandler(h *Hub) *Handler {
	return &Handler{
		hub: h,
	}
}

type CreateRoomReq struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RoomActiveClients struct {
	ID          string   `json:"id"`
	Clients     []string `json:"clients"`
	ClientCount int      `json:"clientCount"`
}

func (h *Handler) CreateRoom(c *gin.Context) {
	var req CreateRoomReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if h.hub.Rooms[req.ID] != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": types.NewError(400, "Room Already Exists")})
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": types.NewError(400, "Name is empty, Name should not be empty")})
		return
	}

	h.hub.Rooms[req.ID] = &Room{
		ID:      req.ID,
		Name:    req.Name,
		Clients: make(map[string]*Client),
	}

	c.JSON(http.StatusOK, req)
}

func (h *Handler) GetRoomActiveClients(c *gin.Context) {
	roomID := c.Param("roomID")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": types.NewError(400, "ID is empty,ID should not be empty")})
		return
	}
	if h.hub.Rooms[roomID] == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": types.NewError(400, "This room Id doesn't exist")})
		return
	}

	clients := make([]string, 0)
	var clientCount int
	for _, client := range h.hub.Rooms[roomID].Clients {
		clients = append(clients, client.Username)
		clientCount++
	}

	response := &RoomActiveClients{
		ID:          roomID,
		Clients:     clients,
		ClientCount: clientCount,
	}
	c.JSON(http.StatusOK, response)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) JoinRoom(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roomID := c.Param("roomId")
	clientID := c.Query("userId")
	username := c.Query("username")

	cl := &Client{
		Conn:     conn,
		Message:  make(chan *Message, 10),
		ID:       clientID,
		RoomID:   roomID,
		Username: username,
	}

	m := &Message{
		Content:  "A new user has joined the room",
		RoomID:   roomID,
		Username: username,
	}

	h.hub.Register <- cl
	h.hub.Broadcast <- m

	go cl.writeMessage()
	cl.readMessage(h.hub)
}

type RoomRes struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Clients []string `json:"clients"`
}

func (h *Handler) GetRooms(c *gin.Context) {
	rooms := make([]RoomRes, 0)

	for _, r := range h.hub.Rooms {
		var clients []string
		var clientCount int
		for _, client := range r.Clients {
			clients = append(clients, client.Username)
			clientCount++
		}
		rooms = append(rooms, RoomRes{
			ID:      r.ID,
			Name:    r.Name,
			Clients: clients,
		})
	}

	c.JSON(http.StatusOK, rooms)
}

func (h *Handler) CloseWSConnection(c *gin.Context) {

	if room := h.hub.Rooms[c.Param("roomID")]; room == nil || room.ID == "" {
		fmt.Println(room, "2")
		c.JSON(http.StatusBadRequest, types.NewError(http.StatusBadRequest, "bad request(roomID)"))
		return
	}
	client := h.hub.Rooms[c.Param("roomID")].Clients[c.Param("userID")]
	if cl := client; cl == nil || cl.ID == "" {
		c.JSON(http.StatusBadRequest, types.NewError(http.StatusBadRequest, "bad request(userID)"))
		return
	}

	client.Conn.Close()
	c.JSON(http.StatusOK, types.NewResponse("response", "User disconnected successfully"))
}
