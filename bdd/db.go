package bdd

import(
	"database/sql"
	"log"
	// Driver pour BDD MySQL : blank import --> https://www.calhoun.io/why-we-import-sql-drivers-with-the-blank-identifier/
	_ "github.com/go-sql-driver/mysql"
)

// Variable globale au package bdd, pilote de communication avec la BDD
var	db *sql.DB

// InitDB : Ouvre une connexion avec la BDD avec 
func InitDB(username string){

	var err error
	db, err = sql.Open("mysql", username + "@tcp(127.0.0.1:3306)/gitus")
	if err != nil {
        log.Panic(err)
	}
	if db != nil{
		log.Println("Connection Established with database gitus")
	}else{
		log.Println("ERROR Connection Established with database gitus")
	}
}

// CloseDB : Ferme la connexion avec la BDD
func CloseDB(){
	db.Close()
}