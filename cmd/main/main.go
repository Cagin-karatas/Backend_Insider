package main

import (
	"log"
	"net/http"

	"github.com/cahitcaginkaratas/backend_insider/internal/database"
	"github.com/cahitcaginkaratas/backend_insider/internal/handlers"
	"github.com/cahitcaginkaratas/backend_insider/internal/models"
	"github.com/gorilla/mux"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Initialize database
	db := database.NewSQLiteDB()
	err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	// Check if teams already exist
	existingTeams, err := db.GetTeams()
	if err != nil {
		log.Fatal(err)
	}

	// Only create teams if none exist
	if len(existingTeams) == 0 {
		// Initialize teams
		teams := []*models.Team{
			models.NewTeam("Manchester City", 90),
			models.NewTeam("Liverpool", 85),
			models.NewTeam("Arsenal", 80),
			models.NewTeam("Chelsea", 75),
		}

		// Save teams to database
		for _, team := range teams {
			err = db.SaveTeam(team)
			if err != nil {
				log.Fatal(err)
			}
		}

		// Create league and generate fixtures
		league := models.NewLeague()
		for _, team := range teams {
			league.AddTeam(team)
		}
		league.GenerateFixtures()

		// Save matches to database
		for _, match := range league.Matches {
			err = db.SaveMatch(&match)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// Initialize API handler
	apiHandler := handlers.NewAPIHandler(db)

	// Set up router
	router := mux.NewRouter()
	router.Use(corsMiddleware)

	// API routes
	router.HandleFunc("/api/teams", apiHandler.GetTeams).Methods("GET")
	router.HandleFunc("/api/matches", apiHandler.GetMatches).Methods("GET")
	router.HandleFunc("/api/league", apiHandler.GetLeagueStats).Methods("GET")
	router.HandleFunc("/api/matches/simulate/{week}", apiHandler.SimulateWeek).Methods("POST")
	router.HandleFunc("/api/matches/simulate-all", apiHandler.SimulateAll).Methods("POST")
	router.HandleFunc("/api/matches/{id}", apiHandler.UpdateMatchResult).Methods("PUT")
	router.HandleFunc("/api/reset", apiHandler.ResetLeague).Methods("POST")

	// Start server
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
