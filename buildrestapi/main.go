package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	_ "github.com/lib/pq"
  "github.com/gorilla/mux"
)

type Article struct {
    ID string `json:"ID"`
    Title string `json:"Title"`
    Abstract string `json:"Abstract"`
    PublishedDate string `json:"PublishedDate"`
    URL string `json:"URL"`
}


const (
  host     = "nyt-db.cdfnbhvkuejj.us-east-2.rds.amazonaws.com"
	port     = 5432
	user     = "username"
	password = "password"
	dbname   = "mydb"
)


func OpenConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

func GETHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

  vars := mux.Vars(r)
  key := vars["query"]
  fmt.Println(key)
  str1 := "%"
  str2 := "%"
  full_key := str1+key+str2
	rows, err := db.Query(`SELECT *From AllArticles WHERE Abstract LIKE $1`, full_key)
	if err != nil {
		panic(err)
	}

	var articles []Article

	for rows.Next() {
		var article Article
		rows.Scan(&article.ID, &article.Title, &article.Abstract, &article.PublishedDate, &article.URL)
		articles = append(articles, article)
	}

	articlesBytes, _ := json.MarshalIndent(articles, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(articlesBytes)

	defer rows.Close()
	defer db.Close()
}



func PUTHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

  vars := mux.Vars(r)
  key := vars["id"]
  fmt.Println(key)


	var article Article

	err := json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement :=`UPDATE AllArticles SET title = $1, abstract = $2, publisheddate = $3, url = $4 WHERE id = $5`
  _, err = db.Exec(sqlStatement, article.Title, article.Abstract, article.PublishedDate, article.URL, key)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	defer db.Close()
}



func POSTHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

	var article Article

	err := json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `INSERT INTO AllArticles (title, abstract, publisheddate, url) VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(sqlStatement, article.Title, article.Abstract, article.PublishedDate, article.URL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	defer db.Close()
}


func main() {
  router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/article/{query}", GETHandler).Methods("GET")
  router.HandleFunc("/article/{id}", PUTHandler).Methods("PUT")
  router.HandleFunc("/article/insert", POSTHandler)
  //https://stackoverflow.com/questions/14081066/gae-golang-gorilla-mux-404-page-not-found
  http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":10000", nil))
}
