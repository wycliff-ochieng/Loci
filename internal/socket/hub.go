package socket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/wycliff-ochieng/internal/models"
)

/*
type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	Send     chan []byte
	UserID   uuid.UUID
	Location *models.GeoPoint
} */

const (
	writeWait  = 10 * time.Second
	pingPeriod = (writeWait * 9) / 10
	pongWait   = 10 * time.Second
)

type Hub struct {
	Clients        map[*Client]bool
	Register       chan *Client
	UnRegister     chan *Client
	Broadcast      chan []byte
	BroadcastLocus chan *models.Locus
	BroadcastReply chan *models.ReplyEvent
}

type Websocket struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:      make(chan []byte),
		Register:       make(chan *Client),
		UnRegister:     make(chan *Client),
		Clients:        make(map[*Client]bool),
		BroadcastLocus: make(chan *models.Locus),
		BroadcastReply: make(chan *models.ReplyEvent),
	}
}

func (h *Hub) Run() {

	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			log.Printf("[HUB] Client registered. Total clients: %d", len(h.Clients))

		case client := <-h.UnRegister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
				log.Printf("[HUB] Client unregistered. Total clients: %d", len(h.Clients))
			}

		case locus := <-h.BroadcastLocus:
			//prepare outbound message
			message := Websocket{
				Type:    "LOCUS_NEW",
				Payload: locus,
			}

			//unmarshal the message
			marshalledMessage, err := json.Marshal(message)
			if err != nil {
				log.Printf("error marshaling message due to: %s", err)
				continue
			}

			for c := range h.Clients {
				if c.Location == nil {
					continue // skip clients with no location
				}

				locusPoint := &models.GeoPoint{
					Lat:  locus.Location.Lat,
					Long: locus.Location.Long,
				}

				//calculate distance btw client and loci location

				dist := models.CalculateDistance(c.Location, locusPoint)

				//set broadcast distance -> distance allowed to see the message
				const broadcastDistance = 5000

				if dist < broadcastDistance {
					select {
					case c.Send <- marshalledMessage:
					default:
						close(c.Send)
						delete(h.Clients, c)
					}
				}
			}
		case reply := <-h.BroadcastReply:
			log.Printf("[HUB] Received Reply Event for Locus: %s. Processing for %d clients...", reply.LocusID, len(h.Clients))

			replyEvent := Websocket{
				Type:    "REPLY_EVENT",
				Payload: reply,
			}

			marshalledReply, err := json.Marshal(replyEvent)
			if err != nil {
				log.Printf("Error marshalling reply  due to: %s", err)
				continue
			}

			//clients location
			for c := range h.Clients {
				if c.Location == nil {
					log.Println("[HUB] Skipping client: Location is nil")
					continue
				}

				replyLocation := &models.GeoPoint{
					Lat:  reply.LocusLocation.Lat,
					Long: reply.LocusLocation.Long,
				}

				dist := models.CalculateDistance(c.Location, replyLocation)
				log.Printf("[HUB] Client Dist: %.2f km. (Limit: 5000km)", dist)

				const replyBroadcastDistance = 5000

				if dist < replyBroadcastDistance {
					select {
					case c.Send <- marshalledReply:
					default:
						close(c.Send)
						delete(h.Clients, c)
					}
				}
			}

		}
	}
}

/*
func (c *Client) WritePump() {

	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				//handle err
			}

			w.Write(message)

			for i := 0; i < len(c.Send); i++ {
				w.Write(<-c.Send)
			}

			w.Close()

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				//handle err
			}
			return
		}

	}

}
*/
func (h *Hub) BroadcastNewLoci(locus *models.Locus) {
	h.BroadcastLocus <- locus
}
