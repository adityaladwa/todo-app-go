package main

import (
    "github.com/adityaladwa/todo-app/internal/server"
    "log"
    "os"
    "os/signal"
    "syscall"
)

func main() {
    configPath := "config/config.yaml"
    if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
        configPath = envPath
    }

    srv, err := server.NewServer(configPath)
    if err != nil {
        log.Fatalf("Failed to initialize server: %v", err)
    }

    // Handle graceful shutdown
    go func() {
        srv.Start()
    }()

    // Wait for interrupt signal to gracefully shutdown the server
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    srv.Shutdown()
}
