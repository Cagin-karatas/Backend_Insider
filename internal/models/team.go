package models

import (
	"time"
)

// Team represents a football team in the league
type Team struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Strength  int       `json:"strength"` // 1-100 scale for team strength
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TeamStats represents the statistics for a team in the league
type TeamStats struct {
	TeamID         uint   `json:"team_id" gorm:"primaryKey"`
	TeamName       string `json:"team_name"`
	Played         int    `json:"played"`
	Won            int    `json:"won"`
	Drawn          int    `json:"drawn"`
	Lost           int    `json:"lost"`
	GoalsFor       int    `json:"goals_for"`
	GoalsAgainst   int    `json:"goals_against"`
	GoalDifference int    `json:"goal_difference"`
	Points         int    `json:"points"`
}

// NewTeam creates a new team instance
func NewTeam(name string, strength int) *Team {
	return &Team{
		Name:     name,
		Strength: strength,
	}
}

// CalculateWinProbability calculates the probability of winning against another team
func (t *Team) CalculateWinProbability(opponent *Team) float64 {
	totalStrength := float64(t.Strength + opponent.Strength)
	return float64(t.Strength) / totalStrength
}
