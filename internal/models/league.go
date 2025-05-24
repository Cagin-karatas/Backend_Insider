package models

import (
	"sort"
)

// League represents the football league
type League struct {
	Teams   []Team      `json:"teams"`
	Matches []Match     `json:"matches"`
	Stats   []TeamStats `json:"stats"`
}

// NewLeague creates a new league instance
func NewLeague() *League {
	return &League{
		Teams:   make([]Team, 0),
		Matches: make([]Match, 0),
		Stats:   make([]TeamStats, 0),
	}
}

// AddTeam adds a team to the league
func (l *League) AddTeam(team *Team) {
	l.Teams = append(l.Teams, *team)
	l.Stats = append(l.Stats, TeamStats{
		TeamID:   team.ID,
		TeamName: team.Name,
	})
}

// GenerateFixtures generates all matches for the league
func (l *League) GenerateFixtures() {
	numTeams := len(l.Teams)
	weeks := (numTeams - 1) * 2
	matchesPerWeek := numTeams / 2

	for week := 1; week <= weeks; week++ {
		for i := 0; i < matchesPerWeek; i++ {
			homeIdx := (week + i) % (numTeams - 1)
			awayIdx := (numTeams - 1 - i) % (numTeams - 1)

			if i == 0 {
				awayIdx = numTeams - 1
			}

			// Swap home and away for second half of season
			if week > weeks/2 {
				homeIdx, awayIdx = awayIdx, homeIdx
			}

			match := NewMatch(week, &l.Teams[homeIdx], &l.Teams[awayIdx])
			l.Matches = append(l.Matches, *match)
		}
	}
}

// UpdateStats updates the league statistics based on match results
func (l *League) UpdateStats() {
	// Reset stats
	for i := range l.Stats {
		l.Stats[i] = TeamStats{
			TeamID:   l.Teams[i].ID,
			TeamName: l.Teams[i].Name,
		}
	}

	// Update stats based on matches
	for _, match := range l.Matches {
		if !match.Played {
			continue
		}

		// Update home team stats
		homeStats := &l.Stats[match.HomeTeamID-1]
		homeStats.Played++
		homeStats.GoalsFor += match.HomeGoals
		homeStats.GoalsAgainst += match.AwayGoals

		// Update away team stats
		awayStats := &l.Stats[match.AwayTeamID-1]
		awayStats.Played++
		awayStats.GoalsFor += match.AwayGoals
		awayStats.GoalsAgainst += match.HomeGoals

		// Update points and results
		if match.HomeGoals > match.AwayGoals {
			homeStats.Won++
			homeStats.Points += 3
			awayStats.Lost++
		} else if match.HomeGoals < match.AwayGoals {
			awayStats.Won++
			awayStats.Points += 3
			homeStats.Lost++
		} else {
			homeStats.Drawn++
			awayStats.Drawn++
			homeStats.Points++
			awayStats.Points++
		}
	}

	// Calculate goal differences and sort
	for i := range l.Stats {
		l.Stats[i].GoalDifference = l.Stats[i].GoalsFor - l.Stats[i].GoalsAgainst
	}

	// Sort by points, then goal difference, then goals scored
	sort.Slice(l.Stats, func(i, j int) bool {
		if l.Stats[i].Points != l.Stats[j].Points {
			return l.Stats[i].Points > l.Stats[j].Points
		}
		if l.Stats[i].GoalDifference != l.Stats[j].GoalDifference {
			return l.Stats[i].GoalDifference > l.Stats[j].GoalDifference
		}
		return l.Stats[i].GoalsFor > l.Stats[j].GoalsFor
	})
}

// GetMatchesByWeek returns all matches for a specific week
func (l *League) GetMatchesByWeek(week int) []Match {
	var weekMatches []Match
	for _, match := range l.Matches {
		if match.Week == week {
			weekMatches = append(weekMatches, match)
		}
	}
	return weekMatches
}

// SimulateWeek simulates all matches for a specific week
func (l *League) SimulateWeek(week int) {
	matches := l.GetMatchesByWeek(week)
	for i := range matches {
		homeTeam := &l.Teams[matches[i].HomeTeamID-1]
		awayTeam := &l.Teams[matches[i].AwayTeamID-1]
		matches[i].Simulate(homeTeam, awayTeam)
	}
	l.UpdateStats()
}

// SimulateAll simulates all remaining matches
func (l *League) SimulateAll() {
	maxWeek := 0
	for _, match := range l.Matches {
		if match.Week > maxWeek {
			maxWeek = match.Week
		}
	}

	for week := 1; week <= maxWeek; week++ {
		l.SimulateWeek(week)
	}
}
