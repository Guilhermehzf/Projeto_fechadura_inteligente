package httpserver

import (
	"log"
	"net/http"

	"smartlock/internal/auth"
	"smartlock/internal/config"
	"smartlock/internal/mqtt"
	"smartlock/internal/state"
)

type Server struct {
	cfg  *config.Config
	http *http.ServeMux
}

func NewServer(cfg *config.Config, st *state.Store, mq *mqtt.Client, as *auth.Service) *Server {
	h := NewHandlers(cfg, st, mq, as)

	mux := http.NewServeMux()

	// pública
	mux.Handle("/login", withCORS(http.HandlerFunc(h.Login)))

	// protegidas
	mux.Handle("/status", withCORS(authMiddleware(as, http.HandlerFunc(h.StatusSimple))))
	mux.Handle("/history", withCORS(authMiddleware(as, http.HandlerFunc(h.History))))
	mux.Handle("/toggle", withCORS(authMiddleware(as, http.HandlerFunc(h.Toggle))))

	return &Server{cfg: cfg, http: mux}
}

func (s *Server) Start() {
	log.Printf("HTTP escutando em http://%s", s.cfg.HTTPAddr)
	if err := http.ListenAndServe(s.cfg.HTTPAddr, s.http); err != nil {
		log.Fatalf("HTTP erro: %v", err)
	}
}
