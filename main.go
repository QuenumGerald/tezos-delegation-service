package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

// Sender struct pour contenir les informations sur l'expéditeur
type Sender struct {
	Address string `json:"address"`
}

// Delegation struct pour contenir les informations sur la délégation
type Delegation struct {
	Timestamp string `json:"timestamp"`
	Amount    int64  `json:"amount"`
	Delegator string `json:"delegator"`
	Block     string `json:"block"`
}

// initDB initialise la base de données
func initDB() *sql.DB {
	database, err := sql.Open("sqlite3", "./delegations.db")
	if err != nil {
		log.Fatal(err)
	}
	// Crée la table "delegations" si elle n'existe pas encore
	statement, err := database.Prepare(`CREATE TABLE IF NOT EXISTS delegations (
		timestamp TEXT,
		amount INT64,
		delegator TEXT,
		block TEXT,
		PRIMARY KEY (timestamp, delegator))`)
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()

	return database
}

// fetchDelegations récupère les délégations à partir d'une API et les stocke dans la base de données
func fetchDelegations(db *sql.DB) {
	for {
		url := "https://api.tzkt.io/v1/operations/delegations?timestamp.gt=2023-09-23T00:00:00Z&limit=10000"
		resp, err := http.Get(url)
		if err != nil {
			log.Println("Error fetching data:", err)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error reading body:", err)
		}
		resp.Body.Close()
		// Décodage du JSON
		var fetched []struct {
			Timestamp string `json:"timestamp"`
			Amount    int64  `json:"amount"`
			Sender    Sender `json:"sender"`
			Block     string `json:"block"`
		}
		if err := json.Unmarshal(body, &fetched); err != nil {
			log.Println("Error unmarshaling:", err)
			continue
		}
		// Insérer les nouvelles délégations dans la base de données
		for _, f := range fetched {
			_, err := db.Exec("INSERT OR IGNORE INTO delegations (timestamp, amount, delegator, block) VALUES (?, ?, ?, ?)",
				f.Timestamp, f.Amount, f.Sender.Address, f.Block)
			if err != nil {
				log.Println("Database insert error:", err)
			}
		}
	}
}

// getDelegations récupère les délégations stockées et les renvoie sous forme de JSON

func getDelegations(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	rows, err := db.Query("SELECT * FROM delegations ORDER BY timestamp DESC")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var delegations []Delegation
	for rows.Next() {
		var d Delegation
		if err := rows.Scan(&d.Timestamp, &d.Amount, &d.Delegator, &d.Block); err != nil {
			log.Fatal(err)
		}
		delegations = append(delegations, d)
	}

	json.NewEncoder(w).Encode(delegations)
}

func main() {
	db := initDB()
	defer db.Close()

	// Définition des routes
	r := mux.NewRouter()
	r.HandleFunc("/xtz/delegations", func(w http.ResponseWriter, r *http.Request) {
		getDelegations(w, r, db)
	}).Methods("GET")

	// Récupération des délégations dans une goroutine
	go fetchDelegations(db)

	fmt.Println("Server is running on port 8000")
	http.Handle("/", r)
	http.ListenAndServe(":8000", nil)
}
