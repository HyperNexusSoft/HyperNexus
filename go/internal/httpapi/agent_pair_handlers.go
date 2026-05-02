package httpapi

import (
	"encoding/json"
	"net/http"
)

func (s *Server) handleAgentPairRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"success": false, "error": "method not allowed"})
		return
	}

	var payload struct {
		Task string `json:"task"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "error": "invalid JSON body"})
		return
	}

	// Try upstream first
	var result any
	upstreamBase, err := s.callUpstreamJSON(r.Context(), "agent.runPairSession", payload, &result)
	if err == nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"data":    result,
			"bridge": map[string]any{
				"upstreamBase": upstreamBase,
				"procedure":    "agent.runPairSession",
			},
		})
		return
	}

	// Fallback to local Go pair orchestrator
	pairResult, fallbackErr := s.pairOrchestrator.RunTask(r.Context(), payload.Task)
	if fallbackErr != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   fallbackErr.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    pairResult,
		"bridge": map[string]any{
			"fallback":  "go-local-pair-orchestrator",
			"procedure": "agent.runPairSession",
		},
	})
}
