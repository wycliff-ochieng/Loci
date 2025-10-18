package socket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
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
}

type Websocket struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func (h *Hub) Run() {

	for {
		select {
		case client := <-h.Register:
			//map[*Client] bool = client -> add client to map
			h.Clients[client] = true

		case client := <-h.UnRegister:
			//check if client exists in the map
			if _, ok := h.Clients[client]; ok {

				//delete if present then close
				delete(h.Clients, client)

				close(client.Send)

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
		}

	}
}

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

func (h *Hub) BroadcastNewLoci(locus *models.Locus) {
	h.BroadcastLocus <- locus
}
