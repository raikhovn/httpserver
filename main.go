package main

import (
	"httpfileserver/server"
	"log"
)

func main() {

	log.Println("Creating instance of the server ...")
	s := server.NewServer(8080, "0.0.0.0")
	log.Println("Registring route handlers ...")
	s.RegisterIndex()
	s.RegisterListDirs()
	s.RegisterListFiles()
	s.RegisterGetFileInfo()
	log.Println("Starting listener ...")
	s.StartListener()
}
