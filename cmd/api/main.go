package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"github.com/permalik/utility/db"
	"github.com/permalik/utility/models"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load .env", err)
	}

	pool := db.InitDB()
	defer pool.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	db.Ping(ctx)

	selectRepos(pool, ctx)
}

func selectRepos(pool *sql.DB, ctx context.Context) ([]byte, error) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	rows, err := pool.QueryContext(ctx, "select * from repos;")
	if err != nil {
		log.Fatal("unable to execute select all", err)
	}
	defer rows.Close()

	var repos []models.JsonRepo
	for rows.Next() {
		var (
			id          int
			owner       string
			name        string
			category    string
			description string
			htmlURL     string
			homepage    string
			topics      string
			createdAt   string
			updatedAt   string
			uid         int
		)
		if err := rows.Scan(
			&id,
			&owner,
			&name,
			&category,
			&description,
			&htmlURL,
			&homepage,
			&topics,
			&createdAt,
			&updatedAt,
			&uid); err != nil {
			log.Fatal(err)
		}
		repo := models.JsonRepo{
			Owner:       owner,
			Name:        name,
			Category:    category,
			Description: description,
			HTMLURL:     htmlURL,
			Homepage:    homepage,
			Topics:      topics,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
			UID:         uid,
		}
		repos = append(repos, repo)
	}

	rerr := rows.Close()
	if rerr != nil {
		log.Fatal(rerr)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	jsonData, err := json.Marshal(repos)
	if err != nil {
		log.Fatalf("error marshaling to json:\n%v", err)
	}
	fmt.Println(string(jsonData))
	return jsonData, nil
}