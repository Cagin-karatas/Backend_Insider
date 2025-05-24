package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cahitcaginkaratas/backend_insider/internal/database"
	"github.com/cahitcaginkaratas/backend_insider/internal/models"
	"github.com/gorilla/mux"
)

// APIHandler handles all API requests
type APIHandler struct {
	db database.Database
}

// NewAPIHandler creates a new API handler
func NewAPIHandler(db database.Database) *APIHandler {
	return &APIHandler{db: db}
}

// GetTeams returns all teams
func (h *APIHandler) GetTeams(w http.ResponseWriter, r *http.Request) {
	teams, err := h.db.GetTeams()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(teams)
}

// GetMatches returns all matches
func (h *APIHandler) GetMatches(w http.ResponseWriter, r *http.Request) {
	matches, err := h.db.GetMatches()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(matches)
}

// GetLeagueStats returns the current league table
func (h *APIHandler) GetLeagueStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.db.GetLeagueStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(stats)
}

// SimulateWeek simulates all matches for a specific week
func (h *APIHandler) SimulateWeek(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	week, err := strconv.Atoi(vars["week"])
	if err != nil {
		http.Error(w, "Invalid week number", http.StatusBadRequest)
		return
	}

	matches, err := h.db.GetMatchesByWeek(week)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	teams, err := h.db.GetTeams()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a map of teams for quick lookup
	teamMap := make(map[uint]*models.Team)
	for i := range teams {
		teamMap[teams[i].ID] = &teams[i]
	}

	// Simulate matches
	for i := range matches {
		homeTeam := teamMap[matches[i].HomeTeamID]
		awayTeam := teamMap[matches[i].AwayTeamID]
		matches[i].Simulate(homeTeam, awayTeam)
		err = h.db.UpdateMatch(&matches[i])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	json.NewEncoder(w).Encode(matches)
}

// SimulateAll simulates all remaining matches
func (h *APIHandler) SimulateAll(w http.ResponseWriter, r *http.Request) {
	matches, err := h.db.GetMatches()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	teams, err := h.db.GetTeams()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a map of teams for quick lookup
	teamMap := make(map[uint]*models.Team)
	for i := range teams {
		teamMap[teams[i].ID] = &teams[i]
	}

	// Find max week
	maxWeek := 0
	for _, match := range matches {
		if match.Week > maxWeek {
			maxWeek = match.Week
		}
	}

	// Simulate all weeks
	for week := 1; week <= maxWeek; week++ {
		weekMatches, err := h.db.GetMatchesByWeek(week)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for i := range weekMatches {
			homeTeam := teamMap[weekMatches[i].HomeTeamID]
			awayTeam := teamMap[weekMatches[i].AwayTeamID]
			weekMatches[i].Simulate(homeTeam, awayTeam)
			err = h.db.UpdateMatch(&weekMatches[i])
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	json.NewEncoder(w).Encode(matches)
}

// UpdateMatchResult updates a match result manually
func (h *APIHandler) UpdateMatchResult(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	matchID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid match ID", http.StatusBadRequest)
		return
	}

	var result models.MatchResult
	err = json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	matches, err := h.db.GetMatches()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Find the match
	var match *models.Match
	for i := range matches {
		if matches[i].ID == uint(matchID) {
			match = &matches[i]
			break
		}
	}

	if match == nil {
		http.Error(w, "Match not found", http.StatusNotFound)
		return
	}

	// Update the match result
	match.UpdateResult(result.HomeGoals, result.AwayGoals)
	err = h.db.UpdateMatch(match)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(match)
}

// ResetLeague resets the database and reinitializes the league
func (h *APIHandler) ResetLeague(w http.ResponseWriter, r *http.Request) {
	err := h.db.ResetDatabase()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get updated data
	stats, err := h.db.GetLeagueStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	matches, err := h.db.GetMatches()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Stats   []models.TeamStats `json:"stats"`
		Matches []models.Match     `json:"matches"`
	}{
		Stats:   stats,
		Matches: matches,
	}

	json.NewEncoder(w).Encode(response)
}
