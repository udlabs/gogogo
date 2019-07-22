package server

import "net/http"

type HttpServer struct {
	Port string
}

func (s *HttpServer) Start(handler http.Handler) {
	http.ListenAndServe(s.Port, handler)
}
