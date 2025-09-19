package state

import (
	"encoding/json"
	"sync"
	"time"

	"smartlock/internal/model"
)

type Store struct {
	mu               sync.RWMutex
	connected        bool
	lastSeen         time.Time
	lastTopic        string
	lastQoS          byte
	lastRetained     bool
	lastPayloadRaw   string
	lastPayload      map[string]interface{}
	trancaEstaAberta *bool
	history          []model.UltimaMensagem
	maxHistory       int
}

func NewStore(maxHistory int) *Store {
	return &Store{
		lastPayload: make(map[string]interface{}),
		maxHistory:  maxHistory,
	}
}

func (s *Store) UpdateFromMQTT(topic string, qos byte, retained bool, payload []byte) {
	now := time.Now().UTC()

	var m map[string]interface{}
	if err := json.Unmarshal(payload, &m); err != nil {
		m = map[string]interface{}{"_non_json_payload": string(payload)}
	}

	var ta *bool
	if v, ok := m["tranca_aberta"]; ok {
		switch b := v.(type) {
		case bool:
			ta = &b
		case float64:
			val := b != 0
			ta = &val
		case string:
			switch b {
			case "true", "1":
				val := true; ta = &val
			case "false", "0":
				val := false; ta = &val
			}
		}
	}

	um := model.UltimaMensagem{
		Topic:      topic,
		QoS:        qos,
		Retained:   retained,
		Timestamp:  now.Format(time.RFC3339),
		PayloadRaw: string(payload),
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.connected = true
	s.lastSeen = now
	s.lastTopic = topic
	s.lastQoS = qos
	s.lastRetained = retained
	s.lastPayloadRaw = string(payload)
	s.lastPayload = m
	s.trancaEstaAberta = ta

	s.history = append(s.history, um)
	if len(s.history) > s.maxHistory {
		s.history = s.history[len(s.history)-s.maxHistory:]
	}
}

func (s *Store) SetConnected(c bool) {
	s.mu.Lock()
	s.connected = c
	s.mu.Unlock()
}

func (s *Store) Snapshot(historyLimit int) (connected bool, lastSeen time.Time, lastTopic string, lastQoS byte, lastRetained bool, lastPayloadRaw string, lastPayload map[string]interface{}, tranca *bool, hist []model.UltimaMensagem) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if historyLimit > len(s.history) {
		historyLimit = len(s.history)
	}
	if historyLimit < 0 {
		historyLimit = 0
	}

	var historyOut []model.UltimaMensagem
	if historyLimit > 0 {
		historyOut = append(historyOut, s.history[len(s.history)-historyLimit:]...)
	}

	payloadCopy := make(map[string]interface{}, len(s.lastPayload))
	for k, v := range s.lastPayload {
		payloadCopy[k] = v
	}

	return s.connected, s.lastSeen, s.lastTopic, s.lastQoS, s.lastRetained, s.lastPayloadRaw, payloadCopy, s.trancaEstaAberta, historyOut
}
