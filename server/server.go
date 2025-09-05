package main

import (
    "log"
    "net"
    "net/http"
    "time"
)

const (
    udpPort      = 9999
    tickRateHz   = 20
    tickDuration = time.Second / tickRateHz
    maxPlayers   = 8
)

type GameServer struct {
    UDPConn *net.UDPConn
    Logger  *log.Logger
}

func NewGameServer(logger *log.Logger) (*GameServer, error) {
    udpAddr := &net.UDPAddr{IP: net.IPv4zero, Port: udpPort}
    conn, err := net.ListenUDP("udp", udpAddr)
    if err != nil {
        return nil, err
    }
    return &GameServer{
        UDPConn: conn,
        Logger:  logger,
    }, nil
}

func (s *GameServer) Start() {
    s.Logger.Printf("UDP listening on %v", s.UDPConn.LocalAddr())
    http.HandleFunc("/ws", wsHandler)

    go func() {
        s.Logger.Println("Starting HTTP (WebSocket) server on :8080")
        if err := http.ListenAndServe(":8080", nil); err != nil {
            s.Logger.Fatalf("http serve: %v", err)
        }
    }()

    go udpReadLoop(s.UDPConn)

    ticker := time.NewTicker(tickDuration)
    defer ticker.Stop()
    for {
        <-ticker.C
        tick(s.UDPConn)
    }
}
