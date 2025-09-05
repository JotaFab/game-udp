package main

import (
	realtimepb "game-server/proto/realtime"
	"log"
	"math"
	"net"

	"google.golang.org/protobuf/proto"
)

func udpReadLoop(conn *net.UDPConn) {
	buf := make([]byte, 2048)
	for {
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Println("udp read:", err)
			continue
		}
		var in realtimepb.ClientInput
		if err := proto.Unmarshal(buf[:n], &in); err != nil {
			log.Println("proto unmarshal:", err)
			continue
		}
		sessionsMu.Lock()
		s, ok := sessions[in.SessionToken]
		sessionsMu.Unlock()
		if !ok {
			log.Printf("unknown token from %v", addr)
			continue
		}
		if s.UDPAddr == nil {
			s.UDPAddr = addr
			log.Printf("registered UDP addr %v for player %s", addr, s.PlayerID)
		}
		if in.Seq <= s.LastSeq {
			continue
		}
		s.LastSeq = in.Seq
		applyInputToState(s.PlayerID, &in)
	}
}

func applyInputToState(playerID string, in *realtimepb.ClientInput) {
	game.mu.Lock()
	defer game.mu.Unlock()
	p, ok := game.Players[playerID]
	if !ok {
		return
	}
	speed := float32(0.1)
	mag := float32(math.Sqrt(float64(in.MoveX*in.MoveX + in.MoveY*in.MoveY)))
	if mag > 0 {
		nx := in.MoveX / mag
		ny := in.MoveY / mag
		p.X += nx * speed
		p.Y += ny * speed
	}
	if p.X > 100 {
		p.X = 100
	}
	if p.X < -100 {
		p.X = -100
	}
	if p.Y > 100 {
		p.Y = 100
	}
	if p.Y < -100 {
		p.Y = -100
	}
}
