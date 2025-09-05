package main

import (
    "log"
    "net"
    "sync"
    realtimepb "game-server/proto/realtime"
    "google.golang.org/protobuf/proto"
)

func tick(conn *net.UDPConn) {
    game.mu.Lock()
    game.Tick++
    tick := game.Tick
    snap := &realtimepb.ServerSnapshot{
        Tick: tick,
    }
    for id, st := range game.Players {
        ps := &realtimepb.PlayerState{
            PlayerId: id,
            X:        st.X,
            Y:        st.Y,
        }
        snap.Players = append(snap.Players, ps)
    }
    game.mu.Unlock()

    data, err := proto.Marshal(snap)
    if err != nil {
        log.Println("proto marshal snapshot:", err)
        return
    }
    sessionsMu.Lock()
    defer sessionsMu.Unlock()
    for _, s := range sessions {
        if s.UDPAddr == nil {
            continue
        }
        if _, err := conn.WriteToUDP(data, s.UDPAddr); err != nil {
            log.Printf("udp write to %v error: %v", s.UDPAddr, err)
        }
    }
}

// Shared types and variables

type PlayerSession struct {
    PlayerID     string
    SessionToken string
    UDPAddr *net.UDPAddr
    LastSeq uint32
}

type GameState struct {
    mu      sync.Mutex
    Tick    uint32
    Players map[string]*PlayerState
}

type PlayerState struct {
    X float32
    Y float32
}

var (
    sessionsMu sync.Mutex
    sessions   = map[string]*PlayerSession{}

    idToTokenMu sync.Mutex
    idToToken   = map[string]string{}

    game = &GameState{
        Players: map[string]*PlayerState{},
    }
)
