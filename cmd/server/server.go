package main

import "merch_service/internal/server"

func main() {
	serv := server.NewServer()
	serv.Run()
}
