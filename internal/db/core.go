package db

import (
	"DC_NewsSender/internal/db/models"

	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ConnectionFormat string = "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable"
)

type postgresDatabase struct {
	Configuration *PostgresDatabaseConfiguration
	Connection    *gorm.DB
}

type PostgresDatabaseConfiguration struct {
	Host         string
	UserName     string
	UserPassword string
	DatabaseName string
	Port         uint16
	ServerPort   uint16
}

// getConnectionString returns connection string from configuration.
func (dbConfig *PostgresDatabaseConfiguration) getConnectionString() string {
	return fmt.Sprintf(ConnectionFormat,
		dbConfig.Host,
		dbConfig.UserName,
		dbConfig.UserPassword,
		dbConfig.DatabaseName,
		dbConfig.Port)
}

// Creates new gorm connection and adds to connections pool.
func (db *postgresDatabase) newConnection() error {
	orm, err := gorm.Open(postgres.Open(db.Configuration.getConnectionString()), &gorm.Config{})
	if err != nil {
		return err
	}

	db.Connection = orm

	return err
}

// InitializePostgresDatabase initializes database.
func InitializePostgresDatabase(config *PostgresDatabaseConfiguration) (*gorm.DB, error) {
	database := &postgresDatabase{
		Configuration: config,
	}

	if err := database.newConnection(); err != nil {
		return nil, err
	}

	if err := database.migrateDatabase(); err != nil {
		return nil, err
	}

	return database.Connection, nil
}

// CleanupConnection closes active connections.
// Should be called with defer in main thread, either should be executed on application close.
func CleanupConnection(orm *gorm.DB) {
	conn, err := orm.DB()
	if err == nil {
		conn.Close()
	}
}

// migrateDatabase сreates database structure
func (db *postgresDatabase) migrateDatabase() error {
	if err := db.Connection.AutoMigrate(&models.Admin{}); err != nil {
		return err
	}
	if err := db.Connection.AutoMigrate(&models.Language{}); err != nil {
		return err
	}
	if err := db.Connection.AutoMigrate(&models.Group{}); err != nil {
		return err
	}
	if err := db.Connection.AutoMigrate(&models.Chat{}); err != nil {
		return err
	}

	return nil
}
