package socket

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/wycliff-ochieng/pkg/middleware"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServerWS(h *Hub, w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Error: upgrading issue due to: %s", http.StatusFailedDependency)
		return
	}

	defer conn.Close()

	userID, err := middleware.GetUserUUIDFromContext(ctx)
	if err != nil {
		http.Error(w, "Error fetching ID from context", http.StatusExpectationFailed)
		return
	}

	//create a client instance
	client := &Client{
		Hub:    h,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		UserID: userID,
		//Location: ,
	}

	//register new client with the hub
	client.Hub.Register <- client

	//start read and write pumps
}
