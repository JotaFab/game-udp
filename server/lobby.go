package server

import (
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
)

type Lobby struct {
	MatchID   string
	PlayerIDs []string
	UDPPort   int
}

var (
	searchingPlayers []string
	searchingMu      sync.Mutex
	lobbies          = make(map[string]*Lobby)
	lobbiesMu        sync.Mutex
	nextPort         = 9000 // starting port for matches
)

func init() {
	rand.Seed(time.Now().UnixNano())
	go matchmakerLoop()
}

func AddPlayerToSearch(playerID string) {
	searchingMu.Lock()
	defer searchingMu.Unlock()
	searchingPlayers = append(searchingPlayers, playerID)
}

func matchmakerLoop() {
	for {
		searchingMu.Lock()
		if len(searchingPlayers) >= 4 {
			players := searchingPlayers[:4]
			searchingPlayers = searchingPlayers[4:]
			searchingMu.Unlock()

			matchID := genMatchID()
			port := getNextPort()
			lobby := &Lobby{
				MatchID:   matchID,
				PlayerIDs: players,
				UDPPort:   port,
			}
			lobbiesMu.Lock()
			lobbies[matchID] = lobby
			lobbiesMu.Unlock()

			go startLobby(lobby)
		} else {
			searchingMu.Unlock()
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func genMatchID() string {
	return fmt.Sprintf("match-%d", rand.Int63())
}

func getNextPort() int {
	lobbiesMu.Lock()
	defer lobbiesMu.Unlock()
	port := nextPort
	nextPort++
	return port
}

func startLobby(lobby *Lobby) {
	addr := net.UDPAddr{Port: lobby.UDPPort, IP: net.ParseIP("0.0.0.0")}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Failed to start lobby on port %d: %v\n", lobby.UDPPort, err)
		return
	}
	fmt.Printf("Lobby %s started on port %d for players %v\n", lobby.MatchID, lobby.UDPPort, lobby.PlayerIDs)
	go udpReadLoop(conn)
	// Add more lobby logic as needed
}
