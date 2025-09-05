package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"

	controlpb "game-server/proto/control"
)

func wsHandler(w http.ResponseWriter, r *http.Request) {

    c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer c.Close()

	_, msg, err := c.ReadMessage()
	if err != nil {
		log.Println("ws read:", err)
		return
	}
	playerID := string(msg)
	log.Printf("WS join: %s", playerID)

	token := genToken(12)
	session := &PlayerSession{PlayerID: playerID, SessionToken: token}
	sessionsMu.Lock()
	sessions[token] = session
	sessionsMu.Unlock()

	idToTokenMu.Lock()
	idToToken[playerID] = token
	idToTokenMu.Unlock()

	game.mu.Lock()
	game.Players[playerID] = &PlayerState{X: 0, Y: 0}
	game.mu.Unlock()

	start := &controlpb.StartGame{
		UdpPort:      int32(udpPort),
		SessionToken: token,
		PlayerId:     playerID,
	}
	data, _ := proto.Marshal(start)
	if err := c.WriteMessage(websocket.BinaryMessage, data); err != nil {
		log.Println("ws write:", err)
		return
	}

	for {
		if _, _, err := c.NextReader(); err != nil {
			log.Println("ws closed for", playerID)
			sessionsMu.Lock()
			delete(sessions, token)
			sessionsMu.Unlock()
			idToTokenMu.Lock()
			delete(idToToken, playerID)
			idToTokenMu.Unlock()
			game.mu.Lock()
			delete(game.Players, playerID)
			game.mu.Unlock()
			return
		}
	}
}
