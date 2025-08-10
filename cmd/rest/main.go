package main

import (
    "abet/internal/app/rest"
)

func main() {
    server := rest.NewHttpServer()
    server.Start()
    defer server.Stop()
}
