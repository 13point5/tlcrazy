package server

import (
	"encoding/json"
	"log"
	"net/http"

	"tlcrazy-backend/internal/ai"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
	}))
	r.Use(middleware.Logger)

	r.Post("/tldraw-tool", s.GenerateToolHandler)

	return r
}

type GenerateToolRequest struct {
	Query string `json:"query"`
}

func (s *Server) GenerateToolHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	body := GenerateToolRequest{}
	err := decoder.Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tool, err := ai.GenTldrawTool(body.Query)
	if err != nil {
		log.Printf("Error generating tool: %s", err)
		w.WriteHeader(500)
		return
	}

	resp, err := json.Marshal(tool)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(resp)
}
