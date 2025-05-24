package models

import (
	"math/rand"
	"time"
)

// Match represents a football match between two teams
type Match struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Week       int       `json:"week"`
	HomeTeam   Team      `json:"home_team" gorm:"foreignKey:HomeTeamID"`
	HomeTeamID uint      `json:"home_team_id"`
	AwayTeam   Team      `json:"away_team" gorm:"foreignKey:AwayTeamID"`
	AwayTeamID uint      `json:"away_team_id"`
	HomeGoals  int       `json:"home_goals"`
	AwayGoals  int       `json:"away_goals"`
	Played     bool      `json:"played"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// MatchResult represents the result of a match
type MatchResult struct {
	HomeTeamID uint `json:"home_team_id"`
	AwayTeamID uint `json:"away_team_id"`
	HomeGoals  int  `json:"home_goals"`
	AwayGoals  int  `json:"away_goals"`
}

// NewMatch creates a new match instance
func NewMatch(week int, homeTeam, awayTeam *Team) *Match {
	return &Match{
		Week:       week,
		HomeTeamID: homeTeam.ID,
		AwayTeamID: awayTeam.ID,
		Played:     false,
	}
}

// Simulate simulates the match result based on team strengths
func (m *Match) Simulate(homeTeam, awayTeam *Team) {
	// Calculate base probabilities
	homeWinProb := homeTeam.CalculateWinProbability(awayTeam)
	awayWinProb := awayTeam.CalculateWinProbability(homeTeam)
	drawProb := 1.0 - homeWinProb - awayWinProb

	// Add some randomness to make it more interesting
	homeAdvantage := 0.1 // 10% home advantage
	homeWinProb += homeAdvantage
	awayWinProb -= homeAdvantage / 2
	drawProb -= homeAdvantage / 2

	// Generate random result
	// This is a simplified simulation - in a real system, you'd want more sophisticated logic
	random := rand.Float64()

	switch {
	case random < homeWinProb:
		m.HomeGoals = rand.Intn(4) + 1
		m.AwayGoals = rand.Intn(m.HomeGoals)
	case random < homeWinProb+awayWinProb:
		m.AwayGoals = rand.Intn(4) + 1
		m.HomeGoals = rand.Intn(m.AwayGoals)
	default:
		m.HomeGoals = rand.Intn(3)
		m.AwayGoals = m.HomeGoals
	}

	m.Played = true
}

// UpdateResult updates the match result manually
func (m *Match) UpdateResult(homeGoals, awayGoals int) {
	m.HomeGoals = homeGoals
	m.AwayGoals = awayGoals
	m.Played = true
}
