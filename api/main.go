// main.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

// --- Estado Global da Aplicação ---
var stateMutex sync.Mutex
var trancaEstaAberta bool = true // Estado inicial é "aberta"

// Estrutura da resposta JSON
type StatusResposta struct {
	TrancaAberta    bool   `json:"tranca_aberta"`
	UltimaAtualizacao string `json:"ultima_atualizacao"`
}

// Handler para a rota GET /status (Consultar)
func statusHandler(w http.ResponseWriter, r *http.Request) {
	stateMutex.Lock()
	defer stateMutex.Unlock()

	resposta := StatusResposta{
		TrancaAberta:    trancaEstaAberta,
		UltimaAtualizacao: time.Now().Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resposta)
}

// Handler para a rota POST /toggle (Mudar o Estado)
func toggleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido, use POST", http.StatusMethodNotAllowed)
		return
	}
	stateMutex.Lock()
	defer stateMutex.Unlock()

	trancaEstaAberta = !trancaEstaAberta
	log.Printf("Estado da tranca alterado para: %v", trancaEstaAberta)

	resposta := StatusResposta{
		TrancaAberta:    trancaEstaAberta,
		UltimaAtualizacao: time.Now().Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resposta)
}

func main() {
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/toggle", toggleHandler)

	log.Println("API rodando em http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}