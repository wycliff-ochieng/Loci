package socket

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/wycliff-ochieng/internal/models"
)

type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	Send     chan []byte
	UserID   uuid.UUID
	mu       sync.Mutex
	Location *models.GeoPoint
}

var (
	ErrInvalidLocation = errors.New("0.0 cant be a valid geolocation")
)

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
		case "USER_UPDATE_LOCATION":

			var loc models.GeoPoint

			var buf bytes.Buffer

			enc := gob.NewEncoder(&buf)

			err := enc.Encode(msg.Payload)
			if err != nil {
				log.Printf("djjd%s", err)

			}

			if err := json.Unmarshal(buf.Bytes(), &loc); err != nil {
				log.Printf("failed to unmarshal due to : %v", err)
				break
			}

			newLocation := &models.GeoPoint{
				Lat:  loc.Lat,
				Long: loc.Long,
			}

			//handle business logic
			if loc.Lat == 0.0 || loc.Long == 0.0 {
				//http.Error(w,"invlaid data input",http.StatusLengthRequired)
				log.Println("input valid address")
			}

			c.mu.Lock()

			c.Location = newLocation

			c.mu.Unlock()
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
				//hub closed the channel
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)

			//add queued chat message to the current websocket message
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
