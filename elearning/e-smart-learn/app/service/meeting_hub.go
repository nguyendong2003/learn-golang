package service

import (
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	meetingWriteWait      = 10 * time.Second
	meetingPongWait       = 60 * time.Second
	meetingPingPeriod     = (meetingPongWait * 9) / 10
	meetingMaxMessageSize = 1024 * 1024
)

const (
	MsgParticipantsSnapshot = "participants_snapshot"
	MsgParticipantJoined    = "participant_joined"
	MsgParticipantLeft      = "participant_left"
	MsgError                = "error"
	MsgPing                 = "ping"
	MsgPong                 = "pong"
	MsgSignal               = "signal"
)

type meetingRoom struct {
	id      string
	clients map[*meetingClient]struct{}
}

type MeetingHub struct {
	mu    sync.RWMutex
	rooms map[string]*meetingRoom
}

type meetingClient struct {
	hub         *MeetingHub
	conn        *websocket.Conn
	send        chan []byte
	roomID      string
	eventID     uuid.UUID
	participant MeetingParticipant
}

type meetingEnvelope struct {
	Type      string      `json:"type"`
	Payload   any         `json:"payload,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	EventID   uuid.UUID   `json:"event_id"`
	Sender    any         `json:"sender,omitempty"`
	Meta      interface{} `json:"meta,omitempty"`
}

type meetingIncoming struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type meetingSignalPayload struct {
	TargetUserID string          `json:"target_user_id"`
	Kind         string          `json:"kind,omitempty"`
	Data         json.RawMessage `json:"data"`
}

func NewMeetingHub() *MeetingHub {
	return &MeetingHub{
		rooms: make(map[string]*meetingRoom),
	}
}

func (h *MeetingHub) HandleConnection(conn *websocket.Conn, roomID string, eventID uuid.UUID, participant MeetingParticipant) {
	client := &meetingClient{
		hub:         h,
		conn:        conn,
		send:        make(chan []byte, 256),
		roomID:      roomID,
		eventID:     eventID,
		participant: participant,
	}

	h.register(client)
	go client.writePump()
	go client.readPump()
}

func (h *MeetingHub) register(client *meetingClient) {
	h.mu.Lock()
	room, ok := h.rooms[client.roomID]
	if !ok {
		room = &meetingRoom{
			id:      client.roomID,
			clients: make(map[*meetingClient]struct{}),
		}
		h.rooms[client.roomID] = room
	}
	room.clients[client] = struct{}{}
	h.mu.Unlock()

	participants := h.getParticipants(client.roomID)
	h.sendToClient(client, meetingEnvelope{
		Type:      MsgParticipantsSnapshot,
		Timestamp: time.Now().UTC(),
		EventID:   client.eventID,
		Payload: map[string]any{
			"participants": participants,
		},
	})

	h.broadcast(client.roomID, meetingEnvelope{
		Type:      MsgParticipantJoined,
		Timestamp: time.Now().UTC(),
		EventID:   client.eventID,
		Payload: map[string]any{
			"participant": client.participant,
		},
	}, client)
}

func (h *MeetingHub) unregister(client *meetingClient) {
	h.mu.Lock()
	room, ok := h.rooms[client.roomID]
	if !ok {
		h.mu.Unlock()
		return
	}
	if _, exists := room.clients[client]; exists {
		delete(room.clients, client)
		close(client.send)
	}
	if len(room.clients) == 0 {
		delete(h.rooms, client.roomID)
	}
	h.mu.Unlock()

	h.broadcast(client.roomID, meetingEnvelope{
		Type:      MsgParticipantLeft,
		Timestamp: time.Now().UTC(),
		EventID:   client.eventID,
		Payload: map[string]any{
			"participant": client.participant,
		},
	}, client)
}

func (h *MeetingHub) getParticipants(roomID string) []MeetingParticipant {
	h.mu.RLock()
	defer h.mu.RUnlock()

	room, ok := h.rooms[roomID]
	if !ok {
		return []MeetingParticipant{}
	}

	participants := make([]MeetingParticipant, 0, len(room.clients))
	for client := range room.clients {
		participants = append(participants, client.participant)
	}

	return participants
}

func (h *MeetingHub) broadcast(roomID string, envelope meetingEnvelope, excluded ...*meetingClient) {
	envelopeBytes, err := json.Marshal(envelope)
	if err != nil {
		return
	}

	excludedSet := make(map[*meetingClient]struct{}, len(excluded))
	for _, c := range excluded {
		excludedSet[c] = struct{}{}
	}

	h.mu.RLock()
	room, ok := h.rooms[roomID]
	if !ok {
		h.mu.RUnlock()
		return
	}

	clients := make([]*meetingClient, 0, len(room.clients))
	for client := range room.clients {
		if _, skip := excludedSet[client]; skip {
			continue
		}
		clients = append(clients, client)
	}
	h.mu.RUnlock()

	for _, client := range clients {
		h.sendRaw(client, envelopeBytes)
	}
}

func (h *MeetingHub) sendToClient(client *meetingClient, envelope meetingEnvelope) {
	b, err := json.Marshal(envelope)
	if err != nil {
		return
	}
	h.sendRaw(client, b)
}

func (h *MeetingHub) sendRaw(client *meetingClient, payload []byte) {
	select {
	case client.send <- payload:
	default:
		h.unregister(client)
		_ = client.conn.Close()
	}
}

func (h *MeetingHub) sendToUser(roomID string, userID uuid.UUID, envelope meetingEnvelope, excluded ...*meetingClient) {
	envelopeBytes, err := json.Marshal(envelope)
	if err != nil {
		return
	}

	h.mu.RLock()
	room, ok := h.rooms[roomID]
	if !ok {
		h.mu.RUnlock()
		return
	}

	clients := make([]*meetingClient, 0, len(room.clients))
	for client := range room.clients {
		if client.participant.UserID != userID {
			continue
		}
		clients = append(clients, client)
	}
	h.mu.RUnlock()

	for _, client := range clients {
		h.sendRaw(client, envelopeBytes)
	}
}

func (c *meetingClient) readPump() {
	defer func() {
		c.hub.unregister(c)
		_ = c.conn.Close()
	}()

	c.conn.SetReadLimit(meetingMaxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(meetingPongWait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(meetingPongWait))
	})

	for {
		_, payload, err := c.conn.ReadMessage()
		if err != nil {
			return
		}

		var incoming meetingIncoming
		if err := json.Unmarshal(payload, &incoming); err != nil {
			c.hub.sendToClient(c, meetingEnvelope{
				Type:      MsgError,
				Timestamp: time.Now().UTC(),
				EventID:   c.eventID,
				Payload: map[string]any{
					"message": "Invalid payload format",
				},
			})
			continue
		}

		switch strings.ToLower(strings.TrimSpace(incoming.Type)) {
		case MsgPing:
			c.hub.sendToClient(c, meetingEnvelope{
				Type:      MsgPong,
				Timestamp: time.Now().UTC(),
				EventID:   c.eventID,
			})
		case MsgSignal:
			var signal meetingSignalPayload
			if err := json.Unmarshal(incoming.Payload, &signal); err != nil {
				continue
			}
			envelope := meetingEnvelope{
				Type:      MsgSignal,
				Timestamp: time.Now().UTC(),
				EventID:   c.eventID,
				Sender:    c.participant,
				Payload: map[string]any{
					"kind": signal.Kind,
					"data": json.RawMessage(signal.Data),
				},
			}

			targetUserID := strings.TrimSpace(signal.TargetUserID)
			if targetUserID == "" {
				c.hub.broadcast(c.roomID, envelope, c)
				continue
			}
			toUserID, err := uuid.Parse(targetUserID)
			if err != nil {
				c.hub.sendToClient(c, meetingEnvelope{
					Type:      MsgError,
					Timestamp: time.Now().UTC(),
					EventID:   c.eventID,
					Payload: map[string]any{
						"message": "Invalid target user ID format",
					},
				})
				continue
			}
			c.hub.sendToUser(c.roomID, toUserID, envelope, c)
		default:
			c.hub.sendToClient(c, meetingEnvelope{
				Type:      MsgError,
				Timestamp: time.Now().UTC(),
				EventID:   c.eventID,
				Payload: map[string]any{
					"message": "Unsupported message type",
				},
			})
		}
	}
}

func (c *meetingClient) writePump() {
	ticker := time.NewTicker(meetingPingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(meetingWriteWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			writer, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			if _, err := writer.Write(msg); err != nil {
				_ = writer.Close()
				return
			}

			pending := len(c.send)
			for i := 0; i < pending; i++ {
				if _, err := writer.Write([]byte{'\n'}); err != nil {
					_ = writer.Close()
					return
				}
				if _, err := writer.Write(<-c.send); err != nil {
					_ = writer.Close()
					return
				}
			}

			if err := writer.Close(); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(meetingWriteWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
