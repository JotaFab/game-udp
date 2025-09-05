package main

import (
    "log"
    "os"
    "game-server/server"
)
// Removed context import as it is not used

func main() {
    logger := log.New(os.Stdout, "[game-server] ", log.LstdFlags)
    server, err := server.NewGameServer(logger)
    if err != nil {
        logger.Fatalf("failed to start server: %v", err)
    }
    defer server.UDPConn.Close()
    server.Start()
}
