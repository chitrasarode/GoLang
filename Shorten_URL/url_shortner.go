package main

import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type Url struct {
	Long_Url  string `json:"long_url"`
	Short_Url string `json:"short_url"`
}

var db *sql.DB

func main() {
	// Open a connection to the SQLite database (creates the database if it doesn't exist)
	db, _ = sql.Open("sqlite3", "./ShortenURL1.db")

	defer db.Close()

	createTable(db)

	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/expand", expandHandler)
	http.HandleFunc("/display", displayHandler)

	fmt.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func createTable(db *sql.DB) {
	createTableSQl :=
		`CREATE TABLE IF NOT EXISTS url (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		longUrl TEXT NOT NULL,
		shortUrl TEXT NOT NULL
	);`
	_, err := db.Exec(createTableSQl)
	if err != nil {
		log.Fatal(err)
	}
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {

	var input Url
	var short Url
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	short.Short_Url = shorten(input.Long_Url)
	fmt.Println("short url : ", short.Short_Url)
	fmt.Println(input.Long_Url, "      ", short.Short_Url)

	insertUrl(db, input.Long_Url, short.Short_Url)

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response
	json.NewEncoder(w).Encode(short)

}

func shorten(input string) string {
	hasher := sha1.New()
	hasher.Write([]byte(input))
	hasherValue := hasher.Sum(nil)
	//fmt.Println(hasherValue)
	hasherString := hex.EncodeToString(hasherValue)
	return hasherString[:12]
}

func insertUrl(db *sql.DB, long_Url, short_Url string) {
	fmt.Println("inserting data.......")
	_, err := db.Exec("INSERT INTO url (longUrl, shortUrl) VALUES (?,?)", long_Url, short_Url)

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("data inserted successfully")
	}
}

func expandHandler(w http.ResponseWriter, r *http.Request) {
	var input Url
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	longURL, found := expand(input.Short_Url)
	if !found {
		fmt.Fprintln(w, "Long url not found")
		return
	}
	fmt.Fprintln(w, "long url : ", longURL)

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response
	json.NewEncoder(w).Encode(Url{Long_Url: longURL, Short_Url: input.Short_Url})
}

func expand(shortURL string) (string, bool) {

	var longURL string
	err := db.QueryRow(`SELECT longUrl FROM url WHERE shortUrl = ?`, shortURL).Scan(&longURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", false
		}
		log.Fatal(err)
		return "", false
	}

	return longURL, true
}

func displayHandler(w http.ResponseWriter, r *http.Request) {
	var urls []Url

	rows, err := db.Query(`SELECT longUrl, shortUrl FROM url`)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var u Url
		err := rows.Scan(&u.Long_Url, &u.Short_Url)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
		urls = append(urls, u)
	}

	err = rows.Err()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response
	json.NewEncoder(w).Encode(urls)
}
