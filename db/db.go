package db

import(
	"database/sql"
	// Driver pour BDD MySQL : blank import --> https://www.calhoun.io/why-we-import-sql-drivers-with-the-blank-identifier/
	_ "github.com/go-sql-driver/mysql"
)

// Variable globale au package bdd, pilote de communication avec la BDD
var	db *sql.DB

// InitDB : Ouvre une connexion avec la BDD avec 
func InitDB(username string)error{

	var err error
	db, err = sql.Open("mysql", username + "@tcp(127.0.0.1:3306)/gitus?parseTime=true")
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