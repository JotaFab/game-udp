package server

import (
	"testing"
	"time"
)

func TestLobbyMatchmaking(t *testing.T) {
	// Clear state before test
	searchingMu.Lock()
	searchingPlayers = nil
	searchingMu.Unlock()
	lobbiesMu.Lock()
	lobbies = make(map[string]*Lobby)
	nextPort = 9000
	lobbiesMu.Unlock()

	// Add 4 players to search
	AddPlayerToSearch("p1")
	AddPlayerToSearch("p2")
	AddPlayerToSearch("p3")
	AddPlayerToSearch("p4")

	// Wait for matchmaker to process
	time.Sleep(1 * time.Second)

	lobbiesMu.Lock()
	defer lobbiesMu.Unlock()
	if len(lobbies) != 1 {
		t.Fatalf("expected 1 lobby, got %d", len(lobbies))
	}
	for _, lobby := range lobbies {
		if len(lobby.PlayerIDs) != 4 {
			t.Errorf("expected 4 players in lobby, got %d", len(lobby.PlayerIDs))
		}
		t.Logf("Lobby created: matchID=%s, port=%d, players=%v", lobby.MatchID, lobby.UDPPort, lobby.PlayerIDs)
	}
}
