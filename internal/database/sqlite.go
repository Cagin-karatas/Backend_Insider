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

// ... rest of the existing code ...
