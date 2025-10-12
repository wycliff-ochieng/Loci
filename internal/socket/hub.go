package socket

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/wycliff-ochieng/internal/models"
	"github.com/wycliff-ochieng/sqlc"
)

type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	Send     chan []byte
	UserID   uuid.UUID
	Location *models.GeoPoint
}

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

func (h *Hub) BroadcastNewLoci(locus *sqlc.Loci) {
	return
}

// readPump pumps messages from the websocket connection to the hub
func (c *Client) ReadPump() {

	defer func() {
		c.Hub.UnRegister <- c //tell hub this client is gone
		c.Conn.Close()        // close the network connection
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %s", err)
			}
			break
		}
		//message was successfully read process it, unmarshal it to figure out its type
		var msg Websocket
		json.Unmarshal(message, &msg)

		switch msg.Type {

		}
	}

}
