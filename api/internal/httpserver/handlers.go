package httpserver

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"golang.org/x/crypto/bcrypt"

	"smartlock/internal/auth"
	"smartlock/internal/config"
	"smartlock/internal/model"
	"smartlock/internal/mqtt"
	"smartlock/internal/state"
	"smartlock/internal/db"
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

// rota login (pública) 
func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "use POST", http.StatusMethodNotAllowed)
        return
    }
    var req model.LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "payload inválido", http.StatusBadRequest)
        return
    }

    // busca no banco
    user, err := db.GetUserByEmail(r.Context(), h.cfg.DB, req.Email)
    if err != nil {
        http.Error(w, "email ou senha errados", http.StatusUnauthorized)
        return
    }

    // compara senha hash (bcrypt)
    if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
        http.Error(w, "email ou senha errados", http.StatusUnauthorized)
        return
    }

    // gera token JWT
    tok, _ := h.authSvc.Generate(user.ID, user.Email)

    resp := model.LoginResponse{
        Success: true,
        Token:   tok,
        User: &struct {
            ID    string `json:"id"`
            Email string `json:"email"`
			Name  string `json:"name"`
        }{ID: user.ID, Email: user.Email, Name: user.Name},
    }
    writeJSON(w, resp)
}



// Status simples (protegido, com long-poll)
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

// Histórico (protegido)
func (h *Handlers) History(w http.ResponseWriter, r *http.Request) {
	limit := 20
	if q := r.URL.Query().Get("limit"); q != "" {
		if n, err := strconv.Atoi(q); err == nil && n >= 0 && n <= h.cfg.MaxHistory {
			limit = n
		}
	}

	ctx := r.Context()
	hist, err := db.GetHistory(ctx, h.cfg.DB, limit)
	if err != nil {
		http.Error(w, "Erro ao consultar histórico", http.StatusInternalServerError)
		return
	}

	out := make([]model.HistoryItem, 0, len(hist))
	for _, rec := range hist {
		out = append(out, model.HistoryItem{
			ID:        rec.ID,
			Action:    rec.Action,
			Timestamp: rec.Timestamp.Format(time.RFC3339),
			Method:    rec.Method,
		})
	}

	writeJSON(w, model.HistoryResponse{History: out})
}
// Toggle (protegido)
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
