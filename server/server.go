package server

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	port string
}

func New() *Server {
	s := &Server{}
	db = NewDB()
	return s
}

func (s *Server) Run() {
	log.Info("Booting up server")
	router := CreateRouter()

	log.WithFields(log.Fields{"port": "3030"}).Info("Running server.")

	log.Fatal(http.ListenAndServe(":3030", router))
}
