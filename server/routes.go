package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Routes struct {
	manager *AgentManager
}

func (ro *Routes) ListAgents(w http.ResponseWriter, r *http.Request) {

	// Set content type header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(fmt.Sprintf("agents connected %d", ro.manager.AgentsConnnected())))
}

func (ro *Routes) SendMessage(w http.ResponseWriter, r *http.Request) {
	// Decode the request body

	type Message struct {
		Message   string `json:"message"`
		AgentID   string `json:"agent_id"`
		Namespace string `json:"namespace"`
	}

	var message Message
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Send the message to the agent
	res, err := ro.manager.SendMessage(message.AgentID, message.Message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set headers and encode JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func NewRoutes(manager *AgentManager) *Routes {
	return &Routes{
		manager: manager,
	}
}
