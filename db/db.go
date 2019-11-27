package db

import(
	"database/sql"
	// Driver pour BDD MySQL : blank import --> https://www.calhoun.io/why-we-import-sql-drivers-with-the-blank-identifier/
	_ "github.com/mattn/go-sqlite3"
)

// Variable globale au package bdd, pilote de communication avec la BDD
var	db *sql.DB

// InitDB : Ouvre une connexion avec la BDD avec 
func InitDB()error{

	var err error
	db, err = sql.Open("sqlite3", "C:\\sqlite\\gitus.db")
	err = db.Ping()
	if err != nil {
        return err
	}
	return nil
}

// CloseDB : Ferme la connexion avec la BDD
func CloseDB(){
	db.Close()
}