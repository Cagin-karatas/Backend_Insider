package database

import (
	"github.com/cahitcaginkaratas/backend_insider/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Database interface defines the methods for database operations
type Database interface {
	InitDB() error
	GetTeams() ([]models.Team, error)
	GetMatches() ([]models.Match, error)
	GetLeagueStats() ([]models.TeamStats, error)
	SaveTeam(team *models.Team) error
	SaveMatch(match *models.Match) error
	UpdateMatch(match *models.Match) error
	GetMatchesByWeek(week int) ([]models.Match, error)
	ResetDatabase() error
}

// SQLiteDB implements the Database interface using SQLite
type SQLiteDB struct {
	db *gorm.DB
}

// NewSQLiteDB creates a new SQLite database instance
func NewSQLiteDB() *SQLiteDB {
	return &SQLiteDB{}
}

// InitDB initializes the database connection and creates tables
func (s *SQLiteDB) InitDB() error {
	db, err := gorm.Open(sqlite.Open("league.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	s.db = db

	// Auto migrate the schema
	err = s.db.AutoMigrate(&models.Team{}, &models.Match{})
	if err != nil {
		return err
	}

	return nil
}

// GetTeams returns all teams
func (s *SQLiteDB) GetTeams() ([]models.Team, error) {
	var teams []models.Team
	err := s.db.Find(&teams).Error
	return teams, err
}

// GetMatches returns all matches
func (s *SQLiteDB) GetMatches() ([]models.Match, error) {
	var matches []models.Match
	err := s.db.Preload("HomeTeam").Preload("AwayTeam").Find(&matches).Error
	return matches, err
}

// GetLeagueStats returns the current league statistics
func (s *SQLiteDB) GetLeagueStats() ([]models.TeamStats, error) {
	var stats []models.TeamStats
	// This is a simplified version - in a real system, you'd want to calculate this from matches
	err := s.db.Raw(`
		SELECT 
			t.id as team_id,
			t.name as team_name,
			COUNT(m.id) as played,
			SUM(CASE WHEN m.home_team_id = t.id AND m.home_goals > m.away_goals THEN 1
				WHEN m.away_team_id = t.id AND m.away_goals > m.home_goals THEN 1
				ELSE 0 END) as won,
			SUM(CASE WHEN m.home_goals = m.away_goals THEN 1 ELSE 0 END) as drawn,
			SUM(CASE WHEN m.home_team_id = t.id AND m.home_goals < m.away_goals THEN 1
				WHEN m.away_team_id = t.id AND m.away_goals < m.home_goals THEN 1
				ELSE 0 END) as lost,
			SUM(CASE WHEN m.home_team_id = t.id THEN m.home_goals ELSE m.away_goals END) as goals_for,
			SUM(CASE WHEN m.home_team_id = t.id THEN m.away_goals ELSE m.home_goals END) as goals_against,
			SUM(CASE WHEN m.home_team_id = t.id THEN m.home_goals - m.away_goals ELSE m.away_goals - m.home_goals END) as goal_difference,
			SUM(CASE 
				WHEN m.home_team_id = t.id AND m.home_goals > m.away_goals THEN 3
				WHEN m.away_team_id = t.id AND m.away_goals > m.home_goals THEN 3
				WHEN m.home_goals = m.away_goals THEN 1
				ELSE 0 END) as points
		FROM teams t
		LEFT JOIN matches m ON (m.home_team_id = t.id OR m.away_team_id = t.id) AND m.played = true
		GROUP BY t.id, t.name
		ORDER BY points DESC, goal_difference DESC, goals_for DESC
	`).Scan(&stats).Error
	return stats, err
}

// SaveTeam saves a team to the database
func (s *SQLiteDB) SaveTeam(team *models.Team) error {
	return s.db.Create(team).Error
}

// SaveMatch saves a match to the database
func (s *SQLiteDB) SaveMatch(match *models.Match) error {
	return s.db.Create(match).Error
}

// UpdateMatch updates a match in the database
func (s *SQLiteDB) UpdateMatch(match *models.Match) error {
	return s.db.Save(match).Error
}

// GetMatchesByWeek returns all matches for a specific week
func (s *SQLiteDB) GetMatchesByWeek(week int) ([]models.Match, error) {
	var matches []models.Match
	err := s.db.Preload("HomeTeam").Preload("AwayTeam").Where("week = ?", week).Find(&matches).Error
	return matches, err
}

// ResetDatabase clears all data and reinitializes the database
func (s *SQLiteDB) ResetDatabase() error {
	// Drop all tables
	err := s.db.Migrator().DropTable(&models.Match{}, &models.Team{})
	if err != nil {
		return err
	}

	// Recreate tables
	err = s.db.AutoMigrate(&models.Team{}, &models.Match{})
	if err != nil {
		return err
	}

	// Initialize teams
	teams := []*models.Team{
		models.NewTeam("Manchester City", 90),
		models.NewTeam("Liverpool", 85),
		models.NewTeam("Arsenal", 80),
		models.NewTeam("Chelsea", 75),
	}

	// Save teams to database
	for _, team := range teams {
		err = s.SaveTeam(team)
		if err != nil {
			return err
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
		err = s.SaveMatch(&match)
		if err != nil {
			return err
		}
	}

	return nil
}
