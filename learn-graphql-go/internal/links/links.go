package links

import (
	"fmt"
	database "learn-graphql-go/internal/pkg/db/migrations/psql"
	"learn-graphql-go/internal/users"
	"log"
	"strconv"
)

// #1
type Link struct {
	ID      string
	Title   string
	Address string
	User    *users.User
}

// #2
func (link Link) Save() (int64, error) {
	// Prepare the INSERT statement with a RETURNING clause for PostgreSQL
	stmt, err := database.Db.Prepare(`
		INSERT INTO Links (Title, Address, UserID)
		VALUES ($1, $2, $3)
		RETURNING id
	`)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	// Execute the statement with the provided link data
	var id int64

	userID, err := strconv.Atoi(link.User.ID)
	if err != nil {
		return 0, fmt.Errorf("invalid user ID: %v", err)
	}

	err = stmt.QueryRow(link.Title, link.Address, userID).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %v", err)
	}

	log.Print("Row inserted!")
	return id, nil
}

func GetAll() ([]Link, error) {
	rows, err := database.Db.Query("SELECT id, title, address FROM Links")
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	var links []Link

	for rows.Next() {
		var link Link
		if err := rows.Scan(&link.ID, &link.Title, &link.Address); err != nil {
			return nil, fmt.Errorf("row scan error: %v", err)
		}
		links = append(links, link)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return links, nil
}
