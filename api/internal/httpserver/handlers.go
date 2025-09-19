package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"smartlock/internal/auth"
	"smartlock/internal/config"
	"smartlock/internal/model"
	"smartlock/internal/mqtt"
	"smartlock/internal/state"
)

type Handlers struct {
	cfg     *config.Config
	store   *state.Store
	mqtt    *mqtt.Client
	authSvc *auth.Service
}

func NewHandlers(cfg *config.Config, st *state.Store, mq *mqtt.Client, as *auth.Service) *Handlers {
	return &Handlers{cfg: cfg, store: st, mqtt: mq, authSvc: as}
}

// Login permanece igual ao que te passei antes
func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "use POST", http.StatusMethodNotAllowed)
		return
	}
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" || req.Password == "" {
		http.Error(w, "payload inválido", http.StatusBadRequest)
		return
	}

	tok, _ := h.authSvc.Generate("user-1", req.Email)

	resp := model.LoginResponse{
		Success: true,
		Token:   tok,
		User: &struct {
			ID    string `json:"id"`
			Email string `json:"email"`
			Name  string `json:"name,omitempty"`
		}{ID: "user-3", Email: req.Email},
	}
	writeJSON(w, resp)
}

// Status simples e protegido
func (h *Handlers) StatusSimple(w http.ResponseWriter, r *http.Request) {
	// Long-poll params
	q := r.URL.Query()
	var since time.Time
	if s := q.Get("since"); s != "" {
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			since = t
		}
	}
	timeout := 25 * time.Second
	if v := q.Get("timeout"); v != "" {
		if d, err := time.ParseDuration(v); err == nil && d > 0 && d <= 60*time.Second {
			timeout = d
		}
	}

	// Espera mudança (ou segue direto se já mudou)
	if !h.store.WaitForUpdate(since, timeout) && !since.IsZero() {
		// timeout sem mudanças → 204
		w.WriteHeader(http.StatusNoContent)
		return
	}

	connected, lastSeen, _, _, _, _, _, tranca, _ := h.store.Snapshot(0)

	isLocked := true
	if tranca != nil {
		isLocked = !*tranca // tranca_aberta=true => isLocked=false
	}
	last := ""
	if !lastSeen.IsZero() {
		last = lastSeen.Format(time.RFC3339)
	}

	writeJSON(w, model.SimpleStatus{
		IsLocked:    isLocked,
		IsConnected: connected,
		LastUpdate:  last, // o cliente usa este valor no próximo “since”
	})
}

// Histórico transformado para HistoryItem[] e protegido
func (h *Handlers) History(w http.ResponseWriter, r *http.Request) {
	limit := 20
	if q := r.URL.Query().Get("limit"); q != "" {
		if n, err := strconv.Atoi(q); err == nil && n >= 0 && n <= h.cfg.MaxHistory {
			limit = n
		}
	}

	_, _, _, _, _, _, _, _, hist := h.store.Snapshot(limit)

	out := make([]model.HistoryItem, 0, len(hist))
	for i, um := range hist {
		// tenta extrair tranca_aberta/method do payload
		method := "auto"
		action := ""

		var m map[string]interface{}
		_ = json.Unmarshal([]byte(um.PayloadRaw), &m)

		if v, ok := m["method"].(string); ok && v != "" {
			method = v
		}

		open := false
		if v, ok := m["tranca_aberta"]; ok {
			switch t := v.(type) {
			case bool:
				open = t
			case string:
				open = t == "true" || t == "1"
			case float64:
				open = t != 0
			}
		}
		if open {
			action = "unlocked"
		} else {
			action = "locked"
		}

		id := fmt.Sprintf("%d-%s", i, um.Timestamp)
		out = append(out, model.HistoryItem{
			ID:        id,
			Action:    action,
			Timestamp: um.Timestamp,
			Method:    method,
		})
	}

	writeJSON(w, model.HistoryResponse{History: out})
}

func (h *Handlers) Toggle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "use POST", http.StatusMethodNotAllowed)
		return
	}
	// publica o comando no MQTT
	if err := h.mqtt.PublishToggle(); err != nil {
		writeJSON(w, map[string]any{
			"success": false,
			"message": "erro ao publicar comando",
		})
		return
	}

	// Prediz o novo estado a partir do último snapshot
	_, _, _, _, _, _, _, tranca, _ := h.store.Snapshot(0)
	// tranca_aberta=true => isLocked=false
	curOpen := false
	if tranca != nil {
		curOpen = *tranca
	}
	newOpen := !curOpen
	newIsLocked := !newOpen

	writeJSON(w, map[string]any{
		"success":  true,
		"isLocked": newIsLocked,
	})
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
