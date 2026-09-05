package main

import (
	"database/sql"
	"fmt"
)

// DBClient wraps our active SQL database connection pool
type DBClient struct {
	Conn *sql.DB
}

// NewDatabaseConnection handles opening a robust connection pool to SQL
func NewDatabaseConnection(dataSourceName string) (*DBClient, error) {
	// Opens database connection pool without external framework dependencies
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Verifies the network path to the database server is actively responsive
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database connection failed ping verification: %w", err)
	}

	return &DBClient{Conn: db}, nil
}

// StreamDatabaseRecords queries a data table and pipes records straight into our pipeline channel
func (client *DBClient) StreamDatabaseRecords(outputChan chan<- UserRecord) error {
	query := "SELECT id, name, email, age, state FROM users"
	
	rows, err := client.Conn.Query(query)
	if err != nil {
		return fmt.Errorf("failed to execute data selection query: %w", err)
	}
	defer rows.Close()

	// Ingests rows sequentially, converting them to type-safe Go structs on the fly
	for rows.Next() {
		var user UserRecord
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Age, &user.State)
		if err != nil {
			continue // Skip single malformed database rows to protect stream health
		}
		
		// Push the database record instantly into the shared concurrent pipeline channel
		outputChan <- user
	}

	return rows.Err()
}
