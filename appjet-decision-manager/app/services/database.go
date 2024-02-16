package services

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Configuration struct {
	Database struct {
		Driver    string `json:"driver"`
		Username  string `json:"username"`
		Password  string `json:"password"`
		Host      string `json:"host"`
		Port      int    `json:"port"`
		DBName    string `json:"dbname"`
		Charset   string `json:"charset"`
		ParseTime bool   `json:"parseTime"`
		Loc       string `json:"loc"`
	} `json:"database"`
	Dependencies map[string]string `json:"dependencies"`
}

var db *gorm.DB

// GetHardcodedConfiguration returns a hardcoded Configuration
func GetHardcodedConfiguration() Configuration {
	return Configuration{
		Database: struct {
			Driver    string `json:"driver"`
			Username  string `json:"username"`
			Password  string `json:"password"`
			Host      string `json:"host"`
			Port      int    `json:"port"`
			DBName    string `json:"dbname"`
			Charset   string `json:"charset"`
			ParseTime bool   `json:"parseTime"`
			Loc       string `json:"loc"`
		}{
			Driver:    "mysql",
			Username:  "root",
			Password:  "root",
			Host:      "database_appjet",
			Port:      3306,
			DBName:    "app-db",
			Charset:   "utf8mb4",
			ParseTime: true,
			Loc:       "Local",
		},
		Dependencies: map[string]string{},
	}
}

func CreateDbConnection() (*gorm.DB, error) {
	config := GetHardcodedConfiguration()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		config.Database.Username,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.DBName,
		config.Database.Charset,
		config.Database.ParseTime,
		config.Database.Loc,
	)

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}) // Assign to package-level db variable
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	return db, nil
}

// CloseDbConnection closes the database connection
func CloseDbConnection() error {
	if db == nil {
		return nil // Database not initialized
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("error getting underlying database connection: %w", err)
	}

	err = sqlDB.Close()
	if err != nil {
		return fmt.Errorf("error closing the database connection: %w", err)
	}

	return nil
}

func GetDBConnection() *gorm.DB {
	return db
}
